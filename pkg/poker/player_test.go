package poker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/richard-to/go-poker/pkg/poker"
)

var _ = Describe("Player", func() {
	var player poker.Player

	BeforeEach(func() {
		player = poker.Player{
			ID:   "1",
			Name: "Player 1",
			HoleCards: [2]*poker.Card{
				{Rank: poker.Ace, Suit: poker.Clubs},
				{Rank: poker.Ten, Suit: poker.Hearts},
			},
			HasFolded: false,
			Chips:     10,
		}
	})

	Describe("PrintHand", func() {
		Context("when player has a hand", func() {
			It("prints the hand", func() {
				hand, err := player.PrintHoleCards()
				Expect(hand).To(Equal("A♣ 10♥"))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when player has no hand", func() {
			It("is an error", func() {
				player.HoleCards = [2]*poker.Card{}
				hand, err := player.PrintHoleCards()
				Expect(hand).To(Equal(""))
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("CanFold", func() {
		Context("when player has not folded", func() {
			It("can fold", func() {
				b := poker.BettingRound{
					Bets:       map[string]int{"1": 0},
					CallAmount: 20,
				}
				Expect(player.CanFold(&b)).To(BeTrue())
			})
		})

		Context("when player has folded", func() {
			It("cannot fold", func() {
				b := poker.BettingRound{
					Bets:       map[string]int{"1": 0},
					CallAmount: 20,
				}
				player.HasFolded = true
				Expect(player.CanFold(&b)).To(BeFalse())
			})
		})

		Context("when player has no chips", func() {
			It("cannot fold", func() {
				b := poker.BettingRound{
					Bets:       map[string]int{"1": 0},
					CallAmount: 20,
				}
				player.Chips = 0
				Expect(player.CanFold(&b)).To(BeFalse())
			})
		})

		Context("when no one has bet", func() {
			It("cannot fold", func() {
				b := poker.BettingRound{
					Bets:       map[string]int{"1": 0},
					CallAmount: 0,
				}
				Expect(player.CanFold(&b)).To(BeFalse())
			})
		})

		Context("when bet has been matched", func() {
			It("cannot fold", func() {
				b := poker.BettingRound{
					Bets:       map[string]int{"1": 20},
					CallAmount: 20,
				}
				Expect(player.CanFold(&b)).To(BeFalse())
			})
		})
	})
})
