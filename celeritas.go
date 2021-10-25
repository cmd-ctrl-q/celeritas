package celeritas

const version = "1.0.0"

type Celeritas struct {
	AppName string

	// If Debug is set to true then the application
	// is running in development mode, otherwise production.
	Debug   bool // dev
	Version string
}

func (c *Celeritas) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := c.Init(pathConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if doesn't exist
		err := c.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}
