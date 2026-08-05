package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/DanielHABDias/SimForge/internal/config"
	"github.com/DanielHABDias/SimForge/internal/finder"
	"github.com/DanielHABDias/SimForge/internal/transfer"
	"github.com/DanielHABDias/SimForge/internal/unlocker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	items, err := finder.Find()
	if err != nil {
		log.Fatalf("Error detecting installations: %v", err)
	}

	if len(items) == 0 {
		fmt.Println("No EA app or Origin installations found.")
		fmt.Println("If you're on Linux, run the game once first so the Wine/Proton prefix is created.")
		waitEnter()
		return
	}

	// If only one installation, use it directly.
	if len(items) == 1 {
		runMenu(&items[0], cfg)
		return
	}

	// Multiple installations: show selection.
	showPrefixMenu(items, cfg)
}

func showPrefixMenu(items []finder.Installation, cfg *config.Config) {
	for {
		fmt.Println("Multiple wine/proton prefixes found!")
		fmt.Println("Which one do you want to manage?")
		fmt.Println()

		for i, item := range items {
			launcherInfo := item.Launcher
			if item.Name != "" {
				launcherInfo = fmt.Sprintf("%s (%s)", item.Name, item.Launcher)
			}
			fmt.Printf("  %s\x1b[33m%d\x1b[0m. %s\n",
				"\x1b[33m",
				i+1,
				launcherInfo,
			)
		}
		fmt.Printf("  \x1b[31mq\x1b[0m. Quit\n")

		choice := readInput("Type the number of your choice and press Enter: ")
		clearScreen()

		if strings.ToLower(choice) == "q" {
			return
		}
		idx, err := strconv.Atoi(choice)
		if err != nil || idx < 1 || idx > len(items) {
			fmt.Println("\x1b[31mInvalid choice! Select a number shown in the menu.\x1b[0m")
			fmt.Println()
			continue
		}
		runMenu(&items[idx-1], cfg)
	}
}

func runMenu(inst *finder.Installation, cfg *config.Config) {
	// Build unlocker manager for this installation.
	mgr, err := unlocker.New(*inst)
	if err != nil {
		fmt.Printf("\x1b[31mError: %v\x1b[0m\n", err)
		waitEnter()
		return
	}

	for {
		clearScreen()

		// Header.
		fmt.Printf("Prefix: %s (%s)\n", inst.Launcher, inst.PrefixPath)
		fmt.Printf("The Sims 4 root: %s\n", orNA(inst.Sims4Root))
		fmt.Printf("The Sims 4 Mods: %s\n", orNA(inst.Sims4Mods))

		fmt.Printf("DLC Unlocker ")
		if mgr.IsInstalled() {
			fmt.Printf("\x1b[32minstalled\x1b[0m")
			if configs := mgr.InstalledGameConfigs(); len(configs) > 0 {
				fmt.Printf("\nGame configs installed: \x1b[36m%s\x1b[0m", strings.Join(configs, ", "))
			} else {
				fmt.Printf("")
			}
		} else {
			fmt.Printf("\x1b[31mnot installed\x1b[0m")
		}
		fmt.Println()

		fmt.Println()
		fmt.Println("  \x1b[36m--- DLC Unlocker ---\x1b[0m")
		fmt.Printf("  \x1b[33m1\x1b[0m. Install EA DLC Unlocker\n")
		fmt.Printf("  \x1b[33m2\x1b[0m. Add/Update game config\n")
		fmt.Printf("  \x1b[33m3\x1b[0m. Open folder with installed configs\n")
		fmt.Printf("  \x1b[33m4\x1b[0m. Open folder with log file\n")
		fmt.Printf("  \x1b[33m5\x1b[0m. Uninstall EA DLC Unlocker\n")

		fmt.Println()
		fmt.Println("  \x1b[36m--- The Sims 4 Content ---\x1b[0m")
		fmt.Printf("  \x1b[33m6\x1b[0m. Transfer DLCs to game root\n")
		fmt.Printf("  \x1b[33m7\x1b[0m. Transfer mods to Mods folder\n")

		fmt.Println()
		fmt.Println("  \x1b[36m--- Configuration ---\x1b[0m")
		fmt.Printf("  \x1b[33m8\x1b[0m. Show/Edit paths (dlc_path, mods_path)\n")

		fmt.Printf("  \x1b[31mq\x1b[0m. Quit\n")

		choice := readInput("Type the number of your choice and press Enter: ")
		clearScreen()

		switch choice {
		case "1":
			handleInstall(mgr)
		case "2":
			handleAddGameConfig(mgr)
		case "3":
			handleOpenConfigs(mgr)
		case "4":
			handleOpenLogs(mgr)
		case "5":
			handleUninstall(mgr)
		case "6":
			handleTransferDLCs(inst, cfg)
		case "7":
			handleTransferMods(inst, cfg)
		case "8":
			handleEditConfig(cfg)
		case "q":
			return
		default:
			fmt.Println("\x1b[31mInvalid choice! Select a number shown in the menu.\x1b[0m")
		}
		fmt.Println()
	}
}

