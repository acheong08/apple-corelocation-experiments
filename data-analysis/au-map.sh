#!/usr/bin/env bash

wl-paste | xxd -r -p > req.bin

printbin req.bin -proto -o /tmp/au.json

cat /tmp/au.json | jq '[."2".[] | {app: ."1", lat: ."4"."1", lon: ."4"."2"}]' > /tmp/au1.json

python3 location-processor.py /tmp/au1.json ./au.json

python3 visualize_locations.py au.json data/apps/map.html
