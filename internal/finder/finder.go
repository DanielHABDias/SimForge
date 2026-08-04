package finder


type Installation struct {

	Platform string

	Launcher string

	Name string

	AppID string

	PrefixPath string

	ClientPath string

	User string

	DllPath string

	ConfigPath string

	LogsPath string
}



func Find() ([]Installation,error){
	return find()
}