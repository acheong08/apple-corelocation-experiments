import Foundation
import CoreLocation

@MainActor
class DataManager: ObservableObject {
    @Published var recentLogs: [LogEntry] = []
    
    private let documentsDirectory: URL
    private let logFileURL: URL
    
    init() {
        documentsDirectory = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask).first!
        logFileURL = documentsDirectory.appendingPathComponent("events.jsonl")
        
        // Create logs directory if it doesn't exist
        let logsDirectory = documentsDirectory.appendingPathComponent("logs")
        try? FileManager.default.createDirectory(at: logsDirectory, withIntermediateDirectories: true)
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
        
        try writeToLogFile(logEntry)
    }
    
    private func writeToLogFile(_ logEntry: LogEntry) throws {
        // Convert log entry to JSON
        let jsonData = try createJSONLEntry(from: logEntry)
        let jsonString = String(data: jsonData, encoding: .utf8)! + "\n"
        
        // Append to file
        if FileManager.default.fileExists(atPath: logFileURL.path) {
            let fileHandle = try FileHandle(forWritingTo: logFileURL)
            fileHandle.seekToEndOfFile()
            fileHandle.write(jsonString.data(using: .utf8)!)
            fileHandle.closeFile()
        } else {
            try jsonString.write(to: logFileURL, atomically: true, encoding: .utf8)
        }
    }
    
    private func createJSONLEntry(from logEntry: LogEntry) throws -> Data {
        let jsonObject: [String: Any] = [
            "eventTypeName": logEntry.eventTypeName,
            "timestamp": ISO8601DateFormatter().string(from: logEntry.timestamp)
        ]
        
        return try JSONSerialization.data(withJSONObject: jsonObject)
    }
    
    // MARK: - Export Functionality
    
    func getLogContent() -> String {
        guard FileManager.default.fileExists(atPath: logFileURL.path),
              let content = try? String(contentsOf: logFileURL, encoding: .utf8) else {
            return "No log data available"
        }
        return content
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