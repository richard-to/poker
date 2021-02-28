package poker

import (
	"sort"
)

// Comparison is used for comparing hands (>, =, <).
type Comparison int

// Comparisons
const (
	LessThan    Comparison = -1
	EqualTo     Comparison = 0
	GreaterThan Comparison = 1
)

// HandRank is a ranking of poker hands.
type HandRank int

// Poker hand ranks.
const (
	HighCard HandRank = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (h HandRank) String() string {
	return [...]string{
		"High Card", "One Pair", "Two Pair", "Three Of A Kind", "Straight", "Flush",
		"Full House", "Four Of A Kind", "Straight Flush", "Royal Flush",
	}[h]
}

// PlayerHand represents a player's best hand among the possible combinations.
type PlayerHand struct {
	ChipsWon int
	Hand     *Hand
	Player   *Player
}

// FindWinningHands finds the players with the best hands.
func FindWinningHands(ps []*Player, t *Table) []PlayerHand {
	winners := make([]PlayerHand, 0)

	// Find players with the same ranks
	for i := range ps {
		if len(winners) == 0 {
			// If no winners have been found, add the current player as the winner by default
			winners = append(winners, PlayerHand{Hand: GetBestHand(ps[i], t), Player: ps[i]})
		} else {
			bestHand := GetBestHand(ps[i], t)
			winningHand := winners[0].Hand
			result := CompareHand(bestHand, winningHand)
			if result == GreaterThan {
				// If we found a better hand, set the current player as the winner
				winners = []PlayerHand{
					{Hand: bestHand, Player: ps[i]},
				}
			} else if result == EqualTo {
				// If we have a tie, the players will share the win
				winners = append(winners, PlayerHand{Hand: bestHand, Player: ps[i]})
			}
		}
	}
	sort.SliceStable(winners, func(i, j int) bool {
		return winners[i].Player.Chips < winners[i].Player.Chips
	})
	return winners
}

// CompareHand compares two hands.
func CompareHand(a *Hand, b *Hand) Comparison {
	if a.Rank > b.Rank {
		return GreaterThan
	}

	if a.Rank < b.Rank {
		return LessThan
	}

	// Tiebreakers assume we are comparing two hands that are the same rank.
	// This is necessary since the number of tiebreakers differs between
	// different types of hands.
	for i := range a.TieBreakers {
		if a.TieBreakers[i] < b.TieBreakers[i] {
			return LessThan
		}
		if a.TieBreakers[i] > b.TieBreakers[i] {
			return GreaterThan
		}
	}

	return EqualTo
}

// GetBestHand gets the player's best hand.
func GetBestHand(p *Player, t *Table) *Hand {
	cards := []Card{
		*p.HoleCards[0],
		*p.HoleCards[1],
		*t.Flop[0],
		*t.Flop[1],
		*t.Flop[2],
		*t.Turn,
		*t.River,
	}

	// The endIndex is exclusive which is why we pass in 2 + 1.
	// We use 2 since the endIndex is calculated as total cards (7) - hand size (5).
	cardCombos := FindCardCombinations(0, 3, cards)

	var cardHand [5]Card
	var bestHand *Hand
	for _, cs := range cardCombos {
		copy(cardHand[:], cs)
		currentHand := CheckHand(cardHand)
		if bestHand == nil {
			// If this is the first hand, set it as the best by default
			bestHand = currentHand
		} else {
			result := CompareHand(currentHand, bestHand)
			if result == GreaterThan {
				// If the current hand is better than the best hand, use it instead
				bestHand = currentHand
			}
		}
	}
	return bestHand
}

// FindCardCombinations finds all possible card combinations.
//
// In No Limit Texas Hold-em, a player can form a hand based on seven cards: two
// hole cards, and five community cards (via flop, turn, and river). Order does
// not matter.
//
// A hand is made out of five cards. This yields 21 combinations [i.e. 7!/(5!*2!)].
//
// This is a recursive function that takes a brute force approach to calculating
// the card combinations. This is fine since the number of combinations is small.
//
// On the initial call specify a startIndex of 0. The endIndex represents the
// the total cards minus the number of cards that a hand consists of. In addition
// we need to add 1 since endIndex is exclusive.
//
// Example: 7 (total cards) - 5 (hand size) + 1 = 3
//
// The algorithm is as follows:
//
// - Loop through each possible starting card
// - For each starting card, find the next possible cards by treating this subset
//   as a subproblem.
// - Once we've received the possible subcombos, append the starting card.
// - Add subcombos to the master list of combinations.
func FindCardCombinations(startIndex int, endIndex int, cs []Card) [][]Card {
	// Store all possible combinations of cards.
	combos := make([][]Card, 0)
	numCards := len(cs)

	// Loop count is needed since i starts from the startIndex. This disrupts the
	// calculation of newStartIndex which needs to from 0.
	loopCount := 0
	for i := startIndex; i < endIndex; i++ {
		// Set new start and end boundaries to work on the remaining
		// subset of cards.
		newStartIndex := startIndex + loopCount + 1
		newEndIndex := endIndex + 1

		// Base case for when we've reached the last card.
		subCombos := [][]Card{
			{cs[i]},
		}

		// If we haven't reached the last card, call FindCardCombinations again
		// to find combinations for the remaining subset.
		if newEndIndex <= numCards {
			subCombos = FindCardCombinations(newStartIndex, newEndIndex, cs)
			// Append the current card to the subcombos.
			for j := range subCombos {
				subCombos[j] = append(subCombos[j], cs[i])
			}
		}

		// Append all subcombos together
		for _, c := range subCombos {
			combos = append(combos, c)
		}

		loopCount++
	}

	return combos
}

// CheckHand checks a player's hand given their cards.
func CheckHand(cs [5]Card) *Hand {
	possibleHands := []IsHand{
		IsRoyalFlush,
		IsStraightFlush,
		IsFourOfAKind,
		IsFullHouse,
		IsFlush,
		IsStraight,
		IsThreeOfAKind,
		IsTwoPair,
		IsOnePair,
		IsHighCard,
	}

	var hand *Hand

	for _, isHand := range possibleHands {
		hand = isHand(cs)
		if hand != nil {
			break
		}
	}
	return hand
}

// Hand is a poker hand.
type Hand struct {
	Rank        HandRank
	TieBreakers []CardRank
}

// CompareLess compares that the hand is less than another hand.
//
// - The result will be true if the hand is less than or equal to the other
// - The result will be false if the hand is greater than the other
func (h *Hand) CompareLess(otherHand *Hand) bool {
	if h.Rank < otherHand.Rank {
		return true
	}

	if h.Rank > otherHand.Rank {
		return false
	}

	// Tiebreakers should only be compared against the same type of hand.
	for i := range h.TieBreakers {
		if h.TieBreakers[i] < otherHand.TieBreakers[i] {
			return true
		}
	}

	return false
}

// ByHand sorts hands by rank.
type ByHand []Hand

func (h ByHand) Len() int           { return len(h) }
func (h ByHand) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h ByHand) Less(i, j int) bool { return h[i].CompareLess(&h[j]) }

