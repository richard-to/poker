package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/richard-to/go-poker/pkg/poker"
)

// General actions
const actionError string = "error"
const actionOnJoin string = "on-join"
const actionOnTakeSeat string = "on-take-seat"
const actionJoin string = "join"
const actionNewMessage string = "new-message"
const actionSendMessage string = "send-message"
const actionTakeSeat string = "take-seat"

// Game actions
const actionBet string = "bet"
const actionCall string = "call"
const actionCheck string = "check"
const actionFold string = "fold"
const actionRaise string = "raise"
const actionUpdateGame string = "update-game"

// Table settings
const defaultChips int = 100
const defaultMinBet int = 2
const minPlayers int = 2
const numPlayers int = 6

const systemUsername string = "System"

// GameStage is an enum for the current round of betting
type GameStage int

// Stages of a game
const (
	Waiting GameStage = iota
	Preflop
	Flop
	Turn
	River
	Showdown
)

func (g GameStage) String() string {
	return [...]string{"Waiting", "Preflop", "Flop", "Turn", "River", "Showdown"}[g]
}

// Event is a JSON message in the game loop.
type Event struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params"`
}

// UserEvent is an event associated with a user
type UserEvent struct {
	UserID string
	Event
}

// GameState is the current state of the poker game
type GameState struct {
	BettingRound *poker.BettingRound
	CurrentSeat  *poker.Seat
	Deck         poker.Deck
	Stage        GameStage
	Table        poker.Table
}

// ProcessEvent process event
func (g *GameState) ProcessEvent(c *Client, e Event) {
	var err error

	if e.Action == actionJoin {
		err = g.HandleJoin(c, e.Params["username"].(string))
	} else if e.Action == actionSendMessage {
		err = g.HandleSendMessage(c, e.Params["username"].(string), e.Params["message"].(string))
	} else if e.Action == actionTakeSeat {
		err = g.HandleTakeSeat(c, e.Params["seatID"].(string))
	} else if e.Action == actionFold {
		err = g.HandleFold(c)
	} else if e.Action == actionCheck {
		err = g.HandleCheck(c)
	} else if e.Action == actionCall {
		err = g.HandleCall(c)
	} else if e.Action == actionRaise {
		raiseAmount := int(e.Params["value"].(float64))
		err = g.HandleRaise(c, raiseAmount)
	} else {
		err = fmt.Errorf("Unknown action encountered: %s", e.Action)
	}

	if err != nil {
		g.HandlePlayerError(c, err)
	}
}

// HandlePlayerError handles player error (not system error)
func (g *GameState) HandlePlayerError(c *Client, err error) error {
	c.send <- createErrorEvent(c.id, err)
	return nil
}

// HandleJoin handles join event
func (g *GameState) HandleJoin(c *Client, username string) error {
	c.username = username
	c.send <- createOnJoinEvent(c.id, c.username)
	c.hub.broadcast <- createNewMessageEvent(
		c.id,
		systemUsername,
		fmt.Sprintf("%s joined the game.", c.username),
	)
	c.send <- createUpdateGameEvent(c.id, g)
	return nil
}

// HandleSendMessage handles send message event
func (g *GameState) HandleSendMessage(c *Client, username string, message string) error {
	c.hub.broadcast <- createNewMessageEvent(c.id, username, message)
	return nil
}

// HandleTakeSeat takes a seat for the user
func (g *GameState) HandleTakeSeat(c *Client, seatID string) error {
	if c.seatID != "" {
		return fmt.Errorf("You can only sit at one seat")
	}

	seat := g.Table.Seats

	var selectedPlayer *poker.Player
	for i := 0; i < seat.Len(); i++ {
		if seat.Player.ID == seatID {
			// It's possible that two players picked the same seat at the same time
			if seat.Player.Status > poker.PlayerVacated {
				return fmt.Errorf("Seat has already been taken")
			}
			selectedPlayer = seat.Player
			break
		}
		seat = seat.Next()
	}

	if selectedPlayer == nil {
		return fmt.Errorf("Invalid seat chosen")
	}

	// Link user with player seat
	selectedPlayer.Name = c.username
	selectedPlayer.Chips = defaultChips
	selectedPlayer.Status = poker.PlayerSittingOut
	c.seatID = selectedPlayer.ID

	c.send <- createOnTakeSeatEvent(c.id, seatID)

	// Try to start a new game if one hasn't started yet.
	if g.Stage == Waiting {
		g.StartNewHand()
	}

	c.hub.broadcast <- createUpdateGameEvent(c.id, g)
	return nil
}

