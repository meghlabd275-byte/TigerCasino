package services

import (
	"encoding/json"
	"fmt"
	"sync"
)

// I18nService handles internationalization
type I18nService struct {
	mu           sync.RWMutex
	translations map[string]map[string]string
	defaultLang  string
	supportedLangs []string
}

// NewI18nService creates a new i18n service
func NewI18nService() *I18nService {
	s := &I18nService{
		translations: make(map[string]map[string]string),
		defaultLang:  "en",
		supportedLangs: []string{},
	}
	s.initializeTranslations()
	return s
}

func (s *I18nService) initializeTranslations() {
	// English
	s.translations["en"] = map[string]string{
		"welcome": "Welcome to TigerCasino",
		"login": "Login",
		"register": "Register",
		"logout": "Logout",
		"deposit": "Deposit",
		"withdraw": "Withdraw",
		"balance": "Balance",
		"games": "Games",
		"slots": "Slots",
		"live_casino": "Live Casino",
		"sports": "Sports",
		"promotions": "Promotions",
		"vip": "VIP",
		"support": "Support",
		"profile": "Profile",
		"settings": "Settings",
		"language": "Language",
		"currency": "Currency",
		"bet": "Bet",
		"win": "Win",
		"loss": "Loss",
		"bonus": "Bonus",
		"free_spins": "Free Spins",
		"jackpot": "Jackpot",
		"crash": "Crash",
		"mines": "Mines",
		"plinko": "Plinko",
		"dice": "Dice",
		"roulette": "Roulette",
		"blackjack": "Blackjack",
		"poker": "Poker",
		"baccarat": "Baccarat",
	}
	s.supportedLangs = append(s.supportedLangs, "en")

	// Spanish
	s.translations["es"] = map[string]string{
		"welcome": "Bienvenido a TigerCasino",
		"login": "Iniciar sesión",
		"register": "Registrarse",
		"logout": "Cerrar sesión",
		"deposit": "Depositar",
		"withdraw": "Retirar",
		"balance": "Saldo",
		"games": "Juegos",
		"slots": "Tragamonedas",
		"live_casino": "Casino en Vivo",
		"sports": "Deportes",
		"promotions": "Promociones",
		"vip": "VIP",
		"support": "Soporte",
		"profile": "Perfil",
		"settings": "Configuración",
		"language": "Idioma",
		"currency": "Moneda",
	}
	s.supportedLangs = append(s.supportedLangs, "es")

	// French
	s.translations["fr"] = map[string]string{
		"welcome": "Bienvenue au TigerCasino",
		"login": "Connexion",
		"register": "S'inscrire",
		"logout": "Déconnexion",
		"deposit": "Dépôt",
		"withdraw": "Retrait",
		"balance": "Solde",
		"games": "Jeux",
		"slots": "Machines à sous",
		"live_casino": "Casino en direct",
		"sports": "Sports",
		"promotions": "Promotions",
		"vip": "VIP",
		"support": "Support",
		"profile": "Profil",
		"settings": "Paramètres",
		"language": "Langue",
		"currency": "Devise",
	}
	s.supportedLangs = append(s.supportedLangs, "fr")

	// German
	s.translations["de"] = map[string]string{
		"welcome": "Willkommen bei TigerCasino",
		"login": "Anmelden",
		"register": "Registrieren",
		"logout": "Abmelden",
		"deposit": "Einzahlung",
		"withdraw": "Auszahlung",
		"balance": "Guthaben",
		"games": "Spiele",
		"slots": "Slots",
		"live_casino": "Live Casino",
		"sports": "Sportwetten",
		"promotions": "Aktionen",
		"vip": "VIP",
		"support": "Support",
		"profile": "Profil",
		"settings": "Einstellungen",
		"language": "Sprache",
		"currency": "Währung",
	}
	s.supportedLangs = append(s.supportedLangs, "de")

	// Portuguese
	s.translations["pt"] = map[string]string{
		"welcome": "Bem-vindo ao TigerCasino",
		"login": "Entrar",
		"register": "Registrar",
		"logout": "Sair",
		"deposit": "Depositar",
		"withdraw": "Sacar",
		"balance": "Saldo",
		"games": "Jogos",
		"slots": "Slots",
		"live_casino": "Cassino ao Vivo",
		"sports": "Esportes",
		"promotions": "Promoções",
		"vip": "VIP",
		"support": "Suporte",
		"profile": "Perfil",
		"settings": "Configurações",
		"language": "Idioma",
		"currency": "Moeda",
	}
	s.supportedLangs = append(s.supportedLangs, "pt")

	// Russian
	s.translations["ru"] = map[string]string{
		"welcome": "Добро пожаловать в TigerCasino",
		"login": "Войти",
		"register": "Регистрация",
		"logout": "Выйти",
		"deposit": "Депозит",
		"withdraw": "Вывод",
		"balance": "Баланс",
		"games": "Игры",
		"slots": "Слоты",
		"live_casino": "Лайв Казино",
		"sports": "Ставки",
		"promotions": "Акции",
		"vip": "ВИП",
		"support": "Поддержка",
		"profile": "Профиль",
		"settings": "Настройки",
		"language": "Язык",
		"currency": "Валюта",
	}
	s.supportedLangs = append(s.supportedLangs, "ru")

	// Japanese
	s.translations["ja"] = map[string]string{
		"welcome": "TigerCasinoへようこそ",
		"login": "ログイン",
		"register": "登録",
		"logout": "ログアウト",
		"deposit": "入金",
		"withdraw": "出金",
		"balance": "残高",
		"games": "ゲーム",
		"slots": "スロット",
		"live_casino": "ライブカジノ",
		"sports": "スポーツ",
		"promotions": "プロモーション",
		"vip": "VIP",
		"support": "サポート",
		"profile": "プロフィール",
		"settings": "設定",
		"language": "言語",
		"currency": "通貨",
	}
	s.supportedLangs = append(s.supportedLangs, "ja")

	// Korean
	s.translations["ko"] = map[string]string{
		"welcome": "TigerCasino에 오신 것을 환영합니다",
		"login": "로그인",
		"register": "회원가입",
		"logout": "로그아웃",
		"deposit": "입금",
		"withdraw": "출금",
		"balance": "잔액",
		"games": "게임",
		"slots": "슬롯",
		"live_casino": "라이브 카지노",
		"sports": "스포츠",
		"promotions": "프로모션",
		"vip": "VIP",
		"support": "지원",
		"profile": "프로필",
		"settings": "설정",
		"language": "언어",
		"currency": "통화",
	}
	s.supportedLangs = append(s.supportedLangs, "ko")

	// Chinese
	s.translations["zh"] = map[string]string{
		"welcome": "欢迎来到TigerCasino",
		"login": "登录",
		"register": "注册",
		"logout": "退出",
		"deposit": "存款",
		"withdraw": "提款",
		"balance": "余额",
		"games": "游戏",
		"slots": "老虎机",
		"live_casino": "真人娱乐场",
		"sports": "体育",
		"promotions": "优惠",
		"vip": "VIP",
		"support": "支持",
		"profile": "个人资料",
		"settings": "设置",
		"language": "语言",
		"currency": "货币",
	}
	s.supportedLangs = append(s.supportedLangs, "zh")

	// Turkish
	s.translations["tr"] = map[string]string{
		"welcome": "TigerCasino'ya Hoş Geldiniz",
		"login": "Giriş Yap",
		"register": "Kayıt Ol",
		"logout": "Çıkış Yap",
		"deposit": "Para Yatırma",
		"withdraw": "Para Çekme",
		"balance": "Bakiye",
		"games": "Oyunlar",
		"slots": "Slotlar",
		"live_casino": "Canlı Casino",
		"sports": "Spor",
		"promotions": "Promosyonlar",
		"vip": "VIP",
		"support": "Destek",
		"profile": "Profil",
		"settings": "Ayarlar",
		"language": "Dil",
		"currency": "Para Birimi",
	}
	s.supportedLangs = append(s.supportedLangs, "tr")

	// Hindi
	s.translations["hi"] = map[string]string{
		"welcome": "TigerCasino में आपका स्वागत है",
		"login": "लॉगिन",
		"register": "रजिस्टर",
		"logout": "लॉगआउट",
		"deposit": "जमा",
		"withdraw": "निकासी",
		"balance": "शेष",
		"games": "गेम्स",
		"slots": "स्लॉट्स",
		"live_casino": "लाइव कैसीनो",
		"sports": "खेल",
		"promotions": "प्रमोशन",
		"vip": "वीआईपी",
		"support": "सहायता",
		"profile": "प्रोफाइल",
		"settings": "सेटिंग्स",
		"language": "भाषा",
		"currency": "मुद्रा",
	}
	s.supportedLangs = append(s.supportedLangs, "hi")

	// Arabic
	s.translations["ar"] = map[string]string{
		"welcome": "مرحباً بك في TigerCasino",
		"login": "تسجيل الدخول",
		"register": "تسجيل",
		"logout": "تسجيل الخروج",
		"deposit": "إيداع",
		"withdraw": "سحب",
		"balance": "الرصيد",
		"games": "الألعاب",
		"slots": "الفتحات",
		"live_casino": "الكازينو المباشر",
		"sports": "الرياضة",
		"promotions": "العروض",
		"vip": "VIP",
		"support": "الدعم",
		"profile": "الملف الشخصي",
		"settings": "الإعدادات",
		"language": "اللغة",
		"currency": "العملة",
	}
	s.supportedLangs = append(s.supportedLangs, "ar")
}