// IsHand is a type that checks a hand type.
type IsHand func(cs [5]Card) *Hand

// IsRoyalFlush checks if hand is a royal flush.
//
// Ties:
//    - No tiebreaker since the ranks would be the same
func IsRoyalFlush(cs [5]Card) *Hand {
	straightFlush := IsStraightFlush(cs)

	if straightFlush == nil {
		return nil
	}

	if straightFlush.TieBreakers[0] != Ace {
		return nil
	}

	return &Hand{
		Rank:        RoyalFlush,
		TieBreakers: []CardRank{},
	}
}

// IsStraightFlush checks if hand is a straight flush.
//
// Ties:
//    - Compare rank of highest card in straight
func IsStraightFlush(cs [5]Card) *Hand {
	flush := IsFlush(cs)
	straight := IsStraight(cs)

	if flush == nil || straight == nil {
		return nil
	}

	return &Hand{
		Rank:        StraightFlush,
		TieBreakers: straight.TieBreakers,
	}
}

// IsFourOfAKind checks if hand is four of a kind.
//
// Ties:
//   - Compare four of a kind rank
//   - Compare kicker rank
func IsFourOfAKind(cs [5]Card) *Hand {
	rankCount := make(map[CardRank]int)
	for _, c := range cs {
		rankCount[c.Rank]++
	}

	isHand := false
	handStrength := Two
	kicker := Two

	for rank, count := range rankCount {
		if count == 4 {
			isHand = true
			handStrength = rank
		}

		if count == 1 {
			kicker = rank
		}
	}

	if isHand == false {
		return nil
	}

	return &Hand{
		Rank:        FourOfAKind,
		TieBreakers: []CardRank{handStrength, kicker},
	}
}

// IsFullHouse checks if a hand is a full house.
//
// Ties:
//   - Compare three of a kind rank
//   - Compare pair rank
func IsFullHouse(cs [5]Card) *Hand {
	threeOfAKind := IsThreeOfAKind(cs)
	pair := IsOnePair(cs)
	if threeOfAKind == nil || pair == nil {
		return nil
	}

	handStrength := threeOfAKind.TieBreakers[0]
	kicker := pair.TieBreakers[0]

	return &Hand{
		Rank:        FullHouse,
		TieBreakers: []CardRank{handStrength, kicker},
	}
}

