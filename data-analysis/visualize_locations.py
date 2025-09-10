#!/usr/bin/env python3
import folium
import json
import statistics
import argparse
import requests
import os
import colorsys


def load_locations(filename):
    with open(filename, "r") as f:
        return json.load(f)


def get_app_icon_url(bundle_id, cache_dir="icon_cache"):
    """
    Fetch app icon URL from Apple's iTunes Search API.
    Returns the 512x512 icon URL or None if not found.
    """
    if not bundle_id:
        return None

    # Create cache directory if it doesn't exist
    os.makedirs(cache_dir, exist_ok=True)
    cache_file = os.path.join(cache_dir, f"{bundle_id.replace('/', '_')}.json")

    # Check cache first
    if os.path.exists(cache_file):
        try:
            with open(cache_file, "r") as f:
                cached_data = json.load(f)
                return cached_data.get("icon_url")
        except (json.JSONDecodeError, IOError):
            pass

    # Fetch from iTunes API
    try:
        url = f"https://itunes.apple.com/lookup?bundleId={bundle_id}"
        response = requests.get(url, timeout=10)
        response.raise_for_status()

        data = response.json()
        if data.get("resultCount", 0) > 0:
            result = data["results"][0]
            icon_url = result.get("artworkUrl512")

            # Cache the result
            cache_data = {"bundle_id": bundle_id, "icon_url": icon_url}
            with open(cache_file, "w") as f:
                json.dump(cache_data, f)

            return icon_url
    except (requests.RequestException, json.JSONDecodeError, KeyError):
        pass

    return None


def timestamp_to_color(ts, min_ts, max_ts):
    """Convert timestamp to a color from blue (oldest) to red (newest)"""
    if max_ts == min_ts:
        return "#FF0000"  # Default to red if all timestamps are the same
    
    # Normalize timestamp to 0-1 range
    normalized = (ts - min_ts) / (max_ts - min_ts)
    
    # Convert to HSV color space (hue from 240° blue to 0° red)
    hue = (1 - normalized) * 240 / 360  # 240° = blue, 0° = red
    saturation = 1.0
    value = 1.0
    
    # Convert HSV to RGB
    r, g, b = colorsys.hsv_to_rgb(hue, saturation, value)
    
    # Convert to hex color
    return f"#{int(r*255):02x}{int(g*255):02x}{int(b*255):02x}"


def create_map(locations):
    if not locations:
        print("No locations found!")
        return None

    # Calculate center point for the map
    lats = [loc["lat"] for loc in locations]
    lons = [loc["lon"] for loc in locations]

    center_lat = statistics.mean(lats)
    center_lon = statistics.mean(lons)
    
    # Find min and max timestamps for color scaling
    timestamps = [loc.get("ts", 0) for loc in locations]
    min_ts = min(timestamps)
    max_ts = max(timestamps)

    # Create the map centered on the data
    m = folium.Map(
        location=[center_lat, center_lon], zoom_start=15, tiles="OpenStreetMap"
    )

    # Add points to the map
    for i, loc in enumerate(locations):
        point_name = loc.get("app", f"Point {i + 1}")
        ts = loc.get("ts", 0)
        color = timestamp_to_color(ts, min_ts, max_ts)
        popup_text = f"{point_name}<br>Lat: {loc['lat']:.6f}<br>Lon: {loc['lon']:.6f}<br>Timestamp: {ts}"

        # Check if we have an app bundle ID to fetch icon
        app_bundle_id = loc.get("app")
        if app_bundle_id and "." in app_bundle_id:  # Bundle IDs contain dots
            # Try to get app icon
            icon_url = get_app_icon_url(app_bundle_id)

            if icon_url:
                # Use custom icon with colored border
                custom_icon = folium.CustomIcon(
                    icon_image=icon_url,
                    icon_size=(21, 21),
                    icon_anchor=(16, 16),
                    popup_anchor=(0, -16),
                )
                folium.Marker(
                    location=[loc["lat"], loc["lon"]],
                    popup=popup_text,
                    icon=custom_icon,
                ).add_to(m)
                
                # Add a colored circle behind the icon to show timestamp progression
                folium.CircleMarker(
                    location=[loc["lat"], loc["lon"]],
                    radius=8,
                    popup=popup_text,
                    color=color,
                    fill=True,
                    fillColor=color,
                    fillOpacity=0.3,
                    weight=2,
                ).add_to(m)
            else:
                # Fallback to circle marker with timestamp color
                folium.CircleMarker(
                    location=[loc["lat"], loc["lon"]],
                    radius=5,
                    popup=popup_text,
                    color=color,
                    fill=True,
                    fillColor=color,
                    fillOpacity=0.8,
                ).add_to(m)
        else:
            # No app data, use regular circle marker with timestamp color
            folium.CircleMarker(
                location=[loc["lat"], loc["lon"]],
                radius=5,
                popup=popup_text,
                color=color,
                fill=True,
                fillColor=color,
                fillOpacity=0.8,
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
