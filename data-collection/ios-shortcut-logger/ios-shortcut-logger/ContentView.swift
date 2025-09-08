import SwiftUI
import UIKit

struct ContentView: View {
    @StateObject private var dataManager = DataManager.shared
    @State private var selectedTab = 0
    
    var body: some View {
        TabView(selection: $selectedTab) {
            SystemStateView()
                .tabItem {
                    Image(systemName: "gear")
                    Text("System")
                }
                .tag(0)
            
            RecentLogsView()
                .tabItem {
                    Image(systemName: "clock")
                    Text("All Logs")
                }
                .tag(1)
        }
        .environmentObject(dataManager)
    }
}

struct SystemStateView: View {
    @EnvironmentObject var dataManager: DataManager
    @State private var isUpdating = false
    @State private var showingShareSheet = false
    @State private var shareItems: [Any] = []
    @State private var airplaneModeOn = false
    
    var body: some View {
        NavigationView {
            VStack(spacing: 20) {
                VStack(alignment: .leading, spacing: 10) {
                    Text("Manual State Updates")
                        .font(.headline)
                    
                    HStack {
                        Text("Airplane Mode:")
                        Spacer()
                        Toggle("", isOn: $airplaneModeOn)
                    }
                    .padding()
                    .background(Color(.systemGray6))
                    .cornerRadius(8)
                }
                
                Button(action: updateStates) {
                    HStack {
                        if isUpdating {
                            ProgressView()
                                .scaleEffect(0.8)
                        } else {
                            Image(systemName: "arrow.clockwise")
                        }
                        Text("Update States")
                    }
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(Color.blue)
                    .foregroundColor(.white)
                    .cornerRadius(10)
                }
                .disabled(isUpdating)
                
                Button(action: exportLogs) {
                    HStack {
                        Image(systemName: "square.and.arrow.up")
                        Text("Share Logs")
                    }
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(Color.green)
                    .foregroundColor(.white)
                    .cornerRadius(10)
                }
                
                Button(action: clearLogs) {
                    HStack {
                        Image(systemName: "trash")
                        Text("Clear Logs")
                    }
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(Color.red)
                    .foregroundColor(.white)
                    .cornerRadius(10)
                }
                
                Spacer()
            }
            .padding()
            .navigationTitle("System State")
        }
        .sheet(isPresented: $showingShareSheet) {
            ActivityViewController(activityItems: shareItems)
        }
    }
    
    private func updateStates() {
        isUpdating = true
        Task {
            do {
                try await updateStatesAndLog(airplaneModeOverride: airplaneModeOn)
                await MainActor.run {
                    isUpdating = false
                }
            } catch {
                await MainActor.run {
                    isUpdating = false
                }
            }
        }
    }
    
    private func exportLogs() {
        let logContent = dataManager.getLogContent()
        shareItems = [logContent]
        showingShareSheet = true
    }
    
    private func clearLogs() {
        Task {
            do {
                try await dataManager.clearLogs()
            } catch {
                print("Failed to clear logs: \(error)")
            }
        }
    }
}

struct RecentLogsView: View {
    @EnvironmentObject var dataManager: DataManager
    @State private var allLogs: [LogEntry] = []
    
    var body: some View {
        NavigationView {
            List(allLogs) { logEntry in
                LogEntryRow(logEntry: logEntry)
            }
            .navigationTitle("All Logs")
            .onAppear {
                loadAllLogs()
            }
            .refreshable {
                loadAllLogs()
            }
            .onChange(of: dataManager.recentLogs) {
                loadAllLogs()
            }
        }
    }
    
    private func loadAllLogs() {
        allLogs = dataManager.getAllLogs()
    }
}

struct LogEntryRow: View {
    let logEntry: LogEntry
    
    var body: some View {
        HStack {
            Text(logEntry.eventTypeName)
                .font(.headline)
            Spacer()
            Text(logEntry.timestamp, style: .time)
                .font(.caption)
                .foregroundColor(.secondary)
        }
        .padding(.vertical, 2)
    }
}

// MARK: - Activity View Controller for Sharing

struct ActivityViewController: UIViewControllerRepresentable {
    let activityItems: [Any]
    
    func makeUIViewController(context: Context) -> UIActivityViewController {
        let controller = UIActivityViewController(activityItems: activityItems, applicationActivities: nil)
        return controller
    }
    
    func updateUIViewController(_ uiViewController: UIActivityViewController, context: Context) {}
}

// MARK: - Update States Function

private func updateStatesAndLog(airplaneModeOverride: Bool? = nil) async throws {
    let changes = await StateTracker.shared.updateSystemStates(airplaneModeOverride: airplaneModeOverride)
    let dataManager = await DataManager.shared
    
    for (stateKey, _, newValue) in changes {
        let eventTypeName: String
        switch stateKey {
        case "AirplaneMode":
            eventTypeName = newValue ? "AirplaneModeOn" : "AirplaneModeOff"
        case "PlugState":
            eventTypeName = newValue ? "PluggedIn" : "PluggedOut"
        default:
            continue
        }
        
        try await dataManager.logEvent(eventTypeName: eventTypeName)
    }
}
