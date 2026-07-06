//
//  GameViewController.swift
//  TigerCasino iOS
//

import UIKit
import WebKit

class GameViewController: UIViewController {

    private var webView: WKWebView!
    private var loadingIndicator: UIActivityIndicatorView!
    private var gameId: String = ""
    
    // Game state
    private var currentMultiplier: Double = 1.0
    private var isGameRunning: Bool = false
    private var balance: Double = 0.0

    override func viewDidLoad() {
        super.viewDidLoad()
        setupUI()
        loadGame()
    }

    private func setupUI() {
        view.backgroundColor = UIColor(hex: "#0D0D0D")
        
        // WebView for game rendering
        let config = WKWebViewConfiguration()
        config.allowsInlineMediaPlayback = true
        config.mediaTypesRequiringUserActionForPlayback = []
        
        webView = WKWebView(frame: view.bounds, configuration: config)
        webView.navigationDelegate = self
        webView.scrollView.bounces = false
        view.addSubview(webView)
        
        // Loading indicator
        loadingIndicator = UIActivityIndicatorView(style: .large)
        loadingIndicator.color = UIColor(hex: "#FFD700")
        loadingIndicator.center = view.center
        loadingIndicator.startAnimating()
        view.addSubview(loadingIndicator)
        
        // Bottom controls
        setupControls()
    }

    private func setupControls() {
        let controlPanel = UIView()
        controlPanel.backgroundColor = UIColor(hex: "#1A1A1A")
        controlPanel.translatesAutoresizingMaskIntoConstraints = false
        view.addSubview(controlPanel)
        
        NSLayoutConstraint.activate([
            controlPanel.leadingAnchor.constraint(equalTo: view.leadingAnchor),
            controlPanel.trailingAnchor.constraint(equalTo: view.trailingAnchor),
            controlPanel.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor),
            controlPanel.heightAnchor.constraint(equalToConstant: 120)
        ])
        
        // Balance label
        let balanceLabel = UILabel()
        balanceLabel.text = "Balance: $\(String(format: "%.2f", balance))"
        balanceLabel.textColor = .white
        balanceLabel.font = UIFont.boldSystemFont(ofSize: 18)
        balanceLabel.translatesAutoresizingMaskIntoConstraints = false
        controlPanel.addSubview(balanceLabel)
        
        // Bet input
        let betField = UITextField()
        betField.placeholder = "Bet Amount"
        betField.keyboardType = .decimalPad
        betField.backgroundColor = UIColor(hex: "#2A2A2A")
        betField.textColor = .white
        betField.layer.cornerRadius = 8
        betField.translatesAutoresizingMaskIntoConstraints = false
        controlPanel.addSubview(betField)
        
        // Bet button
        let betButton = UIButton(type: .system)
        betButton.setTitle("BET", for: .normal)
        betButton.backgroundColor = UIColor(hex: "#FFD700")
        betButton.setTitleColor(.black, for: .normal)
        betButton.titleLabel?.font = UIFont.boldSystemFont(ofSize: 16)
        betButton.layer.cornerRadius = 8
        betButton.translatesAutoresizingMaskIntoConstraints = false
        betButton.addTarget(self, action: #selector(placeBet), for: .touchUpInside)
        controlPanel.addSubview(betButton)
        
        // Cashout button
        let cashoutButton = UIButton(type: .system)
        cashoutButton.setTitle("CASHOUT", for: .normal)
        cashoutButton.backgroundColor = UIColor(hex: "#00FF00")
        cashoutButton.setTitleColor(.black, for: .normal)
        cashoutButton.titleLabel?.font = UIFont.boldSystemFont(ofSize: 16)
        cashoutButton.layer.cornerRadius = 8
        cashoutButton.translatesAutoresizingMaskIntoConstraints = false
        cashoutButton.addTarget(self, action: #selector(cashOut), for: .touchUpInside)
        controlPanel.addSubview(cashoutButton)
        
        NSLayoutConstraint.activate([
            balanceLabel.topAnchor.constraint(equalTo: controlPanel.topAnchor, constant: 16),
            balanceLabel.leadingAnchor.constraint(equalTo: controlPanel.leadingAnchor, constant: 16),
            
            betField.topAnchor.constraint(equalTo: balanceLabel.bottomAnchor, constant: 12),
            betField.leadingAnchor.constraint(equalTo: controlPanel.leadingAnchor, constant: 16),
            betField.widthAnchor.constraint(equalToConstant: 100),
            betField.heightAnchor.constraint(equalToConstant: 40),
            
            betButton.centerYAnchor.constraint(equalTo: betField.centerYAnchor),
            betButton.leadingAnchor.constraint(equalTo: betField.trailingAnchor, constant: 12),
            betButton.widthAnchor.constraint(equalToConstant: 80),
            betButton.heightAnchor.constraint(equalToConstant: 40),
            
            cashoutButton.centerYAnchor.constraint(equalTo: betField.centerYAnchor),
            cashoutButton.leadingAnchor.constraint(equalTo: betButton.trailingAnchor, constant: 12),
            cashoutButton.trailingAnchor.constraint(equalTo: controlPanel.trailingAnchor, constant: -16),
            cashoutButton.heightAnchor.constraint(equalToConstant: 40),
        ])
    }

    private func loadGame() {
        // Load game from WebView
        let gameURL = URL(string: "https://api.tigercasino.com/game/\(gameId)")!
        webView.load(URLRequest(url: gameURL))
    }

    @objc private func placeBet() {
        // Call Rust game engine for instant processing
        let result = GameEngine.shared.placeBet(
            gameId: gameId,
            amount: 10.0,
            userId: getCurrentUserId()
        )
        
        if result.success {
            isGameRunning = true
            startGameLoop()
        }
    }

    @objc private func cashOut() {
        guard isGameRunning else { return }
        
        let payout = currentMultiplier * 10.0
        let result = GameEngine.shared.cashOut(
            gameId: gameId,
            payout: payout
        )
        
        if result.success {
            balance += payout
            isGameRunning = false
            showPayoutAlert(amount: payout)
        }
    }

    private func startGameLoop() {
        // High-speed game loop at 60fps
        Timer.scheduledTimer(withTimeInterval: 1.0/60.0, repeats: true) { [weak self] timer in
            guard let self = self, self.isGameRunning else {
                timer.invalidate()
                return
            }
            
            // Update multiplier
            self.currentMultiplier = GameEngine.shared.getCurrentMultiplier(gameId: self.gameId)
            self.updateUI()
        }
    }

    private func updateUI() {
        // Update multiplier display
    }

    private func showPayoutAlert(amount: Double) {
        let alert = UIAlertController(
            title: "You Won!",
            message: "Payout: $\(String(format: "%.2f", amount))",
            preferredStyle: .alert
        )
        alert.addAction(UIAlertAction(title: "Continue", style: .default))
        present(alert, animated: true)
    }

    private func getCurrentUserId() -> String {
        return UserDefaults.standard.string(forKey: "userId") ?? ""
    }
}

// MARK: - WKNavigationDelegate
extension GameViewController: WKNavigationDelegate {
    func webView(_ webView: WKWebView, didFinish navigation: WKNavigation!) {
        loadingIndicator.stopAnimating()
        loadingIndicator.isHidden = true
    }
}
