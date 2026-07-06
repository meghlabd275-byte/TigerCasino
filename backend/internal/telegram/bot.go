package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// TelegramBot represents a Telegram Bot for TigerCasino
type TelegramBot struct {
	Token           string
	WebhookURL      string
	APIURL          string
	AllowedUpdates  []string
	Commands        map[string]CommandHandler
	SessionStore    map[int64]*UserSession
	Mutex           sync.RWMutex
	WSHub           interface{} // WebSocket hub for real-time updates
}

// CommandHandler defines the function signature for command handlers
type CommandHandler func(update *Update) *Response

// UserSession represents a user's session in Telegram
type UserSession struct {
	UserID       int64
	Username     string
	State        string // "main", "deposit", "withdraw", "game"
	GameType     string
	LastActivity time.Time
}

// Update represents an incoming Telegram update
type Update struct {
	UpdateID int64 `json:"update_id"`
	Message  *Message `json:"message"`
	Callback *CallbackQuery `json:"callback_query"`
}

// Message represents a Telegram message
type Message struct {
	MessageID int64  `json:"message_id"`
	From      *User  `json:"from"`
	Chat      *Chat  `json:"chat"`
	Text      string `json:"text"`
	Entities  []MessageEntity `json:"entities"`
}

// User represents a Telegram user
type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// Chat represents a Telegram chat
type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// MessageEntity represents a text entity
type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url"`
}

// CallbackQuery represents a callback query
type CallbackQuery struct {
	ID       string   `json:"id"`
	From     *User    `json:"from"`
	Message  *Message `json:"message"`
	Data     string   `json:"data"`
}

// Response represents a Telegram API response
type Response struct {
	Ok          bool                   `json:"ok"`
	ErrorCode   int                    `json:"error_code"`
	Description string                 `json:"description"`
	Result      interface{}            `json:"result"`
}

// NewTelegramBot creates a new Telegram bot instance
func NewTelegramBot(token string) *TelegramBot {
	bot := &TelegramBot{
		Token:          token,
		APIURL:         "https://api.telegram.org",
		Commands:       make(map[string]CommandHandler),
		SessionStore:   make(map[int64]*UserSession),
		AllowedUpdates: []string{"message", "callback_query"},
	}
	
	// Register default commands
	bot.registerCommands()
	
	return bot
}

// Start starts the Telegram bot webhook server
func (b *TelegramBot) Start(port string) {
	http.HandleFunc("/"+b.Token, b.handleUpdate)
	
	// Set webhook
	if err := b.SetWebhook(); err != nil {
		log.Printf("Failed to set webhook: %v", err)
	}
	
	log.Printf("Telegram bot started on port %s", port)
	log.Printf("Webhook URL: https://your-domain.com/%s", b.Token)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// SetWebhook sets the webhook for the bot
func (b *TelegramBot) SetWebhook() error {
	webhookURL := os.Getenv("TELEGRAM_WEBHOOK_URL")
	if webhookURL == "" {
		webhookURL = b.WebhookURL
	}
	
	if webhookURL == "" {
		return fmt.Errorf("webhook URL not set")
	}
	
	_, err := b.Request("setWebhook", map[string]string{
		"url":                  webhookURL + "/" + b.Token,
		"allowed_updates":      `["message","callback_query"]`,
		"drop_pending_updates": "true",
	})
	
	return err
}

// Request makes an API request to Telegram
func (b *TelegramBot) Request(method string, params map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	
	url := fmt.Sprintf("%s/bot%s/%s", b.APIURL, b.Token, method)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	
	if !response.Ok {
		return nil, fmt.Errorf("Telegram API error: %s", response.Description)
	}
	
	result, ok := response.Result.(map[string]interface{})
	if !ok {
		return make(map[string]interface{}), nil
	}
	
	return result, nil
}

// handleUpdate handles incoming Telegram updates
func (b *TelegramBot) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Error decoding update: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// Handle message
	if update.Message != nil {
		b.handleMessage(update.Message)
	}
	
	// Handle callback query
	if update.Callback != nil {
		b.handleCallback(update.Callback)
	}
	
	w.WriteHeader(http.StatusOK)
}

