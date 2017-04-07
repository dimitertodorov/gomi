package omi

import (
	"encoding/xml"
	"fmt"
)

type CharData struct {
	string `xml:",innerxml"`
}

func (c *CharData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	*c = CharData{v}
	return nil
}

func (n CharData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct{
		S string `xml:",innerxml"`
	}{
		S: fmt.Sprintf("<![CDATA[%v]]>", n.string),
	}, start)
}
