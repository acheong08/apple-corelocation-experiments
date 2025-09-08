import SwiftUI
import AppIntents

@main
struct ios_shortcut_loggerApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
    
    init() {
        AppDependencyManager.shared.add(dependency: ShortcutsProvider())
    }
}
