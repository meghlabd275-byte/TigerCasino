package com.tigercasino

import android.os.Bundle
import android.webkit.WebView
import android.webkit.WebSettings
import android.webkit.WebChromeClient
import android.webkit.WebViewClient
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch

class MainActivity : AppCompatActivity() {

    private lateinit var webView: WebView
    private var currentMultiplier = 1.0
    private var isGameRunning = false
    private var balance = 0.0
    private var currentBet = 10.0

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        // Initialize Rust security module
        SecurityModule.initialize(this)
        
        setupWebView()
        loadMainPage()
    }

    private fun setupWebView() {
        webView = WebView(this)
        setContentView(webView)
        
        val settings: WebSettings = webView.settings
        settings.javaScriptEnabled = true
        settings.domStorageEnabled = true
        settings.allowFileAccess = true
        settings.allowContentAccess = true
        settings.mediaPlaybackRequiresUserGesture = false
        settings.cacheMode = WebSettings.LOAD_NO_CACHE
        
        // High-performance settings
        settings.useWideViewPort = true
        settings.loadWithOverviewMode = true
        settings.builtInZoomControls = false
        settings.displayZoomControls = false
        
        webView.webChromeClient = WebChromeClient()
        webView.webViewClient = object : WebViewClient() {
            override fun onPageFinished(view: WebView?, url: String?) {
                super.onPageFinished(view, url)
                hideLoading()
            }
        }
    }

    private fun loadMainPage() {
        showLoading()
        webView.loadUrl("https://tigercasino.com")
    }

    // Game functions
    fun placeBet(amount: Double) {
        currentBet = amount
        
        lifecycleScope.launch(Dispatchers.IO) {
            val result = GameEngine.placeBet(
                gameId = getCurrentGameId(),
                amount = amount,
                userId = getCurrentUserId()
            )
            
            if (result.success) {
                isGameRunning = true
                runOnUiThread {
                    startGameLoop()
                }
            }
        }
    }

    fun cashOut() {
        if (!isGameRunning) return
        
        val payout = currentMultiplier * currentBet
        
        lifecycleScope.launch(Dispatchers.IO) {
            val result = GameEngine.cashOut(
                gameId = getCurrentGameId(),
                payout = payout
            )
            
            if (result.success) {
                isGameRunning = false
                balance += payout
                runOnUiThread {
                    showPayoutToast(payout)
                }
            }
        }
    }

    private fun startGameLoop() {
        val handler = android.os.Handler(android.os.Looper.getMainLooper())
        val runnable = object : Runnable {
            override fun run() {
                if (isGameRunning) {
                    currentMultiplier = GameEngine.getCurrentMultiplier(getCurrentGameId())
                    updateGameUI()
                    handler.postDelayed(this, 16) // 60fps
                }
            }
        }
        handler.post(runnable)
    }

    private fun updateGameUI() {
        webView.evaluateJavascript("updateMultiplier($currentMultiplier)", null)
    }

    private fun showPayoutToast(amount: Double) {
        Toast.makeText(this, "You won ${String.format("%.2f", amount)}", Toast.LENGTH_LONG).show()
    }

    private fun showLoading() {}
    private fun hideLoading() {}
    private fun getCurrentGameId() = "crash_001"
    private fun getCurrentUserId() = "user_123"

    override fun onBackPressed() {
        if (webView.canGoBack()) webView.goBack() else super.onBackPressed()
    }
}
