package cfClearance

import (
	"context"
	"errors"
	"gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/cfClearance/cfclient"
	"gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/cfClearance/validate"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
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

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// Ignore certificate errors (for use with proxy testing)
		chromedp.Flag("ignore-certificate-errors", "1"),
		// User-Agent MUST match what your tooling uses
		chromedp.UserAgent(agent),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create the chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// Challenges should be solved in ~5 seconds but can be slower. Timeout at 30.
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Listen for the Cloudflare cookie
	cookieReceiverChan := make(chan string, 1)
	defer close(cookieReceiverChan)

	// Fetch the login page and wait until CF challenge is solved.
	err := chromedp.Run(ctx,
		chromedp.Navigate(target),
		chromedp.WaitNotPresent(`Checking your browser`, chromedp.BySearch),
		extractCookie(cookieReceiverChan),
	)
	if err != nil {
		if err == context.DeadlineExceeded {
			return client, errors.New("Context deadline exceeded trying to grab cookie using chromedp")
		}
		return client, err
	}

	// block the program until the cloud flare cookie is received, or .WaitVisible times out looking for login-pane
	cfToken := <-cookieReceiverChan

	log.Printf("[*] Grabbed Cloudflare token: %s", cfToken)

	// Finally, build up the cookie jar with the required token
	cookieURL, cookies := cfclient.BakeCookies(target, cfToken)
	client.Jar.SetCookies(cookieURL, cookies)

	return client, nil
}

func extractCookie(c chan string) chromedp.Action {
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
