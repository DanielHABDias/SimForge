//go:build windows

package finder


import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)



func find()([]Installation,error){


	clients := []struct{
		Name string
		Key string
	}{
		{
			Name:"EA App",
			Key:`SOFTWARE\Electronic Arts\EA Desktop`,
		},
		{
			Name:"Origin",
			Key:`SOFTWARE\WOW6432Node\Origin`,
		},
		{
			Name:"Origin",
			Key:`SOFTWARE\Origin`,
		},
	}



	var result []Installation



	for _,client := range clients{


		path,err:=getRegistryPath(client.Key)

		if err!=nil{
			continue
		}


		result=append(result,
			createWindowsInstallation(
				client.Name,
				path,
			),
		)
	}


	return result,nil
}



func getRegistryPath(
	keyPath string,
)(string,error){


	key,err:=registry.OpenKey(
		registry.LOCAL_MACHINE,
		keyPath,
		registry.QUERY_VALUE,
	)

	if err!=nil{
		return "",err
	}


	defer key.Close()


	return key.GetStringValue("ClientPath")
}



func createWindowsInstallation(
	name string,
	client string,
)Installation{


	appData:=os.Getenv("APPDATA")
	local:=os.Getenv("LOCALAPPDATA")

	common:=filepath.Join(
		"anadius",
		"EA DLC Unlocker v2",
	)


	return Installation{

		Platform:"windows",

		Launcher:name,


		ClientPath:client,


		DllPath:
			filepath.Join(
				client,
				"version.dll",
			),


		ConfigPath:
			filepath.Join(
				appData,
				common,
			),


		LogsPath:
			filepath.Join(
				local,
				common,
			),


		User:
			os.Getenv("USERNAME"),
	}
}