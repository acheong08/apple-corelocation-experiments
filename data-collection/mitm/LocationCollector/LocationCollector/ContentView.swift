//
//  ContentView.swift
//  LocationCollector
//
//  Created by Antonio on 03/09/2025.
//

import SwiftUI

struct ContentView: View {
    @StateObject private var vpnManager = VPNManager.shared
    @StateObject private var certificateManager = CertificateManager.shared
    @StateObject private var certificateServer = CertificateServer.shared
    @State private var isInstalling = false
    @State private var showingError = false
    @State private var errorMessage = ""
    
    var body: some View {
        VStack(spacing: 30) {
            Image(systemName: "shield.lefthalf.filled")
                .imageScale(.large)
                .foregroundStyle(.tint)
                .font(.system(size: 60))
            
            Text("Location Collector VPN")
                .font(.title)
                .fontWeight(.bold)
            
            VStack(spacing: 10) {
                Text("Status: \(vpnManager.connectionStatus)")
                    .font(.headline)
                    .foregroundColor(vpnManager.isConnected ? .green : .gray)
                
                if vpnManager.isConnected {
                    Text("🎯 Monitoring Apple Location Services")
                        .font(.subheadline)
                        .foregroundColor(.orange)
                        .padding(.top, 5)
                }
                
                // Certificate status
                HStack {
                    Text("Root CA:")
                        .font(.subheadline)
                    Text(certificateManager.isRootCAGenerated ? "✅ Generated" : "❌ Not Generated")
                        .font(.subheadline)
                        .foregroundColor(certificateManager.isRootCAGenerated ? .green : .red)
                }
                .padding(.vertical, 5)
                
                if !certificateManager.isRootCAGenerated {
                    Button(action: {
                        certificateManager.generateRootCA()
                    }) {
                        Text("Generate Root Certificate")
                            .font(.title3)
                            .fontWeight(.medium)
                            .foregroundColor(.white)
                            .padding()
                            .frame(maxWidth: .infinity)
                            .background(Color.orange)
                            .cornerRadius(8)
                    }
                } else if certificateServer.isRunning && !certificateServer.serverURL.isEmpty {
                    VStack(spacing: 10) {
                        Text("🌐 Certificate Server Running")
                            .font(.subheadline)
                            .foregroundColor(.green)
                        
                        Text("Open in Safari:")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        
                        VStack(spacing: 6) {
                            Text(certificateServer.serverURL)
                                .font(.system(.body, design: .monospaced))
                                .foregroundColor(.blue)
                                .padding(.horizontal, 12)
                                .padding(.vertical, 8)
                                .background(Color(.systemGray6))
                                .cornerRadius(6)
                                .onTapGesture {
                                    UIPasteboard.general.string = certificateServer.serverURL
                                }
                            
                            Text("(or visit mitm.it - redirected to local server)")
                                .font(.caption2)
                                .foregroundColor(.secondary)
                                .italic()
                        }
                        
                        Text("Tap URL to copy")
                            .font(.caption2)
                            .foregroundColor(.secondary)
                    }
                } else if certificateManager.isRootCAGenerated {
                    Text("📋 Certificate ready - connect VPN to serve")
                        .font(.subheadline)
                        .foregroundColor(.orange)
                }
                
                if !vpnManager.isInstalled {
                    Button(action: {
                        installProfile()
                    }) {
                        Text(isInstalling ? "Installing..." : "Install VPN Profile")
                            .font(.title2)
                            .fontWeight(.semibold)
                            .foregroundColor(.white)
                            .padding()
                            .frame(maxWidth: .infinity)
                            .background(Color.blue)
                            .cornerRadius(10)
                    }
                    .disabled(isInstalling)
                } else {
                    Button(action: {
                        vpnManager.toggleVPN()
                    }) {
                        Text(vpnManager.isConnected ? "Disconnect VPN" : "Connect VPN")
                            .font(.title2)
                            .fontWeight(.semibold)
                            .foregroundColor(.white)
                            .padding()
                            .frame(maxWidth: .infinity)
                            .background(vpnManager.isConnected ? Color.red : Color.blue)
                            .cornerRadius(10)
                    }
                    .disabled(vpnManager.connectionStatus.contains("ing"))
                }
            }
            
            Text("Connect the VPN, then visit the certificate URL in Safari to download and install the root certificate for HTTPS interception.")
                .font(.caption)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
                .padding(.horizontal)
            
            Spacer()
        }
        .padding()
        .alert("Installation Failed", isPresented: $showingError) {
            Button("OK", role: .cancel) { }
        } message: {
            Text(errorMessage)
        }
        .onAppear {
            // Generate certificate on first launch if needed
            if !certificateManager.isRootCAGenerated {
                certificateManager.generateRootCA()
            }
        }
    }
    
    private func installProfile() {
        isInstalling = true
        
        vpnManager.installProfile { result in
            DispatchQueue.main.async {
                isInstalling = false
                switch result {
                case .success:
                    break // UI will update automatically via @Published properties
                case .failure(let error):
                    errorMessage = error.localizedDescription
                    showingError = true
                }
            }
        }
    }
}

#Preview {
    ContentView()
}
