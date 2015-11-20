package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name     string
	Contents []byte
}

type Folder struct {
	Name    string
	Files   []File
	Folders []Folder
	Root    bool
}

func extractZip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		fmt.Printf("%s, %s\n", path, filepath.Dir(path))
		os.MkdirAll(filepath.Dir(path), 0755)

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

// createZip
func createZip(sourceFolder, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(sourceFolder)
	if err != nil {
		panic(err)
	}

	if !info.IsDir() {
		return errors.New("CreateZip: Source is not a folder!")
	}

	filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, sourceFolder+string(os.PathSeparator))

		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func extractWordXML(docxPath string) (string, error) {
	reader, err := zip.OpenReader(docxPath)
	if err != nil {
		return "", err
	}

	for _, file := range reader.File {
		if file.Name != "word/document.xml" {
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return "", err
		}
		defer fileReader.Close()

		xmlAsBytes, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return "", err
		}

		return string(xmlAsBytes), nil
	}
	return "", errors.New("No WordML data found..")
}

func repackageWordXML(templatePath, targetPath, wordXML string) {
	// unpack docx to temp
	extractZip(templatePath, "./temp")
	xmlFile, err := os.Create("temp/word/document.xml")
	if err != nil {
		panic(err)
	}

	// replace xml
	_, err = xmlFile.WriteString(wordXML)
	// repack docx
	createZip("temp", "HelloWorld-edited.docx")
	os.RemoveAll("temp")
}
