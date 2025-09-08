import AppIntents
import CoreLocation
import UIKit
import Foundation

// MARK: - System State Detection

struct SystemStateDetector {
    static func isDevicePluggedIn() -> Bool {
        UIDevice.current.isBatteryMonitoringEnabled = true
        defer { UIDevice.current.isBatteryMonitoringEnabled = false }
        
        let batteryState = UIDevice.current.batteryState
        return batteryState == .charging || batteryState == .full
    }
}

// MARK: - State Tracker

actor StateTracker {
    static let shared = StateTracker()
    
    private var states: [String: Bool] = [:]
    private let userDefaults = UserDefaults.standard
    private let statesKey = "StateTrackerStates"
    
    private init() {
        // Load states synchronously in init
        if let savedStates = UserDefaults.standard.object(forKey: "StateTrackerStates") as? [String: Bool] {
            states = savedStates
        }
    }
    
    private func saveStates() {
        userDefaults.set(states, forKey: statesKey)
    }
    
    func toggle(_ key: String) -> Bool {
        let currentState = states[key] ?? false
        let newState = !currentState
        states[key] = newState
        saveStates()
        return newState
    }
    
    func getCurrentState(_ key: String) -> Bool {
        return states[key] ?? false
    }
    
    func setState(_ key: String, to value: Bool) {
        states[key] = value
        saveStates()
    }
    
    func clearAllStates() {
        states.removeAll()
        saveStates()
    }
    
    func updateSystemStates(airplaneModeOverride: Bool? = nil) async -> [(String, Bool, Bool)] {
        var changes: [(String, Bool, Bool)] = []
        
        // Handle airplane mode with manual override (no reliable detection available)
        if let airplaneModeState = airplaneModeOverride {
            let recordedAirplaneMode = states["AirplaneMode"] ?? false
            if airplaneModeState != recordedAirplaneMode {
                states["AirplaneMode"] = airplaneModeState
                changes.append(("AirplaneMode", recordedAirplaneMode, airplaneModeState))
            }
        }
        
        // Check plug state (detectable via battery state)
        let currentPlugState = SystemStateDetector.isDevicePluggedIn()
        let recordedPlugState = states["PlugState"] ?? false
        if currentPlugState != recordedPlugState {
            states["PlugState"] = currentPlugState
            changes.append(("PlugState", recordedPlugState, currentPlugState))
        }
        
        // Reset app states to false (default to closed, not detectable)
        let appStates = ["FirstPartyMap"]
        for appState in appStates {
            if states[appState] != false {
                states[appState] = false
                changes.append((appState, true, false))
            }
        }
        
        // Reset any third-party map states to false
        let thirdPartyKeys = states.keys.filter { $0.hasPrefix("ThirdPartyMap") }
        for key in thirdPartyKeys {
            if states[key] != false {
                states[key] = false
                changes.append((key, true, false))
            }
        }
        
        // Save states if there were any changes
        if !changes.isEmpty {
            saveStates()
        }
        
        return changes
    }
}

// MARK: - Consolidated Event Intents

struct LogFirstPartyMapToggleIntent: AppIntent {
    static let title: LocalizedStringResource = "Toggle First Party Map"
    static let description = IntentDescription("Toggle Apple Maps open/close state")

    func perform() async throws -> some IntentResult {
        let isOpen = await StateTracker.shared.toggle("FirstPartyMap")
        let eventTypeName = isOpen ? "FirstPartyMapOpen" : "FirstPartyMapClose"
        
        try await logEvent(eventTypeName: eventTypeName)
        return .result()
    }
}

struct LogThirdPartyMapToggleIntent: AppIntent {
    static let title: LocalizedStringResource = "Toggle Third Party Map"
    static let description = IntentDescription("Toggle third party map app open/close state")

    func perform() async throws -> some IntentResult {
        let isOpen = await StateTracker.shared.toggle("ThirdPartyMap")
        let eventTypeName = isOpen ? "ThirdPartyMapOpen" : "ThirdPartyMapClose"
        
        try await logEvent(eventTypeName: eventTypeName)
        return .result()
    }
}

struct LogPlugToggleIntent: AppIntent {
    static let title: LocalizedStringResource = "Toggle Plug State"
    static let description = IntentDescription("Toggle device plug in/out state")

