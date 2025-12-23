# Universal Magic Link (Web + iOS)

This project supports a single HTTPS magic link that signs in users on desktop browsers and opens the iOS app on phones via Apple Universal Links.

## Link format

- Universal link: `https://<domain>/auth/magic?token=<TOKEN>`
- Backend generates the link via `MakeMagicLink()` and verifies it at `/api/v1/auth/verifyMagicLink`.
- Frontend handles web flow on the route `/auth/magic`.

## One-Time Code (OTP)

- Alongside the magic link, the backend issues a 6-digit code.
- Email contains both the clickable magic link and the code.
- iOS app can verify via `POST /api/v1/auth/verifyMagicCode` with JSON `{ "code": "123456" }`.
- The code expires at the same time as the magic link (15 minutes) and is consumed upon successful verification.

## iOS (Apple Universal Links)

To open the iOS app when installed, configure Universal Links:

1. Host the AASA file on your domain (no extension, JSON only):
   - `/.well-known/apple-app-site-association`
   - `/apple-app-site-association`
   Both are added in this repo under Frontend/public for dev builds.

2. Fill in your real identifiers in the AASA content:
   ```json
   {
     "applinks": {
       "apps": [],
       "details": [
         {
           "appID": "<TEAMID>.com.clubs.ios",
           "paths": [
             "/auth/magic*"
           ]
         }
       ]
     }
   }
   ```
   - TEAMID: Your Apple Developer Team ID.
   - BUNDLEID: Your app bundle identifier (e.g., com.example.clubs).

3. iOS app configuration:
  - Associated Domains entitlements are configured at `ios/ios/ios.entitlements`.
  - Replace `YOUR-DOMAIN-HERE` with your actual domain (e.g., `applinks:clubs.example.com`).
  - The app handles incoming universal links in `ios/ios/iosApp.swift` by parsing `token` from `/auth/magic` and calling the backend `verifyMagicLink` endpoint to obtain session tokens, then routing in-app.

4. Fallback behavior:
   - If the app is not installed, iOS opens the HTTPS link in Safari, where the web flow at `/auth/magic` signs the user in.

## Notes

- Keep the link path consistent (`/auth/magic`) across web and iOS.
- For Android, consider App Links (Digital Asset Links) using a similar approach.
- Do not serve the AASA file with a `Content-Type: application/json`; Apple accepts `application/pkcs7-mime` or no type, but most hosts work with JSON. Ensure it’s accessible over HTTPS.

## Debugging

- Use `deviceconsole` or Xcode logs to confirm Universal Link association.
- Visit `https://<domain>/.well-known/apple-app-site-association` to verify content.
- On iOS, uninstall/reinstall the app if changes to AASA or entitlements don’t take effect.
