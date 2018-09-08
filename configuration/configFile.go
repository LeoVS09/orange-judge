package configuration

type ConfigFile struct {
	Port             int               `json:"port"`
	Directories      ConfigDirectories `json:"directories"`
	TestListFileName string            `json:"testListFileName"`
}

type ConfigDirectories struct {
	Uploaded string `json:"uploaded"`
	Compiled string `json:"compiled"`
	Test     string `json:"test"`
}
