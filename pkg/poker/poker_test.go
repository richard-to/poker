package poker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/richard-to/go-poker/pkg/poker"
)

var _ = Describe("Pot - GetTotal", func() {
	Context("when the pot is empty", func() {
		It("has a total pot of 0", func() {
			pot := poker.NewPot()
			Expect(pot.GetTotal()).To(Equal(0))
		})
	})

	Context("when the pot has chips", func() {
		It("has a total pot equal to the number of chips", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
				},
				{
					ID:   "2",
					Name: "Player 2",
				},
				{
					ID:   "3",
					Name: "Player 3",
				},
			}

			pot := poker.NewPot()
			pot.Bets[ps[0]] = 20
			pot.Bets[ps[1]] = 25
			pot.Bets[ps[2]] = 31

			Expect(pot.GetTotal()).To(Equal(76))
		})
	})
})

var _ = Describe("Pot - GetSidePots", func() {
	Context("when there is one main pot", func() {
		It("creates one side pot", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
				},
				{
					ID:   "2",
					Name: "Player 2",
				},
				{
					ID:   "3",
					Name: "Player 3",
				},
			}

			pot := poker.NewPot()
			pot.Bets[ps[0]] = 20
			pot.Bets[ps[1]] = 20
			pot.Bets[ps[2]] = 20

			sidePots := pot.GetSidePots()

			expectedSidePots := []*poker.SidePot{
				{
					Players: []*poker.Player{ps[0], ps[1], ps[2]},
					Total:   60,
					MaxBet:  20,
				},
			}
			Expect(sidePots).To(Equal(expectedSidePots))
		})
	})

	Context("when there is one main pot and a player has folded", func() {
		It("creates one side pot", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
				},
				{
					ID:        "2",
					Name:      "Player 2",
					HasFolded: true,
				},
				{
					ID:   "3",
					Name: "Player 3",
				},
			}

			pot := poker.NewPot()
			pot.Bets[ps[0]] = 20
			pot.Bets[ps[1]] = 5
			pot.Bets[ps[2]] = 20

			sidePots := pot.GetSidePots()

			expectedSidePots := []*poker.SidePot{
				{
					Players: []*poker.Player{ps[0], ps[2]},
					Total:   45,
					MaxBet:  20,
				},
			}
			Expect(sidePots).To(Equal(expectedSidePots))
		})
	})

	Context("when a player is all in", func() {
		It("creates two side pots", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
				},
				{
					ID:   "2",
					Name: "Player 2",
				},
				{
					ID:   "3",
					Name: "Player 3",
				},
			}

			pot := poker.NewPot()
			pot.Bets[ps[0]] = 20
			pot.Bets[ps[1]] = 6
			pot.Bets[ps[2]] = 20

			sidePots := pot.GetSidePots()

			expectedSidePots := []*poker.SidePot{
				{
					Players: []*poker.Player{ps[1], ps[0], ps[2]},
					Total:   18,
					MaxBet:  6,
				},
				{
					Players: []*poker.Player{ps[0], ps[2]},
					Total:   28,
					MaxBet:  14,
				},
			}
			Expect(sidePots).To(Equal(expectedSidePots))
		})
	})

	Context("when two players are all in with different stack sizes", func() {
		It("creates three side pots", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
				},
				{
					ID:   "2",
					Name: "Player 2",
				},
				{
					ID:   "3",
					Name: "Player 3",
				},
				{
					ID:   "4",
					Name: "Player 4",
				},
			}

			pot := poker.NewPot()
			pot.Bets[ps[0]] = 20
			pot.Bets[ps[1]] = 6
			pot.Bets[ps[2]] = 15
			pot.Bets[ps[3]] = 20

			sidePots := pot.GetSidePots()

			expectedSidePots := []*poker.SidePot{
				{
					Players: []*poker.Player{ps[1], ps[2], ps[0], ps[3]},
					Total:   24,
					MaxBet:  6,
				},
				{
					Players: []*poker.Player{ps[2], ps[0], ps[3]},
					Total:   27,
					MaxBet:  9,
				},
				{
					Players: []*poker.Player{ps[0], ps[3]},
					Total:   10,
					MaxBet:  5,
				},
			}
			Expect(sidePots).To(Equal(expectedSidePots))
		})
	})

	Context("when two players are all in with the same stack size", func() {
		It("creates two side pots", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
				},
				{
					ID:   "2",
					Name: "Player 2",
				},
				{
					ID:   "3",
					Name: "Player 3",
				},
				{
					ID:   "4",
					Name: "Player 4",
				},
			}

			pot := poker.NewPot()
			pot.Bets[ps[0]] = 20
			pot.Bets[ps[1]] = 6
			pot.Bets[ps[2]] = 6
			pot.Bets[ps[3]] = 20

			sidePots := pot.GetSidePots()

			expectedSidePots := []*poker.SidePot{
				{
					Players: []*poker.Player{ps[1], ps[2], ps[0], ps[3]},
					Total:   24,
					MaxBet:  6,
				},
				{
					Players: []*poker.Player{ps[0], ps[3]},
					Total:   28,
					MaxBet:  14,
				},
			}
			Expect(sidePots).To(Equal(expectedSidePots))
		})
	})

	Context("when two players are all in with the different stack sizes and multiple players have folded", func() {
		It("creates two side pots", func() {
			ps := []*poker.Player{
				{
					ID:   "1",
					Name: "Player 1",
				},
				{
					ID:   "2",
					Name: "Player 2",
				},
				{
					ID:   "3",
					Name: "Player 3",
				},
				{
					ID:        "4",
					Name:      "Player 4",
					HasFolded: true,
				},
				{
					ID:   "5",
					Name: "Player 5",
				},
				{
					ID:        "6",
					Name:      "Player 6",
					HasFolded: true,
				},
			}

			pot := poker.NewPot()
			pot.Bets[ps[0]] = 20
			pot.Bets[ps[1]] = 6
			pot.Bets[ps[2]] = 15
			pot.Bets[ps[3]] = 10
			pot.Bets[ps[4]] = 20
			pot.Bets[ps[5]] = 18

			sidePots := pot.GetSidePots()

			expectedSidePots := []*poker.SidePot{
				{
					Players: []*poker.Player{ps[1], ps[2], ps[0], ps[4]},
					Total:   36,
					MaxBet:  6,
				},
				{
					Players: []*poker.Player{ps[2], ps[0], ps[4]},
					Total:   40,
					MaxBet:  9,
				},
				{
					Players: []*poker.Player{ps[0], ps[4]},
					Total:   13,
					MaxBet:  5,
				},
			}
			Expect(sidePots).To(Equal(expectedSidePots))
		})
	})
})
