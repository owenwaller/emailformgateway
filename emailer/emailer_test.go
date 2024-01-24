package emailer

import (
	"fmt"
	"os"
	"testing"

	"github.com/owenwaller/emailformgateway/config"
)

func TestFullPathToTemplates(t *testing.T) {
	var osSep = string(os.PathSeparator)
	var dir = "Templates" + osSep + "Email" + osSep
	var filename = "customer-email.txt"
	var expected = "Templates" + osSep + "Email" + osSep + "customer-email.txt"
	var result = config.BuildTemplateFilename(dir, filename)

	if result != expected {
		t.Fatalf("Did not not get the expected full path filename of the template. Expected %v but got %v\n", expected, result)
	}

	dir = "Templates" + osSep + "Email"
	result = config.BuildTemplateFilename(dir, filename)

	if result != expected {
		t.Fatalf("Did not get the expected full path filename of the template. Expected %v but got %v\n", expected, result)
	}
}

func TestEmailTemplateMapData(t *testing.T) {
	// create map based on the field names
	var td config.EmailTemplateData
	td.FormData = make(map[string]string)
	td.FormData["Name"] = "Joe Blogs"
	td.FormData["Email"] = "Joe@blogs.com"
	td.FormData["Subject"] = "the feedback subject"
	td.FormData["Feedback"] = "this is the feedback"
	tplate := "My name and email address are {{.FormData.Name}} and {{.FormData.Email}}. " +
		"The subject of my feedback is {{.FormData.Subject}}. " +
		"My feedback message is {{.FormData.Feedback}}."
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

func TestSendEmail(t *testing.T) {
	// read a config
	var c config.Config
	var td config.EmailTemplateData
	td.FormData = make(map[string]string)
	td.FormData["Name"] = "Email Tester"
	td.FormData["Email"] = "owenwaller@gmail.com"
	td.FormData["Subject"] = "the feedback subject"
	td.FormData["Feedback"] = "this is the feedback"

	// first get all the templates
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("PWD: %q\n", cwd)

	// now try with a real file - via the ENV
	var filename = os.Getenv("TEST_CONFIG_FILE")
	if filename == "" {
		t.Fatalf("Required enviromental variable \"TEST_CONFIG_FILE\" not set.\nIt should be the absolute path of the config file.")
	}
	err = config.ReadConfig(filename, &c)
	if err != nil {
		t.Fatalf("unexpected error reading config %v\n", err)
	}

	// first get all the templates

	c.Templates.CustomerTextFileName = config.BuildTemplateFilename(cwd, c.Templates.CustomerText)
	c.Templates.CustomerHtmlFileName = config.BuildTemplateFilename(cwd, c.Templates.CustomerHtml)
	c.Templates.SystemTextFileName = config.BuildTemplateFilename(cwd, c.Templates.SystemText)
	c.Templates.SystemHtmlFileName = config.BuildTemplateFilename(cwd, c.Templates.SystemHtml)

	var expectedCtt = config.BuildTemplateFilename(cwd, "customer-email-text.template")
	var expectedCht = config.BuildTemplateFilename(cwd, "customer-email-html.template")
	var expectedStt = config.BuildTemplateFilename(cwd, "system-email-text.template")
	var expectedSht = config.BuildTemplateFilename(cwd, "system-email-html.template")

	if c.Templates.CustomerTextFileName != expectedCtt {
		t.Fatalf("Did not get expected filename. Expected: %v Got %v\n", expectedCtt, c.Templates.CustomerTextFileName)
	}
	if c.Templates.CustomerHtmlFileName != expectedCht {
		t.Fatalf("Did not get expected filename. Expected: %v Got %v\n", expectedCht, c.Templates.CustomerHtmlFileName)
	}
	if c.Templates.SystemTextFileName != expectedStt {
		t.Fatalf("Did not get expected filename. Expected: %v Got %v\n", expectedStt, c.Templates.SystemTextFileName)
	}
	if c.Templates.SystemHtmlFileName != expectedSht {
		t.Fatalf("Did not get expected filename. Expected: %v Got %v\n", expectedSht, c.Templates.SystemHtmlFileName)
	}

	err = SendEmail(td, c.Smtp, c.Auth, c.Addresses,
		c.Subjects, c.Templates)
	if err != nil {
		t.Fatalf("unexpected error sending email %v\n", err)
	}

}
