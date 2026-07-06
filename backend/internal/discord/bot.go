package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// DiscordBot represents a Discord bot for TigerCasino
type DiscordBot struct {
	Token      string
	ChannelID  string
	WebhookURL string
	Client     *http.Client
	Commands   map[string]CommandHandler
	Mutex      sync.RWMutex
}

// CommandHandler defines the function signature for command handlers
type CommandHandler func(interaction *Interaction) *Embed

// Discord API endpoints
const (
	DiscordAPI = "https://discord.com/api/v10"
)

// NewDiscordBot creates a new Discord bot instance
func NewDiscordBot(token, channelID string) *DiscordBot {
	bot := &DiscordBot{
		Token:     token,
		ChannelID: channelID,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		Commands: make(map[string]CommandHandler),
	}
	
	// Register slash commands
	bot.registerCommands()
	
	return bot
}

// StartWebhookServer starts the webhook server for Discord interactions
func (b *DiscordBot) StartWebhookServer(port string) {
	http.HandleFunc("/webhook", b.handleWebhook)
	log.Printf("Discord webhook server started on port %s", port)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// SetWebhook sets up the Discord webhook
func (b *DiscordBot) SetWebhook(webhookURL string) error {
	b.WebhookURL = webhookURL
	return nil
}

// Request makes an API request to Discord
func (b *DiscordBot) Request(method, endpoint string, body []byte) ([]byte, error) {
	url := DiscordAPI + endpoint
	
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bot "+b.Token)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := json.NewDecoder(resp.Body).Decode(nil)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}
	
	return respBody, nil
}

// ============ Message Sending ============

// SendMessage sends a message to the Discord channel
func (b *DiscordBot) SendMessage(content string, embeds ...*Embed) error {
	if b.WebhookURL == "" {
		return fmt.Errorf("webhook URL not set")
	}
	
	payload := map[string]interface{}{
		"content": content,
	}
	
	if len(embeds) > 0 {
		var embedMaps []map[string]interface{}
		for _, e := range embeds {
			embedMaps = append(embedMaps, e.ToMap())
		}
		payload["embeds"] = embedMaps
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", b.WebhookURL, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := b.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord API error: %d", resp.StatusCode)
	}
	
	return nil
}

// SendEmbed sends a rich embed to the Discord channel
func (b *DiscordBot) SendEmbed(embed *Embed) error {
	return b.SendMessage("", embed)
}

// SendGameUpdate sends a game update to Discord
func (b *DiscordBot) SendGameUpdate(gameName, status string, multiplier float64) error {
	embed := NewEmbed().
		SetTitle("🎰 Game Update").
		AddField("Game", gameName).
		AddField("Status", status).
		SetColor(0xFF6B35)
	
	if multiplier > 0 {
		embed.AddField("Multiplier", fmt.Sprintf("%.2fx", multiplier))
	}
	
	return b.SendEmbed(embed)
}

// SendBigWin sends a big win notification
func (b *DiscordBot) SendBigWin(username string, amount float64, game string) error {
	embed := NewEmbed().
		SetTitle("🎉 BIG WIN!").
		AddField("Player", username).
		AddField("Amount", fmt.Sprintf("$%.2f", amount)).
		AddField("Game", game).
		SetColor(0xFFD700).
		SetFooter("TigerCasino")
	
	return b.SendEmbed(embed)
}

// SendJackpot sends a jackpot notification
func (b *DiscordBot) SendJackpot(username, jackpotType string, amount float64) error {
	embed := NewEmbed().
		SetTitle("💎 JACKPOT!").
		AddField("Player", username).
		AddField("Jackpot", jackpotType).
		AddField("Amount", fmt.Sprintf("$%.2f", amount)).
		SetColor(0xFFD700).
		SetFooter("TigerCasino")
	
	return b.SendEmbed(embed)
}

// SendPromo sends a promotion notification
func (b *DiscordBot) SendPromo(title, description string) error {
	embed := NewEmbed().
		SetTitle("🎁 "+title).
		SetDescription(description).
		SetColor(0x00D26A).
		SetFooter("TigerCasino")
	
	return b.SendEmbed(embed)
}

