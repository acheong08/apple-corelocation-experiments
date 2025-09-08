import json
import struct


def int_to_double(i):
    return struct.unpack("<d", struct.pack("<Q", i))[0]


with open("./parsed_locations.json") as f:
    data = json.load(f)

result = [
    {"lat": int_to_double(item["lat"]), "lon": int_to_double(item["lon"])}
    for item in data
]

print(json.dumps(result, indent=2))
