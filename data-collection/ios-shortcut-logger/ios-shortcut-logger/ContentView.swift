import SwiftUI

struct ContentView: View {
    @StateObject private var dataManager = DataManager.shared
    @State private var selectedTab = 0
    
    var body: some View {
        TabView(selection: $selectedTab) {
            OutputLocationsView()
                .tabItem {
                    Image(systemName: "folder")
                    Text("Output")
                }
                .tag(0)
            
            RecentLogsView()
                .tabItem {
                    Image(systemName: "clock")
                    Text("Recent Logs")
                }
                .tag(1)
        }
        .environmentObject(dataManager)
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
