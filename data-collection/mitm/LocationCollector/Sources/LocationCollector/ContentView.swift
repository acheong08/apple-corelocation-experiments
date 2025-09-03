import SwiftUI
import NetworkExtension

public struct ContentView: View {
    @StateObject private var vpnManager = VPNManager.shared
    @State private var showingAlert = false
    @State private var alertMessage = ""
    
    public init() {}
    
    public var body: some View {
        NavigationView {
            VStack(spacing: 20) {
                // Status Section
                VStack(spacing: 10) {
                    Text("VPN Status")
                        .font(.headline)
                    
                    Text(vpnManager.status.description)
                        .font(.title2)
                        .foregroundColor(statusColor)
                        .padding()
                        .background(statusColor.opacity(0.1))
                        .cornerRadius(10)
                }
                
                Divider()
                
                // Controls Section
                VStack(spacing: 15) {
                    if vpnManager.tunnel == nil {
                        Button("Install VPN Profile") {
                            Task {
                                await installProfile()
                            }
                        }
                        .buttonStyle(.borderedProminent)
                    } else {
                        HStack(spacing: 15) {
                            if vpnManager.isConnected {
                                Button("Disconnect") {
                                    Task {
                                        await vpnManager.disconnect()
                                    }
                                }
                                .buttonStyle(.bordered)
                            } else {
                                Button("Connect") {
                                    Task {
                                        await connectVPN()
                                    }
                                }
                                .buttonStyle(.borderedProminent)
                            }
                            
                            Button("Remove Profile") {
                                Task {
                                    await removeProfile()
                                }
                            }
                            .buttonStyle(.bordered)
                            .foregroundColor(.red)
                        }
                    }
                }
                
                Spacer()
                
                // Info Section
                VStack(alignment: .leading, spacing: 10) {
                    Text("About This VPN")
                        .font(.headline)
                    
                    Text("This is a simple pass-through VPN implementation that demonstrates:")
                        .font(.caption)
                    
                    VStack(alignment: .leading, spacing: 5) {
                        Text("• Installing VPN profiles")
                        Text("• Packet Tunnel Provider basics")
                        Text("• NetworkExtension framework")
                    }
                    .font(.caption)
                    .foregroundColor(.secondary)
                }
                .padding()
                .background(Color(.systemGray6))
                .cornerRadius(10)
            }
            .padding()
            .navigationTitle("Simple VPN")
            .alert("VPN Manager", isPresented: $showingAlert) {
                Button("OK") { }
            } message: {
                Text(alertMessage)
            }
        }
    }
    
    private var statusColor: Color {
        switch vpnManager.status {
        case .connected:
            return .green
        case .connecting, .disconnecting, .reasserting:
            return .orange
        case .disconnected:
            return .gray
        case .invalid:
            return .red
        @unknown default:
            return .gray
        }
    }
    
    private func installProfile() async {
        do {
            try await vpnManager.installProfile()
            alertMessage = "VPN profile installed successfully!"
            showingAlert = true
        } catch {
            alertMessage = "Failed to install VPN profile: \(error.localizedDescription)"
            showingAlert = true
        }
    }
    
    private func removeProfile() async {
        do {
            try await vpnManager.removeProfile()
            alertMessage = "VPN profile removed successfully!"
            showingAlert = true
        } catch {
            alertMessage = "Failed to remove VPN profile: \(error.localizedDescription)"
            showingAlert = true
        }
    }
    
    private func connectVPN() async {
        do {
            try await vpnManager.connect()
        } catch {
            alertMessage = "Failed to connect VPN: \(error.localizedDescription)"
            showingAlert = true
        }
    }
}
