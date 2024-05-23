# Experiment with Apple's public WPS service
Reference: https://www.cs.umd.edu/~dml/papers/wifi-surveillance-sp24.pdf

How it works is explained in Apple's disclosure to congress: https://web.archive.org/web/20101208141602/https://markey.house.gov/docs/applemarkeybarton7-12-10.pdf

I ran `mitmproxy` to find the URL used and later found [iSniff-GPS](https://github.com/hubert3/iSniff-GPS) on GitHub which was used as reference for the protobuf.

When requesting location services, MacOS/IOS sends a list of nearby BSSIDs to Apple, which then responds with GPS Long/Lat/Altitude of other nearby BSSIDs. The GPS location of the device is computed locally based on the signal strength of nearby BSSIDs.

Work in progress:
- Figuring out the values of the protobuf by reverse engineering the `CoreLocationProtobuf` binary (extracted from `/System/Volumes/Preboot/Cryptexes/OS/System/Library/dyld/`). You can see the protobuf parsing code [here](./CoreLocationProtobuf.c)

To do:
- See whether GPS location of device is uploaded at any point
- If Apple uses device data to add new BSSIDs to their database, try to add fake data
- Chinese data is not available from outside. Attempt to proxy into China and map Chinese population centers based on density of access points
