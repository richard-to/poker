package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/richard-to/go-poker/pkg/poker"
)

// General actions
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

// GameStage is an enum the current round of betting
type GameStage int

// Rounds of a game
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
	if e.Action == actionJoin {
		c.username = e.Params["username"].(string)
		c.send <- createOnJoinEvent(c.id, c.username)
		c.hub.broadcast <- createNewMessageEvent(
			c.id,
			systemUsername,
			fmt.Sprintf("%s joined the game.", c.username),
		)
		c.send <- createUpdateGameEvent(c.id, g)
		return
	}

	if e.Action == actionSendMessage {
		c.hub.broadcast <- createNewMessageEvent(
			c.id,
			e.Params["username"].(string),
			e.Params["message"].(string),
		)
		return
	}

	if e.Action == actionTakeSeat {
		err := c.gameState.TakeSeat(c, e.Params["seatID"].(string))
		if err != nil {
			fmt.Println(err)
			// TODO Handle error
			return
		}
		c.send <- createOnTakeSeatEvent(c.id, e.Params["seatID"].(string))

		// Try to start a new game
		if c.gameState.Stage == Waiting {
			StartGame(c.gameState)
		}

		c.hub.broadcast <- createUpdateGameEvent(c.id, g)
	}

	if e.Action == actionFold {
		err := c.gameState.FoldAction(c)
		if err != nil {
			// TODO: Handle error
			return
		}
		for {
			err = c.gameState.NextGameState(c)
			if err != nil {
				// TODO: Handle error
				return
			}
			if c.gameState.CurrentSeat.Player.Status == poker.PlayerActive && c.gameState.CurrentSeat.Player.Chips > 0 {
				break
			}
		}
		c.hub.broadcast <- createUpdateGameEvent(c.id, g)
		return
	}

	if e.Action == actionCheck {
		err := c.gameState.CheckAction(c)
		if err != nil {
			// TODO: Handle error
			return
		}
		for {
			err = c.gameState.NextGameState(c)
			if err != nil {
				// TODO: Handle error
				return
			}
			if c.gameState.CurrentSeat.Player.Status == poker.PlayerActive && c.gameState.CurrentSeat.Player.Chips > 0 {
				break
			}
		}
		c.hub.broadcast <- createUpdateGameEvent(c.id, g)
		return
	}

	if e.Action == actionCall {
		err := c.gameState.CallAction(c)
		if err != nil {
			// TODO: Handle error
			return
		}
		for {
			err = c.gameState.NextGameState(c)
			if err != nil {
				// TODO: Handle error
				return
			}
			if c.gameState.CurrentSeat.Player.Status == poker.PlayerActive && c.gameState.CurrentSeat.Player.Chips > 0 {
				break
			}
		}
		c.hub.broadcast <- createUpdateGameEvent(c.id, g)
		return
	}

	if e.Action == actionRaise {
		raiseAmount := int(e.Params["value"].(float64))
		err := c.gameState.RaiseAction(c, raiseAmount)
		if err != nil {
			// TODO: Handle error
			return
		}
		for {
			err = c.gameState.NextGameState(c)
			if err != nil {
				// TODO: Handle error
				return
			}
			if c.gameState.CurrentSeat.Player.Status == poker.PlayerActive && c.gameState.CurrentSeat.Player.Chips > 0 {
				break
			}
		}
		c.hub.broadcast <- createUpdateGameEvent(c.id, g)
		return
	}
}

// GetPlayers gets players in the game
func (g *GameState) GetPlayers() []*poker.Player {
	ps := make([]*poker.Player, g.Table.Seats.Len())
	seat := g.Table.Seats
	for i := 0; i < seat.Len(); i++ {
		ps[i] = seat.Player
		seat = seat.Next()
	}
	return ps
}

