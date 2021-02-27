import deepFreeze from 'deep-freeze-strict'

export const PlayerLocation = deepFreeze({
  BOTTOM: 'bottom',
  TOP: 'top',
})

export const Event = deepFreeze({
  CALL: 'call',
  CHECK: 'check',
  FOLD: 'fold',
  JOIN: 'join',
  NEW_MESSAGE: 'new-message',
  ON_JOIN: 'on-join',
  RAISE: 'raise',
  SEND_MESSAGE: 'send-message',
  UPDATE_GAME: 'update-game',
})

export const Stage = deepFreeze({
  PREFLOP: 'Preflop',
  FLOP: 'Flop',
  TURN: 'Turn',
  RIVER: 'River',
  SHOWDOWN: 'Showdown',
})
