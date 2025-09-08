import SwiftUI
import UIKit

struct ContentView: View {
    @StateObject private var dataManager = DataManager.shared
    @State private var selectedTab = 0
    @State private var showingUpdateAlert = false
    @State private var updateMessage = ""
    
    var body: some View {
        TabView(selection: $selectedTab) {
            SystemStateView()
                .tabItem {
                    Image(systemName: "gear")
                    Text("System")
                }
                .tag(0)
            
            OutputLocationsView()
                .tabItem {
                    Image(systemName: "folder")
                    Text("Output")
                }
                .tag(1)
            
            RecentLogsView()
                .tabItem {
                    Image(systemName: "clock")
                    Text("Recent Logs")
                }
                .tag(2)
        }
        .environmentObject(dataManager)
        .alert("State Update", isPresented: $showingUpdateAlert) {
            Button("OK") { }
        } message: {
            Text(updateMessage)
        }
    }
}

struct SystemStateView: View {
    @EnvironmentObject var dataManager: DataManager
    @State private var isUpdating = false
    @State private var showingShareSheet = false
    @State private var shareItems: [Any] = []
    @State private var airplaneModeOn = false
    @State private var showingAirplaneModeAlert = false
    
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
                        Text("Export Logs")
                    }
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(Color.green)
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
        let logFiles = dataManager.getAllLogFiles()
        if !logFiles.isEmpty {
            shareItems = logFiles
            showingShareSheet = true
        }
    }
}



struct OutputLocationsView: View {
    @EnvironmentObject var dataManager: DataManager
    @State private var showingAddLocation = false
    
    var body: some View {
        NavigationView {
            List {
                ForEach(dataManager.outputLocations) { location in
                    OutputLocationRow(location: location)
                }
                .onDelete(perform: deleteLocations)
            }
            .navigationTitle("Output Locations")
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Add") {
                        showingAddLocation = true
                    }
                }
            }
            .sheet(isPresented: $showingAddLocation) {
                AddOutputLocationView()
            }
        }
    }
    
    private func deleteLocations(offsets: IndexSet) {
        for index in offsets {
            dataManager.removeOutputLocation(dataManager.outputLocations[index])
        }
    }
}

struct OutputLocationRow: View {
    let location: OutputLocation
    @EnvironmentObject var dataManager: DataManager
    
    var body: some View {
        HStack {
            VStack(alignment: .leading, spacing: 4) {
                Text(location.name)
                    .font(.headline)
                Text(location.filePath)
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            Spacer()
            
            Toggle("", isOn: .constant(location.isActive))
                .onChange(of: location.isActive) {
                    // Note: This would need proper updating logic
                }
        }
    }
}

struct AddOutputLocationView: View {
    @EnvironmentObject var dataManager: DataManager
    @Environment(\.dismiss) var dismiss
    
    @State private var name = ""
    @State private var filePath = ""
    
    var body: some View {
        NavigationView {
            Form {
                TextField("Name", text: $name)
                TextField("File Path", text: $filePath)
                    .autocapitalization(.none)
                    .placeholder(when: filePath.isEmpty) {
                        Text("logs/events.jsonl")
                            .foregroundColor(.secondary)
                    }
            }
            .navigationTitle("New Output Location")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Save") {
                        let location = OutputLocation(
                            name: name,
                            filePath: filePath.isEmpty ? "logs/events.jsonl" : filePath
                        )
                        dataManager.addOutputLocation(location)
                        dismiss()
                    }
                    .disabled(name.isEmpty)
                }
            }
        }
    }
}

struct RecentLogsView: View {
    @EnvironmentObject var dataManager: DataManager
    
    var body: some View {
        NavigationView {
            List(dataManager.recentLogs) { logEntry in
                LogEntryRow(logEntry: logEntry)
            }
            .navigationTitle("Recent Logs")
        }
    }
}

struct LogEntryRow: View {
    let logEntry: LogEntry
    
    var body: some View {
        VStack(alignment: .leading, spacing: 4) {
            HStack {
                Text(logEntry.eventTypeName)
                    .font(.headline)
                Spacer()
                Text(logEntry.timestamp, style: .time)
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            
            if !logEntry.data.isEmpty {
                ForEach(Array(logEntry.data.keys.sorted()), id: \.self) { key in
                    if let value = logEntry.data[key] {
                        HStack {
                            Text(key + ":")
                                .font(.caption)
                                .foregroundColor(.secondary)
                            Text(stringValue(for: value))
                                .font(.caption)
                        }
                    }
                }
            }
            
            if logEntry.location != nil {
                Text("📍 Location included")
                    .font(.caption2)
                    .foregroundColor(.blue)
            }
        }
        .padding(.vertical, 2)
    }
    
    private func stringValue(for logValue: LogValue) -> String {
        switch logValue {
        case .text(let string):
            return string
        case .number(let double):
            return String(double)
        case .boolean(let bool):
            return bool ? "true" : "false"
        case .date(let date):
            return DateFormatter.localizedString(from: date, dateStyle: .short, timeStyle: .short)
        }
    }
}

extension View {
    func placeholder<Content: View>(
        when shouldShow: Bool,
        alignment: Alignment = .leading,
        @ViewBuilder placeholder: () -> Content) -> some View {
        
        ZStack(alignment: alignment) {
            placeholder().opacity(shouldShow ? 1 : 0)
            self
        }
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
        case "FirstPartyMap":
            eventTypeName = newValue ? "FirstPartyMapOpen" : "FirstPartyMapClose"
        case let key where key.hasPrefix("ThirdPartyMap"):
            eventTypeName = newValue ? "ThirdPartyMapOpen" : "ThirdPartyMapClose"
        default:
            continue
        }
        
        var data: [String: LogValue] = [:]
        if stateKey.hasPrefix("ThirdPartyMap") {
            let appName = String(stateKey.dropFirst("ThirdPartyMap_".count))
            data["appName"] = .text(appName)
        }
        
        try await dataManager.logEvent(
            eventTypeName: eventTypeName,
            data: data,
            location: nil
        )
    }
}
