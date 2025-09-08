import Foundation
import SwiftUI

struct LogEntry: Codable, Identifiable {
    let id: UUID
    let eventTypeName: String
    let timestamp: Date
    let data: [String: LogValue]
    let location: LocationData?
    
    init(eventTypeName: String, data: [String: LogValue], location: LocationData? = nil) {
        self.id = UUID()
        self.eventTypeName = eventTypeName
        self.timestamp = Date()
        self.data = data
        self.location = location
    }
}

enum LogValue: Codable {
    case text(String)
    case number(Double)
    case boolean(Bool)
    case date(Date)
    
    func toJSONValue() -> Any {
        switch self {
        case .text(let string):
            return string
        case .number(let double):
            return double
        case .boolean(let bool):
            return bool
        case .date(let date):
            return ISO8601DateFormatter().string(from: date)
        }
    }
}

struct LocationData: Codable {
    let latitude: Double
    let longitude: Double
    let altitude: Double?
    let accuracy: Double
    let timestamp: Date
}