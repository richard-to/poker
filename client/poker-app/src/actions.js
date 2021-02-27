import actionTypes from './actionTypes'
import { Event } from './enums'

const onJoinGame = (dispatch, params) => {
  dispatch({
    type: actionTypes.SERVER.ON_JOIN,
    userID: params.userID,
    username: params.username,
  })
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

const updateGame = (dispatch, params) => {
  dispatch({
    type: actionTypes.GAME.UPDATE,
    gameState: params,
  })
}

const joinGame = (client, username) => {
  client.send(JSON.stringify({
    action: Event.JOIN,
    params: {
      username,
    },
  }))
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

export {
  joinGame,
  newMessage,
  onJoinGame,
  sendMessage,
  sendPlayerAction,
  updateGame,
}
