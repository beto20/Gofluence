package java

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	POM_XML  = "pom.xml"
	ROOT_DIR = "."
)

type Document struct {
	Name    string
	Version string
}

func ReadJavaProject(prefix string) {
	currentDir, _ := os.Getwd()
	var documents []Document

	directories := getDirectoriesByPrefix(ROOT_DIR, prefix)

	for _, d := range directories {
		doc := getXmlData(d)
		changeDirectory(currentDir)
		documents = append(documents, doc)
	}

	for _, ds := range documents {
		fmt.Println("Name: %s", ds.Name)
		fmt.Println("Version: %s", ds.Version)
	}
}

func getXmlData(dirname string) Document {
	changeDirectory(dirname)
	return getXmlMetadata()
}

func changeDirectory(dirname string) {
	err := os.Chdir(dirname)
	if err != nil {
		fmt.Println("Failed to change directory")
	}
}

func getDirectoriesByPrefix(dir string, prefix string) []string {
	var filtered []string
	files, err := os.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			filtered = append(filtered, file.Name())
		}
	}

	return filtered
}

func getXmlMetadata() Document {
	xmlFile, err := os.Open(POM_XML)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	var project Project
	err = xml.Unmarshal(byteValue, &project)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
	}

	return Document{
		Name:    project.ArtifactId,
		Version: project.Version,
	}
}
