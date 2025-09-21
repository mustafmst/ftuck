# FTUCK (file-tuck)

Pronouncing name of this tool remember that `t` is silent ;)

FTUCK is a more user friendly stow written in go for certain purpose.
This app is a small tool for keeping all configuration that You need in one neat repo for easy sync and install on system.

## Logging

FTUCK uses structured logging with configurable levels and formats. See [docs/logging.md](docs/logging.md) for detailed information.

Quick examples:
```bash
# Default logging
./ftuck init

# Debug logging
FTUCK_LOG_LEVEL=debug ./ftuck init

# JSON format for log aggregation
FTUCK_LOG_FORMAT=json ./ftuck sync
```

## TBD

Right now in development. More info later.

## Installation

## Using FTUCK
