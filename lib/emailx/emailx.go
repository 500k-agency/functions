package emailx

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/500k-agency/function/data"
)

var (
	//ErrInvalidFormat returns when email's format is invalid
	ErrInvalidFormat = errors.New("invalid format")
	//ErrUnresolvableHost returns when validator couldn't resolve email's host
	ErrUnresolvableHost = errors.New("unresolvable host")

	userRegexp = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp = regexp.MustCompile(`^[^\s]+\.[^\s]+$`)
	// As per RFC 5332 secion 3.2.3: https://tools.ietf.org/html/rfc5322#section-3.2.3
	// Dots are not allowed in the beginning, end or in occurances of more than 1 in the email address
	userDotRegexp = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")

	userPartsRegexp = regexp.MustCompile("^(.)(.*)(.@.*)$")
)

type Email struct {
	value string // full email value
	User  string
	Host  string
}

// Validate checks format of a given email and resolves its host name.
func Validate(email string) error {
	_, err := New(email)
	return err
}

func New(email string) (*Email, error) {
	email = Normalize(email)

	if len(email) < 6 || len(email) > 254 {
		return nil, ErrInvalidFormat
	}

	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return nil, ErrInvalidFormat
	}

	user := email[:at]
	host := email[at+1:]

	if len(user) > 64 {
		return nil, ErrInvalidFormat
	}
	if userDotRegexp.MatchString(user) || !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return nil, ErrInvalidFormat
	}

	e := &Email{value: email, User: user, Host: host}
	if err := e.Validate(); err != nil {
		return e, err
	}
	return e, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Validate() error {
	switch e.Host {
	case "localhost", "example.com":
		return nil
	}
	if _, err := net.LookupMX(e.Host); err != nil {
		if _, err := net.LookupIP(e.Host); err != nil {
			// Only fail if both MX and A records are missing - any of the
			// two is enough for an email to be deliverable
			return ErrUnresolvableHost
		}
	}
	return nil
}

// ValidateFast checks format of a given email.
func ValidateFast(email string) error {
	if len(email) < 6 || len(email) > 254 {
		return ErrInvalidFormat
	}

	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return ErrInvalidFormat
	}

	user := email[:at]
	host := email[at+1:]

	if len(user) > 64 {
		return ErrInvalidFormat
	}
	if userDotRegexp.MatchString(user) || !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return ErrInvalidFormat
	}

	return nil
}

func Mask(email string) string {
	ex, err := New(email)
	if err != nil {
		return ""
	}
	subs := userPartsRegexp.FindStringSubmatch(email)
	if len(ex.User) < 5 {
		return fmt.Sprintf("%s%s%s", subs[1], strings.Repeat("*", len(ex.User)-1), subs[3])
	}
	return fmt.Sprintf("%s%s%s", subs[1], strings.Repeat("*", data.Min(8, len(subs[2]))), subs[3])
}

// Normalize normalizes email address.
func Normalize(email string) string {
	// Trim whitespaces.
	email = strings.TrimSpace(email)

	// Trim extra dot in hostname.
	email = strings.TrimRight(email, ".")

	// Lowercase.
	email = strings.ToLower(email)

	return email
}
