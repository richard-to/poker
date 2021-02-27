package poker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/richard-to/go-poker/pkg/poker"
)

var _ = Describe("FindWinningHands", func() {
	Context("when player 1 has the winning hand", func() {
		It("returns player 1 as the winner", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
					HoleCards: [2]*poker.Card{
						{Rank: poker.Ace, Suit: poker.Diamonds},
						{Rank: poker.Jack, Suit: poker.Hearts},
					},
					HasFolded: false,
					Chips:     10,
				},
				{
					ID:   "2",
					Name: "Player 2",
					HoleCards: [2]*poker.Card{
						{Rank: poker.Ace, Suit: poker.Clubs},
						{Rank: poker.Ten, Suit: poker.Clubs},
					},
					HasFolded: false,
					Chips:     10,
				},
			}

			t := poker.Table{
				Flop: [3]*poker.Card{
					{Rank: poker.Ace, Suit: poker.Hearts},
					{Rank: poker.Jack, Suit: poker.Clubs},
					{Rank: poker.Two, Suit: poker.Clubs},
				},
				Turn:  &poker.Card{Rank: poker.Seven, Suit: poker.Clubs},
				River: &poker.Card{Rank: poker.Ace, Suit: poker.Diamonds},
			}

			winner := poker.FindWinningHands(ps, &t)

			Expect(winner).To(Equal([]poker.PlayerHand{
				{
					Player: ps[0],
					Hand: &poker.Hand{
						Rank: poker.FullHouse,
						TieBreakers: []poker.CardRank{
							poker.Ace,
							poker.Jack,
						},
					},
				},
			}))
		})
	})

	Context("when player 2 has the winning hand", func() {
		It("returns player 2 as the winner", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
					HoleCards: [2]*poker.Card{
						{Rank: poker.King, Suit: poker.Clubs},
						{Rank: poker.Four, Suit: poker.Clubs},
					},
					HasFolded: false,
					Chips:     10,
				},
				{
					ID:   "2",
					Name: "Player 2",
					HoleCards: [2]*poker.Card{
						{Rank: poker.Ace, Suit: poker.Clubs},
						{Rank: poker.Ten, Suit: poker.Clubs},
					},
					HasFolded: false,
					Chips:     10,
				},
			}

			t := poker.Table{
				Flop: [3]*poker.Card{
					{Rank: poker.Ace, Suit: poker.Hearts},
					{Rank: poker.Jack, Suit: poker.Clubs},
					{Rank: poker.Two, Suit: poker.Clubs},
				},
				Turn:  &poker.Card{Rank: poker.Seven, Suit: poker.Clubs},
				River: &poker.Card{Rank: poker.Ace, Suit: poker.Diamonds},
			}

			winner := poker.FindWinningHands(ps, &t)

			Expect(winner).To(Equal([]poker.PlayerHand{
				{
					Player: ps[1],
					Hand: &poker.Hand{
						Rank: poker.Flush,
						TieBreakers: []poker.CardRank{
							poker.Ace,
							poker.Jack,
							poker.Ten,
							poker.Seven,
							poker.Two,
						},
					},
				},
			}))
		})
	})

	Context("when three players have the winning hand", func() {
		It("returns player 1, 2, and 3 as the winner", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
					HoleCards: [2]*poker.Card{
						{Rank: poker.King, Suit: poker.Diamonds},
						{Rank: poker.Four, Suit: poker.Hearts},
					},
					HasFolded: false,
					Chips:     10,
				},
				{
					ID:   "2",
					Name: "Player 2",
					HoleCards: [2]*poker.Card{
						{Rank: poker.Ace, Suit: poker.Hearts},
						{Rank: poker.Ten, Suit: poker.Clubs},
					},
					HasFolded: false,
					Chips:     10,
				},
				{
					ID:   "3",
					Name: "Player 3",
					HoleCards: [2]*poker.Card{
						{Rank: poker.Ace, Suit: poker.Spades},
						{Rank: poker.Jack, Suit: poker.Hearts},
					},
					HasFolded: false,
					Chips:     10,
				},
			}

			t := poker.Table{
				Flop: [3]*poker.Card{
					{Rank: poker.Three, Suit: poker.Hearts},
					{Rank: poker.Three, Suit: poker.Clubs},
					{Rank: poker.Three, Suit: poker.Spades},
				},
				Turn:  &poker.Card{Rank: poker.Ace, Suit: poker.Clubs},
				River: &poker.Card{Rank: poker.Three, Suit: poker.Diamonds},
			}

			winner := poker.FindWinningHands(ps, &t)

			Expect(winner).To(Equal([]poker.PlayerHand{
				{
					Player: ps[0],
					Hand: &poker.Hand{
						Rank: poker.FourOfAKind,
						TieBreakers: []poker.CardRank{
							poker.Three,
							poker.Ace,
						},
					},
				},
				{
					Player: ps[1],
					Hand: &poker.Hand{
						Rank: poker.FourOfAKind,
						TieBreakers: []poker.CardRank{
							poker.Three,
							poker.Ace,
						},
					},
				},
				{
					Player: ps[2],
					Hand: &poker.Hand{
						Rank: poker.FourOfAKind,
						TieBreakers: []poker.CardRank{
							poker.Three,
							poker.Ace,
						},
					},
				},
			}))
		})
	})
})