// IsFlush checks if hand is a flush.
//
// Ties:
//   - Compare cards in descending order
func IsFlush(cs [5]Card) *Hand {
	suit := cs[0].Suit
	for _, c := range cs {
		if c.Suit != suit {
			return nil
		}
	}

	sort.Sort(sort.Reverse(ByCard(cs[:])))
	tieBreakers := make([]CardRank, len(cs))
	for i := range cs {
		tieBreakers[i] = cs[i].Rank
	}

	return &Hand{
		Rank:        Flush,
		TieBreakers: tieBreakers,
	}
}

// IsStraight checks if hand is a straight.
//
// Ties:
//    - Compare rank of highest card in straight
func IsStraight(cs [5]Card) *Hand {
	sort.Sort(ByCard(cs[:]))

	// Handle Ace-5 edge case first
	if cs[0].Rank == Two && cs[1].Rank == Three && cs[2].Rank == Four && cs[3].Rank == Five && cs[4].Rank == Ace {
		return &Hand{
			Rank:        Straight,
			TieBreakers: []CardRank{Five},
		}
	}

	for i := 0; i < len(cs)-1; i++ {
		if cs[i].Rank+1 != cs[i+1].Rank {
			return nil
		}
	}
	// Tiebreaker is the high card for the straight
	handStrength := cs[4].Rank
	return &Hand{
		Rank:        Straight,
		TieBreakers: []CardRank{handStrength},
	}
}

// IsThreeOfAKind checks if hand is three of a kind.
//
// Ties:
//    - Compare ranks of three of a kind
//    - Compare remaining cards in descending order
func IsThreeOfAKind(cs [5]Card) *Hand {
	sort.Sort(sort.Reverse(ByCard(cs[:])))

	rankCount := make(map[CardRank]int)
	for _, c := range cs {
		rankCount[c.Rank]++
	}

	isHand := false
	tieBreakers := make([]CardRank, 0)
	for rank, count := range rankCount {
		if count == 3 {
			isHand = true
			tieBreakers = append(tieBreakers, rank)
		}
	}

	if isHand == false {
		return nil
	}

	// Add kickers
	for _, c := range cs {
		if rankCount[c.Rank] != 3 {
			tieBreakers = append(tieBreakers, c.Rank)
		}
	}

	return &Hand{
		Rank:        ThreeOfAKind,
		TieBreakers: tieBreakers,
	}
}

// IsTwoPair checks if hand is two pair.
//
// Ties:
//    - Compare ranks of pairs in descending order
//    - Compare remaining cards in descending order
func IsTwoPair(cs [5]Card) *Hand {
	rankCount := make(map[CardRank]int)
	for _, c := range cs {
		rankCount[c.Rank]++
	}

	pairs := 0
	tieBreakers := make([]CardRank, 0)
	for rank, count := range rankCount {
		if count == 2 {
			pairs++
			tieBreakers = append(tieBreakers, rank)
		}
	}

	if pairs != 2 {
		return nil
	}

	// Swap pairs so that the highest pair comes first
	if tieBreakers[0] < tieBreakers[1] {
		tempCard := tieBreakers[0]
		tieBreakers[0] = tieBreakers[1]
		tieBreakers[1] = tempCard
	}

	// Add kicker
	for _, c := range cs {
		if rankCount[c.Rank] != 2 {
			tieBreakers = append(tieBreakers, c.Rank)
		}
	}

	return &Hand{
		Rank:        TwoPair,
		TieBreakers: tieBreakers,
	}
}

// IsOnePair checks if hand is a pair.
//
// Ties:
//    - Compare rank of pair
//    - Compare remaining cards in descending order
func IsOnePair(cs [5]Card) *Hand {
	sort.Sort(sort.Reverse(ByCard(cs[:])))

	rankCount := make(map[CardRank]int)
	for _, c := range cs {
		rankCount[c.Rank]++
	}

	pairs := 0
	tieBreakers := make([]CardRank, 0)
	for rank, count := range rankCount {
		if count == 2 {
			pairs++
			tieBreakers = append(tieBreakers, rank)
		}
	}

	if pairs != 1 {
		return nil
	}

	for _, c := range cs {
		if rankCount[c.Rank] != 2 {
			tieBreakers = append(tieBreakers, c.Rank)
		}
	}

	return &Hand{
		Rank:        OnePair,
		TieBreakers: tieBreakers,
	}
}

// IsHighCard checks if hand is a high card.
//
// Ties:
//   - Compare cards in descending order
func IsHighCard(cs [5]Card) *Hand {
	sort.Sort(sort.Reverse(ByCard(cs[:])))

	tieBreakers := make([]CardRank, len(cs))
	for i := range cs {
		tieBreakers[i] = cs[i].Rank
	}

	return &Hand{
		Rank:        HighCard,
		TieBreakers: tieBreakers,
	}
}
