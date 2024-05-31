# Experiment with Apple's public WPS service
Reference: https://www.cs.umd.edu/~dml/papers/wifi-surveillance-sp24.pdf

How it works is explained in Apple's disclosure to congress: https://web.archive.org/web/20101208141602/https://markey.house.gov/docs/applemarkeybarton7-12-10.pdf

Feel free to poke around the code. Most relevant part is the [protobuf](./pb) and the stuff in [lib](./lib). Experimental CLIs found in [cmd](./cmd), the main ones being the demo api and `wloc`.

# WLOC

URL: https://gs-loc.apple.com/clls/wloc

I ran `mitmproxy` to find the URL used and later found [iSniff-GPS](https://github.com/hubert3/iSniff-GPS) on GitHub which was used as reference for the protobuf. The field names were uncovered by disassembling `CoreLocationProtobuf.framework` on MacOS. The relevant C code can be found [here](./CoreLocationProtobuf.c).

When requesting location services, MacOS/IOS sends a list of nearby BSSIDs to Apple, which then responds with GPS Long/Lat/Altitude of other nearby BSSIDs. The GPS location of the device is computed locally based on the signal strength of nearby BSSIDs.

Apple collects information from iPhones such as speed, activity type (walking/driving/etc), cell provider, and a whole bunch of other data which is used to build their database. This seems to be sent when a phone encounters a BSSID not in the existing database and excludes certain MAC address vendors known not to be stationary (e.g. IOS/MacOS hotspots).

## Wifi Tile

URL: https://gspe85-ssl.ls.apple.com/wifi_request_tile

This is a new discovery while running MITM on an iPhone for an extended period of time. An endpoint from Apple takes a "X-tilekey" and returns all the BSSIDs and GPS locations of access points in the given region. Investigation is ongoing. 

Here's an example color coded cluster between keys 81644851 and 81644861 (Cardiff).

<img alt="Map of access points returned by the API" src="https://github.com/acheong08/apple-corelocation-experiments/assets/36258159/a7e3f898-b632-4d0d-a277-bb36281cf578" width=400>

It seems each key denotes a single network of access points. This information seems to be collected from within the networks. The API is labelled as `wifiQualityTileURL` in the code base.

The tile key is a [morton encoded number](https://en.wikipedia.org/wiki/Z-order_curve) with what appears to be their own coordinate system (not based on GPS). I have used linear regression to successfully convert between GPS (long/lat) to their coordinates ([code](./cmd/morton/main.go)).

Update: Linear regression is not the correct solution. It gets worse as you go north/south.

Update 2: This took more work than I'd like to admit. First, I went through all references to tile keys on GitHub. There are behavior discrepancies between [gojuno/go.morton](https://github.com/gojuno/go.morton) and what Apple uses (In the GeoServices private framework). I ended up finding the implementation by [heremaps](https://github.com/heremaps/here-data-sdk-typescript/blob/d9c39622b2306cb00803a493ea134e341716b96d/%40here/olp-sdk-core/lib/utils/TileKey.ts#L76) to match and translated that into Golang. Then, based on the output, noticed that the xyz looked similar to OpenStreetMap's tiles. I used the firefox debugger on [leafletjs](https://leafletjs.com/) to find the code used to generate the tiles from coordinates. Based on mentions of pixels and other keywords, I found [buckhx/tiles](https://github.com/buckhx/tiles) in this [8 year old Reddit post](https://www.reddit.com/r/golang/comments/4iki5d/map_tiling_library_for_go/). So to chain it together: tileKey → morton unpack → OSM tiles → pixel data → long/lat.

## Interactive demo

`go run ./cmd/demo-api` and head to http://127.0.0.1:1974. 

Click on any spot on the map and wait for a bit. It will plot nearby devices in a few seconds.

How it works: It first uses a spiral pattern to find the closest valid tile (limited to 20 to fail fast). Once it finds a starting point, it finds all the nearby access points using the WLOC API. It then takes the closest access point and tries again until there are no closer access points.

### Impact

In the `umd.edu` paper, they were forced to brute force BSSIDs over the course of 2 years to slowly build up their network. This relies on chance and isolated blocks might not be easily found. With the discovery of the tile API, we are able to create a starting point anywhere in the world and begin exploring from there (given the parameters for which an AP is available over that API). I have tested multiple regions (e.g. Gaza, London, Los Angeles) and it seems to work in most populated cities. However, there appears to be less available networks in certain countries (e.g. Brazil), possibly due to privacy laws.


## To do
- If Apple uses device data to add new BSSIDs to their database, try to add fake data
- Chinese data is not available from outside. Attempt to proxy into China and map Chinese population centers based on density of access points
- Create UI to select point on map → use wifi tile to find nearest available chunk → use wloc to find all BSSIDs in area.