// SendLeaderboardUpdate sends a leaderboard update
func (b *DiscordBot) SendLeaderboardUpdate(entries []LeaderboardEntry) error {
	desc := "🏆 **Top Players Today**\n\n"
	for i, e := range entries {
		emoji := ""
		switch i {
		case 0: emoji = "🥇"
		case 1: emoji = "🥈"
		case 2: emoji = "🥉"
		default: emoji = fmt.Sprintf("%d.", i+1)
		}
		desc += fmt.Sprintf("%s **%s** - $%.2f\n", emoji, e.Username, e.Profit)
	}
	
	embed := NewEmbed().
		SetTitle("🏆 Leaderboard Update").
		SetDescription(desc).
		SetColor(0xFFD700)
	
	return b.SendEmbed(embed)
}

// ============ Webhook Handling ============

// handleWebhook handles incoming Discord webhooks
func (b *DiscordBot) handleWebhook(w http.ResponseWriter, r *http.Request) {
	var interaction Interaction
	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// Handle ping
	if interaction.Type == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"type": 1}`))
		return
	}
	
	// Handle command
	if interaction.Data.Name != "" {
		if handler, ok := b.Commands[interaction.Data.Name]; ok {
			embed := handler(&interaction)
			if embed != nil {
				response := InteractionResponse{
					Type: 4,
					Data: &InteractionResponseData{
						Embeds: []map[string]interface{}{embed.ToMap()},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}
	
	w.WriteHeader(http.StatusOK)
}

// registerCommands registers Discord slash commands
func (b *DiscordBot) registerCommands() {
	// Register commands via API
	commands := []map[string]interface{}{
		{
			"name": "balance",
			"description": "Check your casino balance",
		},
		{
			"name": "deposit",
			"description": "Get deposit address",
		},
		{
			"name": "games",
			"description": "List available games",
		},
		{
			"name": "vip",
			"description": "Check VIP status",
		},
		{
			"name": "help",
			"description": "Get help and support",
		},
	}
	
	// In production, would register with Discord API
	_ = commands
}

// ============ Embed Builder ============

// Embed represents a Discord embed
type Embed struct {
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      *EmbedFooter  `json:"footer,omitempty"`
	Image       *EmbedImage   `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Author      *EmbedAuthor  `json:"author,omitempty"`
	Fields      []EmbedField  `json:"fields,omitempty"`
}

// EmbedField represents a field in an embed
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// EmbedFooter represents the footer of an embed
type EmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// EmbedImage represents an image in an embed
type EmbedImage struct {
	URL string `json:"url"`
}

// EmbedThumbnail represents a thumbnail in an embed
type EmbedThumbnail struct {
	URL string `json:"url"`
}

