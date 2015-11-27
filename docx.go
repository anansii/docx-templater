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

func mergeText(someParent etree.Element) {
	// find all paragraphs
	// find all runs in each paragraphs
	// merge text from adjacent runs with same style
	for _, paragraph := range someParent.FindElements("//w:p") {
		var lastRun, thisRun *etree.Element
		for _, run := range paragraph.SelectElements("w:r") {
			lastRun, thisRun = thisRun, run
			if runsHaveSameStyleAndContainText(lastRun, thisRun) {
				mergeRuns(lastRun, thisRun)
			}
		}
	}
}

func mergeRuns(firstRun, secondRun *etree.Element) {
	for _, child := range secondRun.ChildElements() {
		if child.Tag == "w:t" {
			firstTextTag := firstRun.SelectElement("w:t")
			if firstTextTag != nil {
				newText := firstTextTag.Text() + child.Text()
				firstTextTag.SetText(newText)
				firstTextTag.Attr = append(firstTextTag.Attr, child.Attr...) // doubles?
			}
		} else if child.Tag == "w:rPr" {
			// it's already there
		} else {
			firstRun.Child = append(firstRun.Child, child)
		}

	}
}

func runsHaveSameStyleAndContainText(firstRun, secondRun *etree.Element) bool {
	if firstRun.Tag != "w:r" || secondRun.Tag != "w:r" {
		Error.Panicf("Expected nodes of type <w:r>, instead got %s and %s", firstRun.Tag, secondRun.Tag)
	}

	if firstRun.SelectElement("w:t") != nil && secondRun.SelectElement("w:t") != nil {
		return getFlatRunStyle(firstRun) == getFlatRunStyle(secondRun)
	}
	return false
}

func getFlatRunStyle(run *etree.Element) string {
	style := run.SelectElement("w:rPr")
	if style == nil {
		return ""
	}
	return flattenElement(run)
}

func flattenElement(element *etree.Element) string {
	xml := etree.CreateDocument(element)
	xml.Indent(0)
	flatTags, err := xml.WriteToString()
	if err != nil {
		Error.Panicf("Error when trying to flatten %s", element.Tag)
	}
	return flatTags
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

	//docx.Indent(2)
	docx.WriteToFile("wordxml-after.txt")

	result, err := docx.WriteToString()
	if err != nil {
		panic(err)
	}
	return result
}
