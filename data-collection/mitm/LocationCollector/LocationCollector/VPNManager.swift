//
//  VPNManager.swift
//  LocationCollector
//
//  Created by Antonio on 03/09/2025.
//

import NetworkExtension
import os.log

class VPNManager: ObservableObject {
    static let shared = VPNManager()
    
    @Published var isConnected = false
    @Published var connectionStatus = "Not Installed"
    @Published var isInstalled = false
    
    private let logger = Logger(subsystem: "dev.duti.LocationCollector", category: "VPNManager")
    private let certificateServer = CertificateServer.shared
    private var manager: NETunnelProviderManager?
    
    private init() {
        loadManager()
        setupStatusObserver()
    }
    
    private func loadManager() {
        NETunnelProviderManager.loadAllFromPreferences { [weak self] managers, error in
            if let error = error {
                self?.logger.error("Failed to load VPN managers: \(error.localizedDescription)")
                return
            }
            
            DispatchQueue.main.async {
                if let existingManager = managers?.first {
                    self?.manager = existingManager
                    self?.isInstalled = true
                    self?.updateConnectionStatus()
                } else {
                    self?.manager = nil
                    self?.isInstalled = false
                    self?.connectionStatus = "Not Installed"
                }
            }
        }
    }
    
    func installProfile(_ completion: @escaping (Result<Void, Error>) -> Void) {
        let manager = NETunnelProviderManager()
        manager.localizedDescription = "Location Collector VPN"
        
        let providerProtocol = NETunnelProviderProtocol()
        providerProtocol.providerBundleIdentifier = "dev.duti.LocationCollector.LocationCollectorVPN"
        providerProtocol.serverAddress = "LocationCollector"
        manager.protocolConfiguration = providerProtocol
        manager.isEnabled = true
        
        manager.saveToPreferences { [weak self] error in
            if let error = error {
                self?.logger.error("Failed to save VPN configuration: \(error.localizedDescription)")
                completion(.failure(error))
                return
            }
            
            // Load the manager again to ensure it's properly configured
            manager.loadFromPreferences { [weak self] error in
                DispatchQueue.main.async {
                    if let error = error {
                        self?.logger.error("Failed to load VPN configuration after save: \(error.localizedDescription)")
                        completion(.failure(error))
                    } else {
                        self?.logger.info("VPN configuration installed successfully")
                        self?.manager = manager
                        self?.isInstalled = true
                        self?.updateConnectionStatus()
                        completion(.success(()))
                    }
                }
            }
        }
    }
    
    private func setupStatusObserver() {
        NotificationCenter.default.addObserver(
            forName: .NEVPNStatusDidChange,
            object: nil,
            queue: .main
        ) { [weak self] _ in
            self?.updateConnectionStatus()
        }
    }
    
    private func updateConnectionStatus() {
        guard let manager = manager else {
            connectionStatus = "Not Installed"
            isConnected = false
            return
        }
        
        switch manager.connection.status {
        case .connecting:
            connectionStatus = "Connecting..."
            isConnected = false
        case .connected:
            connectionStatus = "Connected"
            isConnected = true
            // Start certificate server when VPN connects
            certificateServer.startServer()
        case .disconnecting:
            connectionStatus = "Disconnecting..."
            isConnected = false
        case .disconnected:
            connectionStatus = "Disconnected"
            isConnected = false
            // Stop certificate server when VPN disconnects
            certificateServer.stopServer()
        case .invalid:
            connectionStatus = "Invalid"
            isConnected = false
        case .reasserting:
            connectionStatus = "Reconnecting..."
            isConnected = false
        @unknown default:
            connectionStatus = "Unknown"
            isConnected = false
        }
        
        logger.info("VPN status changed to: \(self.connectionStatus)")
    }
    
    func toggleVPN() {
        guard let manager = manager else {
            logger.error("VPN manager not available - profile not installed")
            return
        }
        
        if manager.connection.status == .connected {
            disconnect()
        } else {
            connect()
        }
    }
    
    private func connect() {
        guard let manager = manager else { return }
        
        do {
            try manager.connection.startVPNTunnel()
            logger.info("Starting VPN tunnel")
        } catch {
            logger.error("Failed to start VPN: \(error.localizedDescription)")
        }
    }
    
    private func disconnect() {
        guard let manager = manager else { return }
        
        manager.connection.stopVPNTunnel()
        logger.info("Stopping VPN tunnel")
    }
    
    private func createManager() {
        let manager = NETunnelProviderManager()
        manager.localizedDescription = "Location Collector VPN"
        
        let providerProtocol = NETunnelProviderProtocol()
        providerProtocol.providerBundleIdentifier = "dev.duti.LocationCollector.LocationCollectorVPN"
        providerProtocol.serverAddress = "LocationCollector"
        manager.protocolConfiguration = providerProtocol
        manager.isEnabled = true
        
        manager.saveToPreferences { [weak self] error in
            if let error = error {
                self?.logger.error("Failed to save VPN configuration: \(error.localizedDescription)")
            } else {
                self?.logger.info("VPN configuration saved successfully")
                self?.manager = manager
                DispatchQueue.main.async {
                    self?.updateConnectionStatus()
                }
            }
        }
    }
}
