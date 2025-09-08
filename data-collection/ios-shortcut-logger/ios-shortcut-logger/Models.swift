import Foundation
import SwiftUI

struct LogEntry: Codable, Identifiable, Equatable {
    let eventTypeName: String
    let timestamp: Date
    
    var id: String { "\(eventTypeName)-\(timestamp.timeIntervalSince1970)" }
    
    init(eventTypeName: String) {
        self.eventTypeName = eventTypeName
        self.timestamp = Date()
    }
    
    init(eventTypeName: String, timestamp: Date) {
        self.eventTypeName = eventTypeName
        self.timestamp = timestamp
    }
}