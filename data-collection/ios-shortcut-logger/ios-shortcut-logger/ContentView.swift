import SwiftUI

struct ContentView: View {
    @StateObject private var dataManager = DataManager.shared
    @State private var selectedTab = 0
    
    var body: some View {
        TabView(selection: $selectedTab) {
            EventTypesView()
                .tabItem {
                    Image(systemName: "list.bullet")
                    Text("Event Types")
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
    }
}

struct EventTypesView: View {
    @EnvironmentObject var dataManager: DataManager
    @State private var showingAddEventType = false
    
    var body: some View {
        NavigationView {
            List {
                ForEach(dataManager.eventTypes) { eventType in
                    EventTypeRow(eventType: eventType)
                }
                .onDelete(perform: deleteEventTypes)
            }
            .navigationTitle("Event Types")
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Add") {
                        showingAddEventType = true
                    }
                }
            }
            .sheet(isPresented: $showingAddEventType) {
                AddEventTypeView()
            }
        }
    }
    
    private func deleteEventTypes(offsets: IndexSet) {
        for index in offsets {
            dataManager.removeEventType(dataManager.eventTypes[index])
        }
    }
}

struct EventTypeRow: View {
    let eventType: EventType
    
    var body: some View {
        VStack(alignment: .leading, spacing: 4) {
            Text(eventType.name)
                .font(.headline)
            Text(eventType.description)
                .font(.caption)
                .foregroundColor(.secondary)
            
            if !eventType.fields.isEmpty {
                Text("Fields: \(eventType.fields.map { $0.name }.joined(separator: ", "))")
                    .font(.caption2)
                    .foregroundColor(.secondary)
            }
        }
        .padding(.vertical, 2)
    }
}

struct AddEventTypeView: View {
    @EnvironmentObject var dataManager: DataManager
    @Environment(\.dismiss) var dismiss
    
    @State private var name = ""
    @State private var description = ""
    @State private var fields: [EventField] = []
    @State private var showingAddField = false
    
    var body: some View {
        NavigationView {
            Form {
                Section("Event Type Details") {
                    TextField("Name", text: $name)
                    TextField("Description", text: $description)
                }
                
                Section("Fields") {
                    ForEach(fields) { field in
                        HStack {
                            VStack(alignment: .leading) {
                                Text(field.name)
                                    .font(.headline)
                                Text("\(field.type.rawValue)\(field.isRequired ? " (required)" : "")")
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            Spacer()
                        }
                    }
                    .onDelete(perform: deleteFields)
                    
                    Button("Add Field") {
                        showingAddField = true
                    }
                }
            }
            .navigationTitle("New Event Type")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Save") {
                        saveEventType()
                    }
                    .disabled(name.isEmpty)
                }
            }
            .sheet(isPresented: $showingAddField) {
                AddFieldView { field in
                    fields.append(field)
                }
            }
        }
    }
    
    private func deleteFields(offsets: IndexSet) {
        fields.remove(atOffsets: offsets)
    }
    
    private func saveEventType() {
        let eventType = EventType(
            name: name,
            description: description,
            fields: fields
        )
        dataManager.addEventType(eventType)
        dismiss()
    }
}

struct AddFieldView: View {
    @Environment(\.dismiss) var dismiss
    let onSave: (EventField) -> Void
    
    @State private var name = ""
    @State private var type = EventField.FieldType.text
    @State private var isRequired = false
    
    var body: some View {
        NavigationView {
            Form {
                TextField("Field Name", text: $name)
                
                Picker("Type", selection: $type) {
                    ForEach(EventField.FieldType.allCases, id: \.self) { fieldType in
                        Text(fieldType.rawValue.capitalized).tag(fieldType)
                    }
                }
                
                Toggle("Required", isOn: $isRequired)
            }
            .navigationTitle("New Field")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Save") {
                        let field = EventField(
                            name: name,
                            type: type,
                            isRequired: isRequired
                        )
                        onSave(field)
                        dismiss()
                    }
                    .disabled(name.isEmpty)
                }
            }
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