// Translate returns translated string
func (s *I18nService) Translate(lang, key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if lang == "" {
		lang = s.defaultLang
	}

	if translations, ok := s.translations[lang]; ok {
		if val, ok := translations[key]; ok {
			return val
		}
	}

	// Fallback to default language
	if translations, ok := s.translations[s.defaultLang]; ok {
		if val, ok := translations[key]; ok {
			return val
		}
	}

	return key
}

// GetSupportedLanguages returns list of supported languages
func (s *I18nService) GetSupportedLanguages() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	langs := make([]string, len(s.supportedLangs))
	copy(langs, s.supportedLangs)
	return langs
}

// SetLanguage sets user's preferred language
func (s *I18nService) SetLanguage(userID, lang string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate language
	valid := false
	for _, l := range s.supportedLangs {
		if l == lang {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("unsupported language: %s", lang)
	}

	// Save to user preferences (simplified)
	return nil
}

// GetAllTranslations returns all translations for a language
func (s *I18nService) GetAllTranslations(lang string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if lang == "" {
		lang = s.defaultLang
	}

	if translations, ok := s.translations[lang]; ok {
		result := make(map[string]string)
		for k, v := range translations {
			result[k] = v
		}
		return result, nil
	}

	return nil, fmt.Errorf("language not found: %s", lang)
}

// GetTranslationsJSON returns translations as JSON
func (s *I18nService) GetTranslationsJSON(lang string) (string, error) {
	translations, err := s.GetAllTranslations(lang)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(translations)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
