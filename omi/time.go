package omi

import (
	"time"
	"encoding/xml"
)

type wsTime struct {
	time.Time `xml:",omitempty"`
}

func (c *wsTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	//const shortForm = "1/2/2006 03:04:05 PM" // yyyymmdd date format
	const shortForm = "2006-01-02T15:04:05-07:00"
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = wsTime{parse}
	return nil
}
//2017-03-29T20:40:16-04:00

type wsCustomTime struct {
	time.Time `xml:",omitempty"`
}

func (c *wsCustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "1/2/2006 04:04:05 PM" // yyyymmdd date format
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = wsCustomTime{parse}
	return nil
}