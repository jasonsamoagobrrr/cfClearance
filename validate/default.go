package validate

import "net/url"

func Url(target string) bool {
	u, err := url.Parse(target)

	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
