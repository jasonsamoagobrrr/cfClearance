
```
import github.com/imayberoot/cfClearance

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

```

And that's it! If you use it against a target that is not protected against Cloudflare, that's fine - it will simply return back a `http.Client` object without any special cookies added to it.
