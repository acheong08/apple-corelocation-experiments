Documenting the various endpoints, which binary calls them, and what the purpose/content is

| Domain | Path | Source | Request | Response |
|---|---|---|---|---|
| gspe1-ssl.ls.apple.com | /pep/gcc | /System/Library/PrivateFrameworks/GeoServices.framework/Versions/A/XPCServices/com.apple.geod.xpc/Contents/MacOS/com.apple.geod | Empty | Country code |
| gspe21-ssl.ls.apple.com, gspe19-ssl.ls.apple.com | * |  | Access key (tied to account) | Map icons, style sheets, etc |
| gspe19-2-ssl.ls.apple.com | /poi_update |  | Access key, tile location, and other map settings | Point of interests, not protobuf |
| gsp-ssl.ls.apple.com | /dispatcher.arpc |  | Uniquely identifiable UUID, Location data, SIM carrier, Locale, and opaque binary blobs | Maps data |
|  | /directions.arpc |  |  | Directions |
| gspe35-ssl.ls.apple.com, configuration.ls.apple.com | * |  | Hardware and software versions | Configuration depending on environment (prod, staging, beta) |
| gsp10-ssl.apple.com | /au, /pds/pd | locationd | App IDs, exact location they were opened, and various other metadata!!! | Acknowledgement |
|  | /hcy/pbcwloc |  | Surrounding BSSIDs, cell provider, location, movement and activity |  |
| gsp64-ssl.ls.apple.com | /hvr/v3/use, /hvr/v2/rtloc | geoanalyticsd, nsurlsessiond | Opaque binary data. Includes open apps, home location (NOT CURRENT). IP address and ports, possible open connections. Some requests are encrypted! (~7.5 entropy with metadata removed) |  |
