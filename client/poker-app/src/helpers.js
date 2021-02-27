export const getCardImage = (card) => {
  if (card === null) {
    return `images/cards/1B.svg`
  }
  const cardRank = ['2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A']
  const cardSuit = ['C', 'D', 'H', 'S']
  return `images/cards/${cardRank[card.rank]}${cardSuit[card.suit]}.svg`
}
