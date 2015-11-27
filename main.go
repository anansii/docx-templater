package main

var (
	name    = "DocX Templater"
	version = "0.0.3"
)

func main() {
	setupLogging()
	setupDebugLogging()

	Msg.Println(name, version)
	docx, err := LoadDocx("HelloWorld.docx")
	if err != nil {
		Error.Println(err)
	}

	cleanWordXML(docx.WordXML)

	docx.WriteToFile("HelloWorld-edited.docx")
}
