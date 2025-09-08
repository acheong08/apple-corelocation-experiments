#!/usr/bin/env python3

import sqlite3
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
import seaborn as sns

def connect_db():
    """Connect to the database"""
    return sqlite3.connect('bssid_tracking.db')

def plot_update_frequency_timeline(conn):
    """Plot update frequency over time"""
    query = """
    SELECT 
        datetime(change_time, 'start of hour') as hour,
        COUNT(*) as update_count,
        COUNT(DISTINCT bssid) as unique_bssids,
        COUNT(DISTINCT old_tile_key) as tiles_affected
    FROM location_changes 
    GROUP BY datetime(change_time, 'start of hour')
    ORDER BY hour
    """
    
    df = pd.read_sql_query(query, conn)
    df['hour'] = pd.to_datetime(df['hour'])
    
    fig, (ax1, ax2, ax3) = plt.subplots(3, 1, figsize=(15, 12))
    
    # Updates per hour
    ax1.plot(df['hour'], df['update_count'], 'b-', linewidth=1)
    ax1.set_title('Location Updates Per Hour')
    ax1.set_ylabel('Updates')
    ax1.grid(True, alpha=0.3)
    
    # Unique BSSIDs affected per hour  
    ax2.plot(df['hour'], df['unique_bssids'], 'g-', linewidth=1)
    ax2.set_title('Unique BSSIDs Updated Per Hour')
    ax2.set_ylabel('Unique BSSIDs')
    ax2.grid(True, alpha=0.3)
    
    # Tiles affected per hour
    ax3.plot(df['hour'], df['tiles_affected'], 'r-', linewidth=1)
    ax3.set_title('Tiles Affected Per Hour')
    ax3.set_ylabel('Tiles')
    ax3.set_xlabel('Time')
    ax3.grid(True, alpha=0.3)
    
    plt.tight_layout()
    plt.savefig('output/update_frequency_timeline.png', dpi=300, bbox_inches='tight')
    print("Saved output/update_frequency_timeline.png")

def plot_batch_updates(conn):
    """Detect and visualize batch updates"""
    query = """
    WITH update_windows AS (
        SELECT 
            datetime(change_time, 'start of minute') as minute,
            COUNT(*) as updates_in_minute,
            COUNT(DISTINCT old_tile_key) as tiles_old,
            COUNT(DISTINCT new_tile_key) as tiles_new,
            COUNT(DISTINCT bssid) as bssids_affected
        FROM location_changes 
        GROUP BY datetime(change_time, 'start of minute')
    )
    SELECT * FROM update_windows 
    WHERE updates_in_minute >= 5
    ORDER BY minute
    """
    
    df = pd.read_sql_query(query, conn)
    df['minute'] = pd.to_datetime(df['minute'])
    
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(15, 10))
    
    # Scatter plot of batch sizes over time
    ax1.scatter(df['minute'], df['updates_in_minute'], alpha=0.6, s=20)
    ax1.set_title('Batch Update Detection (≥5 updates per minute)')
    ax1.set_ylabel('Updates per Minute')
    ax1.grid(True, alpha=0.3)
    
    # Histogram of batch sizes
    ax2.hist(df['updates_in_minute'], bins=50, alpha=0.7, edgecolor='black')
    ax2.set_title('Distribution of Batch Sizes')
    ax2.set_xlabel('Updates per Minute')
    ax2.set_ylabel('Frequency')
    ax2.grid(True, alpha=0.3)
    
    plt.tight_layout()
    plt.savefig('output/batch_updates_analysis.png', dpi=300, bbox_inches='tight')
    print("Saved output/batch_updates_analysis.png")
    
    # Print statistics
    print(f"Found {len(df)} minute-windows with ≥5 updates")
    print(f"Largest batch: {df['updates_in_minute'].max()} updates")
    print(f"Average batch size: {df['updates_in_minute'].mean():.1f} updates")

