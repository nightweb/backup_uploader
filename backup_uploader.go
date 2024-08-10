package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	debug          bool
	syncFolderPath string
	credPath       string
	driveId        string
	direction      string
	fileMask       string
	mkdir          bool
	folderPath     string
	uploadFilePath string
)

func init() {
	// Add flags for verbose output
	flag.BoolVar(&debug, "vv", false, "Enable verbose output (debug mode)")
	flag.BoolVar(&debug, "debug", false, "Enable verbose output (debug mode)")
	flag.StringVar(&syncFolderPath, "sync_folder", "", "Path to the local folder to sync")
	flag.StringVar(&credPath, "c",
		filepath.Join(os.Getenv("HOME"), ".backup_uploader", "google", "credentials.json"),
		"Path to credentials.json file")
	flag.StringVar(&driveId, "drive_id", "", "Google Drive ID (for Shared Drives)")
	flag.StringVar(&direction, "direction", "to_drive", "Sync direction: to_drive (default) or from_drive")
	flag.StringVar(&fileMask, "file_mask", "*", "File mask to filter files for syncing")
	flag.BoolVar(&mkdir, "mkdir", false, "Create directories if they do not exist")
	flag.StringVar(&folderPath, "folder_path", "", "Path to the target folder on Google Drive")
	flag.StringVar(&uploadFilePath, "upload_file", "", "Path to a single file to upload to Google Drive")
}

// Returns the MD5 hash of a file
// @param filePath: Path to the file
// @return: MD5 hash as a string
func computeMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Function to find or create a folder by path on Google Drive
// @param srv: Google Drive service instance
// @param rootFolderID: ID of the root folder to start search from
// @param folderPath: Path to the folder to find or create
// @param mkdir: Create missing folders if true
// @return: ID of the found or created folder
func findOrCreateFolderByPath(srv *drive.Service, rootFolderID, folderPath string, mkdir bool) (string, error) {
	folders := strings.Split(folderPath, "/")
	currentFolderID := rootFolderID

	for _, folder := range folders {
		query := fmt.Sprintf("name = '%s' and '%s' in parents and mimeType = 'application/vnd.google-apps.folder' and trashed=false", folder, currentFolderID)
		r, err := srv.Files.List().
			Q(query).
			Fields("files(id, name)").
			SupportsAllDrives(true).
			IncludeItemsFromAllDrives(true).
			Do()

		if err != nil {
			return "", fmt.Errorf("error finding folder '%s': %v", folder, err)
		}

		if len(r.Files) == 0 {
			if mkdir {
				// Создаем папку, если не найдена
				newFolder := &drive.File{
					Name:     folder,
					Parents:  []string{currentFolderID},
					MimeType: "application/vnd.google-apps.folder",
				}
				createdFolder, err := srv.Files.Create(newFolder).SupportsAllDrives(true).Fields("id").Do()
				if err != nil {
					return "", fmt.Errorf("error creating folder '%s': %v", folder, err)
				}
				currentFolderID = createdFolder.Id
				fmt.Printf("Folder '%s' created with ID: %s\n", folder, currentFolderID)
			} else {
				return "", fmt.Errorf("folder '%s' not found and creation not allowed (use -mkdir to create missing folders)", folder)
			}
		} else {
			currentFolderID = r.Files[0].Id
		}
	}

	return currentFolderID, nil
}

// Function to download or update a file in a specified folder on Google Drive
// @param srv: Google Drive service instance
// @param folderID: ID of the target folder on Google Drive
// @param filePath: Path to the file to upload
// @param driveFile: File metadata from Google Drive (nil if file does not exist)
// @return: Error if any
func uploadFile(srv *drive.Service, folderID, filePath string, driveFile *drive.File) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to get file info: %v", err)
	}

	if driveFile != nil {
		// Если файл уже существует на Google Drive, обновляем его содержимое
		fmt.Printf("Updating file '%s' on Google Drive...\n", fileInfo.Name())
		// Обновляем только содержимое файла, без изменения других полей
		updatedFile, err := srv.Files.Update(driveFile.Id, &drive.File{}).Media(file).SupportsAllDrives(true).Do()
		if err != nil {
			return fmt.Errorf("unable to update file: %v", err)
		}
		fmt.Printf("File '%s' successfully updated. File ID: %s\n", updatedFile.Name, updatedFile.Id)
	} else {
		// Если файл не существует на Google Drive, создаем его
		fmt.Printf("Uploading new file '%s' to Google Drive...\n", fileInfo.Name())
		driveFile := &drive.File{
			Name:    fileInfo.Name(),
			Parents: []string{folderID},
		}
		uploadedFile, err := srv.Files.Create(driveFile).Media(file).SupportsAllDrives(true).Do()
		if err != nil {
			return fmt.Errorf("unable to upload file: %v", err)
		}
		fmt.Printf("File '%s' successfully uploaded. File ID: %s\n", uploadedFile.Name, uploadedFile.Id)
	}

	return nil
}

