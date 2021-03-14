import PropTypes from 'prop-types'

import { Event, PlayerStatus, Stage } from './enums'

const actions = PropTypes.arrayOf(PropTypes.oneOf([
  Event.CALL,
  Event.CHECK,
  Event.FOLD,
  Event.RAISE,
]))

const card = PropTypes.shape({
  rank: PropTypes.number.isRequired,
  suit: PropTypes.number.isRequired,
})

const player = PropTypes.shape({
  chips: PropTypes.number,
  chipsInPot: PropTypes.number,
  holeCards: PropTypes.arrayOf(card),
  hasFolded: PropTypes.bool,
  isActive: PropTypes.bool,
  isDealer: PropTypes.bool,
  id: PropTypes.string,
  name: PropTypes.string,
  muted: PropTypes.bool.isRequired,
  status: PropTypes.oneOf([
    PlayerStatus.VACATED,
    PlayerStatus.SITTING_OUT,
    PlayerStatus.ACTIVE,
  ]),
})

const stage = PropTypes.oneOf([
  Stage.WAITING,
  Stage.PREFLOP,
  Stage.FLOP,
  Stage.TURN,
  Stage.RIVER,
  Stage.SHOWDOWN,
])

const propTypes = {
  actions,
  card,
  player,
  stage,
}

export default propTypes
