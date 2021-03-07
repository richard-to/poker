import { has } from 'lodash'
import Peer from 'simple-peer'

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

const newPeer = (dispatch, params, ws, peers) => {
  if (has(peers, params.peerID)) {
    return
  }

  const peer = new Peer({initiator: true})
  peer.on('signal', data => {
    sendSignal(ws, params.peerID, data)
  })

  dispatch({
    type: actionTypes.WEBRTC.NEW_PEER,
    peer,
    peerID: params.peerID,
  })

}

const onReceiveSignal = (dispatch, params, ws, peers) => {
  if (params.signalData.type === "offer") {
    const peer = new Peer()
    peer.on('signal', data => {
      sendSignal(ws, params.peerID, data)
    })
    dispatch({
      type: actionTypes.WEBRTC.NEW_PEER,
      peer,
      peerID: params.peerID,
    })
    peer.signal(params.signalData)
  } else if (has(peers, params.peerID)) {
    peers[params.peerID].peer.signal(params.signalData)
  }
}

const sendSignal = (client, peerID, signalData) => {
  client.send(JSON.stringify({
    action: Event.SEND_SIGNAL,
    params: {
      peerID,
      signalData: signalData,
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

const onTakeSeat = (dispatch, params) => {
  navigator.mediaDevices.getUserMedia({
    video: true,
    audio: false,
  })
  .then(stream => {
    dispatch({
      type: actionTypes.SERVER.ON_TAKE_SEAT,
      seatID: params.seatID,
      userStream: stream,
    })
  })
  .catch(() => {}) // TODO Handle error
}

const sendMessage = (client, username, message) => {
  client.send(JSON.stringify({
    action: Event.SEND_MESSAGE,
    params: {
      username: username,
      message: message,
    },
  }))
}

const sendPlayerAction = (client, action, params = {}) => {
  client.send(JSON.stringify({
    action: action,
    params: params,
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

const updateGame = (dispatch, params) => {
  dispatch({
    type: actionTypes.GAME.UPDATE,
    gameState: params,
  })
}

export {
  // Errors
  error,

  // WebRTC
  newPeer,
  onReceiveSignal,
  sendSignal,

  // Game
  joinGame,
  newMessage,
  onJoinGame,
  onTakeSeat,
  sendMessage,
  sendPlayerAction,
  takeSeat,
  updateGame,
}