func handleInstall(mgr *unlocker.Manager) {
	fmt.Println("Installing...")
	if err := mgr.Install(); err != nil {
		fmt.Printf("\x1b[31mError: %v\x1b[0m\n", err)
		waitEnter()
		return
	}
	fmt.Println("\x1b[32mDLC Unlocker installed!\x1b[0m")
	waitEnter()
}

func handleAddGameConfig(mgr *unlocker.Manager) {
	names, err := mgr.GameConfigNames()
	if err != nil || len(names) == 0 {
		fmt.Println("\x1b[31mGame configs missing, you didn't extract all files.\x1b[0m")
		waitEnter()
		return
	}

	for {
		fmt.Println("Game configs:")
		for i, name := range names {
			fmt.Printf("  \x1b[33m%d\x1b[0m. %s\n", i+1, name)
		}
		fmt.Printf("  \x1b[31mb\x1b[0m. Back\n")

		choice := readInput("Type the number of your choice and press Enter: ")
		clearScreen()

		if strings.ToLower(choice) == "b" {
			fmt.Println("No game config selected.")
			return
		}

		idx, err := strconv.Atoi(choice)
		if err != nil || idx < 1 || idx > len(names) {
			fmt.Println("\x1b[31mInvalid choice! Select a number shown in the menu.\x1b[0m")
			fmt.Println()
			continue
		}

		selected := names[idx-1]
		fmt.Printf("\x1b[33m%s\x1b[0m config selected.\n", selected)
		if err := mgr.AddGameConfig(selected); err != nil {
			fmt.Printf("\x1b[31mError: %v\x1b[0m\n", err)
			waitEnter()
			return
		}
		fmt.Println("\x1b[32mGame config copied!\x1b[0m")
		waitEnter()
		return
	}
}

func handleOpenConfigs(mgr *unlocker.Manager) {
	if err := mgr.OpenConfigsFolder(); err != nil {
		fmt.Printf("\x1b[31m%v\x1b[0m\n", err)
		waitEnter()
		return
	}
	fmt.Println("\x1b[32mConfigs folder opened!\x1b[0m")
	waitEnter()
}

func handleOpenLogs(mgr *unlocker.Manager) {
	if err := mgr.OpenLogsFolder(); err != nil {
		fmt.Printf("\x1b[31m%v\x1b[0m\n", err)
		waitEnter()
		return
	}
	fmt.Println("\x1b[32mLogs folder opened!\x1b[0m")
	waitEnter()
}

func handleUninstall(mgr *unlocker.Manager) {
	fmt.Println("Uninstalling...")
	if err := mgr.Uninstall(); err != nil {
		fmt.Printf("\x1b[31mError: %v\x1b[0m\n", err)
		waitEnter()
		return
	}
	fmt.Println("\x1b[32mEA DLC Unlocker uninstalled!\x1b[0m")
	waitEnter()
}

// handleTransferDLCs copies DLCs from the configured source into the game root,
// extracting any archives found.
func handleTransferDLCs(inst *finder.Installation, cfg *config.Config) {
	src := resolveSourcePath(cfg.DLCPath, "Where are the downloaded DLCs?")
	if src == "" {
		return
	}

	if inst.Sims4Root == "" {
		inst.Sims4Root = promptAdditive("Could not auto-detect The Sims 4 game folder. Enter its path: ")
		if inst.Sims4Root == "" {
			fmt.Println("\x1b[31mNo game root provided. Aborting.\x1b[0m")
			waitEnter()
			return
		}
	}

	count, err := transfer.CountFiles(src)
	if err != nil {
		fmt.Printf("\x1b[31mError counting files: %v\x1b[0m\n", err)
		waitEnter()
		return
	}
	fmt.Printf("\x1b[33mFound %d file(s) in %s\x1b[0m\n", count, src)

	fmt.Printf("Transferring DLCs to %s...\n", inst.Sims4Root)
	res, err := transfer.CopyDLCs(src, inst.Sims4Root)
	if err != nil {
		fmt.Printf("\x1b[31mError: %v\x1b[0m\n", err)
		waitEnter()
		return
	}
	fmt.Printf("\x1b[32mDone! Copied %d file(s), extracted %d archive(s).\x1b[0m\n",
		res.CopiedFiles, res.ExtractedZips)
	waitEnter()
}

