
# Backup Uploader

## Description

`backup_uploader` is a tool for synchronizing files between a local directory and Google Drive. It allows you to upload and update files on Google Drive while ensuring that files are compared using hash values to avoid duplication and reduce unnecessary traffic.

## Features

- Synchronize a local directory with Google Drive (`to_drive`).
- Synchronize Google Drive with a local directory (`from_drive`).
- Automatically create missing folders on Google Drive when uploading.
- Compare files based on MD5 hash values.
- Update only files that have been modified.
- Support for debug mode with detailed output (`-vv` or `--debug`).
- Upload a single file directly to Google Drive with the `-upload_file` flag.

## Installation

1. Install Go (https://golang.org/doc/install).
2. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/backup_uploader.git
   cd backup_uploader
   ```

3. Build the executable:
   ```sh
   go build -o backup_uploader
   ```

## Usage

### Main Flags

- `-c` : Path to the `credentials.json` file for Google API (default: `~/.backup_uploader/google/credentials.json`).
- `-sync_folder` : Path to the local directory to sync.
- `-folder_path` : Path to the target folder on Google Drive.
- `-drive_id` : Google Drive ID (for Shared Drives, optional).
- `-file_mask` : File mask to filter files for syncing (e.g., `*.txt`).
- `-mkdir` : Flag to create directories on Google Drive if they do not exist.
- `-direction` : Sync direction (`to_drive` or `from_drive`).
- `-vv`, `--debug` : Enable verbose output (debug mode).
- `-upload_file` : Path to a single file to upload to Google Drive.

### Examples

1. **Synchronize a local directory with Google Drive**:

   ```sh
   ./backup_uploader -c path/to/credentials.json -sync_folder /path/to/local/folder -folder_path "TargetFolder" -drive_id YOUR_DRIVE_ID -mkdir -direction to_drive -file_mask '*.txt' -vv
   ```

2. **Synchronize Google Drive with a local directory**:

   ```sh
   ./backup_uploader -c path/to/credentials.json -sync_folder /path/to/local/folder -folder_path "TargetFolder" -drive_id YOUR_DRIVE_ID -mkdir -direction from_drive -file_mask '*.txt' -vv
   ```

3. **Upload a single file to Google Drive**:

   ```sh
   ./backup_uploader -c path/to/credentials.json -upload_file /path/to/file.txt -folder_path "TargetFolder" -vv
   ```

### Debugging

When running the script with the `-vv` or `--debug` flag, detailed debugging information is output, including:

- A list of files and their hashes on Google Drive.
- A list of local files and their hashes.
- Messages about which files were updated, skipped, or uploaded.

### Testing

To test the script's functionality, you can use the following command:

```sh
go test -v
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
