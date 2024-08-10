
# Backup Uploader

## Описание

`backup_uploader` - это инструмент для синхронизации файлов между локальной директорией и Google Drive. Он позволяет загружать и обновлять файлы на Google Drive, обеспечивая при этом проверку хэшей файлов для избежания дублирования и излишнего трафика.

## Возможности

- Синхронизация локальной директории с Google Drive (`to_drive`).
- Синхронизация Google Drive с локальной директорией (`from_drive`).
- Автоматическое создание отсутствующих папок на Google Drive при загрузке.
- Сравнение файлов на основе MD5-хэшей.
- Обновление только тех файлов, которые были изменены.
- Поддержка режима отладки с подробным выводом информации (`-vv` или `--debug`).

## Установка

1. Установите Go (https://golang.org/doc/install).
2. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/yourusername/backup_uploader.git
   cd backup_uploader
   ```

3. Соберите исполняемый файл:
   ```sh
   go build -o backup_uploader
   ```

## Использование

### Основные параметры

- `-c` : Путь к файлу `credentials.json` для Google API (если не указан, используется `~/.backup_uploader/google/credentials.json`).
- `-sync_folder` : Путь к локальной директории, которую нужно синхронизировать.
- `-folder_path` : Путь к целевой папке на Google Drive.
- `-drive_id` : ID Google Drive (для Shared Drives, необязательно).
- `-file_mask` : Маска файлов для синхронизации (например, `*.txt`).
- `-mkdir` : Флаг для создания папок на Google Drive, если они отсутствуют.
- `-direction` : Направление синхронизации (`to_drive` или `from_drive`).
- `-vv`, `--debug` : Включение режима отладки с подробным выводом информации.

### Примеры использования

1. **Синхронизация локальной директории с Google Drive**:

   ```sh
   ./backup_uploader -c path/to/credentials.json -sync_folder /path/to/local/folder -folder_path "TargetFolder" -drive_id YOUR_DRIVE_ID -mkdir -direction to_drive -file_mask '*.txt' -vv
   ```

2. **Синхронизация Google Drive с локальной директорией**:

   ```sh
   ./backup_uploader -c path/to/credentials.json -sync_folder /path/to/local/folder -folder_path "TargetFolder" -drive_id YOUR_DRIVE_ID -mkdir -direction from_drive -file_mask '*.txt' -vv
   ```

### Отладка

При запуске скрипта с флагом `-vv` или `--debug` выводится отладочная информация, которая включает в себя:

- Список файлов и их хэшей на Google Drive.
- Список локальных файлов и их хэшей.
- Сообщения о том, какие файлы были обновлены, пропущены или загружены заново.

### Пример тестирования

Для тестирования функциональности скрипта можно использовать следующую команду:

```sh
go test -v
```

## Лицензия

Этот проект лицензирован под лицензией MIT - см. файл [LICENSE](LICENSE) для подробностей.
