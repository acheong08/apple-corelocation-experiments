# Technical Report on Apple's Location Privacy Claims and Observed Behavior

## Methodology, Notes, and Data

### Recording traffic

`mitmweb --web-host {private_ip} --mode wireguard --set allow_hosts="mitm\.it|.*\.apple\.com" --listen-host {public_ip} -w apple.flow --set view_filter="~http mitm\.it|.*\.apple\.com"`

Wireguard is preferred as it allows all traffic to be routed in mobile devices, whereas HTTPS and SOCKS proxies occasionally do not get applied to systems services. You can issue more limited certificates via `./data-collection/mitm/gencerts.sh` to ensure the mitm server doesn't even have the capacity to decrypt your unrelated traffic.

### Polling for changes

Code is in `./data-collection/*.go`. Essentially, 10 evenly spread samples with 5 adjacent tile keys are chosen from a seed database collected by `./cmd/seedcrawl/`. Then, we poll the tile API every 5 minutes and record changes in location.

### Reverse engineering

**Finding the files**

The primary bundle identifiers observed in traffic related to location are `locationd`, `geod`, and `geoanalyticsd`.

| Binary        | Path                                                                                                                              |
| ------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| geoanalyticsd | `/System/Library/PrivateFrameworks/GeoAnalytics.framework/geoanalyticsd`                                                          |
| locationd     | `/usr/libexec/locationd`                                                                                                          |
| geod          | `/System/Library/PrivateFrameworks/GeoServices.framework/Versions/A/XPCServices/com.apple.geod.xpc/Contents/MacOS/com.apple.geod` |

The filesystem for MacOS and IOS are roughly the same.

To obtain the IOS root filesystem, [ipsw](https://github.com/blacktop/ipsw) can be used to download and extract update files.

```sh
ipsw download ipsw --version 18.6.2 --device iPhone15,2
ipsw extract --files --pattern ".*" iPhone15,2_18.6.2_22G100_Restore.ipsw
```

You will find that many of the private frameworks are missing. As of macOS Big Sur, instead of shipping the system libraries with macOS, Apple ships a generated cache of all built in dynamic libraries and excludes the originals.

To extract them, [dyld-shared-cache-extractor](https://github.com/keith/dyld-shared-cache-extractor) can be used for both IOS and MacOS. Decompilations were done with Ghidra to look for relevant symbols and logic.

**Debugging**

The size of these binaries are immense, and it is difficult to find exactly when and where APIs are called. While you can use LLDB on MacOS, the more interesting target is IOS due to the amount of information it sends.

The solution to this is disabling SIP (System Integrity Protection) on MacOS and running an IOS simulator within it. Use `ps aux | grep <process name>` and you should see a binary with a path within the XCode simulator. Simply `sudo lldb -p <pid>` to attach.
To break at HTTP requests, use `br set -n "-[NSURLSessionTask resume]"`. You can then see the URL `po [$x0 originalRequest]`.

The next step is to figure out the address in Ghidra. To do this, you must calculate some offsets.
To find the base load address, run `image list <binary name (e.g. geoanalyticsd)>`. Then, take the top frame address from the backtrace which matches your binary. `p/x <frame_runtime_address> - <base_load_address>` returns a file offset. Ghidra expects a virtual address, which can be calculated by adding `0x100000000`.

As an example, `p/x 0x1048ddda0 - 0x1048d4000` = `0x9da0`. `p/x 0x100000000 + 0x9da0` = `0x100009da0`.

## Recorded traffic patterns, data leakage, and privacy concerns

| Domain                                              | Path                       | Source                       | Request                                                                                                                                                                                 | Response                                                     |
| --------------------------------------------------- | -------------------------- | ---------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------ |
| gspe1-ssl.ls.apple.com                              | /pep/gcc                   | com.apple.geod               | Empty                                                                                                                                                                                   | Country code                                                 |
| gspe21-ssl.ls.apple.com, gspe19-ssl.ls.apple.com    | \*                         |                              | Access key (tied to account)                                                                                                                                                            | Map icons, style sheets, etc                                 |
| gspe19-2-ssl.ls.apple.com                           | /poi_update                |                              | Access key, tile location, and other map settings                                                                                                                                       | Point of interests, not protobuf                             |
| gsp-ssl.ls.apple.com                                | /dispatcher.arpc           |                              | Random UUID, Location data, SIM carrier, Locale, and opaque binary blobs, merchant and transaction information                                                                          | Maps data                                                    |
|                                                     | /directions.arpc           |                              |                                                                                                                                                                                         | Directions                                                   |
|                                                     | /ab.arpc                   |                              | Consistent UID                                                                                                                                                                          | AB testing configurations                                    |
| gspe35-ssl.ls.apple.com, configuration.ls.apple.com | \*                         |                              | Hardware and software versions                                                                                                                                                          | Configuration depending on environment (prod, staging, beta) |
| gsp10-ssl.apple.com                                 | /au, /pds/pd               | locationd                    | App IDs, exact location they were opened, and various other metadata!!!                                                                                                                 | Acknowledgement                                              |
|                                                     | /hcy/pbcwloc               |                              | Surrounding BSSIDs, cell provider, location, movement and activity                                                                                                                      |                                                              |
| gsp64-ssl.ls.apple.com                              | /hvr/v3/use, /hvr/v2/rtloc | geoanalyticsd, nsurlsessiond | Opaque binary data. Includes open apps, home location (NOT CURRENT). IP address and ports, possible open connections. Some requests are encrypted! (~7.5 entropy with metadata removed) |                                                              |

There are a couple main points of concern:

- The `/au` endpoint receives batched data of exactly what apps you had open, where, and when. This allows Apple to know your exact movement and behavioral patterns over the course of a day and is extremely identifiable.
- `/dispatcher.arpc` occasionally receives information on transactions you made. When all privacy settings are enabled, data is still stored until you disable them, where they are then sent in bulk to Apple. This means that you cannot have brief periods of privacy as they will be sent weeks later if privacy settings are disabled at any time, even accidentally (e.g. updates).
- `/hvr/v3/use` receives a wide range of opaque binary data. I have observed various personally identifiable information including my home location which is set in the Contacts app, various IP addresses and ports I'm connected to, and open apps.

TODO: Read Apple privacy policy and see if that's documented anywhere

## Location database update patterns
