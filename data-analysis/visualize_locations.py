#!/usr/bin/env python3
import folium
import json
import statistics
import argparse


def load_locations(filename):
    with open(filename, "r") as f:
        return json.load(f)


def create_map(locations):
    if not locations:
        print("No locations found!")
        return None

    # Calculate center point for the map
    lats = [loc["lat"] for loc in locations]
    lons = [loc["lon"] for loc in locations]

    center_lat = statistics.mean(lats)
    center_lon = statistics.mean(lons)

    # Create the map centered on the data
    m = folium.Map(
        location=[center_lat, center_lon], zoom_start=15, tiles="OpenStreetMap"
    )

    # Add points to the map
    for i, loc in enumerate(locations):
        point_name = loc.get("app", f"Point {i + 1}")
        folium.CircleMarker(
            location=[loc["lat"], loc["lon"]],
            radius=3,
            popup=f"{point_name}<br>Lat: {loc['lat']:.6f}<br>Lon: {loc['lon']:.6f}",
            color="red",
            fill=True,
            fillColor="red",
            fillOpacity=0.6,
        ).add_to(m)

    return m


def main():
    parser = argparse.ArgumentParser(description="Visualize location data on a map")
    parser.add_argument("input_file", help="Input JSON file containing location data")
    parser.add_argument("output_file", help="Output HTML file for the map")

    args = parser.parse_args()

    locations = load_locations(args.input_file)
    print(f"Loaded {len(locations)} locations")

    map_obj = create_map(locations)
    if map_obj:
        map_obj.save(args.output_file)
        print(f"Map saved to {args.output_file}")
        print(f"Open {args.output_file} in your browser to view the map")


if __name__ == "__main__":
    main()

