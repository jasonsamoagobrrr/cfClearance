package cfClearance

import (
	"errors"
	"gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/cfClearance/cfclient"
	"gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/cfClearance/browser"
	"gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/cfClearance/validate"
	"log"
	"net/http"
)

// MakeCfClient should be called directly, using a URL string as the target
// and a second string for the User-Agent to pass. User-Agent must match what
// you use in your tooling for subsequent requests, per Cloudflare.
// It will return an http.Client that has the required Cloudflare cookie
//(cf_clearance) if the site is protected.
func MakeCfClient(target string, agent string) (*http.Client, error) {
	// Start with a fresh http.Client pointer we'll later return
	client := cfclient.Create()

	// Validate the target URL
	if validate.Url(target) == false {
		return client, errors.New("Could not parse the target URL")
	}

	// Check if target is even protected by Cloudflare. If not, just return the
	// client as-is.
	if validate.CloudflareExists(target, client) == false {
		log.Println("[*] Target not protected by Cloudflare.")
		return client, nil
	}

	log.Println("[!] Target is protected by Cloudflare, bypassing...")

	return browser.GetCloudFlareClearanceCookie(client, agent, target)

}