// TakeSeat takes a seat for the user
func (g *GameState) TakeSeat(c *Client, seatID string) error {
	if c.seatID != "" {
		return fmt.Errorf("You can only sit at one seat")
	}

	seat := g.Table.Seats
	for i := 0; i < seat.Len(); i++ {
		if seat.Player.ID != seatID {
			seat = seat.Next()
			continue
		}
		if seat.Player.Status > poker.PlayerVacated {
			return fmt.Errorf("Seat has already been taken")
		}

		seat.Player.Name = c.username
		seat.Player.Chips = defaultChips
		seat.Player.Status = poker.PlayerSittingOut
		c.seatID = seat.Player.ID

		return nil
	}

	return fmt.Errorf("Invalid seat chosen")
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

// NextGameState gets the next game state
func (g *GameState) NextGameState(c *Client) error {
	var err error

	g.CurrentSeat, err = poker.GetNextActiveSeat(g.CurrentSeat)
	if err != nil {
		return err
	}

	everyoneHasFolded := true
	nextSeat := g.CurrentSeat.Next()
	for i := 0; i < nextSeat.Len()-1; i++ {
		if nextSeat.Player.Status == poker.PlayerActive && nextSeat.Player.HasFolded == false {
			everyoneHasFolded = false
			break
		}
		nextSeat = nextSeat.Next()
	}

	if everyoneHasFolded {
		c.hub.broadcast <- createNewMessageEvent(
			c.id,
			systemUsername,
			fmt.Sprintf("%s won the hand.", g.CurrentSeat.Player.Name),
		)
		g.CurrentSeat.Player.Chips += g.Table.Pot.GetTotal()
		StartGame(g)
		c.hub.broadcast <- createNewMessageEvent(c.id, systemUsername, "Starting new hand.")
		return nil
	}

	everyoneAllInOrFolded := true
	nextSeat = g.CurrentSeat.Next()
	for i := 0; i < nextSeat.Len()-1; i++ {
		if nextSeat.Player.Status == poker.PlayerActive && (nextSeat.Player.HasFolded == false && nextSeat.Player.Chips > 0) {
			everyoneAllInOrFolded = false
			break
		}
		nextSeat = nextSeat.Next()
	}

	if everyoneAllInOrFolded {
		poker.DealFlop(&g.Deck, &g.Table)
		poker.DealTurn(&g.Deck, &g.Table)
		poker.DealRiver(&g.Deck, &g.Table)
		subPots := g.Table.Pot.GetSidePots()
		for _, subPot := range subPots {
			winningHands := poker.FindWinningHands(subPot.Players, &g.Table)
			chipsWon := subPot.Total / len(winningHands) // TODO: Handle remainder
			for _, ph := range winningHands {
				c.hub.broadcast <- createNewMessageEvent(
					c.id,
					systemUsername,
					fmt.Sprintf("%s won the hand.", ph.Player.Name), // TODO: Handle sub pot win
				)
				ph.Player.Chips += chipsWon
			}
		}
		StartGame(g)
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
			subPots := g.Table.Pot.GetSidePots()
			numSubPots := len(subPots)
			for i, subPot := range subPots {
				winningHands := poker.FindWinningHands(subPot.Players, &g.Table)
				chipsWon := subPot.Total / len(winningHands)
				remainderChipsWon := subPot.Total % len(winningHands)

				potText := "main pot"
				if i > 0 {
					if numSubPots == 2 {
						potText = "side pot"
					} else {
						potText = fmt.Sprintf("side pot %d", i)
					}
				}
				for _, ph := range winningHands {
					playerChipsWon := chipsWon
					if remainderChipsWon > 0 {
						remainderChipsWon--
						playerChipsWon++
					}
					c.hub.broadcast <- createNewMessageEvent(
						c.id,
						systemUsername,
						fmt.Sprintf(
							"%s wins ℝ%d %s with %s.",
							ph.Player.Name,
							playerChipsWon,
							potText,
							strings.ToLower(ph.Hand.Rank.String()),
						),
					)
					ph.Player.Chips += chipsWon
				}
			}
			StartGame(g)
			c.hub.broadcast <- createNewMessageEvent(c.id, systemUsername, "Starting new hand.")
		} else {
			// TODO: Error
		}
	}

	return nil
}

// FoldAction folds
func (g *GameState) FoldAction(c *Client) error {
	err := g.CurrentSeat.Player.Fold()
	if err != nil {
		return err
	}
	c.hub.broadcast <- createNewMessageEvent(
		c.id,
		systemUsername,
		fmt.Sprintf("%s folds.", g.CurrentSeat.Player.Name),
	)
	return nil
}

// CheckAction checks
func (g *GameState) CheckAction(c *Client) error {
	err := g.CurrentSeat.Player.Check(g.BettingRound)
	if err != nil {
		return err
	}
	c.hub.broadcast <- createNewMessageEvent(
		c.id,
		systemUsername,
		fmt.Sprintf("%s checks.", g.CurrentSeat.Player.Name),
	)
	return nil
}

// CallAction calls
func (g *GameState) CallAction(c *Client) error {
	err := g.CurrentSeat.Player.Call(&g.Table, g.BettingRound)
	if err != nil {
		return err
	}
	c.hub.broadcast <- createNewMessageEvent(
		c.id,
		systemUsername,
		fmt.Sprintf("%s calls.", g.CurrentSeat.Player.Name),
	)
	return nil
}

// RaiseAction raises
func (g *GameState) RaiseAction(c *Client, raiseAmount int) error {
	err := g.CurrentSeat.Player.Raise(&g.Table, g.BettingRound, raiseAmount)
	if err != nil {
		return err
	}
	c.hub.broadcast <- createNewMessageEvent(
		c.id,
		systemUsername,
		fmt.Sprintf("%s raises to %d.", g.CurrentSeat.Player.Name, raiseAmount),
	)
	return nil
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

// StartGame starts a new game
func StartGame(g *GameState) {
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
		// Change active player status to sitting out.
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
