#!/usr/bin/env python3

import sqlite3
import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd
import numpy as np

def main():
    conn = sqlite3.connect("bssid_tracking.db")
    
    # Execute the SQL query from analyze_unique_locations.sql
    query = """
    -- Analyze the number of unique locations per BSSID
    -- This helps identify BSSIDs that oscillate between multiple locations
    
    -- Main query: Count unique locations per BSSID
    WITH location_counts AS (
        SELECT 
            bssid,
            COUNT(DISTINCT ROUND(old_lat, 6) || ',' || ROUND(old_long, 6)) as unique_old_locations,
            COUNT(DISTINCT ROUND(new_lat, 6) || ',' || ROUND(new_long, 6)) as unique_new_locations,
            COUNT(*) as total_changes
        FROM location_changes 
        GROUP BY bssid
    ),
    all_locations AS (
        -- Get all unique locations (both old and new) per BSSID
        SELECT 
            bssid,
            ROUND(old_lat, 6) as lat,
            ROUND(old_long, 6) as long
        FROM location_changes
        UNION
        SELECT 
            bssid,
            ROUND(new_lat, 6) as lat,
            ROUND(new_long, 6) as long
        FROM location_changes
    ),
    unique_location_counts AS (
        SELECT 
            bssid,
            COUNT(DISTINCT lat || ',' || long) as total_unique_locations
        FROM all_locations
        GROUP BY bssid
    )
    
    -- Main results
    SELECT 
        u.total_unique_locations,
        l.total_changes,
        COUNT(u.bssid) as bssid_count,
        ROUND(CAST(l.total_changes AS FLOAT) / u.total_unique_locations, 2) as changes_per_location
    FROM unique_location_counts u
    JOIN location_counts l ON u.bssid = l.bssid
    GROUP BY u.total_unique_locations, l.total_changes
    ORDER BY u.total_unique_locations DESC, l.total_changes DESC;
    """
    
    # Execute query and load into DataFrame
    df = pd.read_sql_query(query, conn)
    conn.close()
    
    if df.empty:
        print("No data found in database")
        return
    
    print(f"Found {len(df)} unique combinations of locations and changes")
    print("\nData sample:")
    print(df.head(10))
    
    # Create visualizations
    plt.style.use('default')
    fig, axes = plt.subplots(1, 2, figsize=(15, 6))
    fig.suptitle('BSSID Location Analysis', fontsize=16, fontweight='bold')
    
    # 1. Scatter plot: Changes per location vs BSSID count
    scatter = axes[0].scatter(
        df['changes_per_location'], 
        df['bssid_count'],
        c=df['total_unique_locations'],
        cmap='viridis',
        alpha=0.7,
        s=60
    )
    axes[0].set_xlabel('Changes per Location')
    axes[0].set_ylabel('BSSID Count')
    axes[0].set_title('BSSID Distribution by Changes per Location')
    axes[0].grid(True, alpha=0.3)
    cbar = plt.colorbar(scatter, ax=axes[0])
    cbar.set_label('Unique Locations')
    
    # 2. Log-scale scatter: Total changes vs BSSID count
    axes[1].scatter(
        df['total_changes'], 
        df['bssid_count'],
        alpha=0.6,
        c=df['total_unique_locations'],
        cmap='plasma'
    )
    axes[1].set_xlabel('Total Changes')
    axes[1].set_ylabel('BSSID Count')
    axes[1].set_title('BSSID Count vs Total Changes')
    axes[1].set_xscale('log')
    axes[1].set_yscale('log')
    axes[1].grid(True, alpha=0.3)
    
    plt.tight_layout()
    
    # Print summary statistics
    print(f"\n=== Summary Statistics ===")
    print(f"Total unique combinations: {len(df)}")
    print(f"Total BSSIDs analyzed: {df['bssid_count'].sum()}")
    print(f"Range of unique locations: {df['total_unique_locations'].min()} - {df['total_unique_locations'].max()}")
    print(f"Range of total changes: {df['total_changes'].min()} - {df['total_changes'].max()}")
    print(f"Range of changes per location: {df['changes_per_location'].min()} - {df['changes_per_location'].max()}")
    
    print(f"\n=== Top 10 Most Volatile BSSIDs (highest changes per location) ===")
    top_volatile = df.nlargest(10, 'changes_per_location')
    for _, row in top_volatile.iterrows():
        print(f"Locations: {row['total_unique_locations']}, Changes: {row['total_changes']}, "
              f"Changes/Location: {row['changes_per_location']}, BSSIDs: {row['bssid_count']}")
    
    print(f"\n=== Most Common Patterns ===")
    top_patterns = df.nlargest(10, 'bssid_count')
    for _, row in top_patterns.iterrows():
        print(f"Locations: {row['total_unique_locations']}, Changes: {row['total_changes']}, "
              f"BSSIDs: {row['bssid_count']}, Changes/Location: {row['changes_per_location']}")
    
    # Save the plot
    plt.savefig("output/unique_locations_analysis.png", dpi=300, bbox_inches="tight")
    print(f"\nPlot saved as 'output/unique_locations_analysis.png'")

if __name__ == "__main__":
    main()