# SimForge Refactor TODO

## 1. Fix compile errors ✅
- [x] Add missing `os` import in `internal/platform/platform.go`
- [x] Remove dead `internal/finder/windows_path.go`

## 2. Refactor `finder` package ✅
- [x] Consolidate `Installation` struct (add `Client`, `DllPath2`, proper `User`)
- [x] Fix `findEA` to locate versioned EA Desktop folder
- [x] Use `steamuser` for Steam prefixes
- [x] Add Origin detection on Linux
- [x] Clean up/format code

## 3. Create `internal/unlocker` package ✅
- [x] `unlocker.go` — Manager, assets location, common operations
- [x] `linux_unlocker.go` — user.reg DllOverrides, DLL copy
- [x] `windows_unlocker.go` — scheduled task, machine.ini
- [x] `assets.go` — locate assets root, source DLL, game configs

## 4. Rewrite `cmd/app/main.go`
- [ ] Interactive menu (list prefixes, install, add/update config, open folders, uninstall, quit)

## 5. Verify
- [ ] `go build ./...` (Linux)
- [ ] `GOOS=windows go build ./...` (cross-check)
- [ ] `go vet ./...`
