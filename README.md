# cfClearance

This small package can help you perform automated actions against web sites using Cloudflare's DDoS protection. You can read more about why we created it and how it works [here](https://gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/red-team-tech-notes/-/tree/master/cloudflare-notes).

You can import this package into your Golang applications so that you can use standard `net/http` requests to interact with a site behind Cloudflare.

## Usage

Let's say you need to scrape the text of `https://example.com`, but your requests fail with a `503` because of the Cloudflare protection mechanism. You can use this package like this:

```
import gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/cfClearance

// Full URL to where you encounter the Cloudflare warning
target := "https://example.com"

// You must use a consistent User-Agent between your tooling and this package
userAgent := "BottyMcBotFace"

// Receive back a http.Client object
cfClient, err := cfClearance.MakeCfClient(target, userAgent)
if err != nil {
  log.Fatal(err)
}

// Continue on with your own tools, ensuring a consistent user agent
req, _ := http.NewRequest("GET", target, nil)
req.Header.Set("User-Agent", userAgent)
resp, err := cfClient.Do(req)

...
```

And that's it! If you use it against a target that is not protected against Cloudflare, that's fine - it will simply return back a `http.Client` object without any special cookies added to it.