package services

// Final missing translations for 3 remaining languages
// Romanian, Hungarian, Danish

func LoadFinalTranslations() map[string]map[string]string {
    translations := make(map[string]map[string]string)
    
    // Romanian
    translations["ro"] = map[string]string{
        "welcome": "Bine ai venit la TigerCasino",
        "welcome_message": "Experimentează emoția jocurilor premium",
        "login": "Autentificare",
        "register": "Înregistrare",
        "logout": "Deconectare",
        "profile": "Profil",
        "settings": "Setări",
        "save": "Salvează",
        "cancel": "Anulează",
        "email": "Email",
        "password": "Parolă",
        "username": "Nume utilizator",
        "wallet": "Portofel",
        "balance": "Sold",
        "deposit": "Depunere",
        "withdraw": "Retragere",
    }
    
    // Hungarian
    translations["hu"] = map[string]string{
        "welcome": "Üdvözöljük a TigerCasino-ban",
        "welcome_message": "Éld meg a prémium játékok izgalmát",
        "login": "Bejelentkezés",
        "register": "Regisztráció",
        "logout": "Kijelentkezés",
        "profile": "Profil",
        "settings": "Beállítások",
        "save": "Mentés",
        "cancel": "Mégse",
        "email": "E-mail",
        "password": "Jelszó",
        "username": "Felhasználónév",
        "wallet": "Tárca",
        "balance": "Egyenleg",
        "deposit": "Befizetés",
        "withdraw": "Kivétel",
    }
    
    // Danish
    translations["da"] = map[string]string{
        "welcome": "Velkommen til TigerCasino",
        "welcome_message": "Oplev spændingen ved premium spil",
        "login": "Log ind",
        "register": "Tilmeld",
        "logout": "Log ud",
        "profile": "Profil",
        "settings": "Indstillinger",
        "save": "Gem",
        "cancel": "Annuller",
        "email": "E-mail",
        "password": "Adgangskode",
        "username": "Brugernavn",
        "wallet": "Pung",
        "balance": "Saldo",
        "deposit": "Indbetaling",
        "withdraw": "Udbetaling",
    }
    
    return translations
}
