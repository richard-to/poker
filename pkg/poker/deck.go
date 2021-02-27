package poker

import (
	"fmt"
	"math/rand"
)

// DeckSize is the number of cards in a deck.
const DeckSize = 52

// CardSuit is an enum for a card's suit.
type CardSuit int

// Card suits
const (
	Clubs CardSuit = iota
	Diamonds
	Hearts
	Spades
)

func (s CardSuit) String() string {
	return [...]string{"Clubs", "Diamonds", "Hearts", "Spades"}[s]
}

// Symbol shows the card's suit in an abbreviated format.
func (s CardSuit) Symbol() string {
	return [...]string{"♣", "♦", "♥", "♠"}[s]
}

// CardRank is an enum for a card's rank.
type CardRank int

// Card ranks
const (
	Two CardRank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (r CardRank) String() string {
	return [...]string{
		"Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King", "Ace",
	}[r]
}

// Symbol shows the card's rank in an abbreviated format.
func (r CardRank) Symbol() string {
	return [...]string{
		"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A",
	}[r]
}

// ByCard sorts cards in ascending order by rank. Suit is not used as a secondary
// sort parameter. This means that the sort is not deterministic.
type ByCard []Card

func (c ByCard) Len() int           { return len(c) }
func (c ByCard) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCard) Less(i, j int) bool { return c[i].Rank < c[j].Rank }

// Card is a playing card with a suit and a value.
type Card struct {
	Rank CardRank `json:"rank"`
	Suit CardSuit `json:"suit"`
}

func (c *Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

// Symbol shows the card's abbreviated output format (e.g. J♣).
func (c *Card) Symbol() string {
	return fmt.Sprintf("%s%s", c.Rank.Symbol(), c.Suit.Symbol())
}

// Deck is a deck of cards.
type Deck struct {
	cards            []Card
	currentCardIndex int
}

// GetNextCard gets the next card from the deck. An error will occur if there
// are no more cards in the deck.
func (d *Deck) GetNextCard() (*Card, error) {
	if d.currentCardIndex >= DeckSize {
		return nil, fmt.Errorf("No more cards left in deck")
	}
	card := d.cards[d.currentCardIndex]
	d.currentCardIndex++
	return &card, nil
}

// NewDeck creates a shuffled deck of cards.
func NewDeck() Deck {
	suits := []CardSuit{Clubs, Diamonds, Hearts, Spades}
	ranks := []CardRank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	cards := make([]Card, DeckSize)

	// Make a deck with a card for each suit and rank pair
	i := 0
	for _, s := range suits {
		for _, r := range ranks {
			cards[i] = Card{
				Rank: r,
				Suit: s,
			}
			i++
		}
	}

	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })

	return Deck{cards: cards, currentCardIndex: 0}
}
