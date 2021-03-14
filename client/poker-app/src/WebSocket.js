import { partial } from 'lodash'
import PropTypes from 'prop-types'
import React, { createContext, useContext, useState, useEffect } from 'react'
import { w3cwebsocket} from "websocket"

import {
  error,
  newMessage,
  onJoinGame,
  onReceiveSignal,
  onTakeSeat,
  joinGame,
  sendMessage,
  sendPlayerAction,
  takeSeat,
  updateGame,
} from './actions'
import { appStore } from './appStore'
import { Event } from './enums'

const BASE_WS_URL = process.env.REACT_APP_WEBSOCKET_URL

const WebSocketContext = createContext(null)

const WebSocketProvider = ({ children }) => {
  let ws

  const [client, setClient] = useState(null)

  const appContext = useContext(appStore)
  const { appState, dispatch } = appContext

  useEffect(() => {
    const _client = w3cwebsocket(BASE_WS_URL)
    _client.onerror = function() {
      error(dispatch, {error: 'Could not connect to the server.'})
    }

    _client.onopen = function() {
      console.log('WebSocket client connected')
    }

    _client.onclose = function() {
      error(dispatch, {error: 'Lost connection to the server.'})
    }

    setClient(_client)
  }, [dispatch])

  if (client) {
    client.onmessage = (payload) => {
      let event = JSON.parse(payload.data)
      if (event.action === Event.ERROR) {
        error(dispatch, event.params)
      } else if (event.action === Event.ON_JOIN) {
        onJoinGame(dispatch, event.params)
      } else if (event.action === Event.ON_TAKE_SEAT) {
        onTakeSeat(dispatch, event.params, client, appState)
      } else if (event.action === Event.NEW_MESSAGE) {
        newMessage(dispatch, event.params)
      } else if (event.action === Event.UPDATE_GAME) {
        updateGame(dispatch, event.params, client, appState)
      } else if (event.action === Event.ON_RECEIVE_SIGNAL) {
        onReceiveSignal(dispatch, event.params, client, appState)
      } else {
        error(dispatch, {error: 'Unknown message received.'})
      }
    }
  }

  ws = {
    client,
    joinGame: partial(joinGame, client),
    sendPlayerAction: partial(sendPlayerAction, client),
    sendMessage: partial(sendMessage, client, appState.username),
    takeSeat: partial(takeSeat, client),
  }

  return (
    <WebSocketContext.Provider value={ws}>
      {children}
    </WebSocketContext.Provider>
  )
}

WebSocketProvider.propTypes = {
  children: PropTypes.node.isRequired,
}

export { WebSocketContext, WebSocketProvider }
