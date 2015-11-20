package main

func main() {
	wordXML, _ := extractWordXML("HelloWorld.docx")
	wordXML = cleanWordXML(wordXML)

	//collapseElement(doc.Element, "w:rPr")

	repackageWordXML("HelloWorld.docx", "HelloWorld-edited.docx", wordXML)

}
