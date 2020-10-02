package client

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"os"
)

func Create() http.Client {
	// If a proxy is defined, skip TLS verification.
	// We do this as it seems likely you are testing via ZAP/Burp/etc
	var tr http.Transport
	if os.Getenv("HTTP_PROXY") != "" || os.Getenv("HTTPS_PROXY") != "" {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		tr.Proxy = http.ProxyFromEnvironment
	}

	// Initialize an empy cookie jar. It will be populated later with Cloudflare cookie
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Transport: &tr,
	}
	client.Jar = cookieJar

	return *client

}
