import Foundation
import CoreLocation

@MainActor
class DataManager: ObservableObject {
    @Published var outputLocations: [OutputLocation] = []
    @Published var recentLogs: [LogEntry] = []
    
    private let documentsDirectory: URL
    private let outputLocationsFile: URL
    
    init() {
        documentsDirectory = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask).first!
        outputLocationsFile = documentsDirectory.appendingPathComponent("outputLocations.json")
        
        loadData()
        createDefaultOutputLocation()
    }
    
    // MARK: - Output Locations Management
    
    func addOutputLocation(_ location: OutputLocation) {
        outputLocations.append(location)
        saveOutputLocations()
    }
    
    func removeOutputLocation(_ location: OutputLocation) {
        outputLocations.removeAll { $0.id == location.id }
        saveOutputLocations()
    }
    
    func updateOutputLocation(_ location: OutputLocation) {
        if let index = outputLocations.firstIndex(where: { $0.id == location.id }) {
            outputLocations[index] = location
            saveOutputLocations()
        }
    }
    
    // MARK: - Logging
    
    func logEvent(eventTypeName: String, data: [String: LogValue], location: LocationData? = nil) throws {
        let logEntry = LogEntry(
            eventTypeName: eventTypeName,
            data: data,
            location: location
        )
        
        recentLogs.insert(logEntry, at: 0)
        if recentLogs.count > 100 {
            recentLogs.removeLast()
        }
        
        try writeToJSONLFiles(logEntry)
    }
    
    private func writeToJSONLFiles(_ logEntry: LogEntry) throws {
        let activeLocations = outputLocations.filter { $0.isActive }
        
        for outputLocation in activeLocations {
            let fileURL = documentsDirectory.appendingPathComponent(outputLocation.filePath)
            
            // Create directory if it doesn't exist
            let directory = fileURL.deletingLastPathComponent()
            try FileManager.default.createDirectory(at: directory, withIntermediateDirectories: true)
            
            // Convert log entry to JSON
            let jsonData = try createJSONLEntry(from: logEntry)
            let jsonString = String(data: jsonData, encoding: .utf8)! + "\n"
            
            // Append to file
            if FileManager.default.fileExists(atPath: fileURL.path) {
                let fileHandle = try FileHandle(forWritingTo: fileURL)
                fileHandle.seekToEndOfFile()
                fileHandle.write(jsonString.data(using: .utf8)!)
                fileHandle.closeFile()
            } else {
                try jsonString.write(to: fileURL, atomically: true, encoding: .utf8)
            }
        }
    }
    
    private func createJSONLEntry(from logEntry: LogEntry) throws -> Data {
        var jsonObject: [String: Any] = [
            "id": logEntry.id.uuidString,
            "eventTypeName": logEntry.eventTypeName,
            "timestamp": ISO8601DateFormatter().string(from: logEntry.timestamp)
        ]
        
        // Add data fields
        var dataObject: [String: Any] = [:]
        for (key, value) in logEntry.data {
            dataObject[key] = value.toJSONValue()
        }
        jsonObject["data"] = dataObject
        
        // Add location if available
        if let location = logEntry.location {
            jsonObject["location"] = [
                "latitude": location.latitude,
                "longitude": location.longitude,
                "altitude": location.altitude as Any,
                "accuracy": location.accuracy,
                "timestamp": ISO8601DateFormatter().string(from: location.timestamp)
            ]
        }
        
        return try JSONSerialization.data(withJSONObject: jsonObject)
    }
    
    // MARK: - Export Functionality
    
    func exportLogs() -> [URL] {
        let activeLocations = outputLocations.filter { $0.isActive }
        var exportURLs: [URL] = []
        
        for outputLocation in activeLocations {
            let fileURL = documentsDirectory.appendingPathComponent(outputLocation.filePath)
            if FileManager.default.fileExists(atPath: fileURL.path) {
                exportURLs.append(fileURL)
            }
        }
        
        return exportURLs
    }
    
    func getAllLogFiles() -> [URL] {
        var logFiles: [URL] = []
        let tempDirectory = FileManager.default.temporaryDirectory
        
        for outputLocation in outputLocations {
            let originalFileURL = documentsDirectory.appendingPathComponent(outputLocation.filePath)
            
            if FileManager.default.fileExists(atPath: originalFileURL.path) {
                do {
                    // Create a temporary copy with a safe filename
                    let fileName = outputLocation.name.replacingOccurrences(of: " ", with: "_") + "_" + originalFileURL.lastPathComponent
                    let tempFileURL = tempDirectory.appendingPathComponent(fileName)
                    
                    // Remove existing temp file if it exists
                    if FileManager.default.fileExists(atPath: tempFileURL.path) {
                        try FileManager.default.removeItem(at: tempFileURL)
                    }
                    
                    // Copy to temp directory
                    try FileManager.default.copyItem(at: originalFileURL, to: tempFileURL)
                    logFiles.append(tempFileURL)
                } catch {
                    print("Error creating temp file for export: \(error)")
                }
            }
        }
        
        return logFiles
    }

    // MARK: - Data Persistence
    
    private func loadData() {
        loadOutputLocations()
    }
    
    private func loadOutputLocations() {
        guard FileManager.default.fileExists(atPath: outputLocationsFile.path),
              let data = try? Data(contentsOf: outputLocationsFile),
              let locations = try? JSONDecoder().decode([OutputLocation].self, from: data) else {
            return
        }
        outputLocations = locations
    }
    
    private func saveOutputLocations() {
        guard let data = try? JSONEncoder().encode(outputLocations) else { return }
        try? data.write(to: outputLocationsFile)
    }
    
    private func createDefaultOutputLocation() {
        if outputLocations.isEmpty {
            let defaultLocation = OutputLocation(
                name: "Default Log",
                filePath: "logs/events.jsonl"
            )
            addOutputLocation(defaultLocation)
        }
    }
}

enum LoggingError: Error {
    case eventTypeNotFound
    case invalidData
    case fileWriteError
}

// MARK: - Shared Instance

extension DataManager {
    static let shared = DataManager()
}