package poker_test

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/richard-to/go-poker/pkg/poker"
)

var _ = Describe("Card", func() {
	Describe("String output", func() {
		var card poker.Card

		BeforeEach(func() {
			card = poker.Card{Rank: poker.Ace, Suit: poker.Clubs}
		})

		It("prints in long format", func() {
			Expect(card.String()).To(Equal("Ace of Clubs"))
		})

		It("prints in short format", func() {
			Expect(card.Symbol()).To(Equal("A♣"))
		})
	})

	Describe("Sorting", func() {
		It("sorts cards in ascending order by rank", func() {
			cards := []poker.Card{
				{Rank: poker.Ace, Suit: poker.Clubs},
				{Rank: poker.Ten, Suit: poker.Hearts},
				{Rank: poker.Queen, Suit: poker.Hearts},
				{Rank: poker.Two, Suit: poker.Spades},
				{Rank: poker.Queen, Suit: poker.Spades},
			}
			expectedRanks := []poker.CardRank{
				poker.Two,
				poker.Ten,
				poker.Queen,
				poker.Queen,
				poker.Ace,
			}
			sort.Sort(poker.ByCard(cards))

			for i := range cards {
				Expect(cards[i].Rank).To(Equal(expectedRanks[i]))
			}
		})
	})
})

var _ = Describe("Deck", func() {
	Describe("GetNextCard", func() {
		It("has 52 unique cards", func() {
			cardCount := 0
			deck := poker.NewDeck()
			cardMap := make(map[poker.Card]int)

			for {
				card, err := deck.GetNextCard()
				cardCount++
				if err != nil {
					Expect(cardCount).To(Equal(poker.DeckSize + 1))
				} else {
					Expect(cardCount <= poker.DeckSize).To(BeTrue())
					cardMap[*card]++
				}

				if cardCount == poker.DeckSize+1 {
					break
				}
			}

			// Cards should be unique. We make the assumption that if all
			// 52 cards are unique, that it covers all combinations of rank
			// and suit.
			for _, count := range cardMap {
				Expect(count).To(Equal(1))
			}
		})
	})
})
