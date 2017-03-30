package omi

import (
	"fmt"
	"encoding/xml"
)

var (
	// ErrNotFound is returned if a provider cannot find a requested item.
	ErrNotFound = fmt.Errorf("item not found")
)

type EventList struct {
	XMLName xml.Name        `xml:"event_list"`
	XmlType string                `json:"type" xml:"type,attr"`
	Event   []Event `json:"event" xml:"event"`
}

type Event struct {
	XMLName             xml.Name        `xml:"event"`
	Xmlns               string        `xml:"xmlns,attr"`
	Id                  string        `json:"id" xml:"id,omitempty"`
	Title               string        `json:"title" xml:"title,omitempty"`
	Key                 string `json:"key,omitempty" xml:"key,omitempty"`
	CloseKeyPattern     CharData `json:"close_key_pattern,omitempty" xml:"close_key_pattern,omitempty"`
	TimeCreated         wsTime        `json:"time_created" xml:"time_created,omitempty"`
	TimeCreatedLabel    wsCustomTime        `json:"time_created_label" xml:"time_created_label,omitempty"`
	CustomAttributeList *CustomAttributeList        `json:",omitempty" xml:",omitempty"`
	RelatedCiHints      *RelatedCiHints        `json:",omitempty" xml:",omitempty"`
	RelatedCi           *RelatedCi                `json:",omitempty" xml:",omitempty"`
}

type AssignedUser struct {
	XMLName   xml.Name        `xml:"assigned_user"`
	Id        string        `json:"id" xml:"id"`
	LoginName string        `json:"login_name,omitempty" xml:"login_name,omitempty"`
	UserName  string        `json:"user_name,omitempty" xml:"user_name,omitempty"`
}

type AssignedGroup struct {
	XMLName xml.Name      `xml:"assigned_group"`
	Id      string        `json:"id" xml:"id"`
	Name    string        `json:"name" xml:"name"`
}

type CustomAttribute struct {
	XMLName xml.Name      `xml:"custom_attribute"`
	XmlType string        `json:"type" xml:"type,attr"`
	Name    string        `json:"name" xml:"name"`
	Value   string        `json:"value" xml:"value"`
}

type CustomAttributeList struct {
	XMLName         xml.Name        `xml:"custom_attribute_list"`
	XmlType         string            `json:"type" xml:"type,attr"`
	CustomAttribute []CustomAttribute `json:"custom_attribute" xml:"custom_attribute"`
}

type RelatedCiHints struct {
	XMLName xml.Name        `xml:"related_ci_hints"`
	XmlType string          `json:"type" xml:"type,attr"`
	Hint    []string        `json:"hint" xml:"hint"`
}

type RelatedCi struct {
	XMLName    xml.Name        `xml:"related_ci,omitempty"`
	XmlType    string          `json:"type,omitempty" xml:"type,attr,omitempty"`
	TargetId   string                 `json:"target_id,omitempty" xml:"target_id,omitempty"`
	TargetType string                 `json:"target_type,omitempty" xml:"target_type,omitempty"`
}



