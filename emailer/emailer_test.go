package emailer

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/owenwaller/emailformgateway/config"
	"github.com/spf13/viper"
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
	//var c config.Config
	var td config.EmailTemplateData
	td.FormData = make(map[string]string)
	td.FormData["Name"] = "Email Tester"

	// take the to address for the test email from an env var so we don't have this in the source code
	// We can't use a viper env var binding for this as this comes from the form data
	toAddr, ok := os.LookupEnv("TEST_CUSTOMER_TO_EMAIL_ADDRESS")
	if !ok {
		t.Fatalf("The environmental TEST_CUSTOMER_TO_EMAIL_ADDRESS is undefined.")
	}
	td.FormData["Email"] = toAddr
	td.FormData["Subject"] = "the feedback subject"
	td.FormData["Feedback"] = "this is the feedback"

	// use a viper env var binding to set the System To address and the teamplates directory
	err := viper.BindEnv("Addresses.SystemTo", "TEST_SYSTEM_TO_EMAIL_ADDRESS")
	if err != nil {
		t.Fatalf("Could not bind to TEST_SYSTEM_TO_EMAIL_ADDRESS env var. Error: %s", err)
	}
	err = viper.BindEnv("Templates.Dir", "TEST_TEMPLATES_DIR")
	if err != nil {
		t.Fatalf("Could not bind to TEST_TEMPLATES_DIR env var. Error: %s", err)
	}

	// now try with a real file - via the ENV
	var filename = os.Getenv("TEST_CONFIG_FILE")
	if filename == "" {
		t.Fatalf("Required environmental variable \"TEST_CONFIG_FILE\" not set.\nIt should be the absolute path of the config file.")
	}

	c, err := config.ReadConfig(filename)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("CustomerTo: %q\n", toAddr)
	log.Printf("SystemTo: %q\n", viper.GetString("Addresses.SystemTo"))
	log.Printf("SystemToName: %q\n", viper.GetString("Addresses.SystemToName"))
	log.Printf("Templates Dir: %q\n", viper.GetString("Templates.Dir"))
	log.Printf("Config File: %q\n", viper.ConfigFileUsed())

	// first get all the templates

	c.Templates.CustomerTextFileName = config.BuildTemplateFilename(c.Templates.Dir, c.Templates.CustomerText)
	c.Templates.CustomerHtmlFileName = config.BuildTemplateFilename(c.Templates.Dir, c.Templates.CustomerHtml)
	c.Templates.SystemTextFileName = config.BuildTemplateFilename(c.Templates.Dir, c.Templates.SystemText)
	c.Templates.SystemHtmlFileName = config.BuildTemplateFilename(c.Templates.Dir, c.Templates.SystemHtml)

	var expectedCtt = config.BuildTemplateFilename(c.Templates.Dir, "customer-email-text.template")
	var expectedCht = config.BuildTemplateFilename(c.Templates.Dir, "customer-email-html.template")
	var expectedStt = config.BuildTemplateFilename(c.Templates.Dir, "system-email-text.template")
	var expectedSht = config.BuildTemplateFilename(c.Templates.Dir, "system-email-html.template")

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

	formatted := time.Now().Format(time.RFC3339)
	// add the time/date to the subject lines
	c.Subjects.Customer = fmt.Sprintf("%s %s", formatted, c.Subjects.Customer)
	c.Subjects.System = fmt.Sprintf("%s %s", formatted, c.Subjects.System)

	err = SendEmail(td, c.Smtp, c.Auth, c.Addresses,
		c.Subjects, c.Templates)
	if err != nil {
		t.Fatalf("unexpected error sending email %v\n", err)
	}

}
