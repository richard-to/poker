package poker

import (
	"container/ring"
	"fmt"
	"sort"
)

// Table represents the state of the current poker game.
type Table struct {
	Seats      *Seat // Start at first seat
	Dealer     *Seat // Start at dealer seat
	SmallBlind *Seat // Start at small blind seat
	BigBlind   *Seat // Start at big blind seat
	MinBet     int
	Pot        *Pot
	Flop       [3]*Card
	Turn       *Card
	River      *Card
}

// GetActivePlayers gets active players at the table
func (t *Table) GetActivePlayers() []*Player {
	players := make([]*Player, 0)
	nextSeat := t.SmallBlind
	for i := 0; i < nextSeat.Len(); i++ {
		if nextSeat.Player.Status == PlayerActive {
			players = append(players, nextSeat.Player)
		}
		nextSeat = nextSeat.Next()
	}
	return players
}

// AwardPot awards the entire pot to a player. This is used when all players have folded
// and there is no need to determine the best hand.
func (t *Table) AwardPot(p *Player) {
	p.Chips += t.Pot.GetTotal()
}

// DetermineWinners determines the winners of the hand
func (t *Table) DetermineWinners() [][]PlayerHand {
	subPots := t.Pot.GetSidePots()
	numSubPots := len(subPots)
	allWinningHands := make([][]PlayerHand, numSubPots)

	// Multiple pots can be created if there are players who are all in with
	// different chip stacks.
	for i, subPot := range subPots {
		winningHands := FindWinningHands(subPot.Players, t)

		// In the case of a tie, divide the pot amongst the winners
		chipsWon := subPot.Total / len(winningHands)
		remainderChipsWon := subPot.Total % len(winningHands)

		// If the pot can be split evenly among all winners, then
		// we will distribute one leftover chip to each player until
		// there are no more chips
		for j := range winningHands {
			playerChipsWon := chipsWon
			if remainderChipsWon > 0 {
				remainderChipsWon--
				playerChipsWon++
			}
			// Keep track of the chips won for logging purposes, such as displaying to chat
			winningHands[j].ChipsWon = playerChipsWon
			// Award winning chips to player
			winningHands[j].Player.Chips += playerChipsWon
		}
		allWinningHands[i] = winningHands
	}

	return allWinningHands
}

// Seat represents a seat at the poker table.
//
// - Seat uses ring.Ring under the hood to create a circular list.
// - There is a bi-directional link between seat/node.
//   - i.e. Seat.node has a Ring which has a Value of Seat.
type Seat struct {
	node   *ring.Ring
	Player *Player
}

// Next gets the next seat at the table.
func (s *Seat) Next() *Seat {
	return s.node.Next().Value.(*Seat)
}

// Move moves n spots from the current seat.
func (s *Seat) Move(n int) *Seat {
	return s.node.Move(n).Value.(*Seat)
}

// Len gets the number of seats at the table.
func (s *Seat) Len() int {
	return s.node.Len()
}

// NewSeat creates a new ring of seats.
func NewSeat(n int) *Seat {
	node := ring.New(n)
	for i := 0; i < node.Len(); i++ {
		node.Value = &Seat{
			node: node,
		}
		node = node.Next()
	}
	return node.Value.(*Seat)
}

// Pot represents the amount of chips in play
type Pot struct {
	Bets map[*Player]int
}

// SidePot represents a side pot
type SidePot struct {
	Players []*Player
	Total   int
	MaxBet  int
}

// NewPot creates a new pot
func NewPot() *Pot {
	return &Pot{
		Bets: make(map[*Player]int),
	}
}

// GetTotal gets the total pots
func (p *Pot) GetTotal() int {
	total := 0
	for _, betAmount := range p.Bets {
		total += betAmount
	}
	return total
}

// GetSidePots splits the pot into multiple pots
func (p *Pot) GetSidePots() []*SidePot {
	activePlayerBets := make(ByPlayerBet, 0)
	foldedPlayerBets := make(ByPlayerBet, 0)
	for player, betAmount := range p.Bets {
		if betAmount <= 0 {
			continue
		}
		if player.HasFolded == true {
			foldedPlayerBets = append(foldedPlayerBets, PlayerBet{Player: player, Total: betAmount})
		} else {
			activePlayerBets = append(activePlayerBets, PlayerBet{Player: player, Total: betAmount})
		}
	}
	sort.Sort(activePlayerBets)
	sidePots := []*SidePot{
		{Players: make([]*Player, 0)},
	}
	for _, playerBet := range activePlayerBets {
		numSidePots := len(sidePots)
		total := playerBet.Total
		for i := 0; i < numSidePots; i++ {
			sidePots[i].Players = append(sidePots[i].Players, playerBet.Player)
			if sidePots[i].Total == 0 {
				sidePots[i].Total += total
				sidePots[i].MaxBet = total
				total = 0
				break
			} else if sidePots[i].MaxBet >= total {
				sidePots[i].Total += total
				total = 0
				break
			} else if sidePots[i].MaxBet < total {
				total = total - sidePots[i].MaxBet
				sidePots[i].Total += sidePots[i].MaxBet
			}
		}
		if total > 0 {
			sidePots = append(sidePots, &SidePot{
				Players: []*Player{playerBet.Player},
				Total:   total,
				MaxBet:  total,
			})
		}
	}

	for _, playerBet := range foldedPlayerBets {
		numSidePots := len(sidePots)
		total := playerBet.Total
		for i := 0; i < numSidePots; i++ {
			if sidePots[i].MaxBet >= total {
				sidePots[i].Total += total
				total = 0
				break
			} else if sidePots[i].MaxBet < total {
				total = total - sidePots[i].MaxBet
				sidePots[i].Total += sidePots[i].MaxBet
			}
		}
	}

	return sidePots
}

