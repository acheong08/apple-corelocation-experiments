#!/usr/bin/env python3

import sqlite3
import matplotlib.pyplot as plt
import pandas as pd

def connect_db():
    """Connect to the database"""
    return sqlite3.connect('bssid_tracking.db')

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
    
    print("Generating temporal analysis visualization...")
    
    try:
        plot_update_intervals(conn)
        
        print("\nVisualization complete!")
        print("Generated file:")
        print("  - output/update_intervals_analysis.png")
        
    finally:
        conn.close()

if __name__ == "__main__":
    main()