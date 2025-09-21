# FTUCK Logging

FTUCK uses structured logging with the standard Go `log/slog` package to provide comprehensive, configurable logging throughout the application.

## Configuration

Logging can be configured using environment variables:

- `FTUCK_LOG_LEVEL`: Controls the minimum log level (debug, info, warn, error). Default: `info`
- `FTUCK_LOG_FORMAT`: Controls the output format (text, json). Default: `text`

## Log Levels

- **debug**: Detailed information for debugging application behavior
- **info**: General information about application operations
- **warn**: Warning messages for unexpected but recoverable situations
- **error**: Error messages for problems that prevent normal operation

## Usage Examples

### Default logging (INFO level, text format)
```bash
./ftuck init
```
Output:
```
time=2025-09-21T18:54:52.081Z level=INFO msg="config file does not exist, creating new one" path=/home/runner/.ftuck.yaml
time=2025-09-21T18:54:52.081Z level=INFO msg="sync file found" path=/tmp/ftuck_test2/.ftucksync.yaml
```

### Debug logging for troubleshooting
```bash
FTUCK_LOG_LEVEL=debug ./ftuck init
```
Output:
```
time=2025-09-21T18:53:31.684Z level=DEBUG msg="starting ftuck application" log_level=0 log_format=text
time=2025-09-21T18:53:31.684Z level=DEBUG msg="attempting to open config file" path=/home/runner/.ftuck.yaml
time=2025-09-21T18:53:31.684Z level=INFO msg="config file does not exist, creating new one" path=/home/runner/.ftuck.yaml
time=2025-09-21T18:53:31.684Z level=DEBUG msg="searching for sync file" sync_file_name=.ftucksync.yaml directory=/tmp/ftuck_test
time=2025-09-21T18:53:31.684Z level=INFO msg="sync file found" path=/tmp/ftuck_test/.ftucksync.yaml
time=2025-09-21T18:53:31.684Z level=DEBUG msg="saving config file" path=/home/runner/.ftuck.yaml
time=2025-09-21T18:53:31.684Z level=DEBUG msg="config file saved successfully" path=/home/runner/.ftuck.yaml
```

### JSON format for log aggregation
```bash
FTUCK_LOG_LEVEL=debug FTUCK_LOG_FORMAT=json ./ftuck addsync -s file1 -t /tmp/link1
```
Output:
```json
{"time":"2025-09-21T18:53:39.522314671Z","level":"DEBUG","msg":"starting ftuck application","log_level":0,"log_format":"json"}
{"time":"2025-09-21T18:53:39.522433444Z","level":"DEBUG","msg":"attempting to open config file","path":"/home/runner/.ftuck.yaml"}
{"time":"2025-09-21T18:53:39.522548279Z","level":"DEBUG","msg":"config file loaded successfully","path":"/home/runner/.ftuck.yaml","sync_file":"/tmp/ftuck_test/.ftucksync.yaml"}
{"time":"2025-09-21T18:53:39.522560822Z","level":"DEBUG","msg":"using configured sync file","sync_file":"/tmp/ftuck_test/.ftucksync.yaml"}
{"time":"2025-09-21T18:53:39.522564549Z","level":"DEBUG","msg":"reading sync file","path":"/tmp/ftuck_test/.ftucksync.yaml"}
{"time":"2025-09-21T18:53:39.522632547Z","level":"INFO","msg":"adding sync definition","source":"new_test","destination":"/tmp/new_link"}
{"time":"2025-09-21T18:53:39.52264496Z","level":"DEBUG","msg":"writing sync schema to file","path":"/tmp/ftuck_test/.ftucksync.yaml","entries":2}
{"time":"2025-09-21T18:53:39.522786185Z","level":"DEBUG","msg":"sync schema written successfully","path":"/tmp/ftuck_test/.ftucksync.yaml"}
```

### Error-only logging for production
```bash
FTUCK_LOG_LEVEL=error ./ftuck sync
```

## Structured Fields

FTUCK uses structured logging with contextual fields:

- `path`: File or directory paths
- `error`: Error messages and details
- `source`/`destination`: Sync operation source and target paths
- `entries`: Number of sync entries
- `sync_file_name`: Name of sync files
- `directory`: Working directories
- And many more contextual fields throughout the application

## Implementation Details

The logging system is implemented in `internal/logging/` and includes:

- **logger.go**: Main logging configuration and initialization
- **logger_test.go**: Comprehensive tests for logging functionality

The logger is initialized at application startup and configured as the default slog logger, so all modules throughout the application use the same configured logger automatically.