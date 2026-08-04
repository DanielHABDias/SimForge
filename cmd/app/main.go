package main

import (
	"fmt"
	"log"

	"github.com/DanielHABDias/SimForge/internal/finder"
)

func main() {

	items, err := finder.Find()

	if err != nil {
		log.Fatal(err)
	}


	for _, item := range items {

		fmt.Println("====================")

		fmt.Println("Platform:", item.Platform)
		fmt.Println("Launcher:", item.Launcher)
		fmt.Println("Name:", item.Name)

		fmt.Println("Prefix:", item.PrefixPath)

		fmt.Println("Client:", item.ClientPath)

		fmt.Println("DLL:", item.DllPath)

		fmt.Println("Config:", item.ConfigPath)

		fmt.Println("Logs:", item.LogsPath)
	}
}