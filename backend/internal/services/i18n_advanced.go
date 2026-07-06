package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

// Supported languages
type Language string

const (
	LanguageEnglish    Language = "en"
	LanguageSpanish    Language = "es"
	LanguageFrench     Language = "fr"
	LanguageGerman     Language = "de"
	LanguagePortuguese Language = "pt"
	LanguageRussian    Language = "ru"
	LanguageJapanese   Language = "ja"
	LanguageKorean     Language = "ko"
	LanguageChinese    Language = "zh"
	LanguageTurkish    Language = "tr"
	LanguageArabic     Language = "ar"
	LanguageHindi      Language = "hi"
)

// Language info
type LanguageInfo struct {
	Code      Language `json:"code"`
	Name      string   `json:"name"`
	NativeName string  `json:"native_name"`
	Direction string   `json:"direction"` // ltr or rtl
	Flag      string   `json:"flag"`
}

// Supported currencies
type Currency struct {
	Code         string  `json:"code"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	DecimalPlaces int    `json:"decimal_places"`
}

// Translation map
type TranslationMap map[string]map[string]string // lang -> key -> value

// Localization service
type LocalizationService struct {
	translations TranslationMap
	mu           sync.RWMutex
	languages    map[Language]LanguageInfo
	currencies  map[string]Currency
	defaultLang Language
}

// NewLocalizationService creates a new localization service
func NewLocalizationService() *LocalizationService {
	s := &LocalizationService{
		translations: make(TranslationMap),
		languages:    make(map[Language]LanguageInfo),
		currencies:  make(map[string]Currency),
		defaultLang: LanguageEnglish,
	}

	// Initialize languages
	s.initLanguages()

	// Initialize currencies
	s.initCurrencies()

	// Load default translations
	s.loadDefaultTranslations()

	return s
}

func (s *LocalizationService) initLanguages() {
	s.languages = map[Language]LanguageInfo{
		LanguageEnglish:    {Code: LanguageEnglish, Name: "English", NativeName: "English", Direction: "ltr", Flag: "🇺🇸"},
		LanguageSpanish:    {Code: LanguageSpanish, Name: "Spanish", NativeName: "Español", Direction: "ltr", Flag: "🇪🇸"},
		LanguageFrench:     {Code: LanguageFrench, Name: "French", NativeName: "Français", Direction: "ltr", Flag: "🇫🇷"},
		LanguageGerman:     {Code: LanguageGerman, Name: "German", NativeName: "Deutsch", Direction: "ltr", Flag: "🇩🇪"},
		LanguagePortuguese: {Code: LanguagePortuguese, Name: "Portuguese", NativeName: "Português", Direction: "ltr", Flag: "🇧🇷"},
		LanguageRussian:    {Code: LanguageRussian, Name: "Russian", NativeName: "Русский", Direction: "ltr", Flag: "🇷🇺"},
		LanguageJapanese:   {Code: LanguageJapanese, Name: "Japanese", NativeName: "日本語", Direction: "ltr", Flag: "🇯🇵"},
		LanguageKorean:     {Code: LanguageKorean, Name: "Korean", NativeName: "한국어", Direction: "ltr", Flag: "🇰🇷"},
		LanguageChinese:    {Code: LanguageChinese, Name: "Chinese", NativeName: "中文", Direction: "ltr", Flag: "🇨🇳"},
		LanguageTurkish:    {Code: LanguageTurkish, Name: "Turkish", NativeName: "Türkçe", Direction: "ltr", Flag: "🇹🇷"},
		LanguageArabic:     {Code: LanguageArabic, Name: "Arabic", NativeName: "العربية", Direction: "rtl", Flag: "🇸🇦"},
		LanguageHindi:      {Code: LanguageHindi, Name: "Hindi", NativeName: "हिन्दी", Direction: "ltr", Flag: "🇮🇳"},
	}
}

func (s *LocalizationService) initCurrencies() {
	s.currencies = map[string]Currency{
		"USD": {Code: "USD", Symbol: "$", Name: "US Dollar", DecimalPlaces: 2},
		"EUR": {Code: "EUR", Symbol: "€", Name: "Euro", DecimalPlaces: 2},
		"GBP": {Code: "GBP", Symbol: "£", Name: "British Pound", DecimalPlaces: 2},
		"JPY": {Code: "JPY", Symbol: "¥", Name: "Japanese Yen", DecimalPlaces: 0},
		"CNY": {Code: "CNY", Symbol: "¥", Name: "Chinese Yuan", DecimalPlaces: 2},
		"KRW": {Code: "KRW", Symbol: "₩", Name: "South Korean Won", DecimalPlaces: 0},
		"BRL": {Code: "BRL", Symbol: "R$", Name: "Brazilian Real", DecimalPlaces: 2},
		"RUB": {Code: "RUB", Symbol: "₽", Name: "Russian Ruble", DecimalPlaces: 2},
		"TRY": {Code: "TRY", Symbol: "₺", Name: "Turkish Lira", DecimalPlaces: 2},
		"INR": {Code: "INR", Symbol: "₹", Name: "Indian Rupee", DecimalPlaces: 2},
		"CAD": {Code: "CAD", Symbol: "C$", Name: "Canadian Dollar", DecimalPlaces: 2},
		"AUD": {Code: "AUD", Symbol: "A$", Name: "Australian Dollar", DecimalPlaces: 2},
		"BTC": {Code: "BTC", Symbol: "₿", Name: "Bitcoin", DecimalPlaces: 8},
		"ETH": {Code: "ETH", Symbol: "Ξ", Name: "Ethereum", DecimalPlaces: 8},
		"USDT": {Code: "USDT", Symbol: "₮", Name: "Tether", DecimalPlaces: 2},
	}
}

func (s *LocalizationService) loadDefaultTranslations() {
	// Load English translations
	s.translations["en"] = map[string]string{
		// General
		"welcome": "Welcome to TigerCasino",
		"welcome_message": "Experience the thrill of premium gaming",
		"login": "Login",
		"register": "Register",
		"logout": "Logout",
		"profile": "Profile",
		"settings": "Settings",
		"save": "Save",
		"cancel": "Cancel",
		"confirm": "Confirm",
		"delete": "Delete",
		"edit": "Edit",
		"search": "Search",
		"loading": "Loading...",
		"error": "Error",
		"success": "Success",
		"warning": "Warning",
		"info": "Information",

		// Auth
		"email": "Email",
		"password": "Password",
		"confirm_password": "Confirm Password",
		"username": "Username",
		"forgot_password": "Forgot Password?",
		"remember_me": "Remember Me",
		"dont_have_account": "Don't have an account?",
		"already_have_account": "Already have an account?",
		"signup_success": "Registration successful!",
		"login_success": "Login successful!",
		"invalid_credentials": "Invalid email or password",
		"email_already_exists": "Email already registered",
		"username_already_exists": "Username already taken",

		// Wallet
		"wallet": "Wallet",
		"balance": "Balance",
		"deposit": "Deposit",
		"withdraw": "Withdraw",
		"transactions": "Transactions",
		"deposit_address": "Deposit Address",
		"copy_address": "Copy Address",
		"address_copied": "Address copied to clipboard",
		"withdrawal_request": "Withdrawal Request",
		"withdrawal_success": "Withdrawal request submitted",
		"min_withdrawal": "Minimum withdrawal:",
		"processing_time": "Processing time:",

		// Games
		"games": "Games",
		"slots": "Slots",
		"live_casino": "Live Casino",
		"table_games": "Table Games",
		"blackjack": "Blackjack",
		"roulette": "Roulette",
		"baccarat": "Baccarat",
		"poker": "Poker",
		"crash": "Crash",
		"dice": "Dice",
		"play_now": "Play Now",
		"free_play": "Demo Play",
		"game_info": "Game Info",
		"bet": "Bet",
		"win": "Win",
		"multiplier": "Multiplier",

		// Sportsbook
		"sportsbook": "Sportsbook",
		"live_betting": "Live Betting",
		"upcoming": "Upcoming",
		"odds": "Odds",
		"place_bet": "Place Bet",
		"potential_payout": "Potential Payout",
		"bet_slip": "Bet Slip",
		"cashout": "Cash Out",
		"cashout_available": "Cash Out Available",
		"parlay": "Parlay",
		"single": "Single",
		"acca": "Accumulator",

		// VIP
		"vip": "VIP",
		"vip_level": "VIP Level",
		"vip_points": "VIP Points",
		"loyalty_tier": "Loyalty Tier",
		"benefits": "Benefits",
		"exclusive_bonuses": "Exclusive Bonuses",
		"personal_manager": "Personal Manager",
		"faster_withdrawals": "Faster Withdrawals",

		// Support
		"support": "Support",
		"contact_us": "Contact Us",
		"faq": "FAQ",
		"live_chat": "Live Chat",
		"email_support": "Email Support",
		"response_time": "Response Time",
		
		// Legal
		"terms_of_service": "Terms of Service",
		"privacy_policy": "Privacy Policy",
		"responsible_gaming": "Responsible Gaming",
		"age_restriction": "18+ Only",
		"gambling_warning": "Gambling involves risk. Please play responsibly.",

		// Admin
		"admin": "Admin",
		"dashboard": "Dashboard",
		"users": "Users",
		"user_management": "User Management",
		"game_management": "Game Management",
		"financial_management": "Financial Management",
		"reports": "Reports",
		"settings": "Settings",
	}

	// Load Spanish translations
	s.translations["es"] = map[string]string{
		"welcome": "Bienvenido a TigerCasino",
		"welcome_message": "Experimenta la emoción del juego premium",
		"login": "Iniciar sesión",
		"register": "Registrarse",
		"logout": "Cerrar sesión",
		"profile": "Perfil",
		"settings": "Configuración",
		"save": "Guardar",
		"cancel": "Cancelar",
		"confirm": "Confirmar",
		"wallet": "Billetera",
		"balance": "Saldo",
		"deposit": "Depositar",
		"withdraw": "Retirar",
		"transactions": "Transacciones",
		"games": "Juegos",
		"slots": "Tragamonedas",
		"live_casino": "Casino en Vivo",
		"blackjack": "Blackjack",
		"roulette": "Ruleta",
		"sportsbook": "Casa de Apuestas",
		"vip": "VIP",
		"support": "Soporte",
		"admin": "Administración",
	}

	// Load more languages...
	for lang := range s.languages {
		if _, ok := s.translations[string(lang)]; !ok {
			s.translations[string(lang)] = s.translations["en"]
		}
	}
}

// LoadTranslationsFromFile loads translations from a JSON file
func (s *LocalizationService) LoadTranslationsFromFile(lang Language, filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.translations[string(lang)] = translations

	return nil
}

// Translate translates a key
func (s *LocalizationService) Translate(key string, lang Language, args ...interface{}) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	translations, ok := s.translations[string(lang)]
	if !ok {
		// Fallback to English
		translations = s.translations["en"]
	}

	value, ok := translations[key]
	if !ok {
		// Fallback to English
		if lang != LanguageEnglish {
			if enValue, enOk := s.translations["en"][key]; enOk {
				return s.formatMessage(enValue, args...)
			}
		}
		return key
	}

	return s.formatMessage(value, args...)
}

func (s *LocalizationService) formatMessage(message string, args ...interface{}) string {
	if len(args) == 0 {
		return message
	}

	result := message
	for i, arg := range args {
		placeholder := fmt.Sprintf("{%d}", i)
		result = strings.Replace(result, placeholder, fmt.Sprintf("%v", arg), -1)
	}

	return result
}

// GetLanguage gets language info
func (s *LocalizationService) GetLanguage(lang Language) (LanguageInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info, ok := s.languages[lang]
	return info, ok
}

// GetLanguages gets all languages
func (s *LocalizationService) GetLanguages() []LanguageInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	infos := make([]LanguageInfo, 0, len(s.languages))
	for _, info := range s.languages {
		infos = append(infos, info)
	}
	return infos
}

// GetCurrencies gets all currencies
func (s *LocalizationService) GetCurrencies() []Currency {
	s.mu.RLock()
	defer s.mu.RUnlock()

	currencies := make([]Currency, 0, len(s.currencies))
	for _, currency := range s.currencies {
		currencies = append(currencies, currency)
	}
	return currencies
}

// GetCurrency gets currency info
func (s *LocalizationService) GetCurrency(code string) (Currency, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	currency, ok := s.currencies[code]
	return currency, ok
}

// FormatCurrency formats an amount with currency
func (s *LocalizationService) FormatCurrency(amount float64, currency string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	curr, ok := s.currencies[currency]
	if !ok {
		return fmt.Sprintf("%.2f %s", amount, currency)
	}

	format := fmt.Sprintf("%%.%df", curr.DecimalPlaces)
	return fmt.Sprintf("%s %s", curr.Symbol, fmt.Sprintf(format, amount))
}

// GetDirection gets text direction for language
func (s *LocalizationService) GetDirection(lang Language) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if info, ok := s.languages[lang]; ok {
		return info.Direction
	}
	return "ltr"
}

// GetSupportedLanguages gets supported languages for user
func (s *LocalizationService) GetSupportedLanguages(userRegion string) []LanguageInfo {
	// Could filter by region in production
	return s.GetLanguages()
}

// DetectLanguage detects language from headers
func (s *LocalizationService) DetectLanguage(acceptLanguage string) Language {
	// Parse Accept-Language header
	parts := strings.Split(acceptLanguage, ",")
	if len(parts) == 0 {
		return s.defaultLang
	}

	langCode := strings.TrimSpace(strings.Split(parts[0], ";")[0])
	lang := Language(strings.ToLower(langCode)[:2])

	if _, ok := s.languages[lang]; ok {
		return lang
	}

	return s.defaultLang
}

// Create language files for all supported languages
func (s *LocalizationService) GenerateLanguageFiles(outputDir string) error {
	for lang := range s.languages {
		data, err := json.MarshalIndent(s.translations[string(lang)], "", "  ")
		if err != nil {
			return err
		}

		filename := filepath.Join(outputDir, fmt.Sprintf("%s.json", lang))
		if err := ioutil.WriteFile(filename, data, 0644); err != nil {
			return err
		}
	}
	return nil
}
