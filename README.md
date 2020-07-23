# boards
A simple task and note app for the command line written in Go.

## Usage
```
boards
boards task|t [@BOARD ...] DESCRIPTION [DESCRIPTION ...]
boards note|n [@BOARD ...] DESCRIPTION [DESCRIPTION ...]
boards delete|d ID [ID ...]
boards complete|c ID [ID ...]
boards uncomplete|u ID [ID ...]
boards edit|e ID DESCRIPTION [DESCRIPTION ...]
boards boards|b ID [ID ...] [@BOARD ...]
boards mark|m ID [ID ...]
boards clear
boards --help|-h

If boards is run without any argument all tasks will be printed grouped by their boards.
Clear will remove all complete tasks.
```

The tasks will be stored in `$XDG_CONFIG_HOME/boards/storage.json` (defaults to `~/.config/boards/storage.json`).

## Build
Run `make` to build `boards` binary.

Run `make install` to install the binary to `/usr/local/bin`. A custom install directory can be set via the `INSTALL_DIR` environment variable.
