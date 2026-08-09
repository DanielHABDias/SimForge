# SimForge TODO

## Task: Add missing menu handlers in cmd/app/main.go

- [x] Add `handleVerifyUnlocker(mgr *unlocker.Manager)` to verify the unlocker install and print a report
- [x] Add `handleCheckPaths(inst *finder.Installation, cfg *config.Config)` to show all detected paths
- [x] Verify build: `go build ./...`
- [x] Verify vet: `go vet ./...`
- [x] Verify tests: `go test ./...`
- [x] Cross-compile Windows: `GOOS=windows go build ./...`