    func perform() async throws -> some IntentResult {
        let isPluggedIn = await StateTracker.shared.toggle("PlugState")
        let eventTypeName = isPluggedIn ? "PluggedIn" : "PluggedOut"
        
        try await logEvent(eventTypeName: eventTypeName)
        return .result()
    }
}

struct LogAlarmGoesOffIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Alarm Goes Off"
    static let description = IntentDescription("Log when an alarm goes off")

    func perform() async throws -> some IntentResult {
        try await logEvent(eventTypeName: "AlarmGoesOff")
        return .result()
    }
}

struct LogAirplaneModeToggleIntent: AppIntent {
    static let title: LocalizedStringResource = "Toggle Airplane Mode"
    static let description = IntentDescription("Toggle airplane mode on/off state")

    func perform() async throws -> some IntentResult {
        let isOn = await StateTracker.shared.toggle("AirplaneMode")
        let eventTypeName = isOn ? "AirplaneModeOn" : "AirplaneModeOff"
        
        try await logEvent(eventTypeName: eventTypeName)
        return .result()
    }
}

struct LogTransactionMadeIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Transaction Made"
    static let description = IntentDescription("Log when a transaction is made")

    func perform() async throws -> some IntentResult {
        try await logEvent(eventTypeName: "TransactionMade")
        return .result()
    }
}

// MARK: - Helper Functions

private func logEvent(eventTypeName: String) async throws {
    let dataManager = await DataManager.shared
    try await dataManager.logEvent(eventTypeName: eventTypeName)
}

// MARK: - Supporting Types

enum LogEventError: Error, LocalizedError {
    case eventTypeNotFound(String)
    case loggingFailed

    var errorDescription: String? {
        switch self {
        case .eventTypeNotFound(let name):
            return "Event type '\(name)' not found"
        case .loggingFailed:
            return "Failed to log event"
        }
    }
}

// MARK: - App Shortcuts Provider

struct ShortcutsProvider: AppShortcutsProvider {
    static var appShortcuts: [AppShortcut] {
        AppShortcut(
            intent: LogFirstPartyMapToggleIntent(),
            phrases: [
                "Toggle first party map with \(.applicationName)",
                "Toggle Apple Maps in \(.applicationName)",
                "Switch Maps app state using \(.applicationName)"
            ],
            shortTitle: "Toggle Map",
            systemImageName: "map"
        )
        
        AppShortcut(
            intent: LogThirdPartyMapToggleIntent(),
            phrases: [
                "Toggle third party map with \(.applicationName)",
                "Toggle third party map app in \(.applicationName)",
                "Switch third party map state using \(.applicationName)"
            ],
            shortTitle: "Toggle 3rd Party Map",
            systemImageName: "location"
        )
        
        AppShortcut(
            intent: LogPlugToggleIntent(),
            phrases: [
                "Toggle plug state with \(.applicationName)",
                "Toggle device charging in \(.applicationName)",
                "Switch plug state using \(.applicationName)"
            ],
            shortTitle: "Toggle Plug",
            systemImageName: "battery.100.bolt"
        )
        
        AppShortcut(
            intent: LogAlarmGoesOffIntent(),
            phrases: [
                "Log alarm goes off with \(.applicationName)",
                "Track alarm activation in \(.applicationName)",
                "Record when alarm sounds using \(.applicationName)"
            ],
            shortTitle: "Log Alarm",
            systemImageName: "alarm"
        )
        
        AppShortcut(
            intent: LogAirplaneModeToggleIntent(),
            phrases: [
                "Toggle airplane mode with \(.applicationName)",
                "Toggle airplane mode in \(.applicationName)",
                "Switch airplane mode using \(.applicationName)"
            ],
            shortTitle: "Toggle Airplane Mode",
            systemImageName: "airplane"
        )
        
        AppShortcut(
            intent: LogTransactionMadeIntent(),
            phrases: [
                "Log transaction made with \(.applicationName)",
                "Track payment in \(.applicationName)",
                "Record transaction using \(.applicationName)"
            ],
            shortTitle: "Log Transaction",
            systemImageName: "creditcard"
        )
    }
}


