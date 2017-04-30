package omi

import (
	"encoding/xml"
	"fmt"
	"time"
)

var (
	wsTimeLayout       = "2006-01-02T15:04:05-07:00"
	wsCustomTimeLayout = "1/2/2006 04:04:05 PM"
)

type wsTime struct {
	time.Time `xml:",omitempty"`
}

func (c *wsTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(wsTimeLayout, v)
	if err != nil {
		return err
	}
	*c = wsTime{parse}
	return nil
}

type wsCustomTime struct {
	time.Time `xml:",omitempty"`
}

func (c *wsCustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(wsCustomTimeLayout, v)
	if err != nil {
		return err
	}
	*c = wsCustomTime{parse}
	return nil
}

func (c wsTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		S string `xml:",innerxml"`
	}{
		S: fmt.Sprintf("%v", c.Time.Format(wsTimeLayout)),
	}, start)
}

func (c wsCustomTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		S string `xml:",innerxml"`
	}{
		S: fmt.Sprintf("%v", c.Time.Format(wsCustomTimeLayout)),
	}, start)
}
