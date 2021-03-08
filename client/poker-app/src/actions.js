import { forEach, has } from 'lodash'
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

const newPeer = (dispatch, params, ws, appState) => {
  if (has(appState.peers, params.peerID)) {
    return
  }

  const peer = new Peer({initiator: true, stream: appState.userStream})
  peer.on('signal', data => {
    sendSignal(ws, params.peerID, data)
  })
  peer.on('connect', () => {
    console.log(`Peer (${params.peerID}) connected`)
  })
  peer.on('stream', stream => {
    console.log(`Peer (${params.peerID}) added new stream`)
    dispatch({
      type: actionTypes.WEBRTC.SET_PEER,
      peer,
      peerID: params.peerID,
      stream,
    })
  })
  peer.on('close', () => {
    console.log(`Peer (${params.peerID}) disconnected`)
    dispatch({
      type: actionTypes.WEBRTC.REMOVE_PEER,
      peerID: params.peerID,
    })
  })

  dispatch({
    type: actionTypes.WEBRTC.SET_PEER,
    peer,
    peerID: params.peerID,
    stream: null,
  })

}

const onReceiveSignal = (dispatch, params, ws, appState) => {
  if (params.signalData.type === "offer") {
    const peer = new Peer({stream: appState.userStream})
    peer.on('signal', data => {
      sendSignal(ws, params.peerID, data)
    })
    peer.on('connect', () => {
      console.log(`Peer (${params.peerID}) connected`)
    })
    peer.on('stream', stream => {
      console.log(`Peer (${params.peerID}) added new stream`)
      dispatch({
        type: actionTypes.WEBRTC.SET_PEER,
        peer,
        peerID: params.peerID,
        stream,
      })
    })
    peer.on('close', () => {
      console.log(`Peer (${params.peerID}) disconnected`)
      dispatch({
        type: actionTypes.WEBRTC.REMOVE_PEER,
        peerID: params.peerID,
      })
    })
    dispatch({
      type: actionTypes.WEBRTC.SET_PEER,
      peer,
      peerID: params.peerID,
      stream: null,
    })
    peer.signal(params.signalData)
  } else if (has(appState.peers, params.peerID)) {
    appState.peers[params.peerID].peer.signal(params.signalData)
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

const onTakeSeat = (dispatch, params, appState) => {
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
    console.log(appState.peers)
    forEach(appState.peers, c => c.peer.addStream(stream))
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
