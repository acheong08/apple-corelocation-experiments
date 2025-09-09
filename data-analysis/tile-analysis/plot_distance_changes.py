#!/usr/bin/env python3

import sqlite3
import matplotlib.pyplot as plt
import numpy as np
from math import radians, cos, sin, asin, sqrt
from datetime import datetime


def haversine(lon1, lat1, lon2, lat2):
    """
    Calculate the great circle distance between two points
    on the earth (specified in decimal degrees)
    Returns distance in meters
    """
    lon1, lat1, lon2, lat2 = map(radians, [lon1, lat1, lon2, lat2])

    dlon = lon2 - lon1
    dlat = lat2 - lat1
    a = sin(dlat / 2) ** 2 + cos(lat1) * cos(lat2) * sin(dlon / 2) ** 2
    c = 2 * asin(sqrt(a))
    r = 6371000  # Radius of earth in meters
    return c * r


def main():
    conn = sqlite3.connect("bssid_tracking.db")
    cursor = conn.cursor()

    # Get all location changes ordered by time
    cursor.execute("""
        SELECT old_lat, old_long, new_lat, new_long 
        FROM location_changes 
        ORDER BY change_time
    """)

    changes = cursor.fetchall()
    conn.close()

    if not changes:
        print("No location changes found in database")
        return

    print(f"Found {len(changes)} location changes")

    # Calculate distances for each change
    distances = []

    for old_lat, old_long, new_lat, new_long in changes:
        distance = haversine(old_long, old_lat, new_long, new_lat)
        distances.append(distance)

    # Batch by 1000 changes and calculate averages
    batch_size = 1000
    batch_averages = []
    batch_indices = []

    for i in range(0, len(distances), batch_size):
        batch_distances = distances[i : i + batch_size]

        avg_distance = np.mean(batch_distances)
        # Use the middle index of the batch
        mid_index = i + len(batch_distances) // 2

        batch_averages.append(avg_distance)
        batch_indices.append(mid_index)

    # Create the plot
    plt.figure(figsize=(12, 8))
    plt.plot(batch_indices, batch_averages, "b-", linewidth=2, marker="o", markersize=4)
    plt.title(
        "Average Distance Moved Over Time\n(Batched by 1000 Location Changes)",
        fontsize=14,
    )
    plt.xlabel("Location Change Entry Number", fontsize=12)
    plt.ylabel("Average Distance (m)", fontsize=12)
    plt.grid(True, alpha=0.3)
    plt.tight_layout()

    # Add statistics
    overall_avg = float(np.mean(distances))
    plt.axhline(
        y=overall_avg,
        color="r",
        linestyle="--",
        alpha=0.7,
        label=f"Overall Average: {overall_avg:.1f} m",
    )
    plt.legend()

    # Print summary statistics
    print(f"\nSummary Statistics:")
    print(f"Total location changes: {len(distances)}")
    print(f"Number of batches: {len(batch_averages)}")
    print(f"Overall average distance: {overall_avg:.2f} m")
    print(f"Min distance: {min(distances):.2f} m")
    print(f"Max distance: {max(distances):.2f} m")
    print(f"Standard deviation: {np.std(distances):.2f} m")

    # Save the plot instead of showing it
    plt.savefig("output/distance_changes_plot.png", dpi=300, bbox_inches="tight")
    print(f"\nPlot saved as 'output/distance_changes_plot.png'")
    print(f"Batch size: {batch_size} location changes per data point")


if __name__ == "__main__":
    main()

