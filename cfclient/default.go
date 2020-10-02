package cfclient

import (
	"context"
	"crypto/tls"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

func Create() *http.Client {
	// If a proxy is defined, skip TLS verification.
	// We do this as it seems likely you are testing via ZAP/Burp/etc
	var tr http.Transport
	if os.Getenv("HTTP_PROXY") != "" || os.Getenv("HTTPS_PROXY") != "" {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		tr.Proxy = http.ProxyFromEnvironment
	}

	// Initialize an empty cookie jar. It will be populated later with Cloudflare cookie
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Transport: &tr,
	}
	client.Jar = cookieJar

	return client

}

func BakeCookies(target string, cfToken string) (*url.URL, []*http.Cookie) {
	u, _ := url.Parse(target)
	d := "." + u.Host
	var cookies []*http.Cookie
	cfCookie := &http.Cookie{
		Name:   "cf_clearance",
		Value:  cfToken,
		Path:   "/",
		Domain: d,
	}
	cookies = append(cookies, cfCookie)
	cookieURL, _ := url.Parse(target)

	return cookieURL, cookies
}

func ExtractCookie(c chan string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return err
		}
		for _, cookie := range cookies {
			if strings.ToLower(cookie.Name) == "cf_clearance" {
				// if we find a proper cookie, put the value on the receiving channel
				c <- cookie.Value
			}
		}
		return nil
	})
}
