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

## Building for Different Platforms

You can build the `backup_uploader` binary for different operating systems and architectures using the following commands:

### For Linux (64-bit)

```sh
GOOS=linux GOARCH=amd64 go build -o backup_uploader_linux
```

### For Windows (64-bit)

```sh
GOOS=windows GOARCH=amd64 go build -o backup_uploader.exe
```

### For macOS (64-bit)

```sh
GOOS=darwin GOARCH=amd64 go build -o backup_uploader_mac
```

### For ARM (e.g., Raspberry Pi)

```sh
GOOS=linux GOARCH=arm go build -o backup_uploader_arm
```

These commands will generate binaries for the specified platforms, which you can then distribute and run on the respective systems.

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

## Testing

To run tests for the Backup Uploader, follow these steps:

1. Make sure you have your `credentials.json` file set up as described above.

2. Set up the necessary environment variables or paths in the test file `backup_uploader_test.go`.

3. Run the tests using the `go test` command:

   ```sh
   go test
   ```

   This will execute the test cases defined in `backup_uploader_test.go` and help ensure the functionality of the Backup Uploader.

## Contributing

If you'd like to contribute to this project, follow these steps:

1. Fork the repository.

2. Create a new branch for your feature or bugfix:

   ```sh
   git checkout -b feature-name
   ```

3. Make your changes and commit them:

   ```sh
   git commit -m "Description of the feature or fix"
   ```

4. Push your changes to your fork:

   ```sh
   git push origin feature-name
   ```

5. Create a Pull Request from your fork's branch to the `main` branch of this repository.

6. Ensure your code passes all tests before submitting your PR.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
