package main

type Loader interface {
	SetParam(interface{}) error
	Load() (string, error)
}

type ConfigByStringsLoader interface {
	Loader
	ConfigByStrings(args []string) error
}

type ConfigByStringsLoaderGenerator func() ConfigByStringsLoader
