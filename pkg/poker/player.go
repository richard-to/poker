package poker

import (
	"errors"
	"fmt"
)

// PlayerStatus is the status of the player
type PlayerStatus int

// Player statuses
const (
	PlayerVacated PlayerStatus = iota
	PlayerSittingOut
	PlayerActive
)

func (p PlayerStatus) String() string {
	return [...]string{"vacated", "sitting-out", "active"}[p]
}

// Player is a player in the poker game. It represents their state for a round
// of poker, which encompasses all betting rounds from preflop, flop, turn, and
// river.
type Player struct {
	Chips     int
	HasFolded bool
	HoleCards [2]*Card
	ID        string
	Name      string
	Status    PlayerStatus
}

// PrintHoleCards gets the player's hand in abbreviated format.
func (p *Player) PrintHoleCards() (string, error) {
	if p.HoleCards[0] == nil || p.HoleCards[1] == nil {
		return "", errors.New("The player does not have any hole cards yet")
	}
	hand := fmt.Sprintf("%s %s", p.HoleCards[0].Symbol(), p.HoleCards[1].Symbol())
	return hand, nil
}

// CanFold checks if the player can fold.
func (p *Player) CanFold(b *BettingRound) bool {
	return (p.Status == PlayerActive &&
		p.HasFolded == false &&
		p.Chips > 0 &&
		b.Bets[p.ID] < b.CallAmount)
}

// Fold folds a players hand.
func (p *Player) Fold(b *BettingRound) error {
	if p.Status != PlayerActive {
		return fmt.Errorf("%s is not in the hand", p.Name)
	}
	if p.Chips <= 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}
	if p.HasFolded {
		return fmt.Errorf("%s has already folded", p.Name)
	}
	if b.Bets[p.ID] >= b.CallAmount {
		return fmt.Errorf("%s has matched or exceeded the current bets", p.Name)
	}
	p.HasFolded = true
	return nil
}

// CanCheck checks if the player can check.
func (p *Player) CanCheck(b *BettingRound) bool {
	return (p.Status == PlayerActive &&
		p.HasFolded == false &&
		p.Chips > 0 &&
		b.Bets[p.ID] == b.CallAmount)
}

// Check passes action to the next player.
func (p *Player) Check(b *BettingRound) error {
	if p.Status != PlayerActive {
		return fmt.Errorf("%s is not in the hand", p.Name)
	}
	if p.Chips == 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}
	if p.HasFolded {
		return fmt.Errorf("%s can't check when folded", p.Name)
	}
	if b.Bets[p.ID] < b.CallAmount {
		return fmt.Errorf("%s can't check when someone has bet", p.Name)
	}
	return nil
}

// CanCall checks if the player can call a bet/raise.
func (p *Player) CanCall(b *BettingRound) bool {
	return (p.Status == PlayerActive &&
		p.HasFolded == false &&
		p.Chips > 0 &&
		b.Bets[p.ID] < b.CallAmount)
}

// Call calls the current bet/raise.
func (p *Player) Call(t *Table, b *BettingRound) error {
	if p.Status != PlayerActive {
		return fmt.Errorf("%s is not in the hand", p.Name)
	}
	if p.Chips <= 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}
	if p.HasFolded {
		return fmt.Errorf("%s can't call when folded", p.Name)
	}

	// A player can't call if they have more chips bet than the call amount
	chipsInPot := b.Bets[p.ID]
	if chipsInPot >= b.CallAmount {
		return fmt.Errorf(
			"%s has more chips wagered (%d) than the call amount (%d)",
			p.Name, chipsInPot, b.CallAmount,
		)
	}

	callAmount := b.CallAmount - chipsInPot
	// If a player does not have enough chips to make a call, then they're
	// all in for the remainder of their chips
	if callAmount > p.Chips {
		callAmount = p.Chips
	}
	t.Pot.Bets[p] += callAmount
	b.Bets[p.ID] += callAmount
	p.Chips -= callAmount

	return nil
}

// CanRaise checks if the player can bet/raise.
//
// A player can bet/raise less than in min raise if they don't have enough chips (i.e. they
// must go all in).
func (p *Player) CanRaise(b *BettingRound) bool {
	return (p.Status == PlayerActive &&
		p.HasFolded == false &&
		p.Chips > 0 &&
		p.Chips >= b.CallAmount-b.Bets[p.ID])
}

// Raise raises the pot.
func (p *Player) Raise(t *Table, b *BettingRound, raiseAmount int) error {
	actionLabel := "bet"
	if b.CallAmount > 0 {
		actionLabel = "raise"
	}
	if p.Status != PlayerActive {
		return fmt.Errorf("%s is not in the hand", p.Name)
	}
	if p.Chips == 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}

	if p.HasFolded {
		return fmt.Errorf("%s can't %s when folded", p.Name, actionLabel)
	}

	chipsInPot := b.Bets[p.ID]
	chipsNeeded := raiseAmount - chipsInPot

	if chipsNeeded > p.Chips {
		return fmt.Errorf(
			"%s does not have enough chips (%d) to %s (%d)",
			p.Name, p.Chips, actionLabel, chipsNeeded,
		)
	}

	minRaiseTo := b.CallAmount + b.RaiseByAmount
	// If a player has enough chips to match the minimum bet/raise, then they must
	// bet/raise at least the minimum.
	if raiseAmount < minRaiseTo && minRaiseTo-chipsInPot < p.Chips {
		return fmt.Errorf(
			"%s's raise (%d) is less than the minimum %s (%d)",
			p.Name, raiseAmount, actionLabel, minRaiseTo,
		)
	}

	// Only increase the min bet/raise amount if the bet/raise was at least a full bet/raise
	if raiseAmount >= minRaiseTo {
		b.RaiseByAmount = raiseAmount - b.CallAmount
	}

	b.CallAmount = raiseAmount
	b.Raiser = p
	t.Pot.Bets[p] += chipsNeeded
	b.Bets[p.ID] += chipsNeeded
	p.Chips -= chipsNeeded

	return nil
}
