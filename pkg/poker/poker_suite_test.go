package poker_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPoker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Poker Suite")
}
