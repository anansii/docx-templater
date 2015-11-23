package main

func main() {
	docx, err := LoadDocx("HelloWorld.docx")
	if err != nil {
		panic(err)
	}

	cleanWordXML(docx.WordXML)

	docx.WriteToFile("HelloWorld-edited.docx")
}
