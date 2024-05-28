# Experiment with Apple's public WPS service
Reference: https://www.cs.umd.edu/~dml/papers/wifi-surveillance-sp24.pdf

How it works is explained in Apple's disclosure to congress: https://web.archive.org/web/20101208141602/https://markey.house.gov/docs/applemarkeybarton7-12-10.pdf

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

The tile key is a [morton encoded number](https://en.wikipedia.org/wiki/Z-order_curve) with what appears to be their own coordinate system (not based on GPS). I have used linear regression to successfully convert between GPS (long/lat) to their coordinates ([code](./cmd/regress/main.go)).

## Usage

**Building**

`go build ./cmd/wloc`

**Running**
```
wloc get - Gets and displays adjacent BSSID locations given an existing BSSID
Flags:

  -bssid value
    	One or more known bssid strings
  -less
    	Only return requested BSSID location

wloc tile - Returns a list of BSSIDs and their associated GPS locations
Flags:

  -key int
    	The tile key used to determine region (default 81644853)
```

Multiple bssids:
`./wloc get -bssid <bssid1> -bssid <bssid2> -bssid <bssid3>...`

Output:
```
BSSID: xx:xx:xx:xx:xx (Vendor) found at Lat: 0.000 Long: 0.000
```

## To do
- If Apple uses device data to add new BSSIDs to their database, try to add fake data
- Chinese data is not available from outside. Attempt to proxy into China and map Chinese population centers based on density of access points

## Changelog
- Unknown fields have been found for protobuf: horizontal_accuracy, altitude, vertical_accuracy (among other less useful ones)
- Added information on collected info
