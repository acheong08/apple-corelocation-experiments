# for i in {xxx..xxx}; do go run ./cmd/wloc tile -key $i | awk '{ print $4,$5 }' > $i.txt; done

import json
import os
import statistics

import folium
import morton

# Create a Folium map
map = folium.Map(
    location=[
        51.506889,
        -3.192196,
    ],
    zoom_start=10,
)


def plot_points(file, color="green"):
    coords = [line.strip().split() for line in open(file, "r").readlines()]
    if len(coords) == 0:
        return
    # coord is the mean of long/lat
    coord = [
        statistics.mean([float(coord[0]) for coord in coords]),
        statistics.mean([float(coord[1]) for coord in coords]),
    ]
    file = int(file[:-4])
    morton_coord = morton.deinterleave_32(file)
    tooltip = str(morton_coord)
    print(json.dumps({"coord": coord, "morton": morton_coord}))
    # Plot each coordinate as a marker
    folium.Marker(coord, tooltip=tooltip, icon=folium.Icon(color=color)).add_to(map)


# List .txt files
files = os.listdir()
files = [file for file in files if file.endswith(".txt")]

# A different color for each
colors = [
    "red",
    "blue",
    "green",
    "purple",
    "orange",
    "darkred",
    "lightred",
    "beige",
    "darkblue",
    "darkgreen",
    "cadetblue",
    "darkpurple",
    "white",
    "pink",
    "lightblue",
    "lightgreen",
    "gray",
    "black",
]

# Plot each file
for i, file in enumerate(files):
    plot_points(file, colors[i % len(colors)])

# Display the map
map.save("map.html")