var _ = Describe("GetBestHand", func() {
	It("returns the best hand", func() {
		p := poker.Player{
			ID:   "1",
			Name: "Player 1",
			HoleCards: [2]*poker.Card{
				{Rank: poker.Ace, Suit: poker.Clubs},
				{Rank: poker.Ten, Suit: poker.Clubs},
			},
			HasFolded: false,
			Chips:     10,
		}

		t := poker.Table{
			Flop: [3]*poker.Card{
				{Rank: poker.Ace, Suit: poker.Hearts},
				{Rank: poker.Jack, Suit: poker.Clubs},
				{Rank: poker.Two, Suit: poker.Clubs},
			},
			Turn:  &poker.Card{Rank: poker.Seven, Suit: poker.Clubs},
			River: &poker.Card{Rank: poker.Ace, Suit: poker.Diamonds},
		}

		bestHand := poker.GetBestHand(&p, &t)
		expectedHand := poker.Hand{
			Rank: poker.Flush,
			TieBreakers: []poker.CardRank{
				poker.Ace,
				poker.Jack,
				poker.Ten,
				poker.Seven,
				poker.Two,
			},
		}
		Expect(bestHand).To(Equal(&expectedHand))
	})
})

var _ = Describe("FindCardCombinations", func() {

	Context("when given seven cards to make a five card hand", func() {
		It("returns 21 combinations", func() {
			cs := []poker.Card{
				{Rank: poker.Ace, Suit: poker.Clubs},
				{Rank: poker.Two, Suit: poker.Clubs},
				{Rank: poker.Three, Suit: poker.Clubs},
				{Rank: poker.Four, Suit: poker.Clubs},
				{Rank: poker.Five, Suit: poker.Clubs},
				{Rank: poker.Six, Suit: poker.Clubs},
				{Rank: poker.Seven, Suit: poker.Clubs},
			}
			combos := poker.FindCardCombinations(0, 3, cs)
			Expect(len(combos)).To(Equal(21))
			Expect(len(combos[0])).To(Equal(5)) // Each combo has five cards
		})
	})

	Context("when given five cards to make a three card hand", func() {
		cs := []poker.Card{
			{Rank: poker.Ace, Suit: poker.Clubs},
			{Rank: poker.Two, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Clubs},
			{Rank: poker.Four, Suit: poker.Clubs},
			{Rank: poker.Five, Suit: poker.Clubs},
		}
		combos := poker.FindCardCombinations(0, 3, cs)

		It("returns 10 combinations", func() {
			Expect(len(combos)).To(Equal(10))
		})

		It("returns the following combinations", func() {
			expectedCombos := [][]poker.Card{
				{
					{Rank: poker.Three, Suit: poker.Clubs},
					{Rank: poker.Two, Suit: poker.Clubs},
					{Rank: poker.Ace, Suit: poker.Clubs},
				},
				{
					{Rank: poker.Four, Suit: poker.Clubs},
					{Rank: poker.Two, Suit: poker.Clubs},
					{Rank: poker.Ace, Suit: poker.Clubs},
				},
				{
					{Rank: poker.Five, Suit: poker.Clubs},
					{Rank: poker.Two, Suit: poker.Clubs},
					{Rank: poker.Ace, Suit: poker.Clubs},
				},
				{
					{Rank: poker.Four, Suit: poker.Clubs},
					{Rank: poker.Three, Suit: poker.Clubs},
					{Rank: poker.Ace, Suit: poker.Clubs},
				},
				{
					{Rank: poker.Five, Suit: poker.Clubs},
					{Rank: poker.Three, Suit: poker.Clubs},
					{Rank: poker.Ace, Suit: poker.Clubs},
				},
				{
					{Rank: poker.Five, Suit: poker.Clubs},
					{Rank: poker.Four, Suit: poker.Clubs},
					{Rank: poker.Ace, Suit: poker.Clubs},
				},
				{
					{Rank: poker.Four, Suit: poker.Clubs},
					{Rank: poker.Three, Suit: poker.Clubs},
					{Rank: poker.Two, Suit: poker.Clubs},
				},
				{
					{Rank: poker.Five, Suit: poker.Clubs},
					{Rank: poker.Three, Suit: poker.Clubs},
					{Rank: poker.Two, Suit: poker.Clubs},
				},
				{
					{Rank: poker.Five, Suit: poker.Clubs},
					{Rank: poker.Four, Suit: poker.Clubs},
					{Rank: poker.Two, Suit: poker.Clubs},
				},
				{
					{Rank: poker.Five, Suit: poker.Clubs},
					{Rank: poker.Four, Suit: poker.Clubs},
					{Rank: poker.Three, Suit: poker.Clubs},
				},
			}
			Expect(combos).To(ConsistOf(expectedCombos))
		})
	})
})