def plot_tile_update_patterns(conn):
    """Analyze tile update patterns"""
    sampled_tiles = [75908699, 75908700, 75908701, 75908702, 75908703,
                     78946920, 78946921, 78946922, 78946923, 78946924,
                     80391377, 80391378, 80391379, 80391380, 80391381,
                     81621497, 81621498, 81621499, 81621500, 81621501,
                     81857498, 81857499, 81857500, 81857501, 81857502]
    
    tile_list = ','.join(map(str, sampled_tiles))
    
    query = f"""
    WITH tile_updates AS (
        SELECT 
            old_tile_key,
            new_tile_key,
            datetime(change_time, 'start of hour') as hour,
            COUNT(*) as update_count
        FROM location_changes 
        WHERE old_tile_key IN ({tile_list}) OR new_tile_key IN ({tile_list})
        GROUP BY old_tile_key, new_tile_key, datetime(change_time, 'start of hour')
    )
    SELECT 
        hour,
        old_tile_key,
        SUM(update_count) as total_updates
    FROM tile_updates
    GROUP BY hour, old_tile_key
    ORDER BY hour, old_tile_key
    """
    
    df = pd.read_sql_query(query, conn)
    df['hour'] = pd.to_datetime(df['hour'])
    
    # Create a pivot table for heatmap
    pivot_df = df.pivot(index='hour', columns='old_tile_key', values='total_updates').fillna(0)
    
    plt.figure(figsize=(20, 12))
    sns.heatmap(pivot_df.T, cmap='YlOrRd', cbar_kws={'label': 'Updates per Hour'})
    plt.title('Tile Update Patterns Over Time (Sampled Tiles Only)')
    plt.xlabel('Time')
    plt.ylabel('Tile Key')
    plt.xticks(rotation=45)
    plt.tight_layout()
    plt.savefig('output/tile_update_heatmap.png', dpi=300, bbox_inches='tight')
    print("Saved output/tile_update_heatmap.png")

def plot_update_intervals(conn):
    """Analyze time intervals between updates"""
    query = """
    WITH time_intervals AS (
        SELECT 
            bssid,
            change_time,
            LAG(change_time) OVER (PARTITION BY bssid ORDER BY change_time) as prev_change,
            (julianday(change_time) - julianday(LAG(change_time) OVER (PARTITION BY bssid ORDER BY change_time))) * 24 * 60 as minutes_between
        FROM location_changes
    )
    SELECT minutes_between
    FROM time_intervals 
    WHERE minutes_between IS NOT NULL AND minutes_between < 1440  -- Less than 24 hours
    """
    
    df = pd.read_sql_query(query, conn)
    
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(15, 6))
    
    # Histogram of intervals (log scale)
    ax1.hist(df['minutes_between'], bins=100, alpha=0.7, edgecolor='black')
    ax1.set_yscale('log')
    ax1.set_title('Distribution of Time Between Updates (Log Scale)')
    ax1.set_xlabel('Minutes Between Updates')
    ax1.set_ylabel('Frequency (Log Scale)')
    ax1.grid(True, alpha=0.3)
    
    # Box plot for better statistics view
    ax2.boxplot(df['minutes_between'])
    ax2.set_title('Time Between Updates (Box Plot)')
    ax2.set_ylabel('Minutes Between Updates')
    ax2.grid(True, alpha=0.3)
    
    plt.tight_layout()
    plt.savefig('output/update_intervals_analysis.png', dpi=300, bbox_inches='tight')
    print("Saved output/update_intervals_analysis.png")
    
    # Print statistics
    print(f"Update interval statistics:")
    print(f"  Mean: {df['minutes_between'].mean():.1f} minutes")
    print(f"  Median: {df['minutes_between'].median():.1f} minutes")
    print(f"  Min: {df['minutes_between'].min():.1f} minutes")
    print(f"  Max: {df['minutes_between'].max():.1f} minutes")

def main():
    conn = connect_db()
    
    print("Generating temporal analysis visualizations...")
    
    try:
        plot_update_frequency_timeline(conn)
        plot_batch_updates(conn)
        plot_tile_update_patterns(conn)
        plot_update_intervals(conn)
        
        print("\nAll visualizations complete!")
        print("Generated files:")
        print("  - output/update_frequency_timeline.png")
        print("  - output/batch_updates_analysis.png") 
        print("  - output/tile_update_heatmap.png")
        print("  - output/update_intervals_analysis.png")
        
    finally:
        conn.close()

if __name__ == "__main__":
    main()