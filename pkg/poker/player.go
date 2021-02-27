package poker

import (
	"errors"
	"fmt"
)

// Player is a player in the poker game. It represents their state for a round
// of poker, which encompasses all betting rounds from preflop, flop, turn, and
// river.
type Player struct {
	Active    bool
	Chips     int
	HasFolded bool
	HoleCards [2]*Card
	ID        string
	Name      string
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
	return p.HasFolded == false && p.Chips > 0 && b.Bets[p.ID] < b.CallAmount
}

// Fold folds a players hand.
func (p *Player) Fold() error {
	if p.Chips <= 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}
	if p.HasFolded {
		return fmt.Errorf("%s has already folded", p.Name)
	}
	p.HasFolded = true
	return nil
}

// CanCheck checks if the player can check.
func (p *Player) CanCheck(b *BettingRound) bool {
	if p.HasFolded || p.Chips <= 0 {
		return false
	}
	return b.Bets[p.ID] == b.CallAmount
}

// Check passes action to the next player.
func (p *Player) Check(b *BettingRound) error {
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
	if p.HasFolded || p.Chips <= 0 {
		return false
	}
	return b.Bets[p.ID] < b.CallAmount
}

// Call calls the current bet/raise.
func (p *Player) Call(t *Table, b *BettingRound) error {
	if p.Chips <= 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}
	if p.HasFolded {
		return fmt.Errorf("%s can't call when folded", p.Name)
	}

	chipsInPot := b.Bets[p.ID]
	if chipsInPot >= b.CallAmount {
		return fmt.Errorf(
			"%s has more chips wagered (%d) than the call amount (%d)",
			p.Name, chipsInPot, b.CallAmount,
		)
	}

	callAmount := b.CallAmount - chipsInPot
	if callAmount > p.Chips {
		callAmount = p.Chips // TODO manage side pots
	}
	t.Pot.Bets[p] += callAmount
	b.Bets[p.ID] += callAmount
	p.Chips -= callAmount

	return nil
}

// CanRaise checks if the player can bet/raise.
func (p *Player) CanRaise(b *BettingRound) bool {
	if p.HasFolded || p.Chips <= 0 || p.Chips <= b.CallAmount-b.Bets[p.ID] {
		return false
	}
	return true
}

// Raise bet/raises the pot.
func (p *Player) Raise(t *Table, b *BettingRound, raiseAmount int) error {
	if p.Chips == 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}
	if p.HasFolded {
		return fmt.Errorf("%s can't bet/raise when folded", p.Name)
	}

	minRaiseTo := b.CallAmount + b.RaiseByAmount

	if raiseAmount < minRaiseTo {
		return fmt.Errorf(
			"%s's raise (%d) is less than the minimum bet/raise (%d)",
			p.Name, raiseAmount, minRaiseTo,
		)
	}

	chipsInPot := b.Bets[p.ID]
	chipsNeeded := raiseAmount - chipsInPot

	if chipsNeeded > p.Chips {
		return fmt.Errorf(
			"%s does not have enough chips (%d) to bet/raise (%d)",
			p.Name, p.Chips, chipsNeeded,
		)
	}
	b.RaiseByAmount = raiseAmount - b.CallAmount
	b.CallAmount = raiseAmount
	b.Raiser = p
	t.Pot.Bets[p] += chipsNeeded
	b.Bets[p.ID] += chipsNeeded
	p.Chips -= chipsNeeded

	return nil
}
