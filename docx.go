package main

import "github.com/beevik/etree"

func removeAttributes(someParent etree.Element, tag string, attributes ...string) {
	for _, element := range someParent.FindElements("//" + tag) {
		for _, attribute := range attributes {
			element.RemoveAttr(attribute)
		}
	}
}

func removeAttributesForEmptyTags(someParent etree.Element, tag string, attributes ...string) {
	for _, element := range someParent.FindElements("//" + tag) {
		for _, attribute := range attributes {
			if element.Text() == "" {
				element.RemoveAttr(attribute)
			}
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

func runsHaveSameStyle(firstNode, secondNode etree.Element) {

}

func cleanWordXML(wordXML string) string {
	docx := etree.NewDocument()
	if err := docx.ReadFromString(wordXML); err != nil {
		panic(err)
	}

	docx.Indent(2)
	docx.WriteToFile("wordxml-before.txt")

	removeElement(docx.Element, "proofErr")
	removeElement(docx.Element, "bookmarkStart")
	removeElement(docx.Element, "bookmarkEnd")

	removeAttributes(docx.Element, "w:p", "w:rsidR", "w:rsidRDefault", "w:rsidRPr")
	removeAttributes(docx.Element, "w:r", "w:rsidR", "w:rsidRDefault", "w:rsidRPr")
	//removeAttributesForEmptyTags(docx.Element, "w:t", "xml:space")

	//for _, element := range docx.FindElements("//w:t") {
	//	fmt.Printf("Attributes: %v, Text: [%s]\n", element.Attr, element.Text())
	//}

	for _, element := range docx.FindElements("//w:t") {
		element.Parent.RemoveElement(element)
	}

	//docx.Indent(2)
	docx.WriteToFile("wordxml-after.txt")

	result, err := docx.WriteToString()
	if err != nil {
		panic(err)
	}
	return result
}
