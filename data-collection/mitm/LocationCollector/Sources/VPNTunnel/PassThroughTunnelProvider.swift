import NetworkExtension
import OSLog

public class PassThroughTunnelProvider: NEPacketTunnelProvider {
    private let logger = Logger(subsystem: "VPNTunnel", category: "Provider")
    
    public override func startTunnel(options: [String: NSObject]?, completionHandler: @escaping (Error?) -> Void) {
        logger.info("Starting pass-through tunnel")
        
        // Create tunnel settings with a dummy remote address
        let settings = NEPacketTunnelNetworkSettings(tunnelRemoteAddress: "127.0.0.1")
        
        // Configure the tunnel to route all traffic through it
        settings.ipv4Settings = {
            let ipv4Settings = NEIPv4Settings(addresses: ["192.168.1.1"], subnetMasks: ["255.255.255.0"])
            ipv4Settings.includedRoutes = [NEIPv4Route.default()]
            return ipv4Settings
        }()
        
        // Set up DNS settings to use system DNS
        settings.dnsSettings = {
            let dnsSettings = NEDNSSettings(servers: ["8.8.8.8", "8.8.4.4"])
            return dnsSettings
        }()
        
        // Apply the tunnel settings
        setTunnelNetworkSettings(settings) { [weak self] error in
            if let error = error {
                self?.logger.error("Failed to set tunnel settings: \(error)")
                completionHandler(error)
                return
            }
            
            self?.logger.info("Tunnel settings applied successfully")
            
            // Start reading packets from the tunnel interface
            self?.startPacketProcessing()
            
            // Signal that the tunnel is ready
            completionHandler(nil)
        }
    }
    
    public override func stopTunnel(with reason: NEProviderStopReason, completionHandler: @escaping () -> Void) {
        logger.info("Stopping pass-through tunnel, reason: \(reason)")
        completionHandler()
    }
    
    private func startPacketProcessing() {
        // Start reading packets from the virtual interface
        readPackets()
    }
    
    private func readPackets() {
        packetFlow.readPacketObjects { [weak self] packets in
            guard let self = self else { return }
            
            // In a real VPN, you would encrypt and send these packets to a server
            // For pass-through, we just immediately write them back to simulate routing
            
            let packetData = packets.map { $0.data }
            let protocols = packets.map { self.protocolNumber(for: $0.data) }
            
            // Log packet info for debugging
            self.logger.debug("Processing \(packets.count) packets")
            
            // For a true pass-through, we would send these to the default interface
            // Here we're just demonstrating the packet tunnel provider functionality
            
            // Write the packets back (this creates a loop but demonstrates the API)
            self.packetFlow.writePackets(packetData, withProtocols: protocols)
            
            // Continue reading
            self.readPackets()
        }
    }
    
    private func protocolNumber(for packet: Data) -> NSNumber {
        guard !packet.isEmpty else {
            return AF_INET as NSNumber
        }
        
        // Check IP version from the first 4 bits
        let ipVersion = (packet[0] & 0xf0) >> 4
        return (ipVersion == 6) ? AF_INET6 as NSNumber : AF_INET as NSNumber
    }
    
    public override func handleAppMessage(_ messageData: Data, completionHandler: ((Data?) -> Void)?) {
        logger.info("Received app message")
        completionHandler?(messageData)
    }
}