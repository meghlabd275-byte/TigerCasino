package com.tigercasino

import android.app.Application
import android.content.Context
import com.tigercasino.security.SecurityModule
import com.tigercasino.network.NetworkClient
import com.tigercasino.storage.PreferencesManager

/**
 * TigerCasino Application Class
 * Initializes all core modules and services
 */
class TigerCasinoApp : Application() {

    override fun onCreate() {
        super.onCreate()
        instance = this
        
        initializeModules()
    }

    private fun initializeModules() {
        // Initialize security module
        SecurityModule.initialize(this)
        
        // Initialize network client
        NetworkClient.initialize(
            baseUrl = BuildConfig.API_BASE_URL,
            apiKey = BuildConfig.API_KEY
        )
        
        // Initialize preferences
        PreferencesManager.initialize(this)
    }

    override fun attachBaseContext(base: Context) {
        super.attachBaseContext(base)
    }

    companion object {
        lateinit var instance: TigerCasinoApp
            private set

        fun getAppContext(): Context = instance.applicationContext
    }
}