// handleMessage handles incoming messages
func (b *TelegramBot) handleMessage(msg *Message) {
	if msg == nil || msg.Text == "" {
		return
	}
	
	// Get or create session
	session := b.getOrCreateSession(msg.From)
	
	// Check if it's a command
	if msg.Text[0] == '/' {
		parts := strings.Split(msg.Text[1:], " ")
		command := strings.ToLower(parts[0])
		
		if handler, ok := b.Commands[command]; ok {
			handler(&Update{
				Message: msg,
			})
		} else {
			b.SendMessage(msg.Chat.ID, "Unknown command. Use /help to see available commands.")
		}
		return
	}
	
	// Handle state-based input
	b.handleStatefulInput(session, msg)
}

// handleCallback handles callback queries
func (b *TelegramBot) handleCallback(cb *CallbackQuery) {
	if cb == nil || cb.Data == "" {
		return
	}
	
	// Get session
	session := b.getOrCreateSession(cb.From)
	
	// Parse callback data
	parts := strings.Split(cb.Data, ":")
	action := parts[0]
	
	switch action {
	case "game":
		b.showGameMenu(cb.Message.Chat.ID, cb.From.ID)
	case "balance":
		b.showBalance(cb.Message.Chat.ID, session.UserID)
	case "deposit":
		b.showDeposit(cb.Message.Chat.ID)
	case "withdraw":
		b.showWithdraw(cb.Message.Chat.ID)
	case "back":
		b.showMainMenu(cb.Message.Chat.ID, cb.From.ID)
	}
	
	// Answer callback query
	b.AnswerCallbackQuery(cb.ID, "")
}

// registerCommands registers all bot commands
func (b *TelegramBot) registerCommands() {
	b.Commands["start"] = b.cmdStart
	b.Commands["help"] = b.cmdHelp
	b.Commands["register"] = b.cmdRegister
	b.Commands["balance"] = b.cmdBalance
	b.Commands["deposit"] = b.cmdDeposit
	b.Commands["withdraw"] = b.cmdWithdraw
	b.Commands["games"] = b.cmdGames
	b.Commands["profile"] = b.cmdProfile
	b.Commands["vip"] = b.cmdVIP
	b.Commands["referral"] = b.cmdReferral
	b.Commands["support"] = b.cmdSupport
}

// Command handlers
func (b *TelegramBot) cmdStart(update *Update) *Response {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	
	welcomeText := fmt.Sprintf(`🎰 *Welcome to TigerCasino!* 🐯

Hello @%s!

I'm your personal TigerCasino assistant. Use me to:
• 📊 Check your balance
• 💰 Make deposits & withdrawals  
• 🎮 Play games directly
• 🎁 Claim bonuses
• 📞 Get support

Ready to play? Use /games to start!

_Responsible gambling: Play responsibly_`, username)
	
	b.SendMessage(chatID, welcomeText)
	b.showMainMenu(chatID, update.Message.From.ID)
	
	return nil
}

func (b *TelegramBot) cmdHelp(update *Update) *Response {
	helpText := `🎰 *TigerCasino Commands*

/start - Start the bot
/help - Show this help message
/register - Register an account
/balance - Check your balance
/deposit - Get deposit address
/withdraw - Request withdrawal
/games - Browse games
/profile - View profile
/vip - VIP status
/referral - Your referral code
/support - Contact support

_Play responsibly_`
	
	b.SendMessage(update.Message.Chat.ID, helpText)
	return nil
}

func (b *TelegramBot) cmdBalance(update *Update) *Response {
	// In production, would fetch from API
	balanceText := `💰 *Your Balance*

_Loading from TigerCasino..._
`
	b.SendMessage(update.Message.Chat.ID, balanceText)
	return nil
}

