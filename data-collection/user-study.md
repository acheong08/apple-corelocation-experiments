# User Study on Privacy of Apple Location Services

## What we want to find out

**How well does deanonymization scale?**
Say there are 10 people on the same IP address. How unique are their set of identifiers (iOS version, apps installed)? What happens when we scale that up to 100?

**How privacy-aware are the general public?**

- Which privacy settings have they got enabled?
- Are descriptions of these settings informative enough?
- Are people surprised by what data is collected?

## What data do we need to collect

- `gsp10-ssl.apple.com` Routing & Traffic, Improve Location Services
- `gsp-ssl.ls.apple.com` Apple Pay Merchant Identification
- `gs-loc.apple.com` Location services

These domains receive the bulk of the analytics data. Other domains are more dual-use and would capture unnecessary authenticated data. We can collect counts and timings of requests (without decrypting) to authenticated Apple services which can then be used as part of deanonymization.

## Privacy considerations and measures

## Ethics: How do we ensure users are informed about what data they're handing over?

**Visualizations**

Using the locally collected data, we can map out all the routes they've taken, apps they've had open, and NFC payments made.
