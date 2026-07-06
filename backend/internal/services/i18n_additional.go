package services

// Additional localization translations for missing languages
// Italian, Polish, Swedish, Dutch, Thai, Vietnamese, Indonesian, Greek, 
// Hindi, Arabic, Bangla, Urdu, Chinese (Simplified), Japanese

func LoadAdditionalTranslations() map[string]map[string]string {
    translations := make(map[string]map[string]string)
    
    // Italian
    translations["it"] = map[string]string{
        // General
        "welcome": "Benvenuto su TigerCasino",
        "welcome_message": "Vivi l'emozione del gioco premium",
        "login": "Accedi",
        "register": "Registrati",
        "logout": "Esci",
        "profile": "Profilo",
        "settings": "Impostazioni",
        "save": "Salva",
        "cancel": "Annulla",
        "confirm": "Conferma",
        "delete": "Elimina",
        "edit": "Modifica",
        "search": "Cerca",
        "loading": "Caricamento...",
        "error": "Errore",
        "success": "Successo",
        
        // Auth
        "email": "Email",
        "password": "Password",
        "username": "Nome utente",
        
        // Wallet
        "wallet": "Portafoglio",
        "balance": "Saldo",
        "deposit": "Deposita",
        "withdraw": "Preleva",
        
        // Games
        "games": "Giochi",
        "slots": "Slot",
        "live_casino": "Casino Live",
        "blackjack": "Blackjack",
        "roulette": "Roulette",
        
        // Sportsbook
        "sportsbook": "Scommesse Sportive",
        "place_bet": "Scommetti",
        
        // VIP
        "vip": "VIP",
        
        // Support
        "support": "Supporto",
        
        // Admin
        "admin": "Amministratore",
    }
    
    // Polish
    translations["pl"] = map[string]string{
        "welcome": "Witamy w TigerCasino",
        "welcome_message": "Doświadcz emocji premium",
        "login": "Zaloguj się",
        "register": "Zarejestruj się",
        "logout": "Wyloguj się",
        "profile": "Profil",
        "settings": "Ustawienia",
        "save": "Zapisz",
        "cancel": "Anuluj",
        "confirm": "Potwierdź",
        "email": "Email",
        "password": "Hasło",
        "username": "Nazwa użytkownika",
        "wallet": "Portfel",
        "balance": "Saldo",
        "deposit": "Depozyt",
        "withdraw": "Wypłata",
        "games": "Gry",
        "slots": "Automaty",
        "live_casino": "Kasyno na żywo",
        "sportsbook": "bukmacher",
    }
    
    // Swedish
    translations["sv"] = map[string]string{
        "welcome": "Välkommen till TigerCasino",
        "welcome_message": "Upplev spänningen",
        "login": "Logga in",
        "register": "Registrera",
        "logout": "Logga ut",
        "profile": "Profil",
        "settings": "Inställningar",
        "save": "Spara",
        "cancel": "Avbryt",
        "email": "E-post",
        "password": "Lösenord",
        "username": "Användarnamn",
        "wallet": "Plånbok",
        "balance": "Saldo",
        "deposit": "Insättning",
        "withdraw": "Uttag",
        "games": "Spel",
        "slots": "Slots",
        "live_casino": "Live Casino",
    }
    
    // Dutch
    translations["nl"] = map[string]string{
        "welcome": "Welkom bij TigerCasino",
        "welcome_message": "Ervaar premium gaming",
        "login": "Inloggen",
        "register": "Registreren",
        "logout": "Uitloggen",
        "profile": "Profiel",
        "settings": "Instellingen",
        "save": "Opslaan",
        "cancel": "Annuleren",
        "email": "E-mail",
        "password": "Wachtwoord",
        "username": "Gebruikersnaam",
        "wallet": "Portemonnee",
        "balance": "Saldo",
        "deposit": "Storten",
        "withdraw": "Opnemen",
        "games": "Spellen",
        "slots": "Gokkasten",
        "live_casino": "Live Casino",
    }
    
    // Thai
    translations["th"] = map[string]string{
        "welcome": "ยินดีต้อนรับสู่ TigerCasino",
        "welcome_message": "สัมผัสประสบการณ์เกมระดับพรีเมียม",
        "login": "เข้าสู่ระบบ",
        "register": "ลงทะเบียน",
        "logout": "ออกจากระบบ",
        "profile": "โปรไฟล์",
        "settings": "การตั้งค่า",
        "save": "บันทึก",
        "cancel": "ยกเลิก",
        "email": "อีเมล",
        "password": "รหัสผ่าน",
        "username": "ชื่อผู้ใช้",
        "wallet": "กระเป๋าเงิน",
        "balance": "ยอดเงิน",
        "deposit": "ฝากเงิน",
        "withdraw": "ถอนเงิน",
        "games": "เกม",
        "slots": "สล็อต",
    }
    
    // Vietnamese
    translations["vi"] = map[string]string{
        "welcome": "Chào mừng đến với TigerCasino",
        "welcome_message": "Trải nghiệm chơi game cao cấp",
        "login": "Đăng nhập",
        "register": "Đăng ký",
        "logout": "Đăng xuất",
        "profile": "Hồ sơ",
        "settings": "Cài đặt",
        "save": "Lưu",
        "cancel": "Hủy",
        "email": "Email",
        "password": "Mật khẩu",
        "username": "Tên người dùng",
        "wallet": "Ví",
        "balance": "Số dư",
        "deposit": "Nạp tiền",
        "withdraw": "Rút tiền",
    }
    
    // Indonesian
    translations["id"] = map[string]string{
        "welcome": "Selamat datang di TigerCasino",
        "welcome_message": "Rasakan pengalaman bermain game premium",
        "login": "Masuk",
        "register": "Daftar",
        "logout": "Keluar",
        "profile": "Profil",
        "settings": "Pengaturan",
        "save": "Simpan",
        "cancel": "Batal",
        "email": "Email",
        "password": "Kata sandi",
        "username": "Nama pengguna",
        "wallet": "Dompet",
        "balance": "Saldo",
        "deposit": "Setor",
        "withdraw": "Tarik",
    }
    
    // Greek
    translations["el"] = map[string]string{
        "welcome": "Καλωσήρθατε στο TigerCasino",
        "welcome_message": "Ζήστε την εμπειρία του premium gaming",
        "login": "Σύνδεση",
        "register": "Εγγραφή",
        "logout": "Αποσύνδεση",
        "profile": "Προφίλ",
        "settings": "Ρυθμίσεις",
        "save": "Αποθήκευση",
        "cancel": "Ακύρωση",
        "email": "Email",
        "password": "Κωδικός",
        "username": "Όνομα χρήστη",
        "wallet": "Πορτοφόλι",
        "balance": "Υπόλοιπο",
        "deposit": "Κατάθεση",
        "withdraw": "Ανάληψη",
    }
    
    // Hindi
    translations["hi"] = map[string]string{
        "welcome": "TigerCasino में आपका स्वागत है",
        "welcome_message": "प्रीमियम गेमिंग का अनुभव लें",
        "login": "लॉगिन करें",
        "register": "रजिस्टर करें",
        "logout": "लॉगआउट करें",
        "profile": "प्रोफ़ाइल",
        "settings": "सेटिंग्स",
        "save": "सहेजें",
        "cancel": "रद्द करें",
        "email": "ईमेल",
        "password": "पासवर्ड",
        "username": "यूज़रनेम",
        "wallet": "वॉलेट",
        "balance": "शेष राशि",
        "deposit": "जमा करें",
        "withdraw": "निकालें",
    }
    
    // Arabic (RTL)
    translations["ar"] = map[string]string{
        "welcome": "مرحباً بك في TigerCasino",
        "welcome_message": "استمتع بتجربة الألعاب المتميزة",
        "login": "تسجيل الدخول",
        "register": "إنشاء حساب",
        "logout": "تسجيل الخروج",
        "profile": "الملف الشخصي",
        "settings": "الإعدادات",
        "save": "حفظ",
        "cancel": "إلغاء",
        "email": "البريد الإلكتروني",
        "password": "كلمة المرور",
        "username": "اسم المستخدم",
        "wallet": "المحفظة",
        "balance": "الرصيد",
        "deposit": "إيداع",
        "withdraw": "سحب",
    }
    
    // Bangla
    translations["bn"] = map[string]string{
        "welcome": "TigerCasino-তে স্বাগতম",
        "welcome_message": "প্রিমিয়াম গেমিংয়ের অভিজ্ঞতা নিন",
        "login": "লগইন করুন",
        "register": "নিবন্ধন করুন",
        "logout": "লগআউট করুন",
        "profile": "প্রোফাইল",
        "settings": "সেটিংস",
        "save": "সংরক্ষণ করুন",
        "cancel": "বাতিল করুন",
        "email": "ইমেইল",
        "password": "পাসওয়ার্ড",
        "username": "ব্যবহারকারীর নাম",
        "wallet": "ওয়ালেট",
        "balance": "ব্যালেন্স",
        "deposit": "জমা করুন",
        "withdraw": "উত্তোলন করুন",
    }
    
    // Urdu
    translations["ur"] = map[string]string{
        "welcome": "TigerCasino میں خوش آمدید",
        "welcome_message": "پریمیم گیمنگ کا تجربہ کریں",
        "login": "لاگ ان کریں",
        "register": "رجسٹر کریں",
        "logout": "لاگ آؤٹ کریں",
        "profile": "پروفائل",
        "settings": "ترتیبات",
        "save": "محفوظ کریں",
        "cancel": "منسوخ کریں",
        "email": "ای میل",
        "password": "پاس ورڈ",
        "username": "صارف نام",
        "wallet": "والیٹ",
        "balance": "بیلینس",
        "deposit": "ڈپوزٹ",
        "withdraw": "ویدرا",
    }
    
    // Chinese Simplified
    translations["zh"] = map[string]string{
        "welcome": "欢迎来到 TigerCasino",
        "welcome_message": "体验顶级游戏",
        "login": "登录",
        "register": "注册",
        "logout": "退出",
        "profile": "个人资料",
        "settings": "设置",
        "save": "保存",
        "cancel": "取消",
        "email": "邮箱",
        "password": "密码",
        "username": "用户名",
        "wallet": "钱包",
        "balance": "余额",
        "deposit": "充值",
        "withdraw": "提现",
    }
    
    // Japanese
    translations["ja"] = map[string]string{
        "welcome": "TigerCasinoへようこそ",
        "welcome_message": "プレミアムゲーム体验",
        "login": "ログイン",
        "register": "登録",
        "logout": "ログアウト",
        "profile": "プロフィール",
        "settings": "設定",
        "save": "保存",
        "cancel": "キャンセル",
        "email": "メール",
        "password": "パスワード",
        "username": "ユーザー名",
        "wallet": "ウォレット",
        "balance": "残高",
        "deposit": "入金",
        "withdraw": "出金",
    }
    
    return translations
}
