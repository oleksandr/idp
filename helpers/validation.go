package helpers

import "regexp"

const (
	// BasicEmailRegExp is a regular expression that matches 99% of the email
	// addresses in use today.
	BasicEmailRegExp = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`

	// RFCEmailRegExp - practical implementation of RFC 2822 omitting the
	// syntax using double quotes (") and square brackets ([]), since not
	// all applications support the syntax using it.
	RFCEmailRegExp = `(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+` +
		`(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*` +
		`@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+` +
		`[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`
)

// ValidBasicEmail validates email using BasicEmailRegExp regexp
func ValidBasicEmail(email string) bool {
	exp, _ := regexp.Compile(BasicEmailRegExp)
	if exp.MatchString(email) {
		return true
	}
	return false
}

// ValidRFC2822Email validates email using RFCEmailRegExp regexp
func ValidRFC2822Email(email string) bool {
	exp, _ := regexp.Compile(RFCEmailRegExp)
	if exp.MatchString(email) {
		return true
	}
	return false
}
