import { find, forEach, keys, omit } from 'lodash'
import Peer from 'simple-peer'
import { v4 as uuidv4 } from 'uuid'

import actionTypes from './actionTypes'
import { Event } from './enums'

// Errors
const error = (dispatch, params) => {
  dispatch({
    type: actionTypes.SERVER.ERROR,
    error: params.error,
  })
}

// WebRTC
const initPeer = (dispatch, ws, initiator, stream, peerID, streamID) => {
  const peer = new Peer({initiator, stream})

  dispatch({
    type: actionTypes.WEBRTC.SET_STREAM,
    peer,
    peerID,
    peerStream: null,
    streamID,
  })

  peer.on('signal', data => {
    sendSignal(ws, peerID, streamID, data)
  })

  peer.on('connect', () => {
    console.log(`Stream (${streamID}) connected`)
  })

  peer.on('stream', peerStream => {
    console.log(`Stream (${streamID}) added new stream`)
    dispatch({
      type: actionTypes.WEBRTC.SET_STREAM,
      peer,
      peerID,
      stream: peerStream,
      streamID,
    })
  })

  peer.on('close', () => {
    console.log(`Stream (${streamID}) disconnected`)
    dispatch({
      type: actionTypes.WEBRTC.REMOVE_STREAM,
      streamID: streamID,
    })
  })
  peer.on('error', err => {
    console.log(`Stream (${streamID}) error: ${err}`)
    dispatch({
      type: actionTypes.WEBRTC.REMOVE_STREAM,
      streamID: streamID,
    })
  })

  return peer
}


const updatePeers = (dispatch, ws, clientSeatMap, userID, seatID, userStream, streams) => {
  if (!seatID || !userStream) {
    return
  }
  keys(clientSeatMap).forEach(clientID => {
    if (clientID === userID) {
      return
    }
    const stream = find(streams, s => s.streamID.startsWith(userID) && s.peerID === clientID)
    if (!stream) {
      const streamID = `${userID}-${uuidv4()}`
      initPeer(dispatch, ws, true, userStream, clientID, streamID)
    }
  })
}

const updateStreamMap = (dispatch, clientSeatMap) => {
  const streamSeatMap = {}
  forEach(clientSeatMap, (seatID, clientID) => {
    if (seatID) {
      streamSeatMap[seatID] = clientID
    }
  })

  dispatch({
    type: actionTypes.WEBRTC.SET_STREAM_SEAT_MAP,
    streamSeatMap,
  })
}

const onReceiveSignal = (dispatch, params, ws, appState) => {
  if (params.signalData.type === "offer") {
    const peer = initPeer(dispatch, ws, false, false, params.peerID, params.streamID)
    peer.signal(params.signalData)
  } else {
    const stream = find(appState.streams, s => s.streamID === params.streamID && s.peerID === params.peerID)
    if (stream) {
      stream.peer.signal(params.signalData)
    }
  }
}

const sendSignal = (client, peerID, streamID, signalData) => {
  client.send(JSON.stringify({
    action: Event.SEND_SIGNAL,
    params: {
      peerID,
      signalData,
      streamID,
    },
  }))
}

// Game

const joinGame = (client, username) => {
  client.send(JSON.stringify({
    action: Event.JOIN,
    params: {
      username,
    },
  }))
}

const newMessage = (dispatch, params) => {
  dispatch({
    type: actionTypes.CHAT.NEW_MESSAGE,
    message: {
      id: params.id,
      message: params.message,
      username: params.username,
    },
  })
}

const onJoinGame = (dispatch, params) => {
  dispatch({
    type: actionTypes.SERVER.ON_JOIN,
    userID: params.userID,
    username: params.username,
  })
}

const onTakeSeat = (dispatch, params, ws, appState) => {
  navigator.mediaDevices.getUserMedia({
    video: true,
    audio: process.env.REACT_APP_ENABLE_AUDIO === '1',
  })
  .then(stream => {
    dispatch({
      type: actionTypes.SERVER.ON_TAKE_SEAT,
      seatID: params.seatID,
      userStream: stream,
    })
    updateStreamMap(dispatch, params.clientSeatMap)
    updatePeers(
      dispatch,
      ws,
      params.clientSeatMap,
      appState.userID,
      params.seatID,
      stream,
      appState.streams,
    )
  })
  .catch(() => {}) // TODO Handle error
}

const sendMessage = (client, username, message) => {
  client.send(JSON.stringify({
    action: Event.SEND_MESSAGE,
    params: {
      username,
      message,
    },
  }))
}

const sendMuteVideo = (client, muted) => {
  client.send(JSON.stringify({
    action: Event.MUTE_VIDEO,
    params: {
      muted,
    },
  }))
}


const sendPlayerAction = (client, action, params = {}) => {
  client.send(JSON.stringify({
    action,
    params,
  }))
}

const takeSeat = (client, seatID) => {
  client.send(JSON.stringify({
    action: Event.TAKE_SEAT,
    params: {
      seatID,
    },
  }))
}

const updateGame = (dispatch, params, ws, appState) => {
  updateStreamMap(dispatch, params.clientSeatMap)
  updatePeers(
    dispatch,
    ws,
    params.clientSeatMap,
    appState.userID,
    appState.seatID,
    appState.userStream,
    appState.streams,
  )
  dispatch({
    type: actionTypes.GAME.UPDATE,
    gameState: omit(params, ['clientSeatMap']),
  })
}

export {
  // Errors
  error,

  // WebRTC
  onReceiveSignal,

  // Game
  joinGame,
  newMessage,
  onJoinGame,
  onTakeSeat,
  sendMessage,
  sendMuteVideo,
  sendPlayerAction,
  takeSeat,
  updateGame,
}
