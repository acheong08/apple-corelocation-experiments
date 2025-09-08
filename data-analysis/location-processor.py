import json
import struct
import sys
import argparse


def int_to_double(i):
    return struct.unpack("<d", struct.pack("<Q", i))[0]


def main():
    parser = argparse.ArgumentParser(description='Process location data from parsed JSON')
    parser.add_argument('input_file', help='Input JSON file containing parsed location data')
    parser.add_argument('output_file', help='Output JSON file for processed locations')
    
    args = parser.parse_args()
    
    with open(args.input_file) as f:
        data = json.load(f)

    result = []
    for item in data:
        processed_item = {}
        for key, value in item.items():
            if key == "lat" or key == "lon":
                processed_item[key] = int_to_double(value)
            else:
                processed_item[key] = value
        result.append(processed_item)

    with open(args.output_file, 'w') as f:
        json.dump(result, f, indent=2)
    
    print(f"Processed {len(result)} locations and saved to {args.output_file}")


if __name__ == "__main__":
    main()
