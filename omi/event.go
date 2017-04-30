package omi

import (
	"encoding/xml"
	"fmt"
	"reflect"
)

var (
	// ErrNotFound is returned if a provider cannot find a requested item.
	ErrNotFound = fmt.Errorf("item not found")
)

type EventList struct {
	XMLName xml.Name `xml:"event_list"`
	XmlType string   `json:"type" xml:"type,attr"`
	Event   []Event  `json:"event" xml:"event"`
	Client  *Client
}

type Event struct {
	XMLName             xml.Name             `json:"-" xml:"event"`
	Xmlns               string               `json:"-" xml:"xmlns,attr"`
	Id                  string               `json:"id,omitempty" xml:"id,omitempty"`
	Title               string               `json:"title,omitempty" xml:"title,omitempty"`
	Description         string               `json:"description,omitempty" xml:"description,omitempty"`
	Key                 string               `json:"key,omitempty" xml:"key,omitempty"`
	State               string               `json:"state,omitempty" xml:"state,omitempty"`
	Severity            string               `json:"severity,omitempty" xml:"severity,omitempty"`
	Category            string               `json:"category,omitempty" xml:"category,omitempty"`
	Application         string               `json:"application,omitempty" xml:"application,omitempty"`
	Object              string               `json:"object,omitempty" xml:"object,omitempty"`
	DuplicateCount      int                  `json:"duplicate_count,omitempty" xml:"duplicate_count,omitempty"`
	CloseKeyPattern     string               `json:"close_key_pattern,omitempty" xml:"close_key_pattern,omitempty"`
	TimeCreated         *wsTime              `json:"time_created,omitempty" xml:"time_created,omitempty"`
	TimeCreatedLabel    *wsCustomTime        `json:"time_created_label,omitempty" xml:"time_created_label,omitempty"`
	CustomAttributeList *CustomAttributeList `json:"custom_attribute_list,omitempty" xml:",omitempty"`
	RelatedCiHints      *RelatedCiHints      `json:"related_ci_hints,omitempty" xml:",omitempty"`
	RelatedCi           *RelatedCi           `json:"related_ci,omitempty" xml:",omitempty"`
	MatchInfo           *MatchInfo           `json:"match_info,omitempty" xml:",omitempty"`
	Client              *Client              `json:"-" xml"-"`
}

var (
	READ_ONLY_FIELDS = [...]string{
		"RelatedCi",
		"CustomAttributeList",
		"TimeCreated",
		"TimeCreatedLabel",
		"DuplicateCount",
	}
	EventXmlns   = "http://www.hp.com/2009/software/opr/data_model"
	DefaultEvent = Event{
		Xmlns: EventXmlns,
	}
	DefaultCustomAttribute = CustomAttribute{
		Xmlns: EventXmlns,
	}
	DefaultAnnotation = Annotation{
		Xmlns: EventXmlns,
	}
	DefaultAnnotationList = AnnotationList{
		Xmlns: EventXmlns,
	}
)

type AssignedUser struct {
	XMLName   xml.Name `json:"-" xml:"assigned_user"`
	Id        string   `json:"id" xml:"id"`
	LoginName string   `json:"login_name,omitempty" xml:"login_name,omitempty"`
	UserName  string   `json:"user_name,omitempty" xml:"user_name,omitempty"`
}

type AssignedGroup struct {
	XMLName xml.Name `json:"-" xml:"assigned_group"`
	Id      string   `json:"id" xml:"id"`
	Name    string   `json:"name" xml:"name"`
}

type CustomAttribute struct {
	XMLName xml.Name `json:"-" xml:"custom_attribute"`
	Xmlns   string   `json:"-" xml:"xmlns,attr"`
	XmlType string   `json:"type" xml:"type,attr"`
	Id      string   `json:"id" xml:"id,omitempty"`
	Name    string   `json:"name" xml:"name"`
	Value   string   `json:"value" xml:"value"`
}

type CustomAttributeList struct {
	XMLName          xml.Name          `json:"-" xml:"custom_attribute_list"`
	XmlType          string            `json:"type" xml:"type,attr"`
	CustomAttributes []CustomAttribute `json:"custom_attribute" xml:"custom_attribute"`
}

type Annotation struct {
	XMLName     xml.Name `json:"-" xml:"annotation"`
	Xmlns       string   `json:"-" xml:"xmlns,attr"`
	XmlType     string   `json:"type,omitempty" xml:"type,attr,omitempty"`
	Id          string   `json:"id,omitempty" xml:"id,omitempty"`
	Author      string   `json:"author" xml:"author"`
	Text        string   `json:"text" xml:"text"`
	TimeCreated *wsTime  `json:"time_created,omitempty" xml:"time_created,omitempty"`
}

type AnnotationList struct {
	XMLName     xml.Name     `json:"-" xml:"annotation_list"`
	Xmlns       string       `json:"-" xml:"xmlns,attr"`
	XmlType     string       `json:"type" xml:"type,attr"`
	Annotations []Annotation `json:"custom_attribute" xml:"custom_attribute"`
}

type RelatedCiHints struct {
	XMLName xml.Name `json:"-" xml:"related_ci_hints"`
	XmlType string   `json:"type" xml:"type,attr"`
	Hint    []string `json:"hint" xml:"hint"`
}

type RelatedCi struct {
	XMLName           xml.Name           `json:"-" xml:"related_ci,omitempty"`
	XmlType           string             `json:"type,omitempty" xml:"type,attr,omitempty"`
	TargetId          string             `json:"target_id,omitempty" xml:"target_id,omitempty"`
	TargetType        string             `json:"target_type,omitempty" xml:"target_type,omitempty"`
	ConfigurationItem *ConfigurationItem `json:"configuration_item,omitempty" xml:",omitempty"`
}

type ConfigurationItem struct {
	XMLName          xml.Name `json:"-" xml:"configuration_item,omitempty"`
	XmlType          string   `json:"type,omitempty" xml:"type,attr,omitempty"`
	Id               string   `json:"id,omitempty" xml:"id,omitempty"`
	Name             string   `json:"name,omitempty" xml:"name,omitempty"`
	Type             string   `json:"type,omitempty" xml:"type,omitempty"`
	TypeLabel        string   `json:"type_label,omitempty" xml:"type_label,omitempty"`
	PrimaryDnsName   string   `json:"primary_dns_name,omitempty" xml:"primary_dns_name,omitempty"`
	DiscoveredOsName string   `json:"discovered_os_name,omitempty" xml:"discovered_os_name,omitempty"`
}

type MatchInfo struct {
	XMLName       xml.Name `json:"-" xml:"match_info,omitempty"`
	PolicyType    string   `json:"policy_type,omitempty" xml:"policy_type,attr,omitempty"`
	PolicyName    string   `json:"policy_name,omitempty" xml:"policy_name,attr,omitempty"`
	ConditionId   string   `json:"condition_id,omitempty" xml:"condition_id,attr,omitempty"`
	ConditionName string   `json:"condition_name,omitempty" xml:"condition_name,attr,omitempty"`
}

func (self *Event) wipe(fieldName string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in wipe %v", r)
		}
	}()
	indir := reflect.ValueOf(self).Elem().FieldByName(fieldName)
	indir.Set(reflect.Zero(indir.Type()))
}

func (self *Event) clean() {
	for _, field := range READ_ONLY_FIELDS {
		self.wipe(field)
	}
}
