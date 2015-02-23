package config

type Config struct {
	Server    ServerData
	LogFile   LogFileData
	Smtp      SmtpData
	Addresses EmailAddressData
	Subjects  EmailSubjectData
	Templates EmailTemplatesData
	Fields    map[string]FieldData
}

type ServerData struct {
	Host string
	Path string
	Port int
}

type LogFileData struct {
	Filename string
	Path     string
	Level    string
}

type SmtpData struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailAddressData struct {
	CustomerFrom     string
	CustomerFromName string
	CustomerReplyTo  string
	SystemTo         string
	SystemToName     string
	SystemFrom       string
	SystemFromName   string
	SystemReplyTo    string
}

type EmailSubjectData struct {
	Customer string
	System   string
}

type EmailTemplatesData struct {
	Dir          string
	CustomerText string
	CustomerHtml string
	SystemText   string
	SystemHtml   string
}

type FieldData struct {
	Name string
	Type string
}

type EmailTemplateData struct {
	Name     string
	Email    string
	Subject  string
	Feedback string
}
