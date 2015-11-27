package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
)

// File stores information about a file inside a docx
type File struct {
	Name     string
	Contents []byte
}

// Docx stores the contents of a docx file
type Docx struct {
	Name    string
	Files   []File
	WordXML string
	Writer  *zip.Writer
}

//LoadDocx extracts a docx into memory and allows to process the WordXML
func LoadDocx(source string) (Docx, error) {
	Trace.Printf("Loading %s\n", source)
	var docx Docx
	reader, err := zip.OpenReader(source)
	if err != nil {
		return docx, err
	}

	docx.Name = filepath.Base(source)
	docx.WordXML = "This should not be here."

	for _, file := range reader.File {
		var f File
		f.Name = file.Name

		fileReader, err := file.Open()
		if err != nil {
			return docx, err
		}
		defer fileReader.Close()

		//bytes := make([]byte, file.UncompressedSize64)
		bytes, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return docx, err
		}
		f.Contents = bytes
		if f.Name == "word/document.xml" {
			docx.WordXML = string(f.Contents)
		} else {
			docx.Files = append(docx.Files, f)
		}

	}
	defer reader.Close()

	return docx, nil
}

// WriteToFile creates a zip with all its contents
func (d *Docx) WriteToFile(target string) error {
	docxFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer docxFile.Close()

	d.Writer = zip.NewWriter(docxFile)
	defer d.Writer.Close()

	for _, f := range d.Files {
		zippedFile, err := d.Writer.Create(f.Name)
		if err != nil {
			return err
		}
		zippedFile.Write(f.Contents)
	}

	zippedFile, err := d.Writer.Create("word/document.xml")
	if err != nil {
		return err
	}
	zippedFile.Write([]byte(d.WordXML))

	return nil
}
