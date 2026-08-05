# SimForge TODO

## 1. Fix compile errors ✅
- [x] Add missing `os` import in `internal/platform/platform.go`
- [x] Remove dead `internal/finder/windows_path.go`

## 2. Refactor `finder` package ✅
- [x] Consolidate `Installation` struct (add `Client`, `DllPath2`, proper `User`)
- [x] Fix `findEA` to locate versioned EA Desktop folder
- [x] Use `steamuser` for Steam prefixes
- [x] Add Origin detection on Linux
- [x] Add `Sims4Root` and `Sims4Mods` detection
- [x] Clean up/format code

## 3. Create `internal/unlocker` package ✅
- [x] `unlocker.go` — Manager, assets location, common operations
- [x] `linux_unlocker.go` — user.reg DllOverrides, DLL copy
- [x] `windows_unlocker.go` — scheduled task, machine.ini
- [x] `assets.go` — locate assets root, source DLL, game configs

## 4. Rewrite `cmd/app/main.go` ✅
- [x] Interactive menu (list prefixes, install, add/update config, open folders, uninstall, quit)

## 5. Configuration & Transfer ✅
- [x] `internal/config` — `config.json` with `dlc_path` (default `assets/dlcs`) and `mods_path`
- [x] `internal/transfer` — copy DLCs to game root, mods to Mods folder, extract archives
- [x] Menu options to transfer DLCs and mods
- [x] Menu option to show/edit paths

## 6. Verify
- [x] `go build ./...` (Linux)
- [x] `GOOS=windows go build ./...` (cross-check)
- [x] `go vet ./...`