// HandleFold folds
func (g *GameState) HandleFold(c *Client) error {
	err := g.CurrentSeat.Player.Fold()
	if err != nil {
		return err
	}
	c.hub.broadcast <- createNewMessageEvent(
		c.id,
		systemUsername,
		fmt.Sprintf("%s folds.", g.CurrentSeat.Player.Name),
	)
	return g.GoToNextGameState(c)
}

// HandleCheck checks
func (g *GameState) HandleCheck(c *Client) error {
	err := g.CurrentSeat.Player.Check(g.BettingRound)
	if err != nil {
		return err
	}
	c.hub.broadcast <- createNewMessageEvent(
		c.id,
		systemUsername,
		fmt.Sprintf("%s checks.", g.CurrentSeat.Player.Name),
	)
	return g.GoToNextGameState(c)
}

// HandleCall calls
func (g *GameState) HandleCall(c *Client) error {
	err := g.CurrentSeat.Player.Call(&g.Table, g.BettingRound)
	if err != nil {
		return err
	}
	c.hub.broadcast <- createNewMessageEvent(
		c.id,
		systemUsername,
		fmt.Sprintf("%s calls.", g.CurrentSeat.Player.Name),
	)
	return g.GoToNextGameState(c)
}

// HandleRaise raises/bets
func (g *GameState) HandleRaise(c *Client, raiseAmount int) error {
	err := g.CurrentSeat.Player.Raise(&g.Table, g.BettingRound, raiseAmount)
	if err != nil {
		return err
	}
	c.hub.broadcast <- createNewMessageEvent(
		c.id,
		systemUsername,
		fmt.Sprintf("%s raises to %d.", g.CurrentSeat.Player.Name, raiseAmount),
	)
	return g.GoToNextGameState(c)
}

// GoToNextGameState moves to the next game state
func (g *GameState) GoToNextGameState(c *Client) error {
	for {
		err := c.gameState.NextGameState(c)
		if err != nil {
			return err
		}
		if c.gameState.CurrentSeat.Player.Status == poker.PlayerActive && c.gameState.CurrentSeat.Player.Chips > 0 {
			break
		}
	}
	c.hub.broadcast <- createUpdateGameEvent(c.id, g)
	return nil
}

// NextGameState gets the next game state
func (g *GameState) NextGameState(c *Client) error {
	var err error

	g.CurrentSeat, err = poker.GetNextActiveSeat(g.CurrentSeat)
	if err != nil {
		return err
	}

	everyoneHasFolded := poker.HasEveryoneFolded(g.CurrentSeat)

	if everyoneHasFolded {
		c.hub.broadcast <- createNewMessageEvent(
			c.id,
			systemUsername,
			fmt.Sprintf("%s won the hand.", g.CurrentSeat.Player.Name),
		)
		g.Table.AwardPot(g.CurrentSeat.Player)
		g.StartNewHand()
		c.hub.broadcast <- createNewMessageEvent(c.id, systemUsername, "Starting new hand.")
		return nil
	}

	everyoneAllInOrFolded := poker.HasEveryoneFoldedOrIsAllIn(g.CurrentSeat)

	if everyoneAllInOrFolded {
		poker.DealFlop(&g.Deck, &g.Table)
		poker.DealTurn(&g.Deck, &g.Table)
		poker.DealRiver(&g.Deck, &g.Table)
		g.DetermineWinners(c)
		g.StartNewHand()
		c.hub.broadcast <- createNewMessageEvent(c.id, systemUsername, "Starting new hand.")
		return nil
	}

	if g.CurrentSeat.Player == g.BettingRound.Raiser {
		g.Stage++
		if g.Stage == Flop {
			poker.DealFlop(&g.Deck, &g.Table)
			g.CurrentSeat, err = poker.GetNextActiveSeat(g.Table.Dealer)
			if err != nil {
				return err
			}
			g.BettingRound, err = poker.NewBettingRound(g.CurrentSeat, 0, g.Table.MinBet)
			if err != nil {
				return err
			}
			c.hub.broadcast <- createNewMessageEvent(c.id, systemUsername, "Dealing flop.")
		} else if g.Stage == Turn {
			poker.DealTurn(&g.Deck, &g.Table)
			g.CurrentSeat, err = poker.GetNextActiveSeat(g.Table.Dealer)
			if err != nil {
				return err
			}
			g.BettingRound, err = poker.NewBettingRound(g.CurrentSeat, 0, g.Table.MinBet)
			if err != nil {
				return err
			}
			c.hub.broadcast <- createNewMessageEvent(c.id, systemUsername, "Dealing turn.")
		} else if g.Stage == River {
			poker.DealRiver(&g.Deck, &g.Table)
			g.CurrentSeat, err = poker.GetNextActiveSeat(g.Table.Dealer)
			if err != nil {
				return err
			}
			g.BettingRound, err = poker.NewBettingRound(g.CurrentSeat, 0, g.Table.MinBet)
			if err != nil {
				return err
			}
			c.hub.broadcast <- createNewMessageEvent(c.id, systemUsername, "Dealing river.")
		} else if g.Stage == Showdown {
			g.DetermineWinners(c)
			g.StartNewHand()
			c.hub.broadcast <- createNewMessageEvent(c.id, systemUsername, "Starting new hand.")
		} else {
			return fmt.Errorf("Invalid game stage encountered: %s", g.Stage.String())
		}
	}

	return nil
}

