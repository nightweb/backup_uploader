package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Функция для поиска или создания папки по пути
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

// Функция для загрузки файла в указанную папку
func uploadFile(srv *drive.Service, folderID, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to get file info: %v", err)
	}

	driveFile := &drive.File{
		Name:    fileInfo.Name(),
		Parents: []string{folderID},
	}

	uploadedFile, err := srv.Files.Create(driveFile).Media(file).SupportsAllDrives(true).Do()
	if err != nil {
		return fmt.Errorf("unable to upload file: %v", err)
	}

	fmt.Printf("File '%s' successfully uploaded. File ID: %s\n", uploadedFile.Name, uploadedFile.Id)
	return nil
}

// Функция для рекурсивного вывода всех файлов и папок
func listFilesAndFolders(srv *drive.Service, folderID string, indent string) error {
	query := fmt.Sprintf("'%s' in parents and trashed=false", folderID)
	r, err := srv.Files.List().
		Q(query).
		Fields("files(id, name, mimeType, parents)").
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Do()

	if err != nil {
		return err
	}

	if len(r.Files) == 0 {
		fmt.Printf("%sNo files or folders found in folder %s.\n", indent, folderID)
		return nil
	}

	for _, file := range r.Files {
		fmt.Printf("%sName: %s, ID: %s, Type: %s\n", indent, file.Name, file.Id, file.MimeType)
		if file.MimeType == "application/vnd.google-apps.folder" {
			fmt.Printf("%sEntering folder %s...\n", indent, file.Name)
			listFilesAndFolders(srv, file.Id, indent+"  ")
		}
	}

	return nil
}

func main() {
	credPath := flag.String("c", filepath.Join(os.Getenv("HOME"), ".backup_uploader", "google", "credentials.json"), "Path to credentials.json file")
	driveId := flag.String("drive_id", "", "Google Drive ID (for Shared Drives)")
	uploadFilePath := flag.String("upload_file", "", "Path to the file to be uploaded")
	targetFolderPath := flag.String("folder_path", "", "Path to the folder to upload the file to")
	mkdir := flag.Bool("mkdir", false, "Create directory if it does not exist")
	listOnly := flag.Bool("list", false, "List all files and folders and exit")
	flag.Parse()

	b, err := ioutil.ReadFile(*credPath)
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

	if *listOnly {
		rootFolderID := "root"
		if *driveId != "" {
			rootFolderID = *driveId // Используем driveId как корень, если он указан
		}

		fmt.Println("Listing all files and folders:")
		err = listFilesAndFolders(srv, rootFolderID, "")
		if err != nil {
			log.Fatalf("Error listing files and folders: %v", err)
		}
		return
	}

	if *uploadFilePath != "" && *targetFolderPath != "" {
		rootFolderID := "root"
		if *driveId != "" {
			rootFolderID = *driveId // Используем driveId как корень, если он указан
		}

		folderID, err := findOrCreateFolderByPath(srv, rootFolderID, *targetFolderPath, *mkdir)
		if err != nil {
			log.Fatalf("Error finding or creating target folder: %v", err)
		}

		err = uploadFile(srv, folderID, *uploadFilePath)
		if err != nil {
			log.Fatalf("Error uploading file: %v", err)
		}
		return
	}

	log.Fatal("No operation specified. Use -list to list files and folders, or specify a file and folder path to upload.")
}
