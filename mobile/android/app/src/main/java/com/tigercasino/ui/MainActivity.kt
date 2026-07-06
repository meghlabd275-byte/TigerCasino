package com.tigercasino.ui

import android.os.Bundle
import android.view.View
import android.webkit.WebChromeClient
import android.webkit.WebView
import android.webkit.WebViewClient
import android.widget.FrameLayout
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import com.tigercasino.R
import com.tigercasino.databinding.ActivityMainBinding
import com.tigercasino.network.ApiClient
import com.tigercasino.security.FingerprintManager
import com.tigercasino.storage.PreferencesManager
import kotlinx.coroutines.launch

/**
 * Main Activity - Primary entry point for the mobile app
 */
class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private var webView: WebView? = null
    private var loadingView: View? = null
    
    private val baseUrl = "https://tigercasino.com"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)
        
        setupWebView()
        checkAuthentication()
    }

    private fun setupWebView() {
        webView = binding.webView.apply {
            settings.apply {
                javaScriptEnabled = true
                domStorageEnabled = true
                databaseEnabled = true
                cacheMode = android.webkit.WebSettings.LOAD_DEFAULT
                
                // Performance settings
                useWideViewPort = true
                loadWithOverviewMode = true
                builtInZoomControls = false
                displayZoomControls = false
                
                // Security settings
                allowFileAccess = false
                allowContentAccess = false
                setSupportZoom(false)
                
                // Media settings
                mediaPlaybackRequiresUserGesture = false
                loadsImagesAutomatically = true
            }

            webViewClient = CustomWebViewClient()
            webChromeClient = CustomChromeClient()
            
            // Add JavaScript interface
            addJavascriptInterface(JSInterface(), "TigerCasino")
        }
    }

    private fun checkAuthentication() {
        val token = PreferencesManager.getAuthToken()
        
        if (token.isNullOrEmpty()) {
            // Show login
            loadUrl("$baseUrl/auth/login")
        } else {
            // Check if token is valid
            lifecycleScope.launch {
                try {
                    val isValid = ApiClient.verifyToken(token)
                    if (isValid) {
                        loadMainApp()
                    } else {
                        loadUrl("$baseUrl/auth/login")
                    }
                } catch (e: Exception) {
                    loadUrl("$baseUrl/auth/login")
                }
            }
        }
    }

    private fun loadMainApp() {
        // Inject auth token
        val token = PreferencesManager.getAuthToken()
        val js = """
            localStorage.setItem('auth_token', '$token');
            window.TigerCasino = { isMobile: true, platform: 'android' };
        """.trimIndent()
        
        webView?.evaluateJavascript(js, null)
        loadUrl(baseUrl)
    }

    private fun loadUrl(url: String) {
        showLoading()
        webView?.loadUrl(url)
    }

    private fun showLoading() {
        loadingView?.visibility = View.VISIBLE
    }

    private fun hideLoading() {
        loadingView?.visibility = View.GONE
    }

    // Handle back navigation
    override fun onBackPressed() {
        if (webView?.canGoBack() == true) {
            webView?.goBack()
        } else {
            super.onBackPressed()
        }
    }

    // JavaScript Interface
    inner class JSInterface {
        @android.webkit.JavascriptInterface
        fun getAuthToken(): String = PreferencesManager.getAuthToken() ?: ""
        
        @android.webkit.JavascriptInterface
        fun saveAuthToken(token: String) {
            PreferencesManager.saveAuthToken(token)
        }
        
        @android.webkit.JavascriptInterface
        fun clearAuthToken() {
            PreferencesManager.clearAuthToken()
        }
        
        @android.webkit.JavascriptInterface
        fun getBalance(): String {
            return PreferencesManager.getBalance().toString()
        }
        
        @android.webkit.JavascriptInterface
        fun updateBalance(amount: Double) {
            PreferencesManager.saveBalance(amount)
        }
        
        @android.webkit.JavascriptInterface
        fun getLanguage(): String {
            return PreferencesManager.getLanguage()
        }
        
        @android.webkit.JavascriptInterface
        fun authenticateWithFingerprint(callback: String) {
            FingerprintManager.authenticate(this@MainActivity) { success ->
                webView?.evaluateJavascript("$callback($success)", null)
            }
        }
        
        @android.webkit.JavascriptInterface
        fun openGame(gameId: String) {
            // Navigate to game screen
            runOnUiThread {
                // Start game activity
            }
        }
        
        @android.webkit.JavascriptInterface
        fun logout() {
            PreferencesManager.clearAuthToken()
            runOnUiThread {
                loadUrl("$baseUrl/auth/login")
            }
        }
    }

    // Custom WebView Client
    inner class CustomWebViewClient : WebViewClient() {
        override fun shouldOverrideUrlLoading(view: WebView?, url: String?): Boolean {
            url?.let {
                // Handle deep links
                if (it.startsWith("tigercasino://")) {
                    handleDeepLink(it)
                    return true
                }
            }
            return false
        }

        override fun onPageFinished(view: WebView?, url: String?) {
            super.onPageFinished(view, url)
            hideLoading()
            
            // Inject mobile-specific JavaScript
            injectMobileJS()
        }
    }

    // Custom Chrome Client
    inner class CustomChromeClient : WebChromeClient() {
        override fun onProgressChanged(view: WebView?, newProgress: Int) {
            super.onProgressChanged(view, newProgress)
            if (newProgress == 100) {
                hideLoading()
            }
        }
    }

    private fun injectMobileJS() {
        val js = """
            (function() {
                window.isMobileApp = true;
                window.platform = 'android';
                window.TigerCasino = window.TigerCasino || {};
            })();
        """.trimIndent()
        webView?.evaluateJavascript(js, null)
    }

    private fun handleDeepLink(url: String) {
        // Parse deep link and navigate
        val uri = android.net.Uri.parse(url)
        val path = uri.path
        
        when (path) {
            "/game" -> {
                val gameId = uri.getQueryParameter("id")
                gameId?.let { openGame(it) }
            }
            "/deposit" -> openDeposit()
            "/withdraw" -> openWithdraw()
        }
    }

    private fun openGame(gameId: String) {
        // Start game activity
    }

    private fun openDeposit() {
        // Open payment activity
    }

    private fun openWithdraw() {
        // Open payment activity
    }
}