// handleTransferMods copies mods from the configured source into the game Mods
// folder, extracting any archives found.
func handleTransferMods(inst *finder.Installation, cfg *config.Config) {
	src := resolveSourcePath(cfg.ModsPath, "Where are the downloaded mods?")
	if src == "" {
		return
	}

	if inst.Sims4Root == "" {
		inst.Sims4Root = promptAdditive("Could not auto-detect The Sims 4 game folder. Enter its path: ")
		if inst.Sims4Root == "" {
			fmt.Println("\x1b[31mNo game root provided. Aborting.\x1b[0m")
			waitEnter()
			return
		}
	}

	// Ensure the Mods folder exists.
	if inst.Sims4Mods == "" {
		inst.Sims4Mods = promptAdditive("Enter the The Sims 4 Mods folder path: ")
		if inst.Sims4Mods == "" {
			fmt.Println("\x1b[31mNo Mods folder provided. Aborting.\x1b[0m")
			waitEnter()
			return
		}
	}
	if err := os.MkdirAll(inst.Sims4Mods, 0o755); err != nil {
		fmt.Printf("\x1b[31mError creating Mods folder: %v\x1b[0m\n", err)
		waitEnter()
		return
	}

	count, err := transfer.CountFiles(src)
	if err != nil {
		fmt.Printf("\x1b[31mError counting files: %v\x1b[0m\n", err)
		waitEnter()
		return
	}
	fmt.Printf("\x1b[33mFound %d file(s) in %s\x1b[0m\n", count, src)

	fmt.Printf("Transferring mods to %s...\n", inst.Sims4Mods)
	res, err := transfer.CopyMods(src, inst.Sims4Mods)
	if err != nil {
		fmt.Printf("\x1b[31mError: %v\x1b[0m\n", err)
		waitEnter()
		return
	}
	fmt.Printf("\x1b[32mDone! Copied %d file(s), extracted %d archive(s).\x1b[0m\n",
		res.CopiedFiles, res.ExtractedZips)
	waitEnter()
}

// handleEditConfig shows the current paths and lets the user change them.
func handleEditConfig(cfg *config.Config) {
	fmt.Println("Current configuration:")
	fmt.Printf("  dlc_path = \x1b[36m%s\x1b[0m\n", cfg.DLCPath)
	fmt.Printf("  mods_path = \x1b[36m%s\x1b[0m\n", cfg.ModsPath)
	fmt.Printf("  (config file: %s)\n", config.Path())
	fmt.Println()

	choice := readInput("Edit dlc_path? [y/N]: ")
	if strings.EqualFold(choice, "y") {
		newPath := readInput(fmt.Sprintf("New dlc_path (current: %s): ", cfg.DLCPath))
		if newPath != "" {
			cfg.DLCPath = newPath
		}
	}

	choice = readInput("Edit mods_path? [y/N]: ")
	if strings.EqualFold(choice, "y") {
		newPath := readInput(fmt.Sprintf("New mods_path (current: %s): ", cfg.ModsPath))
		if newPath != "" {
			cfg.ModsPath = newPath
		}
	}

	if err := cfg.Save(); err != nil {
		fmt.Printf("\x1b[31mError saving config: %v\x1b[0m\n", err)
		waitEnter()
		return
	}
	fmt.Println("\x1b[32mConfiguration saved!\x1b[0m")
	waitEnter()
}

// resolveSourcePath returns the configured path if it exists, otherwise asks
// the user for a valid path.
func resolveSourcePath(configured, prompt string) string {
	if configured != "" && directoryExistsPath(configured) {
		return configured
	}

	for {
		current := " (configured: " + configured + ")"
		if configured == "" {
			current = ""
		}
		p := readInput(prompt + current + ": ")
		if p == "" {
			return ""
		}
		if directoryExistsPath(p) {
			return p
		}
		fmt.Printf("\x1b[31mPath not found: %s\x1b[0m\n", p)
	}
}

func directoryExistsPath(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func promptAdditive(prompt string) string {
	return readInput(prompt)
}

func orNA(s string) string {
	if s == "" {
		return "\x1b[31mnot detected\x1b[0m"
	}
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", s)
}

// ----- helpers -----

func readInput(prompt string) string {
	fmt.Print(prompt)
	var s string
	_, _ = fmt.Scanln(&s)
	return strings.TrimSpace(s)
}

func clearScreen() {
	// ANSI clear
	fmt.Print("\x1b[H\x1b[2J")
}

func waitEnter() {
	fmt.Println()
	fmt.Print("Press Enter to continue...")
	_, _ = fmt.Scanln()
}

func init() {
	// Handle Ctrl+C gracefully.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		os.Exit(0)
	}()
}
