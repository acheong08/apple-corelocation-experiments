import AppIntents
import CoreLocation
import CoreTelephony
import UIKit
import Foundation

// MARK: - System State Detection

struct SystemStateDetector {
    static func isAirplaneModeEnabled() -> Bool {
        let networkInfo = CTTelephonyNetworkInfo()
        let carrier = networkInfo.serviceSubscriberCellularProviders?.first?.value
        return carrier?.carrierName == nil
    }
    
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
    
    private init() {}
    
    func toggle(_ key: String) -> Bool {
        let currentState = states[key] ?? false
        let newState = !currentState
        states[key] = newState
        return newState
    }
    
    func getCurrentState(_ key: String) -> Bool {
        return states[key] ?? false
    }
    
    func setState(_ key: String, to value: Bool) {
        states[key] = value
    }
    
    func updateSystemStates() async -> [(String, Bool, Bool)] {
        var changes: [(String, Bool, Bool)] = []
        
        // Check airplane mode (detectable via telephony)
        let currentAirplaneMode = SystemStateDetector.isAirplaneModeEnabled()
        let recordedAirplaneMode = states["AirplaneMode"] ?? false
        if currentAirplaneMode != recordedAirplaneMode {
            states["AirplaneMode"] = currentAirplaneMode
            changes.append(("AirplaneMode", recordedAirplaneMode, currentAirplaneMode))
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
        
        return changes
    }
}

// MARK: - Consolidated Event Intents

struct LogFirstPartyMapToggleIntent: AppIntent {
    static let title: LocalizedStringResource = "Toggle First Party Map"
    static let description = IntentDescription("Toggle Apple Maps open/close state")

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        let isOpen = await StateTracker.shared.toggle("FirstPartyMap")
        let eventTypeName = isOpen ? "FirstPartyMapOpen" : "FirstPartyMapClose"
        
        try await logEvent(
            eventTypeName: eventTypeName,
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogThirdPartyMapToggleIntent: AppIntent {
    static let title: LocalizedStringResource = "Toggle Third Party Map"
    static let description = IntentDescription("Toggle third party map app open/close state")

    @Parameter(title: "App Name")
    var appName: String?

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        let stateKey = "ThirdPartyMap_\(appName ?? "Unknown")"
        let isOpen = await StateTracker.shared.toggle(stateKey)
        let eventTypeName = isOpen ? "ThirdPartyMapOpen" : "ThirdPartyMapClose"
        
        var data: [String: LogValue] = [:]
        if let appName = appName {
            data["appName"] = .text(appName)
        }
        try await logEvent(
            eventTypeName: eventTypeName,
            data: data,
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogPlugToggleIntent: AppIntent {
    static let title: LocalizedStringResource = "Toggle Plug State"
    static let description = IntentDescription("Toggle device plug in/out state")

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        let isPluggedIn = await StateTracker.shared.toggle("PlugState")
        let eventTypeName = isPluggedIn ? "PluggedIn" : "PluggedOut"
        
        try await logEvent(
            eventTypeName: eventTypeName,
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogAlarmGoesOffIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Alarm Goes Off"
    static let description = IntentDescription("Log when an alarm goes off")

    @Parameter(title: "Alarm Name")
    var alarmName: String?

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        var data: [String: LogValue] = [:]
        if let alarmName = alarmName {
            data["alarmName"] = .text(alarmName)
        }
        try await logEvent(
            eventTypeName: "AlarmGoesOff",
            data: data,
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogAirplaneModeToggleIntent: AppIntent {
    static let title: LocalizedStringResource = "Toggle Airplane Mode"
    static let description = IntentDescription("Toggle airplane mode on/off state")

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        let isOn = await StateTracker.shared.toggle("AirplaneMode")
        let eventTypeName = isOn ? "AirplaneModeOn" : "AirplaneModeOff"
        
        try await logEvent(
            eventTypeName: eventTypeName,
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogTransactionMadeIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Transaction Made"
    static let description = IntentDescription("Log when a transaction is made")

    @Parameter(title: "Amount")
    var amount: Double?

    @Parameter(title: "Merchant")
    var merchant: String?

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        var data: [String: LogValue] = [:]
        if let amount = amount {
            data["amount"] = .number(amount)
        }
        if let merchant = merchant {
            data["merchant"] = .text(merchant)
        }
        try await logEvent(
            eventTypeName: "TransactionMade",
            data: data,
            includeLocation: includeLocation
        )
        return .result()
    }
}

// MARK: - Helper Functions

private func logEvent(
    eventTypeName: String,
    data: [String: LogValue] = [:],
    includeLocation: Bool = false
) async throws {
    let dataManager = await DataManager.shared

    // Get location if requested
    var locationData: LocationData? = nil
    if includeLocation {
        locationData = await getCurrentLocation()
    }

    // Log the event
    try await dataManager.logEvent(
        eventTypeName: eventTypeName,
        data: data,
        location: locationData
    )
}

private func getCurrentLocation() async -> LocationData? {
    let locationManager = CLLocationManager()

    let authStatus = locationManager.authorizationStatus
    if authStatus == .notDetermined {
        locationManager.requestWhenInUseAuthorization()
    }

    guard authStatus == .authorizedWhenInUse || authStatus == .authorizedAlways
    else {
        return nil
    }

    do {
        let location = try await withCheckedThrowingContinuation {
            continuation in
            let delegate = LocationDelegate(continuation: continuation)
            locationManager.delegate = delegate
            locationManager.requestLocation()
        }

        return LocationData(
            latitude: location.coordinate.latitude,
            longitude: location.coordinate.longitude,
            altitude: location.altitude,
            accuracy: location.horizontalAccuracy,
            timestamp: location.timestamp
        )
    } catch {
        return nil
    }
}

// MARK: - Supporting Types

enum LogEventError: Error, LocalizedError {
    case eventTypeNotFound(String)
    case invalidDataFormat
    case loggingFailed
    case locationPermissionDenied

    var errorDescription: String? {
        switch self {
        case .eventTypeNotFound(let name):
            return "Event type '\(name)' not found"
        case .invalidDataFormat:
            return "Invalid JSON data format"
        case .loggingFailed:
            return "Failed to log event"
        case .locationPermissionDenied:
            return "Location permission denied"
        }
    }
}

class LocationDelegate: NSObject, CLLocationManagerDelegate {
    private let continuation: CheckedContinuation<CLLocation, Error>

    init(continuation: CheckedContinuation<CLLocation, Error>) {
        self.continuation = continuation
    }

    func locationManager(
        _ manager: CLLocationManager,
        didUpdateLocations locations: [CLLocation]
    ) {
        if let location = locations.first {
            continuation.resume(returning: location)
        }
    }

    func locationManager(
        _ manager: CLLocationManager,
        didFailWithError error: Error
    ) {
        continuation.resume(throwing: error)
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


