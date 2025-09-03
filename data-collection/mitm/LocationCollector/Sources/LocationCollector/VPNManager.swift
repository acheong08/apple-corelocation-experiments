import Foundation
@preconcurrency import NetworkExtension
import OSLog

@MainActor
public final class VPNManager: ObservableObject {
    @Published public private(set) var isConnected = false
    @Published public private(set) var status: NEVPNStatus = .invalid
    @Published public private(set) var tunnel: NETunnelProviderManager?
    
    public static let shared = VPNManager()
    
    private let logger = Logger(subsystem: "LocationCollector", category: "VPN")
    nonisolated(unsafe) private var statusObserver: NSObjectProtocol?
    
    private init() {
        Task {
            await loadExistingConfiguration()
            setupStatusObserver()
        }
    }
    
    public func loadExistingConfiguration() async {
        do {
            let managers = try await NETunnelProviderManager.loadAllFromPreferences()
            tunnel = managers.first
            updateStatus()
            logger.info("Loaded existing VPN configuration")
        } catch {
            logger.error("Failed to load VPN configuration: \(error)")
        }
    }
    
    public func installProfile() async throws {
        let manager = createTunnelManager()
        
        try await manager.saveToPreferences()
        
        // Reload to get the saved configuration
        try await manager.loadFromPreferences()
        
        tunnel = manager
        updateStatus()
        
        logger.info("VPN profile installed successfully")
    }
    
    public func removeProfile() async throws {
        guard let tunnel = tunnel else {
            throw VPNError.noConfiguration
        }
        
        try await tunnel.removeFromPreferences()
        self.tunnel = nil
        updateStatus()
        
        logger.info("VPN profile removed successfully")
    }
    
    public func connect() async throws {
        guard let tunnel = tunnel else {
            throw VPNError.noConfiguration
        }
        
        try tunnel.connection.startVPNTunnel()
        logger.info("VPN connection started")
    }
    
    public func disconnect() async {
        tunnel?.connection.stopVPNTunnel()
        logger.info("VPN connection stopped")
    }
    
    private func createTunnelManager() -> NETunnelProviderManager {
        let manager = NETunnelProviderManager()
        manager.localizedDescription = "Simple VPN"
        
        let protocolConfig = NETunnelProviderProtocol()
        // This should match your app's bundle identifier + ".VPNTunnel"
        protocolConfig.providerBundleIdentifier = "com.example.LocationCollector.VPNTunnel"
        protocolConfig.serverAddress = "127.0.0.1" // Dummy address for pass-through
        
        manager.protocolConfiguration = protocolConfig
        manager.isEnabled = true
        
        return manager
    }
    
    private func setupStatusObserver() {
        statusObserver = NotificationCenter.default.addObserver(
            forName: .NEVPNStatusDidChange,
            object: nil,
            queue: .main
        ) { [weak self] _ in
            Task { @MainActor in
                self?.updateStatus()
            }
        }
    }
    
    private func updateStatus() {
        status = tunnel?.connection.status ?? .invalid
        isConnected = status == .connected
    }
    
    deinit {
        if let observer = statusObserver {
            NotificationCenter.default.removeObserver(observer)
        }
    }
}

public enum VPNError: Error, LocalizedError {
    case noConfiguration
    case connectionFailed
    
    public var errorDescription: String? {
        switch self {
        case .noConfiguration:
            return "No VPN configuration found"
        case .connectionFailed:
            return "VPN connection failed"
        }
    }
}

extension NEVPNStatus: @retroactive CustomStringConvertible {
    public var description: String {
        switch self {
        case .invalid: return "Invalid"
        case .disconnected: return "Disconnected"
        case .connecting: return "Connecting"
        case .connected: return "Connected"
        case .reasserting: return "Reasserting"
        case .disconnecting: return "Disconnecting"
        @unknown default: return "Unknown"
        }
    }
}