# SimForge
Upgrade your The Sims 4 installation with mods and DLCs

## Features

- Detect EA App / Origin installations (Windows and Wine/Proton prefixes on Linux).
- Install/uninstall the EA DLC Unlocker and manage game configs.
- **Transfer DLCs** from a source folder into the game root.
- **Transfer mods** from a source folder into the game `Mods` folder.
- Automatically **extract archives** (zip, tar, tar.gz, tgz, gz) found in the source folders.
- Configurable source paths.

## Configuration

The program uses a `config.json` file (created automatically in the working
directory on first run). It stores the source folders where you keep your
downloaded files:

```json
{
  "dlc_path": "assets/dlcs",
  "mods_path": "assets/mods"
}
```

- `dlc_path` — folder containing the downloaded DLC files. Default: `assets/dlcs`.
- `mods_path` — folder containing the downloaded mod files. Default: `assets/mods`.

Both paths can be changed from the menu (`Show/Edit paths`) or by editing
`config.json` directly.

## How it works

1. The program reads the configured source paths.
2. If a configured path does not exist, the program asks you to type a valid one.
3. It reports how many files were found in the source folder.
4. It locates the The Sims 4 game root and its `Mods` folder:
   - DLCs are copied/extracted into the **game root**.
   - Mods are copied/extracted into the **`Mods` folder**.
5. Any compressed file (zip, tar, tar.gz, gz) inside the copied collection is
   extracted automatically.

## Build

```sh
go build ./...
```

## Run

```sh
go run ./cmd/app