var _ = Describe("CheckHand", func() {
	Context("when a hand is found", func() {
		It("returns the highest possible hand", func() {
			cs := [5]poker.Card{
				{Rank: poker.Nine, Suit: poker.Clubs},
				{Rank: poker.Jack, Suit: poker.Clubs},
				{Rank: poker.Queen, Suit: poker.Clubs},
				{Rank: poker.Eight, Suit: poker.Clubs},
				{Rank: poker.Ten, Suit: poker.Clubs},
			}
			expectedHand := poker.Hand{
				Rank: poker.StraightFlush,
				TieBreakers: []poker.CardRank{
					poker.Queen,
				},
			}
			Expect(poker.CheckHand(cs)).To(Equal(&expectedHand))
		})
	})

	Context("when no hand is found", func() {
		It("returns high card", func() {
			cs := [5]poker.Card{
				{Rank: poker.Jack, Suit: poker.Clubs},
				{Rank: poker.Three, Suit: poker.Clubs},
				{Rank: poker.Six, Suit: poker.Spades},
				{Rank: poker.Ten, Suit: poker.Spades},
				{Rank: poker.Queen, Suit: poker.Diamonds},
			}
			expectedHand := poker.Hand{
				Rank: poker.HighCard,
				TieBreakers: []poker.CardRank{
					poker.Queen,
					poker.Jack,
					poker.Ten,
					poker.Six,
					poker.Three,
				},
			}
			Expect(poker.CheckHand(cs)).To(Equal(&expectedHand))
		})
	})
})

var _ = Describe("Hand.Compare", func() {
	Context("when a hand is greater", func() {
		It("returns false", func() {
			a := poker.Hand{
				Rank:        poker.TwoPair,
				TieBreakers: []poker.CardRank{poker.Ace, poker.Jack, poker.Two},
			}
			b := poker.Hand{
				Rank:        poker.OnePair,
				TieBreakers: []poker.CardRank{poker.Ace, poker.Jack, poker.Five, poker.Two},
			}
			Expect(a.CompareLess(&b)).To(BeFalse())
		})
	})

	Context("when a hand is lesser", func() {
		It("returns true", func() {
			a := poker.Hand{
				Rank:        poker.OnePair,
				TieBreakers: []poker.CardRank{poker.Ace, poker.Jack, poker.Five, poker.Two},
			}
			b := poker.Hand{
				Rank:        poker.TwoPair,
				TieBreakers: []poker.CardRank{poker.Ace, poker.Jack, poker.Two},
			}
			Expect(a.CompareLess(&b)).To(BeTrue())
		})
	})

	Context("when a hand is lesser via tiebreaker", func() {
		It("returns true", func() {
			a := poker.Hand{
				Rank:        poker.TwoPair,
				TieBreakers: []poker.CardRank{poker.Ace, poker.Jack, poker.Two},
			}
			b := poker.Hand{
				Rank:        poker.TwoPair,
				TieBreakers: []poker.CardRank{poker.Ace, poker.Jack, poker.Three},
			}
			Expect(a.CompareLess(&b)).To(BeTrue())
		})
	})

	Context("when a hand is equal", func() {
		It("returns false", func() {
			a := poker.Hand{
				Rank:        poker.TwoPair,
				TieBreakers: []poker.CardRank{poker.Ace, poker.Jack, poker.Two},
			}
			b := poker.Hand{
				Rank:        poker.TwoPair,
				TieBreakers: []poker.CardRank{poker.Ace, poker.Jack, poker.Two},
			}
			Expect(a.CompareLess(&b)).To(BeFalse())
		})
	})
})

