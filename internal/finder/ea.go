package finder

import (
	"os"
	"path/filepath"
)


func findEA(prefixPath string) *EAInstallation {
	// Possíveis caminhos do EA App dentro do Wine/Proton
	possibleClients := []string{

		filepath.Join(
			prefixPath,
			"drive_c",
			"Program Files",
			"Electronic Arts",
			"EA Desktop",
			"EA Desktop",
		),


		filepath.Join(
			prefixPath,
			"drive_c",
			"Program Files",
			"Electronic Arts",
			"EA Desktop",
		),


		filepath.Join(
			prefixPath,
			"drive_c",
			"Program Files",
			"Electronic Arts",
			"EADesktop",
		),


		filepath.Join(
			prefixPath,
			"drive_c",
			"Program Files (x86)",
			"Electronic Arts",
			"EA Desktop",
		),
	}



	var clientPath string



	for _,path := range possibleClients {

		if directoryExists(path){

			clientPath = path
			break

		}
	}



	if clientPath == "" {
		return nil
	}



	user :=
		findWineUser(
			prefixPath,
		)



	configPath :=
		filepath.Join(
			prefixPath,
			"drive_c",
			"users",
			user,
			"AppData",
			"Roaming",
			"anadius",
			"EA DLC Unlocker v2",
		)



	logPath :=
		filepath.Join(
			prefixPath,
			"drive_c",
			"users",
			user,
			"AppData",
			"Local",
			"anadius",
			"EA DLC Unlocker v2",
		)



	return &EAInstallation{

		ClientPath: clientPath,

		DllPath:
			filepath.Join(
				clientPath,
				"version.dll",
			),


		ConfigPath: configPath,

		LogPath: logPath,

		User: user,
	}
}



func findWineUser(
	prefixPath string,
) string {


	usersPath :=
		filepath.Join(
			prefixPath,
			"drive_c",
			"users",
		)



	users,err :=
		os.ReadDir(
			usersPath,
		)



	if err != nil {
		return "user"
	}



	for _,user := range users {


		if user.IsDir() {


			name :=
				user.Name()


			// ignora usuários padrões do Wine
			if name != "Public" &&
				name != "Default" &&
				name != "defaultuser0" {


				return name
			}
		}
	}

	return "user"
}