// EmbedAuthor represents the author of an embed
type EmbedAuthor struct {
	Name    string `json:"name"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// NewEmbed creates a new embed
func NewEmbed() *Embed {
	return &Embed{
		Fields: make([]EmbedField, 0),
	}
}

// SetTitle sets the embed title
func (e *Embed) SetTitle(title string) *Embed {
	e.Title = title
	return e
}

// SetDescription sets the embed description
func (e *Embed) SetDescription(desc string) *Embed {
	e.Description = desc
	return e
}

// SetURL sets the embed URL
func (e *Embed) SetURL(url string) *Embed {
	e.URL = url
	return e
}

// SetColor sets the embed color
func (e *Embed) SetColor(color int) *Embed {
	e.Color = color
	return e
}

// SetFooter sets the embed footer
func (e *Embed) SetFooter(text string) *Embed {
	e.Footer = &EmbedFooter{Text: text}
	return e
}

// AddField adds a field to the embed
func (e *Embed) AddField(name, value string) *Embed {
	e.Fields = append(e.Fields, EmbedField{
		Name:   name,
		Value:  value,
		Inline: false,
	})
	return e
}

// AddInlineField adds an inline field to the embed
func (e *Embed) AddInlineField(name, value string) *Embed {
	e.Fields = append(e.Fields, EmbedField{
		Name:   name,
		Value:  value,
		Inline: true,
	})
	return e
}

// ToMap converts the embed to a map
func (e *Embed) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	
	if e.Title != "" {
		m["title"] = e.Title
	}
	if e.Description != "" {
		m["description"] = e.Description
	}
	if e.URL != "" {
		m["url"] = e.URL
	}
	if e.Color != 0 {
		m["color"] = e.Color
	}
	if e.Footer != nil {
		m["footer"] = map[string]string{"text": e.Footer.Text}
	}
	if len(e.Fields) > 0 {
		var fields []map[string]interface{}
		for _, f := range e.Fields {
			fields = append(fields, map[string]interface{}{
				"name":   f.Name,
				"value":  f.Value,
				"inline": f.Inline,
			})
		}
		m["fields"] = fields
	}
	
	return m
}

// ============ Discord API Types ============

// Interaction represents a Discord interaction
type Interaction struct {
	ID        string      `json:"id"`
	Type      int         `json:"type"`
	Data      CommandData `json:"data"`
	Member    *Member     `json:"member"`
	ChannelID string      `json:"channel_id"`
}

// CommandData represents command data
type CommandData struct {
	Name string `json:"name"`
}

// Member represents a Discord member
type Member struct {
	User         User   `json:"user"`
	Nick        string `json:"nick"`
	Roles       []string `json:"roles"`
	JoinedAt    string `json:"joined_at"`
}

// User represents a Discord user
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Discriminator string `json:"discriminator"`
}

// InteractionResponse represents a Discord interaction response
type InteractionResponse struct {
	Type int                  `json:"type"`
	Data *InteractionResponseData `json:"data"`
}

// InteractionResponseData represents interaction response data
type InteractionResponseData struct {
	Content string               `json:"content,omitempty"`
	Embeds []map[string]interface{} `json:"embeds,omitempty"`
	Flags  int                  `json:"flags,omitempty"`
}

// LeaderboardEntry represents a leaderboard entry
type LeaderboardEntry struct {
	Username string  `json:"username"`
	Profit   float64 `json:"profit"`
	Rank     int     `json:"rank"`
}

// ============ Discord Webhook Client ============

// WebhookClient provides methods for Discord webhooks
type WebhookClient struct {
	WebhookURL string
	Client    *http.Client
}

// NewWebhookClient creates a new webhook client
func NewWebhookURL(webhookURL string) *WebhookClient {
	return &WebhookClient{
		WebhookURL: webhookURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send sends a message via webhook
func (c *WebhookClient) Send(content string, embeds ...*Embed) error {
	payload := map[string]interface{}{
		"content": content,
	}
	
	if len(embeds) > 0 {
		var embedMaps []map[string]interface{}
		for _, e := range embeds {
			embedMaps = append(embedMaps, e.ToMap())
		}
		payload["embeds"] = embedMaps
	}
	
	jsonData, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", c.WebhookURL, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("webhook error: %d", resp.StatusCode)
	}
	
	return nil
}

// ============ Notification Functions ============

// SendBigWinNotification sends a big win notification to Discord
func SendBigWinNotification(webhookURL, username string, amount float64, game string) error {
	client := NewWebhookURL(webhookURL)
	
	embed := NewEmbed().
		SetTitle("🎉 BIG WIN!").
		AddField("Player", username).
		AddField("Amount", fmt.Sprintf("$%.2f", amount)).
		AddField("Game", game).
		SetColor(0xFFD700)
	
	return client.Send("", embed)
}

// SendPromoNotification sends a promo notification to Discord
func SendPromoNotification(webhookURL, title, description string) error {
	client := NewWebhookURL(webhookURL)
	
	embed := NewEmbed().
		SetTitle("🎁 "+title).
		SetDescription(description).
		SetColor(0x00D26A)
	
	return client.Send("", embed)
}

// SendLeaderboardNotification sends leaderboard to Discord
func SendLeaderboardNotification(webhookURL string, entries []LeaderboardEntry) error {
	client := NewWebhookURL(webhookURL)
	
	desc := "🏆 **Today's Top Players**\n\n"
	for i, e := range entries {
		if i >= 10 {
			break
		}
		emoji := ""
		switch i {
		case 0: emoji = "🥇"
		case 1: emoji = "🥈"
		case 2: emoji = "🥉"
		default: emoji = fmt.Sprintf("%d.", i+1)
		}
		desc += fmt.Sprintf("%s **%s** - $%.2f\n", emoji, e.Username, e.Profit)
	}
	
	embed := NewEmbed().
		SetTitle("🏆 Leaderboard").
		SetDescription(desc).
		SetColor(0xFFD700)
	
	return client.Send("", embed)
}

// InitDiscordEnv initializes Discord from environment variables
func InitDiscordEnv() (string, string) {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	return botToken, webhookURL
}

// CheckMentions checks if a message contains dangerous mentions
func CheckMentions(content string) bool {
	dangerous := []string{"@everyone", "@here", "@all"}
	lower := strings.ToLower(content)
	for _, d := range dangerous {
		if strings.Contains(lower, d) {
			return true
		}
	}
	return false
}