// DetermineWinners determines who won the hand and awards chips to the winner
func (g *GameState) DetermineWinners(c *Client) {
	allWinningHands := g.Table.DetermineWinners()
	for i, winningHandsByPot := range allWinningHands {
		potText := "main pot"
		if i > 0 {
			if len(allWinningHands) == 2 {
				potText = "side pot"
			} else {
				potText = fmt.Sprintf("side pot %d", i)
			}
		}
		for _, ph := range winningHandsByPot {
			c.hub.broadcast <- createNewMessageEvent(
				c.id,
				systemUsername,
				fmt.Sprintf(
					"%s wins ℝ%d %s with %s.",
					ph.Player.Name,
					ph.ChipsWon,
					potText,
					strings.ToLower(ph.Hand.Rank.String()),
				),
			)
		}
	}
}

// NewGameState creates a new game state
func NewGameState() *GameState {
	// Initialize vacated seats
	seats := poker.NewSeat(numPlayers)
	for i := 0; i < seats.Len(); i++ {
		seats.Player = &poker.Player{
			ID:     strconv.Itoa(i + 1),
			Status: poker.PlayerVacated,
		}
		seats = seats.Next()
	}

	return &GameState{
		BettingRound: nil,
		CurrentSeat:  seats,
		Deck:         poker.NewDeck(),
		Table: poker.Table{
			MinBet: defaultMinBet,
			Pot:    poker.NewPot(),
			Seats:  seats,
		},
		Stage: Waiting,
	}
}

// StartNewHand starts a new hand
func (g *GameState) StartNewHand() {
	deck := poker.NewDeck()
	seats := g.Table.Seats

	// Reset player hands
	for i := 0; i < seats.Len(); i++ {
		seats.Player.HoleCards = [2]*poker.Card{}
		seats.Player.HasFolded = false
		seats = seats.Next()
	}

	// Get active players for the next game
	for i := 0; i < seats.Len(); i++ {
		if seats.Player.Status > poker.PlayerVacated {
			if seats.Player.Chips == 0 {
				seats.Player.Status = poker.PlayerSittingOut
			} else {
				seats.Player.Status = poker.PlayerActive
			}
		}
		seats = seats.Next()
	}

	activePlayerCount := poker.CountSeatsByPlayerStatus(seats, poker.PlayerActive)

	if activePlayerCount < minPlayers {
		// Change active player status to sitting out if we don't have enough players
		for i := 0; i < seats.Len(); i++ {
			if seats.Player.Status == poker.PlayerActive {
				seats.Player.Status = poker.PlayerSittingOut
			}
			seats = seats.Next()
		}
		return
	}

	if g.Table.Dealer == nil {
		g.Table.Dealer = g.Table.Seats
	}

	dealer, err := poker.GetNextActiveSeat(g.Table.Dealer)
	if err != nil {
		panic(err)
	}

	smallBlind, err := poker.GetNextActiveSeat(dealer)
	if err != nil {
		panic(err)
	}

	// In a head to head match, the dealer is the small blind
	if activePlayerCount == 2 {
		smallBlind = dealer
	}

	bigBlind, err := poker.GetNextActiveSeat(smallBlind)
	if err != nil {
		panic(err)
	}

	table := poker.Table{
		BigBlind:   bigBlind,
		Dealer:     dealer,
		MinBet:     defaultMinBet,
		Pot:        poker.NewPot(),
		Seats:      seats,
		SmallBlind: smallBlind,
	}

	poker.DealHands(&deck, &table)

	currentSeat, err := poker.GetNextActiveSeat(bigBlind)
	if err != nil {
		panic(err)
	}

	preflopRound, err := poker.NewBettingRound(currentSeat, table.MinBet, table.MinBet)
	if err != nil {
		panic(err)
	}

	poker.TakeSmallBlind(&table, preflopRound)
	poker.TakeBigBlind(&table, preflopRound)

	g.BettingRound = preflopRound
	g.CurrentSeat = currentSeat
	g.Deck = deck
	g.Stage = Preflop
	g.Table = table
}

