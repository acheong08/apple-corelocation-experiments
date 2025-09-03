# Data collection and possible attacks

## Threats against user privacy

- **Apple**: Is able to collect location data tied to IP address, IOS version, and phone model. However, due to the high volume, deanonymization may be difficult.
- **ISP / passive MITM**: Able to view frequency of requests via SNI sniffing and DNS requests. If the regularity of requests is predictable, an ISP may be able to identify the presence of a user within the network.

For each given request to wloc, the request includes your IOS version, iPhone model, locale, and the binary making the request as obvious points of fingerprinting. There are also exact version numbers of `locationd`, `CFNetwork`, and `Darwin` as part of the user agent string. Overall, there is ~111 bytes of unnecessary fingerprintable data.
The IP address can also be correlated with other requests made to Apple which contain account information (e.g. iMessage, Apple Wallet, etc.). This then becomes a problem of how many non-unique Apple devices share an IP space such that an exact location is not known.

For an ISP, it is already possible to track a user's location based on cell towers. There is not an increased attack surface here.
