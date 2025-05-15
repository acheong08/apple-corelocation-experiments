# Experiment with Apple's public WPS service
Reference: https://www.cs.umd.edu/~dml/papers/wifi-surveillance-sp24.pdf

How it works is explained in Apple's disclosure to congress: https://web.archive.org/web/20101208141602/https://markey.house.gov/docs/applemarkeybarton7-12-10.pdf

Feel free to poke around the code. Most relevant part is the [protobuf](./pb) and the stuff in [lib](./lib). Experimental CLIs found in [cmd](./cmd), the main ones being the demo api and `wloc`.

## WLOC

URL: https://gs-loc.apple.com/clls/wloc

I ran `mitmproxy` to find the URL used and later found [iSniff-GPS](https://github.com/hubert3/iSniff-GPS) on GitHub which was used as reference for the protobuf. The field names were uncovered by disassembling `CoreLocationProtobuf.framework` on MacOS. The relevant C code can be found [here](./CoreLocationProtobuf.c).

When requesting location services, MacOS/IOS sends a list of nearby BSSIDs to Apple, which then responds with GPS Long/Lat/Altitude of other nearby BSSIDs. The GPS location of the device is computed locally based on the signal strength of nearby BSSIDs.

Apple collects information from iPhones such as speed, activity type (walking/driving/etc), cell provider, and a whole bunch of other data which is used to build their database. This seems to be sent when a phone encounters a BSSID not in the existing database and excludes certain MAC address vendors known not to be stationary (e.g. IOS/MacOS hotspots).

### Cell towers

Using the same API as above, we can also request cell tower information. `MCC`, `MNC`, `CellId`, and `TacID` is sent off and used to find the surrounding cell towers. It seems to have more data than opencellid.org but missing UMTS and GSM (only LTE is available).

## Wifi Tile

URL: https://gspe85-ssl.ls.apple.com/wifi_request_tile

This is a new discovery while running MITM on an iPhone for an extended period of time. An endpoint from Apple takes a "X-tilekey" and returns all the BSSIDs and GPS locations of access points in the given region. Investigation is ongoing. 

Here's an example color coded cluster between keys 81644851 and 81644861 (Cardiff).

<img alt="Map of access points returned by the API" src="https://github.com/acheong08/apple-corelocation-experiments/assets/36258159/a7e3f898-b632-4d0d-a277-bb36281cf578" width=400>

It seems each key denotes a single network of access points. This information seems to be collected from within the networks. The API is labelled as `wifiQualityTileURL` in the code base.

The tile key is a [morton encoded number](https://en.wikipedia.org/wiki/Z-order_curve) with what appears to be their own coordinate system (not based on GPS). I have used linear regression to successfully convert between GPS (long/lat) to their coordinates ([code](./cmd/morton/main.go)).

Update: Linear regression is not the correct solution. It gets worse as you go north/south.

Update 2: This took more work than I'd like to admit. First, I went through all references to tile keys on GitHub. There are behavior discrepancies between [gojuno/go.morton](https://github.com/gojuno/go.morton) and what Apple uses (In the GeoServices private framework). I ended up finding the implementation by [heremaps](https://github.com/heremaps/here-data-sdk-typescript/blob/d9c39622b2306cb00803a493ea134e341716b96d/%40here/olp-sdk-core/lib/utils/TileKey.ts#L76) to match and translated that into Golang. Then, based on the output, noticed that the xyz looked similar to OpenStreetMap's tiles. I used the firefox debugger on [leafletjs](https://leafletjs.com/) to find the code used to generate the tiles from coordinates. Based on mentions of pixels and other keywords, I found [buckhx/tiles](https://github.com/buckhx/tiles) in this [8 year old Reddit post](https://www.reddit.com/r/golang/comments/4iki5d/map_tiling_library_for_go/). So to chain it together: tileKey → morton unpack → OSM tiles → pixel data → long/lat.

## China

Perhaps for data sovereignty reasons, Chinese data is isolated from the main API. However, we are still able to access it from outside.

URLS:
- "https://gs-loc.apple.com" -> "https://gs-loc-cn.apple.com"
- "https://gspe85-ssl.ls.apple.com" -> "https://gspe85-cn-ssl.ls.apple.com"

A quick `dig` and `whois` shows that these are hosted within China at a Unicom IP. `CNAME` records show that Akamai is used for DNS and Kingsoft Cloud as a CDN.

To swap out the APIs, simply add `-china` to any of the CLIs in `cmd`.

Credits to [JaneCCP](https://github.com/JaneCCP) for finding this info.

## Interactive demo

`go run ./cmd/demo-api` and head to http://127.0.0.1:1974. 

Click on any spot on the map and wait for a bit. It will plot nearby devices in a few seconds.

How it works: It first uses a spiral pattern to find the closest valid tile (limited to 20 to fail fast). Once it finds a starting point, it finds all the nearby access points using the WLOC API. It then takes the closest access point and tries again until there are no closer access points.

## Mass data collection

It is relatively simple to collect data via the tile API. The working code is [here](https://github.com/acheong08/apple-corelocation-experiments/tree/main/cmd/seedcrawl). You can collect around 9 million records by going through every tile (on land). Some work was done to detect if a coordinate is in water (to skip) or in China (to choose the right API). You can find some details [here](https://github.com/acheong08/apple-corelocation-experiments/tree/main/lib/shapefiles). 

Source for China's shapefile: [GaryBikini/ChinaAdminDivisonSHP](https://github.com/GaryBikini/ChinaAdminDivisonSHP/). This was [forked](https://github.com/acheong08/ChinaAdminDivisonSHP/) to remove special administration regions which are part of the international API.

Source for water polygons [here](https://osmdata.openstreetmap.de/data/water-polygons.html)

Once that data has been collected, we can begin using the wloc API to fetch and explore nearby sections. My code for that is incomplete but a friend was able to fetch around 1 billion records. His source code can be found here: https://codeberg.org/joelkoen/wtfps/.

Seed data: https://tmp.duti.dev/seeds.db.zst

<img src="https://github.com/user-attachments/assets/8da21d51-a506-4c32-94b7-b3ae853d65ab" alt="Grafana plot of collected seeds" height=400></img>

## Submitting data

To trigger data submission, turn on iPhone Analytics, Routing and Traffic, and Improve Maps in Location privacy settings. Then visit Maps, set a location, and start directions.

The endpoint is `https://gsp10-ssl.apple.com/hcy/pbcwloc`, labelled as `_CLPCellWifiCollectionRequest` in `CoreLocationProtobuf`.

It seems to only allow updating of existing records rather than collecting new records. More investigation is ongoing. You can find my attempt at uploading fake data [here](https://github.com/acheong08/apple-corelocation-experiments/tree/main/cmd/fakeloc).

## Spoofing your location

Using the info here, we can easily spoof your iPhone's location. You can find a server implementation [here](https://github.com/acheong08/apple-corelocation-experiments/tree/main/cmd/spoofed). After starting the server, run `mitmproxy` with [this](https://github.com/acheong08/apple-corelocation-experiments/blob/main/wloc.py) script which forwards wloc requests to your own server. Remember to change the IP address.

<img src="https://github.com/user-attachments/assets/a6ccdbc4-4317-4db8-8f3f-303e2e5a7c2f" alt="screenshot of maps with spoofed location" height="400"></img>

## Ichnaea

This is the format used by Mozilla/Google/etc for their location service. We can imitate that with [this](https://github.com/acheong08/apple-corelocation-experiments/blob/main/cmd/ichnaea/main.go). I have tested this with geoclue and it works just fine.

## To do
- If Apple uses device data to add new BSSIDs to their database, try to add fake data
- Ichnaea compatibility and add direct support in geoclue ([proposal](https://gitlab.freedesktop.org/geoclue/geoclue/-/issues/193))
