#!/usr/bin/env bash

# Check if input file is provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <req.bin> <output_directory>"
    exit 1
fi

INPUT_FILE="$1"
OUTPUT_DIR="$2"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Step 1: Run printbin on the input file to generate JSON
printbin "$INPUT_FILE" -proto -o /tmp/pds.json

# Step 2: Use jq to extract lat/lon data
jq '[."2"[] | {lat: .["1"], lon: .["2"]}]' /tmp/pds.json > /tmp/pds_loc.json

# Step 3: Run location-processor to generate final locations.json
python3 location-processor.py /tmp/pds_loc.json "$OUTPUT_DIR/locations.json"

# Step 4: Run visualize_locations to create the map
python3 visualize_locations.py "$OUTPUT_DIR/locations.json" "$OUTPUT_DIR/map.html"

echo "Pipeline complete! Files generated:"
echo "- $OUTPUT_DIR/locations.json"
echo "- $OUTPUT_DIR/map.html"