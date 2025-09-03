# Simple VPN Implementation

This project demonstrates a basic VPN client and packet tunnel provider implementation using NetworkExtension framework.

## What This Demonstrates

- **VPN Profile Installation**: How to programmatically install and manage VPN configurations
- **Packet Tunnel Provider**: A minimal implementation that passes traffic through
- **NetworkExtension Framework**: Basic usage of Apple's VPN APIs

## Key Components

### VPNManager
Located in `Sources/LocationCollector/VPNManager.swift`, this class handles:
- Installing/removing VPN profiles 
- Starting/stopping VPN connections
- Monitoring connection status

### PassThroughTunnelProvider  
Located in `Sources/VPNTunnel/PassThroughTunnelProvider.swift`, this provides:
- A basic packet tunnel provider implementation
- Pass-through functionality (no actual server)
- Demonstrates packet reading/writing APIs

## Usage

1. Build and run the app
2. Tap "Install VPN Profile" to create the VPN configuration
3. Tap "Connect" to start the tunnel (this will show in Settings > VPN)
4. The tunnel provider will process packets but pass them through unchanged

## Important Notes

- This is a simplified example for learning purposes
- No actual VPN server is required - it's a pass-through implementation  
- You'll need proper entitlements for NetworkExtension in a real app
- Bundle identifiers need to match between app and extension

## Based On

This implementation is inspired by the excellent article series:
- [VPN, Part 2: Packet Tunnel Provider](https://kean.blog/post/packet-tunnel-provider) by Alex Grebenyuk