// Function to synchronize a local folder with a folder on Google Drive
// @param srv: Google Drive service instance
// @param localDir: Path to the local folder to sync
// @param remoteFolderID: ID of the target folder on Google Drive
// @param direction: Sync direction: "to_drive" or "from_drive"
// @param fileMask: File mask to filter files for syncing
// @param mkdir: Create missing folders if true
// @return: Error if any
func syncFolder(srv *drive.Service, localDir, remoteFolderID, direction, fileMask string, mkdir bool) error {
	localFilesMap := make(map[string]string)
	driveFilesMap, err := getAllFilesFromDrive(srv, remoteFolderID, "")
	if err != nil {
		return fmt.Errorf("unable to get files from Google Drive: %v", err)
	}

	// Обход всех файлов в локальной директории
	err = filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relativePath := strings.TrimPrefix(path, localDir)
			relativePath = strings.TrimPrefix(relativePath, string(filepath.Separator))

			// Проверка на соответствие маске файла
			if fileMask != "" && !matchesFileMask(info.Name(), fileMask) {
				return nil
			}

			localHash, err := computeMD5(path)
			if err != nil {
				return fmt.Errorf("unable to compute hash for file '%s': %v", path, err)
			}

			// Сохраняем хэш локального файла в localFilesMap
			localFilesMap[relativePath] = localHash

			// Сравнение полного относительного пути
			driveFile, exists := driveFilesMap[relativePath]
			if exists {
				if driveFile.Md5Checksum == localHash {
					fmt.Printf("File '%s' is up-to-date on Google Drive. Skipping upload.\n", relativePath)
					return nil // Файл не изменился, пропускаем его
				} else {
					fmt.Printf("File '%s' has changed. Updating on Google Drive...\n", relativePath)
					err = uploadFile(srv, remoteFolderID, path, driveFile)
					if err != nil {
						return fmt.Errorf("unable to update file '%s': %v", path, err)
					}
				}
			} else {
				// Если файл отсутствует на Google Drive, загружаем его
				fmt.Printf("File '%s' does not exist on Google Drive. Uploading...\n", relativePath)
				remoteFileDir := filepath.Dir(relativePath)
				folderID := remoteFolderID

				if remoteFileDir != "." && remoteFileDir != "" {
					folderID, err = findOrCreateFolderByPath(srv, remoteFolderID, remoteFileDir, mkdir)
					if err != nil {
						return fmt.Errorf("unable to create or find folder '%s': %v", remoteFileDir, err)
					}
				}

				err = uploadFile(srv, folderID, path, nil)
				if err != nil {
					return fmt.Errorf("unable to upload file '%s': %v", path, err)
				}
			}
		}

		return nil
	})

	// Output debug information about local files and their hashes
	debugPrintln("Local files:")
	for path, hash := range localFilesMap {
		debugPrintf("  %s: %s\n", path, hash)
	}

	if err != nil {
		return err
	}

	// Проверка на файлы, которые есть на Google Drive, но нет локально (только если direction == "from_drive")
	if direction == "from_drive" {
		for remoteFileName, driveFile := range driveFilesMap {
			if _, exists := localFilesMap[remoteFileName]; !exists {
				localFilePath := filepath.Join(localDir, remoteFileName)
				fmt.Printf("File '%s' exists on Google Drive but not locally. Downloading to '%s'...\n", remoteFileName, localFilePath)

				// Создание директорий на локальном диске, если они отсутствуют
				localFileDir := filepath.Dir(localFilePath)
				if err := os.MkdirAll(localFileDir, os.ModePerm); err != nil {
					return fmt.Errorf("unable to create local directories '%s': %v", localFileDir, err)
				}

				err := downloadFile(srv, driveFile.Id, localFilePath)
				if err != nil {
					return fmt.Errorf("unable to download file '%s': %v", remoteFileName, err)
				}
			}
		}
	}

	return nil
}

// Функция для загрузки файла с Google Drive на локальный диск
func downloadFile(srv *drive.Service, fileID, localPath string) error {
	// Открываем файл на запись
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("unable to create local file '%s': %v", localPath, err)
	}
	defer file.Close()

	// Загружаем содержимое файла с Google Drive
	res, err := srv.Files.Get(fileID).SupportsAllDrives(true).Download()
	if err != nil {
		return fmt.Errorf("unable to download file with ID '%s': %v", fileID, err)
	}
	defer res.Body.Close()

	// Копируем содержимое файла в локальный файл
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return fmt.Errorf("unable to save file '%s': %v", localPath, err)
	}

	fmt.Printf("File '%s' successfully downloaded.\n", localPath)
	return nil
}

// Функция для фильтрации файлов по маске
func matchesFileMask(fileName, fileMask string) bool {
	matched, err := filepath.Match(fileMask, fileName)
	if err != nil {
		log.Fatalf("Invalid file mask: %v", err)
	}
	return matched
}

