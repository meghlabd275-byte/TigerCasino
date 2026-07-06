//
//  TigerCasino iOS App
//  High-performance mobile gaming
//

import UIKit

@main
class AppDelegate: UIResponder, UIApplicationDelegate {

    var window: UIWindow?

    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        
        // Initialize core services
        initializeServices()
        
        // Configure appearance
        configureAppearance()
        
        return true
    }

    private func initializeServices() {
        // Initialize Rust security module via JNI
        SecurityModule.initialize()
        
        // Initialize game engine
        GameEngine.shared.initialize()
        
        // Setup WebSocket connection for real-time gaming
        WebSocketService.shared.connect()
    }

    private func configureAppearance() {
        // Dark theme for casino gaming
        UINavigationBar.appearance().barTintColor = UIColor(hex: "#0D0D0D")
        UINavigationBar.appearance().tintColor = UIColor(hex: "#FFD700")
        UINavigationBar.appearance().titleTextAttributes = [.foregroundColor: UIColor.white]
    }

    // MARK: UISceneSession Lifecycle

    func application(_ application: UIApplication, configurationForConnecting connectingSceneSession: UISceneSession, options: UIScene.ConnectionOptions) -> UISceneConfiguration {
        return UISceneConfiguration(name: "Default Configuration", sessionRole: connectingSceneSession.role)
    }

    func application(_ application: UIApplication, didDiscardSceneSessions sceneSessions: Set<UISceneSession>) {
    }
}
