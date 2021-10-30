package celeritas

type initPaths struct {
	rootPath    string
	folderNames []string
}

type cookieConfig struct {
	name     string
	lifetime string

	// does it persist between browswer closes
	persist string

	// is the cookie encrypted
	secure string

	// the domain the cookie is associated with
	domain string
}