func getAllFilesFromDrive(srv *drive.Service, folderID, currentPath string) (map[string]*drive.File, error) {
	filesMap := make(map[string]*drive.File)

	query := fmt.Sprintf("'%s' in parents and trashed=false", folderID)
	pageToken := ""
	for {
		req := srv.Files.List().
			Q(query).
			Fields("nextPageToken, files(id, name, md5Checksum, mimeType)").
			SupportsAllDrives(true).
			IncludeItemsFromAllDrives(true).
			PageToken(pageToken)

		r, err := req.Do()
		if err != nil {
			return nil, err
		}

		for _, file := range r.Files {
			fullPath := filepath.Join(currentPath, file.Name)
			if file.MimeType == "application/vnd.google-apps.folder" {
				// Рекурсивно обходим поддиректории
				subDirFiles, err := getAllFilesFromDrive(srv, file.Id, fullPath)
				if err != nil {
					return nil, err
				}
				for k, v := range subDirFiles {
					filesMap[k] = v
				}
			} else {
				filesMap[fullPath] = file
			}
		}

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}

	// debug information
	debugPrintln("Files on Google Drive:")
	for name, file := range filesMap {
		debugPrintf("  %s: %s\n", name, file.Md5Checksum)
	}

	return filesMap, nil
}

func debugPrintln(args ...interface{}) {
	if debug {
		fmt.Println(args...)
	}
}

func debugPrintf(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format, args...)
	}
}

func syncSingleFile(srv *drive.Service, localFilePath, remoteFolderID, fileMask string, mkdir bool) error {
	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		return fmt.Errorf("unable to get file info: %v", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("provided path is a directory, expected a file")
	}

	localHash, err := computeMD5(localFilePath)
	if err != nil {
		return fmt.Errorf("unable to compute hash for file '%s': %v", localFilePath, err)
	}

	relativePath := fileInfo.Name()
	driveFilesMap, err := getAllFilesFromDrive(srv, remoteFolderID, "")
	if err != nil {
		return fmt.Errorf("unable to get files from Google Drive: %v", err)
	}

	driveFile, exists := driveFilesMap[relativePath]
	if exists {
		if driveFile.Md5Checksum == localHash {
			fmt.Printf("File '%s' is up-to-date on Google Drive. Skipping upload.\n", relativePath)
			return nil // Файл не изменился, пропускаем его
		} else {
			fmt.Printf("File '%s' has changed. Updating on Google Drive...\n", relativePath)
			err = uploadFile(srv, remoteFolderID, localFilePath, driveFile)
			if err != nil {
				return fmt.Errorf("unable to update file '%s': %v", localFilePath, err)
			}
		}
	} else {
		// Если файл отсутствует на Google Drive, загружаем его
		fmt.Printf("File '%s' does not exist on Google Drive. Uploading...\n", relativePath)
		err = uploadFile(srv, remoteFolderID, localFilePath, nil)
		if err != nil {
			return fmt.Errorf("unable to upload file '%s': %v", localFilePath, err)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	if debug {
		fmt.Println("Debug mode enabled")
	}
	flag.Parse()

	if uploadFilePath == "" && (syncFolderPath == "" || folderPath == "") {
		log.Fatal("Both -sync_folder and -folder_path must be specified, unless using -upload_file.")
	}

	b, err := ioutil.ReadFile(credPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	ctx := context.Background()
	client := config.Client(ctx)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	rootFolderID := "root"
	if driveId != "" {
		rootFolderID = driveId // Используем driveId как корень, если он указан
	}

	// Only one file upload is supported
	if uploadFilePath != "" {
		fmt.Printf("Uploading file: %s to Google Drive folder: %s\n", uploadFilePath, folderPath)

		folderID, err := findOrCreateFolderByPath(srv, rootFolderID, folderPath, mkdir)
		if err != nil {
			log.Fatalf("Failed to create or find folder: %v", err)
		}

		err = syncSingleFile(srv, uploadFilePath, folderID, fileMask, mkdir)
		if err != nil {
			log.Fatalf("Failed to upload file: %v", err)
		}
	} else if syncFolderPath != "" && folderPath != "" {
		fmt.Printf("Syncing folder: %s to Google Drive folder: %s\n", syncFolderPath, folderPath)

		folderID, err := findOrCreateFolderByPath(srv, rootFolderID, folderPath, mkdir)
		if err != nil {
			log.Fatalf("Error finding or creating target folder: %v", err)
		}

		err = syncFolder(srv, syncFolderPath, folderID, direction, fileMask, mkdir)
		if err != nil {
			log.Fatalf("Error syncing folder: %v", err)
		}
	} else {
		log.Fatal("No operation specified. Use -sync_folder to sync a local folder with Google Drive or -upload_file to upload a single file.")
	}
}