var _ = Describe("IsRoyalFlush", func() {
	Context("when a hand is a royal flush", func() {
		cs := [5]poker.Card{
			{Rank: poker.Ace, Suit: poker.Clubs},
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Queen, Suit: poker.Clubs},
			{Rank: poker.King, Suit: poker.Clubs},
			{Rank: poker.Ten, Suit: poker.Clubs},
		}
		hand := poker.IsRoyalFlush(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.RoyalFlush))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is not a royal flush", func() {
		cs := [5]poker.Card{
			{Rank: poker.Nine, Suit: poker.Clubs},
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Queen, Suit: poker.Clubs},
			{Rank: poker.King, Suit: poker.Clubs},
			{Rank: poker.Ten, Suit: poker.Clubs},
		}
		hand := poker.IsRoyalFlush(cs)

		It("is nil", func() {
			Expect(hand).To(BeNil())
		})
	})
})

var _ = Describe("IsStraightFlush", func() {
	Context("when a hand is a straight flush", func() {
		cs := [5]poker.Card{
			{Rank: poker.Nine, Suit: poker.Clubs},
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Queen, Suit: poker.Clubs},
			{Rank: poker.King, Suit: poker.Clubs},
			{Rank: poker.Ten, Suit: poker.Clubs},
		}
		hand := poker.IsStraightFlush(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.StraightFlush))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.King,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is not a straight flush", func() {
		cs := [5]poker.Card{
			{Rank: poker.Nine, Suit: poker.Clubs},
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Queen, Suit: poker.Clubs},
			{Rank: poker.Ace, Suit: poker.Clubs},
			{Rank: poker.Ten, Suit: poker.Clubs},
		}
		hand := poker.IsStraightFlush(cs)

		It("is nil", func() {
			Expect(hand).To(BeNil())
		})
	})
})

var _ = Describe("IsFourOfAKind", func() {
	Context("when a hand is four of a kind", func() {
		cs := [5]poker.Card{
			{Rank: poker.Ace, Suit: poker.Clubs},
			{Rank: poker.Jack, Suit: poker.Hearts},
			{Rank: poker.Ace, Suit: poker.Diamonds},
			{Rank: poker.Ace, Suit: poker.Spades},
			{Rank: poker.Ace, Suit: poker.Hearts},
		}
		hand := poker.IsFourOfAKind(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.FourOfAKind))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.Ace,
				poker.Jack,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is not four of a kind", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Jack, Suit: poker.Diamonds},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Jack, Suit: poker.Hearts},
		}
		hand := poker.IsFourOfAKind(cs)

		It("is nil", func() {
			Expect(hand).To(BeNil())
		})
	})
})

var _ = Describe("IsFullHouse", func() {
	Context("when a hand is a full house", func() {
		cs := [5]poker.Card{
			{Rank: poker.Ace, Suit: poker.Clubs},
			{Rank: poker.Jack, Suit: poker.Hearts},
			{Rank: poker.Ace, Suit: poker.Diamonds},
			{Rank: poker.Ace, Suit: poker.Spades},
			{Rank: poker.Jack, Suit: poker.Clubs},
		}
		hand := poker.IsFullHouse(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.FullHouse))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.Ace,
				poker.Jack,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is not a full house", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Three, Suit: poker.Diamonds},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Queen, Suit: poker.Clubs},
		}
		hand := poker.IsFullHouse(cs)

		It("is nil", func() {
			Expect(hand).To(BeNil())
		})
	})
})

var _ = Describe("IsFlush", func() {
	Context("when a hand is a flush", func() {
		cs := [5]poker.Card{
			{Rank: poker.Ace, Suit: poker.Clubs},
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Six, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Clubs},
			{Rank: poker.Four, Suit: poker.Clubs},
		}
		hand := poker.IsFlush(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.Flush))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.Ace,
				poker.Jack,
				poker.Six,
				poker.Four,
				poker.Three,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is not a flush", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Clubs},
			{Rank: poker.Six, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Queen, Suit: poker.Clubs},
		}
		hand := poker.IsFlush(cs)

		It("is nil", func() {
			Expect(hand).To(BeNil())
		})
	})
})

