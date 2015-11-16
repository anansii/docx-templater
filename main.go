package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	//"os"
	"errors"
	//"fmt"
	"github.com/beevik/etree"
)

func extractWordXML(wordFile string) (string, error) {
	reader, err := zip.OpenReader(wordFile)
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

func removeAttributes(someParent etree.Element, tag string, attributes ...string) {
	for _, element := range someParent.FindElements("//" + tag) {
		for _, attribute := range attributes {
			element.RemoveAttr(attribute)
		}
	}
}

func removeElement(someParent etree.Element, tag string) {
	for _, element := range someParent.FindElements("//" + tag) {
		element.Parent.RemoveElement(element)
	}
}

func collapseElement(someParent etree.Element, tag string) {
	for _, element := range someParent.FindElements("//" + tag) {
		element.Parent.Child = append(element.Parent.Child, element.Child...)
		element.Parent.RemoveElement(element)
	}
}

func main() {
	wordXML, _ := extractWordXML("HelloWorld.docx")

	doc := etree.NewDocument()
	if err := doc.ReadFromString(wordXML); err != nil {
		panic(err)
	}

	doc.Indent(2)
	doc.WriteToFile("wordxml-before.txt")

	removeElement(doc.Element, "proofErr")
	removeElement(doc.Element, "bookmarkStart")
	removeElement(doc.Element, "bookmarkEnd")

	removeAttributes(doc.Element, "w:p", "w:rsidR", "w:rsidRDefault", "w:rsidRPr")
	removeAttributes(doc.Element, "w:r", "w:rsidR", "w:rsidRDefault", "w:rsidRPr")

	for _, element := range doc.FindElements("//w:t") {
		fmt.Printf("Attributes: %v, Text: [%s]\n", element.Attr, element.Text())

	}

	//collapseElement(doc.Element, "w:rPr")
	doc.WriteToFile("wordxml-after.txt")
}
