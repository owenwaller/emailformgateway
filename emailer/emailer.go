package emailer

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/owenwaller/emailformgateway/config"
)

func buildTemplateFilename(dir, filename string) string {
	return filepath.Join(dir, filename)
}

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
