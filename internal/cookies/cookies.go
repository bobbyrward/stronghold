package cookies

import (
	"net/http"
)

type (
	Cookie = http.Cookie
)

func FindCookieByName(cookies []*http.Cookie, name string) (*http.Cookie, bool) {
	var foundCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == name {
			foundCookie = cookie
			break
		}
	}

	return foundCookie, foundCookie != nil
}
