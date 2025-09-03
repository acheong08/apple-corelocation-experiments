//
//  PacketTunnelProvider.swift
//  LocationCollectorVPN
//
//  Created by Antonio on 03/09/2025.
//

import NetworkExtension
import Network
import os.log

class PacketTunnelProvider: NEPacketTunnelProvider {
    
    private let logger = Logger(subsystem: "dev.duti.LocationCollector.VPN", category: "PacketTunnel")
    
    // Target domains for Apple location services
    private let targetDomains = [
        "gs-loc.apple.com",
        "ls.apple.com"
    ]
    
    // Certificate server domain redirect
    private let certificateDomain = "mitm.it"
    private let certificateServerIP = "10.0.0.1"
    
    override func startTunnel(options: [String : NSObject]?, completionHandler: @escaping (Error?) -> Void) {
        logger.info("Starting VPN tunnel")
        
        let settings = NEPacketTunnelNetworkSettings(tunnelRemoteAddress: "127.0.0.1")
        
        // Configure IPv4 settings
        let ipv4Settings = NEIPv4Settings(addresses: ["10.0.0.1"], subnetMasks: ["255.255.255.0"])
        ipv4Settings.includedRoutes = [NEIPv4Route.default()]
        settings.ipv4Settings = ipv4Settings
        
        // Configure DNS settings to use system DNS  
        let dnsSettings = NEDNSSettings(servers: ["8.8.8.8", "8.8.4.4"])
        dnsSettings.matchDomains = [""]
        settings.dnsSettings = dnsSettings
        
        setTunnelNetworkSettings(settings) { [weak self] error in
            if let error = error {
                self?.logger.error("Failed to set tunnel network settings: \(error.localizedDescription)")
                completionHandler(error)
            } else {
                self?.logger.info("VPN tunnel started successfully")
                #if DEBUG
                self?.testSNIDetection()
                #endif
                self?.startPacketForwarding()
                completionHandler(nil)
            }
        }
    }
    
    override func stopTunnel(with reason: NEProviderStopReason, completionHandler: @escaping () -> Void) {
        logger.info("Stopping VPN tunnel with reason: \(reason.rawValue)")
        completionHandler()
    }
    
    override func handleAppMessage(_ messageData: Data, completionHandler: ((Data?) -> Void)?) {
        if let handler = completionHandler {
            handler(messageData)
        }
    }
    
    private func startPacketForwarding() {
        logger.info("Starting packet forwarding")
        
        packetFlow.readPackets { [weak self] packets, protocols in
            guard let self = self else { return }
            
            var packetsToForward: [Data] = []
            var protocolsToForward: [NSNumber] = []
            
            // Process each packet
            for (index, packet) in packets.enumerated() {
                let protocolNumber = protocols[index]
                
                // Check if this is a packet we need to intercept
                if let targetDomain = self.extractTargetDomain(from: packet, protocol: protocolNumber) {
                    self.logger.info("🎯 INTERCEPTED: Target domain detected: \(targetDomain)")
                    
                    // For now, we'll still forward the packet but log that we detected it
                    // TODO: In future checkpoints, this is where we'll redirect to local proxy
                    packetsToForward.append(packet)
                    protocolsToForward.append(protocolNumber)
                } else {
                    // Regular packet - forward as normal
                    packetsToForward.append(packet)
                    protocolsToForward.append(protocolNumber)
                }
            }
            
            // Forward packets
            if !packetsToForward.isEmpty {
                self.packetFlow.writePackets(packetsToForward, withProtocols: protocolsToForward)
            }
            
            // Continue reading packets
            self.startPacketForwarding()
        }
    }
    
    private func extractTargetDomain(from packet: Data, protocol protocolNumber: NSNumber) -> String? {
        guard packet.count >= 20 else { return nil }
        
        let ipVersion = (packet[0] & 0xF0) >> 4
        guard ipVersion == 4 else { return nil }
        
        let ipHeaderLength = Int(packet[0] & 0x0F) * 4
        let protocol = packet[9]
        
        // Only process TCP packets (protocol 6)
        guard protocol == 6, packet.count >= ipHeaderLength + 20 else { return nil }
        
        let tcpHeaderOffset = ipHeaderLength
        let tcpHeaderLength = Int(packet[tcpHeaderOffset + 12] >> 4) * 4
        let tcpPayloadOffset = ipHeaderLength + tcpHeaderLength
        
        guard packet.count > tcpPayloadOffset else { return nil }
        
        let destPort = UInt16(packet[tcpHeaderOffset + 2]) << 8 | UInt16(packet[tcpHeaderOffset + 3])
        
        // Check if this is HTTPS traffic (port 443)
        guard destPort == 443 else { return nil }
        
        let tcpPayload = packet.subdata(in: tcpPayloadOffset..<packet.count)
        
        // Parse TLS and extract SNI
        return extractSNI(from: tcpPayload)
    }
    
