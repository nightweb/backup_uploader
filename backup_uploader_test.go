package main

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// Функция для инициализации Google Drive Service для тестов
func initDriveService(t *testing.T) *drive.Service {
	credPath := "path/to/credentials.json" // Укажите путь к вашему файлу учетных данных
	b, err := ioutil.ReadFile(credPath)
	if err != nil {
		t.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		t.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	ctx := context.Background()
	client := config.Client(ctx)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		t.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return srv
}

func TestUploadFile(t *testing.T) {
	srv := initDriveService(t)

	// Создаем временный файл для загрузки
	tmpFile, err := ioutil.TempFile("", "upload_test_file")
	if err != nil {
		t.Fatalf("Unable to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("This is a test file for upload.")
	if err != nil {
		t.Fatalf("Unable to write to temporary file: %v", err)
	}

	folderID := "your-test-folder-id" // Укажите тестовый идентификатор папки

	err = uploadFile(srv, folderID, tmpFile.Name())
	if err != nil {
		t.Errorf("Failed to upload file: %v", err)
	}
}

func TestFindOrCreateFolderByPath(t *testing.T) {
	srv := initDriveService(t)

	rootFolderID := "your-test-drive-id" // Укажите тестовый идентификатор диска
	folderPath := "Test/Path/For/Creation"

	folderID, err := findOrCreateFolderByPath(srv, rootFolderID, folderPath, true)
	if err != nil {
		t.Errorf("Failed to find or create folder: %v", err)
	}

	if !strings.HasPrefix(folderID, "folder:") {
		t.Errorf("Invalid folder ID returned: %s", folderID)
	}
}
