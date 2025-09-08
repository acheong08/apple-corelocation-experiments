import Foundation
import SwiftUI

struct EventType: Codable, Identifiable, Hashable {
    let id: UUID
    let name: String
    let description: String
    let fields: [EventField]
    let createdAt: Date
    
    init(name: String, description: String, fields: [EventField] = []) {
        self.id = UUID()
        self.name = name
        self.description = description
        self.fields = fields
        self.createdAt = Date()
    }
}

struct EventField: Codable, Identifiable, Hashable {
    let id: UUID
    let name: String
    let type: FieldType
    let isRequired: Bool
    
    init(name: String, type: FieldType, isRequired: Bool = false) {
        self.id = UUID()
        self.name = name
        self.type = type
        self.isRequired = isRequired
    }
    
    enum FieldType: String, Codable, CaseIterable {
        case text = "text"
        case number = "number"
        case boolean = "boolean"
        case date = "date"
    }
}

struct LogEntry: Codable, Identifiable {
    let id: UUID
    let eventTypeId: UUID
    let eventTypeName: String
    let timestamp: Date
    let data: [String: LogValue]
    let location: LocationData?
    
    init(eventTypeId: UUID, eventTypeName: String, data: [String: LogValue], location: LocationData? = nil) {
        self.id = UUID()
        self.eventTypeId = eventTypeId
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

struct OutputLocation: Codable, Identifiable {
    let id: UUID
    let name: String
    let filePath: String
    let isActive: Bool
    let createdAt: Date
    
    init(name: String, filePath: String, isActive: Bool = true) {
        self.id = UUID()
        self.name = name
        self.filePath = filePath
        self.isActive = isActive
        self.createdAt = Date()
    }
}