func (b *TelegramBot) cmdDeposit(update *Update) *Response {
	depositText := `💵 *Deposit*

Your deposit address:

*BTC:* ` + "```" + `bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh```" + `

_Network: Bitcoin (BTC)_
_Minimum: 0.0001 BTC_

Copy the address above and send BTC. Deposits are processed automatically after 1 confirmation.`
	
	b.SendMessage(update.Message.Chat.ID, depositText)
	return nil
}

func (b *TelegramBot) cmdWithdraw(update *Update) *Response {
	withdrawText := `💸 *Withdraw*

To withdraw:
1. Enter your wallet address
2. Enter amount
3. Confirm

_Processing time: 5-30 minutes_
_Fee: 1%_
`
	b.SendMessage(update.Message.Chat.ID, withdrawText)
	return nil
}

func (b *TelegramBot) cmdGames(update *Update) *Response {
	b.showGameMenu(update.Message.Chat.ID, update.Message.From.ID)
	return nil
}

func (b *TelegramBot) cmdProfile(update *Update) *Response {
	profileText := `👤 *Your Profile*

_Loading from TigerCasino..._
`
	b.SendMessage(update.Message.Chat.ID, profileText)
	return nil
}

func (b *TelegramBot) cmdVIP(update *Update) *Response {
	vipText := `👑 *VIP Program*

Your VIP Status:
• Level: Bronze
• Points: 0
• Rakeback: 0%

_Play more to level up!_

Bronze → Silver → Gold → Platinum → Diamond → VIP
`
	b.SendMessage(update.Message.Chat.ID, vipText)
	return nil
}

func (b *TelegramBot) cmdReferral(update *Update) *Response {
	refText := `📢 *Referral Program*

Your referral code: *TIGER123*

Earn 20% commission on your friends' deposits!

Share your code and earn free crypto! 🎉
`
	b.SendMessage(update.Message.Chat.ID, refText)
	return nil
}

func (b *TelegramBot) cmdSupport(update *Update) *Response {
	supportText := `📞 *Support*

• Email: support@tigercasino.com
• Live Chat: Available 24/7
• Response Time: < 1 minute

_We're here to help!_
`
	b.SendMessage(update.Message.Chat.ID, supportText)
	return nil
}

// SendMessage sends a message to a chat
func (b *TelegramBot) SendMessage(chatID int64, text string) error {
	_, err := b.Request("sendMessage", map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	})
	return err
}

// SendInlineKeyboard sends a message with inline keyboard
func (b *TelegramBot) SendInlineKeyboard(chatID int64, text string, keyboard [][]KeyboardButton) error {
	_, err := b.Request("sendMessage", map[string]interface{}{
		"chat_id": chatID,
		"text": text,
		"parse_mode": "Markdown",
		"reply_markup": map[string]interface{}{
			"inline_keyboard": keyboard,
		},
	})
	return err
}

// KeyboardButton represents an inline keyboard button
type KeyboardButton struct {
	Text string `json:"text"`
	URL  string `json:"url,omitempty"`
	Data string `json:"callback_data,omitempty"`
}

// AnswerCallbackQuery answers a callback query
func (b *TelegramBot) AnswerCallbackQuery(callbackID string, text string) error {
	_, err := b.Request("answerCallbackQuery", map[string]interface{}{
		"callback_query_id": callbackID,
		"text":             text,
	})
	return err
}

// Helper functions
func (b *TelegramBot) getOrCreateSession(user *User) *UserSession {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	
	if session, ok := b.SessionStore[user.ID]; ok {
		session.LastActivity = time.Now()
		return session
	}
	
	session := &UserSession{
		UserID:       user.ID,
		Username:     user.Username,
		State:        "main",
		LastActivity: time.Now(),
	}
	
	b.SessionStore[user.ID] = session
	return session
}

func (b *TelegramBot) showMainMenu(chatID int64, userID int64) {
	menu := [][]KeyboardButton{
		{{Text: "🎮 Play Games", Data: "game:show"}, {Text: "💰 Balance", Data: "balance:show"}},
		{{Text: "💵 Deposit", Data: "deposit:show"}, {Text: "💸 Withdraw", Data: "withdraw:show"}},
		{{Text: "👑 VIP Status", Data: "vip:show"}, {Text: "📢 Referral", Data: "ref:show"}},
	}
	
	b.SendInlineKeyboard(chatID, "🎰 *Main Menu*\n\nChoose an option:", menu)
}

func (b *TelegramBot) showGameMenu(chatID int64, userID int64) {
	games := [][]KeyboardButton{
		{{Text: "🎲 Crash", Data: "game:crash"}, {Text: "💣 Mines", Data: "game:mines"}},
		{{Text: "🎯 Plinko", Data: "game:plinko"}, {Text: "🎲 Dice", Data: "game:dice"}},
		{{Text: "🎰 Slots", Data: "game:slots"}, {Text: "♠️ Poker", Data: "game:poker"}},
		{{Text: "« Back", Data: "back:main"}},
	}
	
	b.SendInlineKeyboard(chatID, "🎮 *Choose a Game*", games)
}

func (b *TelegramBot) showBalance(chatID int64, userID int64) {
	b.SendMessage(chatID, "💰 *Your Balance*\n\n_Loading..._")
}

func (b *TelegramBot) showDeposit(chatID int64) {
	b.SendMessage(chatID, "💵 *Deposit*\n\n_Select cryptocurrency:_")
}

func (b *TelegramBot) showWithdraw(chatID int64) {
	b.SendMessage(chatID, "💸 *Withdraw*\n\n_Select cryptocurrency:_")
}

func (b *TelegramBot) handleStatefulInput(session *UserSession, msg *Message) {
	switch session.State {
	case "awaiting_deposit_address":
		// Handle deposit address input
		b.SendMessage(msg.Chat.ID, "Address saved! Your deposit is being processed.")
		session.State = "main"
	case "awaiting_withdraw_amount":
		// Handle withdrawal amount
		b.SendMessage(msg.Chat.ID, "Withdrawal request submitted!")
		session.State = "main"
	default:
		b.showMainMenu(msg.Chat.ID, msg.From.ID)
	}
}

// TelegramBotClient provides a client for interacting with TigerCasino API
type TelegramBotClient struct {
	APIURL string
}

// NewTelegramBotClient creates a new Telegram bot client
func NewTelegramBotClient(apiURL string) *TelegramBotClient {
	return &TelegramBotClient{APIURL: apiURL}
}

// NotifyDeposit notifies a user about a deposit via Telegram
func (c *TelegramBotClient) NotifyDeposit(telegramID int64, amount float64, currency string) error {
	// In production, would use Telegram Bot API to send message
	log.Printf("Notifying user %d about deposit: %.2f %s", telegramID, amount, currency)
	return nil
}

// NotifyWithdrawal notifies a user about a withdrawal via Telegram
func (c *TelegramBotClient) NotifyWithdrawal(telegramID int64, amount float64, currency string) error {
	log.Printf("Notifying user %d about withdrawal: %.2f %s", telegramID, amount, currency)
	return nil
}

// NotifyWin notifies a user about a big win via Telegram
func (c *TelegramBotClient) NotifyWin(telegramID int64, amount float64, game string) error {
	log.Printf("Notifying user %d about win in %s: %.2f", telegramID, game, amount)
	return nil
}

// NotifyPromo notifies a user about a promotion via Telegram
func (c *TelegramBotClient) NotifyPromo(telegramID int64, promoName string) error {
	log.Printf("Notifying user %d about promo: %s", telegramID, promoName)
	return nil
}