    private func extractSNI(from tlsData: Data) -> String? {
        guard tlsData.count >= 5 else { return nil }
        
        // Check if this is a TLS handshake record (type 22)
        guard tlsData[0] == 0x16 else { return nil }
        
        // Skip TLS record header (5 bytes)
        let handshakeData = tlsData.subdata(in: 5..<tlsData.count)
        guard handshakeData.count >= 4 else { return nil }
        
        // Check if this is a Client Hello (type 1)
        guard handshakeData[0] == 0x01 else { return nil }
        
        // Parse Client Hello message
        let clientHelloLength = Int(handshakeData[1]) << 16 | Int(handshakeData[2]) << 8 | Int(handshakeData[3])
        guard handshakeData.count >= 4 + clientHelloLength else { return nil }
        
        let clientHello = handshakeData.subdata(in: 4..<(4 + clientHelloLength))
        return parseSNIFromClientHello(clientHello)
    }
    
    private func parseSNIFromClientHello(_ clientHello: Data) -> String? {
        guard clientHello.count >= 38 else { return nil }
        
        var offset = 34 // Skip version (2) + random (32)
        
        // Skip session ID
        guard offset < clientHello.count else { return nil }
        let sessionIdLength = Int(clientHello[offset])
        offset += 1 + sessionIdLength
        
        // Skip cipher suites
        guard offset + 1 < clientHello.count else { return nil }
        let cipherSuitesLength = Int(clientHello[offset]) << 8 | Int(clientHello[offset + 1])
        offset += 2 + cipherSuitesLength
        
        // Skip compression methods
        guard offset < clientHello.count else { return nil }
        let compressionMethodsLength = Int(clientHello[offset])
        offset += 1 + compressionMethodsLength
        
        // Parse extensions
        guard offset + 1 < clientHello.count else { return nil }
        let extensionsLength = Int(clientHello[offset]) << 8 | Int(clientHello[offset + 1])
        offset += 2
        
        let extensionsEnd = offset + extensionsLength
        
        while offset + 3 < extensionsEnd && offset + 3 < clientHello.count {
            let extensionType = Int(clientHello[offset]) << 8 | Int(clientHello[offset + 1])
            let extensionLength = Int(clientHello[offset + 2]) << 8 | Int(clientHello[offset + 3])
            offset += 4
            
            // SNI extension type is 0
            if extensionType == 0 && offset + extensionLength <= clientHello.count {
                return parseSNIExtension(clientHello.subdata(in: offset..<(offset + extensionLength)))
            }
            
            offset += extensionLength
        }
        
        return nil
    }
    
    private func parseSNIExtension(_ sniData: Data) -> String? {
        guard sniData.count >= 5 else { return nil }
        
        let listLength = Int(sniData[0]) << 8 | Int(sniData[1])
        guard sniData.count >= 2 + listLength else { return nil }
        
        var offset = 2
        
        while offset + 2 < sniData.count {
            let nameType = sniData[offset]
            let nameLength = Int(sniData[offset + 1]) << 8 | Int(sniData[offset + 2])
            offset += 3
            
            // Host name type is 0
            if nameType == 0 && offset + nameLength <= sniData.count {
                let hostname = String(data: sniData.subdata(in: offset..<(offset + nameLength)), encoding: .utf8)
                
                // Check if this hostname matches our target domains
                if let hostname = hostname, isTargetDomain(hostname) {
                    return hostname
                }
            }
            
            offset += nameLength
        }
        
        return nil
    }
    
    private func isTargetDomain(_ hostname: String) -> Bool {
        let lowercaseHostname = hostname.lowercased()
        
        for domain in targetDomains {
            if lowercaseHostname == domain || lowercaseHostname.hasSuffix("." + domain) {
                return true
            }
        }
        
        return false
    }
    
    // MARK: - Testing utilities
    #if DEBUG
    func testSNIDetection() {
        logger.info("🧪 Testing SNI detection functionality")
        
        // Test target domain matching
        let testDomains = [
            "gs-loc.apple.com",
            "ls.apple.com", 
            "sub.ls.apple.com",
            "other.apple.com",
            "google.com"
        ]
        
        for domain in testDomains {
            let isTarget = isTargetDomain(domain)
            logger.info("🧪 Domain '\(domain)': \(isTarget ? "TARGET ✅" : "ignored ❌")")
        }
    }
    #endif
    
    private func formatIPv4Address(_ packet: Data, offset: Int) -> String {
        return "\(packet[offset]).\(packet[offset+1]).\(packet[offset+2]).\(packet[offset+3])"
    }
}
