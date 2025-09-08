import AppIntents
import CoreLocation
import Foundation

// MARK: - Hardcoded Event Intents

struct LogFirstPartyMapOpenIntent: AppIntent {
    static let title: LocalizedStringResource = "Log First Party Map Open"
    static let description = IntentDescription("Log when Apple Maps opens")

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        try await logEvent(
            eventTypeName: "FirstPartyMapOpen",
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogFirstPartyMapCloseIntent: AppIntent {
    static let title: LocalizedStringResource = "Log First Party Map Close"
    static let description = IntentDescription("Log when Apple Maps closes")

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        try await logEvent(
            eventTypeName: "FirstPartyMapClose",
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogThirdPartyMapOpenIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Third Party Map Open"
    static let description = IntentDescription(
        "Log when a third party map app opens"
    )

    @Parameter(title: "App Name")
    var appName: String?

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        var data: [String: LogValue] = [:]
        if let appName = appName {
            data["appName"] = .text(appName)
        }
        try await logEvent(
            eventTypeName: "ThirdPartyMapOpen",
            data: data,
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogThirdPartyMapCloseIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Third Party Map Close"
    static let description = IntentDescription(
        "Log when a third party map app closes"
    )

    @Parameter(title: "App Name")
    var appName: String?

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        var data: [String: LogValue] = [:]
        if let appName = appName {
            data["appName"] = .text(appName)
        }
        try await logEvent(
            eventTypeName: "ThirdPartyMapClose",
            data: data,
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogPluggedInIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Plugged In"
    static let description = IntentDescription(
        "Log when device is plugged in to charge"
    )

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        try await logEvent(
            eventTypeName: "PluggedIn",
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogPluggedOutIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Plugged Out"
    static let description = IntentDescription(
        "Log when device is unplugged from charging"
    )

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        try await logEvent(
            eventTypeName: "PluggedOut",
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

struct LogAirplaneModeOnIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Airplane Mode On"
    static let description = IntentDescription(
        "Log when airplane mode is turned on"
    )

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        try await logEvent(
            eventTypeName: "AirplaneModeOn",
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogAirplaneModeOffIntent: AppIntent {
    static let title: LocalizedStringResource = "Log Airplane Mode Off"
    static let description = IntentDescription(
        "Log when airplane mode is turned off"
    )

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        try await logEvent(
            eventTypeName: "AirplaneModeOff",
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogAppOpenIntent: AppIntent {
    static let title: LocalizedStringResource = "Log App Open"
    static let description = IntentDescription(
        "Log when an installed app is opened"
    )

    @Parameter(title: "App Name")
    var appName: String

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        let data: [String: LogValue] = ["appName": .text(appName)]
        try await logEvent(
            eventTypeName: "AppOpen",
            data: data,
            includeLocation: includeLocation
        )
        return .result()
    }
}

struct LogAppCloseIntent: AppIntent {
    static let title: LocalizedStringResource = "Log App Close"
    static let description = IntentDescription(
        "Log when an installed app is closed"
    )

    @Parameter(title: "App Name")
    var appName: String

    @Parameter(title: "Include Location", default: false)
    var includeLocation: Bool

    func perform() async throws -> some IntentResult {
        let data: [String: LogValue] = ["appName": .text(appName)]
        try await logEvent(
            eventTypeName: "AppClose",
            data: data,
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

struct AppShortcutsProvider: AppShortcutsProvider {
    static var appShortcuts: [AppShortcut] {
        AppShortcut(
            intent: LogFirstPartyMapOpenIntent(),
            phrases: [
                "Log first party map open with \(.applicationName)",
                "Track Apple Maps opening in \(.applicationName)",
                "Record when Maps app opens using \(.applicationName)"
            ],
            shortTitle: "Log Map Open",
            systemImageName: "map"
        )
        
        AppShortcut(
            intent: LogFirstPartyMapCloseIntent(),
            phrases: [
                "Log first party map close with \(.applicationName)",
                "Track Apple Maps closing in \(.applicationName)",
                "Record when Maps app closes using \(.applicationName)"
            ],
            shortTitle: "Log Map Close",
            systemImageName: "map.fill"
        )
        
        AppShortcut(
            intent: LogThirdPartyMapOpenIntent(),
            phrases: [
                "Log third party map open with \(.applicationName)",
                "Track third party map opening in \(.applicationName)",
                "Record third party map app opening using \(.applicationName)"
            ],
            shortTitle: "Log 3rd Party Map Open",
            systemImageName: "location"
        )
        
        AppShortcut(
            intent: LogThirdPartyMapCloseIntent(),
            phrases: [
                "Log third party map close with \(.applicationName)",
                "Track third party map closing in \(.applicationName)",
                "Record third party map app closing using \(.applicationName)"
            ],
            shortTitle: "Log 3rd Party Map Close",
            systemImageName: "location.fill"
        )
        
        AppShortcut(
            intent: LogPluggedInIntent(),
            phrases: [
                "Log plugged in with \(.applicationName)",
                "Track device charging in \(.applicationName)",
                "Record when device is plugged in using \(.applicationName)"
            ],
            shortTitle: "Log Plugged In",
            systemImageName: "battery.100.bolt"
        )
        
        AppShortcut(
            intent: LogPluggedOutIntent(),
            phrases: [
                "Log plugged out with \(.applicationName)",
                "Track device unplugging in \(.applicationName)",
                "Record when device is unplugged using \(.applicationName)"
            ],
            shortTitle: "Log Plugged Out",
            systemImageName: "battery.100"
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
            intent: LogAirplaneModeOnIntent(),
            phrases: [
                "Log airplane mode on with \(.applicationName)",
                "Track airplane mode activation in \(.applicationName)",
                "Record airplane mode turning on using \(.applicationName)"
            ],
            shortTitle: "Log Airplane On",
            systemImageName: "airplane"
        )
        
        AppShortcut(
            intent: LogAirplaneModeOffIntent(),
            phrases: [
                "Log airplane mode off with \(.applicationName)",
                "Track airplane mode deactivation in \(.applicationName)",
                "Record airplane mode turning off using \(.applicationName)"
            ],
            shortTitle: "Log Airplane Off",
            systemImageName: "airplane.departure"
        )
        
        AppShortcut(
            intent: LogAppOpenIntent(),
            phrases: [
                "Log app open with \(.applicationName)",
                "Track app opening in \(.applicationName)",
                "Record when app launches using \(.applicationName)"
            ],
            shortTitle: "Log App Open",
            systemImageName: "apps.iphone"
        )
        
        AppShortcut(
            intent: LogAppCloseIntent(),
            phrases: [
                "Log app close with \(.applicationName)",
                "Track app closing in \(.applicationName)",
                "Record when app closes using \(.applicationName)"
            ],
            shortTitle: "Log App Close",
            systemImageName: "apps.iphone.badge.minus"
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


