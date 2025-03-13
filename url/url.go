package url

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/idna"
)

func NormalizeWithOptions(urlToNormalize string, options ...AllowOption) (*url.URL, error) {
	checks := optionsToChecks(options)

	normalizedURL, err := normalizeURLString(urlToNormalize, *checks)
	if err != nil {
		return nil, fmt.Errorf("normalize string URL: %w", err)
	}

	return normalizedURL, nil
}

func normalizeURLString(urlToNormalize string, checks allowChecks) (*url.URL, error) {
	// can't use value from here because it removes fragment from path
	if _, err := url.ParseRequestURI(urlToNormalize); err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	parsed, err := url.Parse(urlToNormalize)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)

	if allowErr := checks.isAllowedScheme(parsed.Scheme); allowErr != nil {
		return nil, fmt.Errorf("schema (%s) not allowed: %w", parsed.Scheme, allowErr)
	}

	if allowErr := checks.isAllowedUser(parsed.User); allowErr != nil {
		return nil, fmt.Errorf("user (%v) not allowed: %w", parsed.User, allowErr)
	}

	asciiHost, err := extractASCIIHost(parsed)
	if err != nil {
		return nil, fmt.Errorf("invalid host: %w", err)
	}

	if allowErr := checks.isAllowedHost(asciiHost); allowErr != nil {
		return nil, fmt.Errorf("host (%s) not allowed: %w", asciiHost, allowErr)
	}

	parsed.Host = removeDefaultPort(asciiHost)
	parsed.Path = normalizePath(parsed.Path)
	parsed.RawQuery = normalizeQueryParams(parsed.Query())

	return parsed, nil
}

func extractASCIIHost(rawURL *url.URL) (string, error) {
	loweredHostname := strings.ToLower(rawURL.Hostname())

	asciiHost, err := idna.Punycode.ToASCII(loweredHostname)
	if err != nil {
		return "", fmt.Errorf("invalid host: %w", err)
	}

	if port := rawURL.Port(); port != "" { // return port if it was there
		asciiHost += ":" + port
	}

	return asciiHost, nil
}

const (
	defaultHTTPPort  = "80"
	defaultHTTPSPort = "443"
)

func removeDefaultPort(host string) string {
	if parts := strings.Split(host, ":"); len(parts) > 1 &&
		(parts[1] == defaultHTTPSPort || parts[1] == defaultHTTPPort) {
		return parts[0]
	}

	return host
}

var filePathRegex = regexp.MustCompile(`.+\..+`)

func normalizePath(urlPath string) string {
	cleanPath := path.Clean(urlPath)
	if cleanPath == "/" || cleanPath == "." {
		return ""
	}

	if filePathRegex.MatchString(cleanPath) {
		cleanPath = removeDirIndex(cleanPath)
	}

	return cleanPath
}

func normalizeQueryParams(params url.Values) string {
	sorted := make(url.Values, len(params))

	for key, values := range params { // sort values
		valuesCopy := make([]string, len(values))
		copy(valuesCopy, values)
		sort.Strings(valuesCopy)
		sorted[key] = valuesCopy
	}

	return sorted.Encode() // encode sort keys
}

//nolint:gochecknoglobals // static value
var defaultDirFiles = []string{
	"index.html", "index.shtml", "index.htm", "Index.html", "Index.htm", "Index.shtml",
	"index.php", "index.php3", "index.phtml", "index.cgi", "index.pl", "index.jsp", "index.wml", "index.pcgi",
	"default.asp",
}

func removeDirIndex(urlPath string) string {
	for _, fileName := range defaultDirFiles {
		suffix := "/" + fileName

		if strings.HasSuffix(urlPath, suffix) {
			urlPath = strings.TrimSuffix(urlPath, suffix)

			break
		}
	}

	return urlPath
}
