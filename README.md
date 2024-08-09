# Backup Uploader

Backup Uploader is a command-line tool for uploading files to Google Drive, with support for creating directories if they don't exist.

## Features

- Upload files to a specified directory on Google Drive or Shared Drives.
- Automatically create directories along the path if they don't exist using the `-mkdir` option.
- List files and folders in Google Drive or Shared Drives.

## Requirements

- Go 1.16 or higher
- Google Drive API credentials

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/backup-uploader.git
   cd backup-uploader
   ```

2. Build the binary:

   ```sh
   go build -o backup_uploader
   ```

3. Obtain Google Drive API credentials and save them as `credentials.json`:

   - Go to [Google Cloud Console](https://console.cloud.google.com/).
   - Create a new project or use an existing one.
   - Enable the Google Drive API for your project.
   - Create credentials for a service account and download the `credentials.json` file.

## Usage

### Upload a File

Upload a file to a specific folder path within Google Drive or a Shared Drive:

```sh
./backup_uploader -c path/to/credentials.json -upload_file /path/to/local/file.txt -folder_path "TargetFolder/SubFolder" -drive_id YOUR_DRIVE_ID
```

### Create Directories Automatically

Use the `-mkdir` option to automatically create directories if they don't exist:

```sh
./backup_uploader -c path/to/credentials.json -upload_file /path/to/local/file.txt -folder_path "TargetFolder/SubFolder" -drive_id YOUR_DRIVE_ID -mkdir
```

### List Files and Folders

List all files and folders in Google Drive or a Shared Drive:

```sh
./backup_uploader -c path/to/credentials.json -drive_id YOUR_DRIVE_ID -list
```

## Options

- `-c path/to/credentials.json` : Path to your Google Drive API credentials file.
- `-upload_file /path/to/local/file.txt` : Path to the file to be uploaded.
- `-folder_path "TargetFolder/SubFolder"` : Path to the folder on Google Drive where the file should be uploaded.
- `-drive_id YOUR_DRIVE_ID` : ID of the Shared Drive (optional if uploading to My Drive).
- `-mkdir` : Create directories if they do not exist.
- `-list` : List all files and folders.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
