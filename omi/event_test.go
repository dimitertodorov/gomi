package omi

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

var (
	eventXml = "../testdata/event.xml"
)

func TestUnmarshalEvent(t *testing.T) {
	var event Event

	externalEvent, err := ioutil.ReadFile(eventXml)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if err := xml.Unmarshal(externalEvent, &event); err != nil {
		t.Fatalf("Error Unmarshalling Event: %s", err)
	}
	assert.Equal(t, "eth0", event.Object, "they should be equal")
	assert.Equal(t, "^Sys_NetworkInterfaceErrorDiagnosis:itslaboml02:eth0:<*>", event.CloseKeyPattern, "CDATA Should unmarshal properly")
}

func TestMarshalEvent(t *testing.T) {
	event := DefaultEvent
	event.Title = "NEW TITLE <*>"
	event.CloseKeyPattern = "potatoes:<*>"
	result, err := xml.Marshal(event)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	assert.NotEqual(t, -1, strings.Index(string(result[:]), "<close_key_pattern>potatoes"), "Should Serialize Event properly")
}

func TestMarshalAnnotation(t *testing.T) {
	ann := DefaultAnnotation
	ann.Author = "Dimiter"
	ann.Text = "TestText"
	result, err := xml.Marshal(ann)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	assert.NotEqual(t, -1, strings.Index(string(result[:]), "<author>Dimiter"), "Should Serialize Annotation properly")

}

func TestDefaultEvent(t *testing.T) {
	var event = DefaultEvent
	assert.Equal(t, event.Xmlns, DefaultEvent.Xmlns)
}

func TestDefaultCustomAttribute(t *testing.T) {
	var ca = DefaultCustomAttribute
	assert.Equal(t, ca.Xmlns, DefaultCustomAttribute.Xmlns)
}
