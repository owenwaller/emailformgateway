package emailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strconv"

	"github.com/owenwaller/email"
	"github.com/owenwaller/emailformgateway/config"
)

func populateTemplate(td config.EmailTemplateData, t string) (string, error) {
	tmpl, err := template.New("email-template").Parse(t)
	if err != nil {
		return "", err
	}
	// create a string writer
	var buf = bytes.NewBufferString("") // buffer implements io.Writer
	err = tmpl.Execute(buf, td)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SendEmail(etd config.EmailTemplateData, smtpData config.SmtpData, addr config.EmailAddressData,
	subject config.EmailSubjectData, templatesData config.EmailTemplatesData) error {

	// now create the templates
	ctt := template.Must(template.ParseFiles(templatesData.CustomerTextFileName))
	cht := template.Must(template.ParseFiles(templatesData.CustomerHtmlFileName))
	stt := template.Must(template.ParseFiles(templatesData.SystemTextFileName))
	sht := template.Must(template.ParseFiles(templatesData.SystemHtmlFileName))
	// now populate the templates - must have set the FormData before this
	var cttbuf = bytes.NewBufferString("") // buffer implements io.Writer
	var chtbuf = bytes.NewBufferString("")
	var sttbuf = bytes.NewBufferString("")
	var shtbuf = bytes.NewBufferString("")
	err := ctt.Execute(cttbuf, etd)
	if err != nil {
		return err
	}
	err = cht.Execute(chtbuf, etd)
	if err != nil {
		return err
	}
	err = stt.Execute(sttbuf, etd)
	if err != nil {
		return err
	}
	err = sht.Execute(shtbuf, etd)
	if err != nil {
		return err
	}

	// now build the emails
	// need to add a reply-to header
	customerEmail := email.NewEmail()
	customerEmail.From = addr.CustomerFromName + "<" + addr.CustomerFrom + ">"
	to := etd.FormData["Name"] + "<" + etd.FormData["Email"] + ">"
	customerEmail.To = []string{to}
	customerEmail.Subject = subject.Customer
	customerEmail.Text = cttbuf.Bytes() // return a []bytes
	customerEmail.HTML = chtbuf.Bytes()
	customerEmail.Headers.Add("Reply-To:", addr.CustomerReplyTo)
	sysEmail := email.NewEmail()

	sysEmail.From = addr.SystemFromName + "<" + addr.SystemFrom + ">"
	to = addr.SystemToName + "<" + addr.SystemTo + ">"
	sysEmail.To = []string{to}
	sysEmail.Subject = subject.System
	sysEmail.Text = sttbuf.Bytes() // return a []bytes
	sysEmail.HTML = shtbuf.Bytes()
	customerEmail.Headers.Add("Reply-To:", addr.SystemReplyTo)

	fmt.Printf("-------\n")
	fmt.Printf("%s\n", customerEmail)
	fmt.Printf("-------\n")
	fmt.Printf("%s\n", sysEmail)
	fmt.Printf("-------\n")
	auth := smtp.PlainAuth("", smtpData.Username, smtpData.Password, smtpData.Host)

	hostname := smtpData.Host + ":" + strconv.Itoa(smtpData.Port)
	err = customerEmail.Send(hostname, auth)
	if err != nil {
		return err
	}
	err = sysEmail.Send(hostname, auth)
	if err != nil {
		return err
	}
	return err

}

//func SendEmail(c Config) error {
//
//}
/*
e := NewEmail()
e.From = "Jordan Wright <test@gmail.com>"
e.To = []string{"test@example.com"}
e.Bcc = []string{"test_bcc@example.com"}
e.Cc = []string{"test_cc@example.com"}
e.Subject = "Awesome Subject"
e.Text = []byte("Text Body is, of course, supported!\n")
e.HTML = []byte("<h1>Fancy Html is supported, too!</h1>\n")
e.Send("smtp.gmail.com:587", smtp.PlainAuth("", e.From, "password123", "smtp.gmail.com"))
*/
