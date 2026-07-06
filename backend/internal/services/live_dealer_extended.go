package services

// Extended Live Dealer Server - 200+ additional tables
type ExtendedLiveDealer struct {
	tables map[string]*ExtendedTable
}

type ExtendedTable struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	GameType  string `json:"game_type"`
	Status    string `json:"status"`
	Players   int    `json:"players"`
	MinBet    float64 `json:"min_bet"`
	MaxBet    float64 `json:"max_bet"`
	Language  string `json:"language"`
	Dealer    string `json:"dealer"`
}

func NewExtendedLiveDealer() *ExtendedLiveDealer {
	s := &ExtendedLiveDealer{tables: make(map[string]*ExtendedTable)}
	s.initializeTables()
	return s
}

func (s *ExtendedLiveDealer) initializeTables() {
	providers := []string{"Evolution", "PragmaticPlay", "Ezugi", "Authentic Gaming", "Vivo Gaming", "AsiaGaming", " SAGaming", "SexyGaming", "TVBet", "7Mojos"}
	languages := []string{"English", "Spanish", "German", "French", "Italian", "Russian", "Turkish", "Chinese", "Japanese", "Korean", "Portuguese", "Hindi", "Arabic", "Thai", "Vietnamese"}
	gameTypes := []string{"Blackjack", "Roulette", "Baccarat", "Poker", "GameShows", "Dice", "Keno", "AndarBahar", "TeenPatti", "War"}
	
	// Generate 200 more tables
	for i := 0; i < 200; i++ {
		provider := providers[i%len(providers)]
		lang := languages[i%len(languages)]
		gameType := gameTypes[i%len(gameTypes)]
		
		s.tables[fmt.Sprintf("ext_live_%d", i)] = &ExtendedTable{
			ID:       fmt.Sprintf("ext_live_%d", i),
			Name:     fmt.Sprintf("%s Premium %d", gameType, i+1),
			Provider: provider,
			GameType: gameType,
			Status:   "online",
			Players:  i % 75,
			MinBet:   0.5 + float64(i%10)*0.5,
			MaxBet:   50000.0 + float64(i%100)*500.0,
			Language: lang,
			Dealer:   fmt.Sprintf("Dealer_%d", i+1),
		}
	}
}

func (s *ExtendedLiveDealer) GetAllTables() []*ExtendedTable {
	tables := make([]*ExtendedTable, 0, len(s.tables))
	for _, t := range s.tables { tables = append(tables, t) }
	return tables
}

func (s *ExtendedLiveDealer) GetTableCount() int { return len(s.tables) }
