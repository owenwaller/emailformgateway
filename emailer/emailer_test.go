package emailer

import (
	"os"
	"testing"

	"github.com/owenwaller/emailformgateway/config"
)

func TestFullPathToTemplates(t *testing.T) {
	var osSep = string(os.PathSeparator)
	var dir = "Templates" + osSep + "Email" + osSep
	var filename = "customer-email.txt"
	var expected = "Templates" + osSep + "Email" + osSep + "customer-email.txt"
	var result = buildTemplateFilename(dir, filename)

	if result != expected {
		t.Fatalf("Did not not get the expected full path filename of the template. Expected %v but got %v\n", expected, result)
	}

	dir = "Templates" + osSep + "Email"
	result = buildTemplateFilename(dir, filename)

	if result != expected {
		t.Fatalf("Did not get the expected full path filename of the template. Expected %v but got %v\n", expected, result)
	}
}

func TestEmailTemplateData(t *testing.T) {
	var td = config.EmailTemplateData{Name: "Joe Blogs",
		Email:    "Joe@blogs.com",
		Subject:  "the feedback subject",
		Feedback: "this is the feedback"}
	tplate := "My name and email address are {{.Name}} and {{.Email}}. " +
		"The subject of my feedback is {{.Subject}}. " +
		"My feedback message is {{.Feedback}}."
	expected := "My name and email address are Joe Blogs and Joe@blogs.com. " +
		"The subject of my feedback is the feedback subject. " +
		"My feedback message is this is the feedback."

	result, err := populateTemplate(td, tplate)
	if err != nil {
		t.Fatalf("Could not populate the test template. Error: %v\n", err)
	}
	if result != expected {
		t.Fatalf("Did not get the expected result of the template expansion. Expected %v but got %v\"", expected, result)
	}
}