// PlayerBet is the player's bet
type PlayerBet struct {
	Player *Player
	Total  int
}

// ByPlayerBet is an array of players bets that can be sorted
type ByPlayerBet []PlayerBet

func (p ByPlayerBet) Len() int      { return len(p) }
func (p ByPlayerBet) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ByPlayerBet) Less(i, j int) bool {
	if p[i].Total < p[j].Total {
		return true
	}

	if p[i].Total > p[j].Total {
		return false
	}

	return p[i].Player.ID < p[j].Player.ID
}

// BettingRound keeps track of chips/bets/raisers for a round of betting (preflop, flop, turn, river).
type BettingRound struct {
	Bets          map[string]int
	CallAmount    int
	Raiser        *Player
	RaiseByAmount int
}

// NewBettingRound makes a new BettingRound to keep track of chips/bets/raisers.
func NewBettingRound(startSeat *Seat, callAmount int, minBetAmount int) (*BettingRound, error) {
	// Initialize bets for all players to zero
	bets := make(map[string]int)
	for i := 0; i < startSeat.Len(); i++ {
		bets[startSeat.Player.ID] = 0
		startSeat = startSeat.Next()
	}

	p := startSeat.Player
	if p.HasFolded == true {
		return nil, fmt.Errorf("%s cannot start the round if they have already folded", p.Name)
	}

	return &BettingRound{
		Bets:          bets,
		CallAmount:    callAmount,
		Raiser:        p,
		RaiseByAmount: minBetAmount,
	}, nil

}

// DealHands deals hands to players.
func DealHands(d *Deck, t *Table) {
	rounds := 0
	numRounds := 2
	activePlayers := t.GetActivePlayers()
	hands := make([][2]*Card, len(activePlayers))
	for rounds < numRounds {
		for i := range hands {
			card, _ := d.GetNextCard()
			hands[i][rounds] = card
		}
		rounds++
	}

	for i := range hands {
		activePlayers[i].HoleCards = hands[i]
	}
}

// TakeSmallBlind takes the small blind and adds it to the pot.
func TakeSmallBlind(t *Table, b *BettingRound) error {
	p := t.SmallBlind.Player
	if p.Chips < t.MinBet {
		return fmt.Errorf("%s does not have enough chips to play", p.Name)
	}

	smallBlind := (t.MinBet / 2)
	p.Chips -= smallBlind
	t.Pot.Bets[p] += smallBlind
	b.Bets[p.ID] = smallBlind

	return nil
}

// TakeBigBlind takes the big blind and adds it to the pot.
func TakeBigBlind(t *Table, b *BettingRound) error {
	p := t.BigBlind.Player
	if p.Chips < t.MinBet {
		return fmt.Errorf("%s does not have enough chips to play", p.Name)
	}

	p.Chips -= t.MinBet
	t.Pot.Bets[p] += t.MinBet
	b.Bets[p.ID] = t.MinBet
	b.CallAmount = t.MinBet
	b.RaiseByAmount = t.MinBet

	return nil
}

// DealFlop deals the flop.
func DealFlop(d *Deck, t *Table) {
	for i := range t.Flop {
		card, _ := d.GetNextCard()
		t.Flop[i] = card
	}
}

// DealTurn deals the turn card.
func DealTurn(d *Deck, t *Table) {
	card, _ := d.GetNextCard()
	t.Turn = card
}

// DealRiver deals the river card.
func DealRiver(d *Deck, t *Table) {
	card, _ := d.GetNextCard()
	t.River = card
}

// GetNextActiveSeat get next active player who has not folded.
func GetNextActiveSeat(currentSeat *Seat) (*Seat, error) {
	nextSeat := currentSeat.Next()
	for i := 0; i < currentSeat.Len(); i++ {
		if nextSeat == currentSeat {
			return nil, fmt.Errorf("Next active seat not found. All players have folded")
		}
		if nextSeat.Player.Status == PlayerActive && nextSeat.Player.HasFolded == false {
			break
		}
		nextSeat = nextSeat.Next()
	}
	return nextSeat, nil
}

// CountSeatsByPlayerStatus counts the number of seats by player status
func CountSeatsByPlayerStatus(seat *Seat, status PlayerStatus) int {
	statusCount := 0
	for i := 0; i < seat.Len(); i++ {
		if seat.Player.Status == status {
			statusCount++
		}
		seat = seat.Next()
	}
	return statusCount
}

// HasEveryoneFolded checks if everyone except the current player has folded
func HasEveryoneFolded(currentSeat *Seat) bool {
	nextSeat := currentSeat.Next()
	for i := 0; i < nextSeat.Len()-1; i++ {
		if nextSeat.Player.Status == PlayerActive && nextSeat.Player.HasFolded == false {
			return false
		}
		nextSeat = nextSeat.Next()
	}
	return true
}

// HasEveryoneFoldedOrIsAllIn checks if everyone except the current player has folded or is all in
func HasEveryoneFoldedOrIsAllIn(currentSeat *Seat) bool {
	nextSeat := currentSeat.Next()
	for i := 0; i < nextSeat.Len()-1; i++ {
		if nextSeat.Player.Status == PlayerActive && (nextSeat.Player.HasFolded == false && nextSeat.Player.Chips > 0) {
			return false
		}
		nextSeat = nextSeat.Next()
	}
	return true
}
