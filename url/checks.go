package url

import (
	"errors"
	"fmt"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

type (
	allowChecks struct {
		isAllowedScheme AllowedSchemeCheck
		isAllowedUser   AllowedUserCheck
		isAllowedHost   AllowedHostCheck
	}

	AllowedSchemeCheck func(string) error
	AllowedUserCheck   func(*url.Userinfo) error
	AllowedHostCheck   func(string) error
)

func defaultAllowChecks() *allowChecks {
	return &allowChecks{
		isAllowedScheme: defaultIsAllowedScheme,
		isAllowedUser:   defaultIsAllowedUser,
		isAllowedHost:   defaultIsAllowedHost,
	}
}

const (
	httpScheme  = "http"
	httpsScheme = "https"
)

var errInvalidScheme = errors.New("invalid scheme")

func defaultIsAllowedScheme(scheme string) error {
	if scheme != httpScheme && scheme != httpsScheme {
		return errInvalidScheme
	}

	return nil
}

var errUserIsNotNil = errors.New("user is not nil")

func defaultIsAllowedUser(usr *url.Userinfo) error {
	if usr != nil {
		return errUserIsNotNil
	}

	return nil
}

func defaultIsAllowedHost(host string) error {
	if _, err := publicsuffix.EffectiveTLDPlusOne(host); err != nil {
		return fmt.Errorf("invalid host: %w", err)
	}

	return nil
}
