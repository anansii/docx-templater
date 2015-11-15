package main

import (
	"archive/zip"
	"io/ioutil"
  //"os"
  "errors"
  //"fmt"
  "github.com/beevik/etree"
)

func extractWordXml(wordFile string) (string, error) {
	reader, err := zip.OpenReader(wordFile)
	if err != nil {
		return "", err
	}

	for _, file := range reader.File {
    if (file.Name != "word/document.xml") {
      continue
    }

    fileReader, err := file.Open();
    if err != nil {
			return "", err
		}
		defer fileReader.Close()

    xmlAsBytes, err := ioutil.ReadAll(fileReader);
    if err != nil {
      return "", err
    }

    return string(xmlAsBytes), nil
	}
  return "", errors.New("No WordML data found..")
}

func removeAttributes(someParent etree.Element, tag string, attributes ...string) {
  for _, element := range someParent.FindElements("//"+ tag) {
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

func collapseElement(someParent *etree.Element, tag string) {
  for _, element := range someParent.FindElements("//" + tag) {
      for _, child := range element.ChildElements() {
          someParent.addChild(child)

      }
      element.Parent.RemoveElement(element)
  }
}


func main() {
  wordXml, _ := extractWordXml("HelloWorld.docx")

  doc := etree.NewDocument()
  if err := doc.ReadFromString(wordXml); err != nil {
    panic(err)
  }

  doc.Indent(2)
  doc.WriteToFile("wordxml-before.txt")

  removeElement(doc.Element, "proofErr")
  removeElement(doc.Element, "bookmarkStart")
  removeElement(doc.Element, "bookmarkEnd")

  removeAttributes(doc.Element, "w:p", "w:rsidR", "w:rsidRDefault", "w:rsidRPr")
  removeAttributes(doc.Element, "w:r", "w:rsidR", "w:rsidRDefault", "w:rsidRPr")

  collapseElement(&doc.Element, "w:rPr")
  doc.WriteToFile("wordxml-after.txt")
}