var _ = Describe("IsStraight", func() {
	Context("when a hand is a straight", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Nine, Suit: poker.Hearts},
			{Rank: poker.Seven, Suit: poker.Spades},
			{Rank: poker.Ten, Suit: poker.Spades},
			{Rank: poker.Eight, Suit: poker.Diamonds},
		}
		hand := poker.IsStraight(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.Straight))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.Jack,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is a Ace-Five straight", func() {
		cs := [5]poker.Card{
			{Rank: poker.Five, Suit: poker.Clubs},
			{Rank: poker.Ace, Suit: poker.Hearts},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Four, Suit: poker.Spades},
			{Rank: poker.Two, Suit: poker.Diamonds},
		}
		hand := poker.IsStraight(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.Straight))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.Five,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is not a straight", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Queen, Suit: poker.Hearts},
			{Rank: poker.Seven, Suit: poker.Spades},
			{Rank: poker.Ten, Suit: poker.Spades},
			{Rank: poker.Eight, Suit: poker.Diamonds},
		}
		hand := poker.IsStraight(cs)

		It("is nil", func() {
			Expect(hand).To(BeNil())
		})
	})
})

var _ = Describe("IsThreeOfAKind", func() {
	Context("when a hand is three of a kind", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Jack, Suit: poker.Hearts},
			{Rank: poker.Six, Suit: poker.Spades},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Jack, Suit: poker.Diamonds},
		}
		hand := poker.IsThreeOfAKind(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.ThreeOfAKind))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.Jack,
				poker.Six,
				poker.Three,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is not two pair", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Clubs},
			{Rank: poker.Six, Suit: poker.Spades},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Queen, Suit: poker.Diamonds},
		}
		hand := poker.IsThreeOfAKind(cs)

		It("is nil", func() {
			Expect(hand).To(BeNil())
		})
	})
})

var _ = Describe("IsTwoPair", func() {
	Context("when a hand is two pair", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Clubs},
			{Rank: poker.Six, Suit: poker.Spades},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Jack, Suit: poker.Diamonds},
		}
		hand := poker.IsTwoPair(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.TwoPair))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.Jack,
				poker.Three,
				poker.Six,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is not two pair", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Clubs},
			{Rank: poker.Six, Suit: poker.Spades},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Queen, Suit: poker.Diamonds},
		}
		hand := poker.IsTwoPair(cs)

		It("is nil", func() {
			Expect(hand).To(BeNil())
		})
	})
})

var _ = Describe("IsOnePair", func() {
	Context("when a hand is a pair", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Clubs},
			{Rank: poker.Six, Suit: poker.Spades},
			{Rank: poker.Three, Suit: poker.Spades},
			{Rank: poker.Queen, Suit: poker.Diamonds},
		}
		hand := poker.IsOnePair(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.OnePair))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.Three,
				poker.Queen,
				poker.Jack,
				poker.Six,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})

	Context("when a hand is not a pair", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Clubs},
			{Rank: poker.Six, Suit: poker.Spades},
			{Rank: poker.King, Suit: poker.Spades},
			{Rank: poker.Queen, Suit: poker.Diamonds},
		}
		hand := poker.IsOnePair(cs)

		It("is nil", func() {
			Expect(hand).To(BeNil())
		})
	})
})

var _ = Describe("IsHighCard", func() {
	Context("when a hand is a high card", func() {
		cs := [5]poker.Card{
			{Rank: poker.Jack, Suit: poker.Clubs},
			{Rank: poker.Three, Suit: poker.Clubs},
			{Rank: poker.Six, Suit: poker.Spades},
			{Rank: poker.Ten, Suit: poker.Spades},
			{Rank: poker.Queen, Suit: poker.Diamonds},
		}
		hand := poker.IsHighCard(cs)

		It("has the right rank", func() {
			Expect(hand.Rank).To(Equal(poker.HighCard))
		})

		It("has the right tiebreakers", func() {
			tieBreakers := []poker.CardRank{
				poker.Queen,
				poker.Jack,
				poker.Ten,
				poker.Six,
				poker.Three,
			}
			Expect(hand.TieBreakers).To(Equal(tieBreakers))
		})
	})
})
