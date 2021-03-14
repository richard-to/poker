import deepFreeze from 'deep-freeze-strict'

export const PlayerLocation = deepFreeze({
  BOTTOM: 'bottom',
  TOP: 'top',
})

export const Event = deepFreeze({
  CALL: 'call',
  CHECK: 'check',
  ERROR: 'error',
  FOLD: 'fold',
  JOIN: 'join',
  MUTE_VIDEO: 'mute-video',
  NEW_MESSAGE: 'new-message',
  NEW_PEER: 'new-peer',
  ON_JOIN: 'on-join',
  ON_RECEIVE_SIGNAL: 'on-receive-signal',
  ON_TAKE_SEAT: 'on-take-seat',
  RAISE: 'raise',
  SEND_MESSAGE: 'send-message',
  SEND_SIGNAL: 'send-signal',
  TAKE_SEAT: 'take-seat',
  UPDATE_GAME: 'update-game',
})

export const PlayerStatus = deepFreeze({
  ACTIVE: 'active',
  SITTING_OUT: 'sitting-out',
  VACATED: 'vacated',
})

export const Stage = deepFreeze({
  WAITING: 'Waiting',
  PREFLOP: 'Preflop',
  FLOP: 'Flop',
  TURN: 'Turn',
  RIVER: 'River',
  SHOWDOWN: 'Showdown',
})
