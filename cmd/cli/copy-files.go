package main

import (
	"embed"
	"io/ioutil"
)

//go:embed templates
var templateFS embed.FS

func copyFileFromTemplate(templatePath, targetFile string) error {
	// TODO: check to ensure file does not already exist

	// read file from file system
	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		exitGraceFully(err)
	}

	err = copyDataToFile(data, targetFile)
	if err != nil {
		exitGraceFully(err)
	}

	return nil
}

func copyDataToFile(data []byte, to string) error {
	err := ioutil.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
