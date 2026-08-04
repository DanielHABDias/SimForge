//go:build linux

package finder


import (
	"os"
	"path/filepath"
)



func find()([]Installation,error){


	var result []Installation


	result =
		append(result,
			findWine()...,
		)


	result =
		append(result,
			findSteam()...,
		)


	result =
		append(result,
			findLutris()...,
		)


	result =
		append(result,
			findBottles()...,
		)



	return result,nil
}



func checkPrefix(
	path string,
	launcher string,
	name string,
)(*Installation){


	ea:=filepath.Join(
		path,
		"drive_c",
		"Program Files",
		"Electronic Arts",
		"EA Desktop",
	)


	if _,err:=os.Stat(ea);err!=nil{
		return nil
	}



	user:=os.Getenv("USER")


	return &Installation{

		Platform:"linux",

		Launcher:launcher,

		Name:name,


		PrefixPath:path,


		ClientPath:ea,


		DllPath:
			filepath.Join(
				ea,
				"version.dll",
			),


		ConfigPath:
			filepath.Join(
				path,
				"drive_c",
				"users",
				user,
				"AppData",
				"Roaming",
				"anadius",
				"EA DLC Unlocker v2",
			),


		LogsPath:
			filepath.Join(
				path,
				"drive_c",
				"users",
				user,
				"AppData",
				"Local",
				"anadius",
				"EA DLC Unlocker v2",
			),
	}
}

func findWine()[]Installation{


	home,_:=os.UserHomeDir()


	prefix:=filepath.Join(
		home,
		".wine",
	)


	if item:=checkPrefix(
		prefix,
		"Wine",
		"Default Wine",
	); item!=nil{

		return []Installation{*item}
	}


	return nil
}

func findLutris()[]Installation{


	home,_:=os.UserHomeDir()


	base:=filepath.Join(
		home,
		"Games",
	)


	var result []Installation


	files,_:=os.ReadDir(base)


	for _,file:=range files{


		p:=filepath.Join(
			base,
			file.Name(),
		)


		if item:=checkPrefix(
			p,
			"Lutris",
			file.Name(),
		);item!=nil{

			result=append(result,*item)
		}

	}


	return result
}

func findBottles()[]Installation{


	home,_:=os.UserHomeDir()


	paths:=[]string{

		".local/share/bottles/bottles",

		".var/app/com.usebottles.bottles/data/bottles/bottles",
	}



	var result []Installation


	for _,base:=range paths{


		base=filepath.Join(home,base)


		files,_:=os.ReadDir(base)


		for _,file:=range files{


			p:=filepath.Join(
				base,
				file.Name(),
			)


			if item:=checkPrefix(
				p,
				"Bottles",
				file.Name(),
			);item!=nil{

				result=append(result,*item)
			}
		}
	}


	return result
}