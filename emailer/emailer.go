package emailer

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"

	//"fmt"
	"html/template"

	"github.com/emersion/go-message/mail"
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

func SendEmail(etd config.EmailTemplateData, smtpData config.SmtpData, authData config.AuthData, addr config.EmailAddressData,
	subject config.EmailSubjectData, templatesData config.EmailTemplatesData) error {

	// write the email we want to send into the customerEmail bytes.Buffer or fail.
	customerEmail, err := newCustomerEmail(etd, addr, subject, templatesData)

	// for the minute...duplicate creating the to and from addresses here (also in newCustomerEmail)
	from := []*mail.Address{{addr.CustomerFromName, addr.CustomerFrom}}
	to := []*mail.Address{{etd.FormData["Name"], etd.FormData["Email"]}}

	toStrs := make([]string, 0)
	for i := range to {
		toStrs = append(toStrs, to[i].Address)
	}
	hostname := smtpData.Host + ":" + strconv.Itoa(smtpData.Port)

	//	do we need auth for this server?
	if authData.Password != "" && authData.Username != "" {
		clientAuth := sasl.NewPlainClient("", authData.Username, authData.Password)
		err = smtp.SendMail(hostname, clientAuth, from[0].String(), toStrs, bytes.NewReader(customerEmail.Bytes()))
	} else {
		// no auth version
		err = smtp.SendMail(hostname, nil, from[0].String(), toStrs, bytes.NewReader(customerEmail.Bytes()))

	}
	if err != nil {
		return err
	}

	return err

}

func newCustomerEmail(etd config.EmailTemplateData, addr config.EmailAddressData,
	subject config.EmailSubjectData, templatesData config.EmailTemplatesData) (bytes.Buffer, error) {
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
		return bytes.Buffer{}, err
	}
	err = cht.Execute(chtbuf, etd)
	if err != nil {
		return bytes.Buffer{}, err
	}
	err = stt.Execute(sttbuf, etd)
	if err != nil {
		return bytes.Buffer{}, err
	}
	err = sht.Execute(shtbuf, etd)
	if err != nil {
		return bytes.Buffer{}, err
	}

	// now build the customer email as a multi part email
	var customerEmail bytes.Buffer
	fmt.Printf("Form Data: %+v\n", etd.FormData)
	from := []*mail.Address{{addr.CustomerFromName, addr.CustomerFrom}}
	to := []*mail.Address{{etd.FormData["Name"], etd.FormData["Email"]}}
	replyTo := []*mail.Address{{addr.CustomerReplyTo, addr.CustomerReplyTo}}
	var h mail.Header
	h.SetDate(time.Now())
	h.SetAddressList("From", from)
	h.SetAddressList("To", to)
	h.SetAddressList("Reply-To", replyTo)
	h.SetSubject(subject.Customer)
	err = h.GenerateMessageIDWithHostname("gophercoders.com") // we need to pass the domain name in somehow.... or do some sort of DNS query??
	if err != nil {
		return bytes.Buffer{}, err
	}

	h.SetContentType("multipart/alternative", nil)
	emailWriter, err := mail.CreateWriter(&customerEmail, h)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer emailWriter.Close()

	htmlWriter, err := emailWriter.CreateInline()
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer htmlWriter.Close()

	var htmlHeader mail.InlineHeader
	htmlHeader.SetContentType("text/html", nil)
	htmlPartWriter, err := htmlWriter.CreatePart(htmlHeader)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer htmlPartWriter.Close()
	_, err = io.Copy(htmlPartWriter, chtbuf)
	if err != nil {
		return bytes.Buffer{}, err
	}

	plainWriter, err := emailWriter.CreateInline()
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer plainWriter.Close()

	var plainHeader mail.InlineHeader
	plainHeader.SetContentType("text/plain", nil)
	plainPartWriter, err := plainWriter.CreatePart(plainHeader)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer plainPartWriter.Close()
	_, err = io.Copy(plainPartWriter, cttbuf)
	if err != nil {
		return bytes.Buffer{}, err
	}

	//log.Println(customerEmail.String())
	return customerEmail, nil
}