func createOnJoinEvent(userID string, username string) UserEvent {
	return UserEvent{
		UserID: userID,
		Event: Event{
			Action: actionOnJoin,
			Params: map[string]interface{}{
				"userID":   userID,
				"username": username,
			},
		},
	}
}

func createOnTakeSeatEvent(userID string, seatID string) UserEvent {
	return UserEvent{
		UserID: userID,
		Event: Event{
			Action: actionOnTakeSeat,
			Params: map[string]interface{}{
				"seatID": seatID,
			},
		},
	}
}

func createUpdateGameEvent(userID string, g *GameState) UserEvent {
	var actionBar map[string]interface{}

	players := make([]map[string]interface{}, 0)
	seats := g.Table.Seats

	if g.Stage == Waiting {
		// Players data
		for i := 0; i < seats.Len(); i++ {
			players = append(players, map[string]interface{}{
				"chips":      seats.Player.Chips,
				"chipsInPot": nil,
				"hasFolded":  seats.Player.HasFolded,
				"holeCards":  seats.Player.HoleCards,
				"id":         seats.Player.ID,
				"isActive":   false,
				"isDealer":   false,
				"name":       seats.Player.Name,
				"status":     seats.Player.Status.String(),
			})
			seats = seats.Next()
		}

		// Actions data
		actionBar = map[string]interface{}{
			"actions":        []string{},
			"callAmount":     0,
			"chipsInPot":     0,
			"maxRaiseAmount": 0,
			"minRaiseAmount": 0,
			"totalChips":     0,
		}
	} else {
		// Players data
		activePlayer := g.CurrentSeat.Player
		for i := 0; i < seats.Len(); i++ {
			players = append(players, map[string]interface{}{
				"chips":      seats.Player.Chips,
				"chipsInPot": g.BettingRound.Bets[seats.Player.ID],
				"hasFolded":  seats.Player.HasFolded,
				"holeCards":  seats.Player.HoleCards,
				"id":         seats.Player.ID,
				"isActive":   seats.Player.ID == activePlayer.ID,
				"isDealer":   seats.Player.ID == g.Table.Dealer.Player.ID,
				"name":       seats.Player.Name,
				"status":     seats.Player.Status.String(),
			})
			seats = seats.Next()
		}

		// Actions data
		callRemainingAmount := g.BettingRound.CallAmount - g.BettingRound.Bets[activePlayer.ID]
		actionBar = map[string]interface{}{
			"actions":        g.GetActions(),
			"callAmount":     g.BettingRound.CallAmount,
			"chipsInPot":     g.BettingRound.Bets[activePlayer.ID],
			"maxRaiseAmount": activePlayer.Chips - callRemainingAmount,
			"minRaiseAmount": g.BettingRound.RaiseByAmount,
			"totalChips":     activePlayer.Chips,
		}
	}

	// Table data
	table := map[string]interface{}{
		"flop":  g.Table.Flop,
		"pot":   g.Table.Pot.GetTotal(),
		"river": g.Table.River,
		"turn":  g.Table.Turn,
	}

	return UserEvent{
		UserID: userID,
		Event: Event{
			Action: actionUpdateGame,
			Params: map[string]interface{}{
				"actionBar": actionBar,
				"players":   players,
				"stage":     g.Stage.String(),
				"table":     table,
			},
		},
	}
}

// GetActions gets the actions available to active player
func (g *GameState) GetActions() []string {
	var actions []string

	if g.CurrentSeat.Player.CanFold(g.BettingRound) {
		actions = append(actions, actionFold)
	}

	if g.CurrentSeat.Player.CanCheck(g.BettingRound) {
		actions = append(actions, actionCheck)
	}

	if g.CurrentSeat.Player.CanCall(g.BettingRound) {
		actions = append(actions, actionCall)
	}

	if g.CurrentSeat.Player.CanRaise(g.BettingRound) {
		actions = append(actions, actionRaise)
	}
	return actions
}

func createNewMessageEvent(userID string, username string, message string) UserEvent {
	return UserEvent{
		UserID: userID,
		Event: Event{
			Action: actionNewMessage,
			Params: map[string]interface{}{
				"id":       uuid.New().String(),
				"message":  message,
				"username": username,
			},
		},
	}
}

func createErrorEvent(userID string, err error) UserEvent {
	return UserEvent{
		UserID: userID,
		Event: Event{
			Action: actionError,
			Params: map[string]interface{}{
				"error": err,
			},
		},
	}
}
