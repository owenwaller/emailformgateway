package validation

import (
	"fmt"
	"html/template"
	"net/mail"
	"regexp"
)

const (
	Control     = 'C'
	Math        = 'M'
	Number      = 'N'
	Letter      = 'L'
	Punctuation = 'P'
	Symbol      = 'S'
	Separator   = 'Z'
)

func RemoveEscapeSequences(s string) string {
	// replace all control/other characters based on the unicode class
	// this also wipes out tabs
	re := regexp.MustCompile(`(?im:\pC)+`)
	s = re.ReplaceAllLiteralString(s, "")
	return s
}

func RemoveEmailHeaders(s string) string {
	re := regexp.MustCompile("(?im:To|From|Bcc|Cc|Reply-To|Sender):")
	s = re.ReplaceAllString(s, "")
	return s
}

func RemoveScriptTagsAndContents(s string) string {
	re := regexp.MustCompile(`(?im:(<|\&lt;)\s*script\s*(>|\&gt;).*?(<|\&lt;)\s*/\s*script\s*(>|\&gt;))`)
	s = re.ReplaceAllLiteralString(s, "")
	return s
}

func EscapeHTML(s string) string {
	return template.HTMLEscapeString(s)
}

func AcceptUnicodeLettersSpacesPunctuation(s string) bool {
	req := map[rune]bool{
		Letter:      true,
		Separator:   true,
		Punctuation: true,
		Control:     false,
		Math:        false,
		Number:      false,
		Symbol:      false}
	ans := acceptUnicode(s, req)
	return ans
}

func AcceptAllUnicodeExceptControl(s string) bool {
	req := map[rune]bool{
		Letter:      true,
		Separator:   true,
		Punctuation: true,
		Control:     false,
		Math:        true,
		Number:      true,
		Symbol:      true}
	ans := acceptUnicode(s, req)
	return ans
}

func ValidateAsEmail(s string) bool {
	//	match.Value = "Mr Blah Blah <" + match.Value + ">"
	_, err := mail.ParseAddress(s)
	if err != nil {
		fmt.Printf("Could not parse email address \"%s\".Error: %s\n", s, err)
		return false
	}
	return true
}

func ValidateAsRestrictedText(s string) bool {
	var accept = AcceptUnicodeLettersSpacesPunctuation(s)
	return accept
}

func ValidateAsUnrestrictedText(s string) bool {
	var accept = AcceptAllUnicodeExceptControl(s)
	return accept
}

func acceptUnicode(s string, cc map[rune]bool) bool {
	present := map[rune]bool{
		Letter:      false,
		Separator:   false,
		Punctuation: false,
		Control:     false,
		Math:        false,
		Number:      false,
		Symbol:      false}

	re := regexp.MustCompile(`(?m:\pL)+`)
	present[Letter] = re.MatchString(s)

	re = regexp.MustCompile(`(?m:[\pZ|` + "\t|\r|\n" + `])+`)
	present[Separator] = re.MatchString(s)

	//re = regexp.MustCompile(`(?m:\pC)+`)
	// anything apart from \pL\pZ\pM\pN\pS\pP + tab + CR + LF
	// Unicode \PC includes Tab LF CR which is not what we want
	re = regexp.MustCompile(`(?m:[^\pL|\pZ|\pM|\pN|\pS|\pP|` + "\t|\r|\n" + `])+`)
	present[Control] = re.MatchString(s)

	re = regexp.MustCompile(`(?m:\pM)+`)
	present[Math] = re.MatchString(s)

	re = regexp.MustCompile(`(?m:\pN)+`)
	present[Number] = re.MatchString(s)

	re = regexp.MustCompile(`(?m:\pS)+`)
	present[Symbol] = re.MatchString(s)

	re = regexp.MustCompile(`(?m:\pP)+`)
	present[Punctuation] = re.MatchString(s)

	var accept = false
	var reject = true

	for k, v := range cc {
		//	fmt.Printf("Present[%c]=%v   cc[%c]=%v\n", k, present[k], k, cc[k])

		if v == true {
			accept = accept || present[k]
		} else {
			reject = reject && !present[k]
		}
	}
	return accept && reject
}
