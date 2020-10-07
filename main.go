package cfClearance

import (
	"errors"
	"gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/cfClearance/browser"
	"gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/cfClearance/cfclient"
	"gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/cfClearance/validate"
	"log"
	"net/http"
)

// ConfigureCfClient should be called directly, using a URL string as the target
// and a second string for the User-Agent to pass. User-Agent must match what
// you use in your tooling for subsequent requests, per Cloudflare.
// Pass in your own http.Client that will receive CloudFlares'
// (cf_clearance) if the site is protected.
func ConfigureClient(client *http.Client, target string, agent string) error {
	// Initialize the client with the things we need to bypass cloudflare
	cfclient.Initialize(client)

	// Validate the target URL
	if validate.Url(target) == false {
		return errors.New("could not parse the target URL")
	}

	// Check if target is even protected by Cloudflare. If not, just return the
	// client as-is.
	if validate.CloudFlareIsPresent(target, client) == false {
		log.Println("[*] Target not protected by Cloudflare.")
		return nil
	}

	log.Println("[!] Target is protected by Cloudflare, bypassing...")

	return browser.GetCloudFlareClearanceCookie(client, agent, target)

}
