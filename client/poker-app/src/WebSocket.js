import { partial } from 'lodash'
import PropTypes from 'prop-types'
import React, { createContext, useContext, useState, useEffect } from 'react'
import { w3cwebsocket} from "websocket"

import {
  onJoinGame,
  onTakeSeat,
  newMessage,
  joinGame,
  sendMessage,
  sendPlayerAction,
  takeSeat,
  updateGame,
} from './actions'
import { appStore } from './appStore'
import { Event } from './enums'

const BASE_WS_URL = 'ws://localhost:8000/ws'

const WebSocketContext = createContext(null)

const WebSocketProvider = ({ children }) => {
  let ws

  const [client, setClient] = useState(null)

  const appContext = useContext(appStore)
  const { appState, dispatch } = appContext

  useEffect(() => {
    const _client = w3cwebsocket(BASE_WS_URL)
    _client.onerror = function() {
      // TODO: connection error
      console.log('Connection error')
    }

    _client.onopen = function() {
      // TODO: client connected
      console.log('WebSocket client connected')
    }

    _client.onclose = function() {
      // TODO: client closed
      console.log('WebSocket client closed')
    }

    _client.onmessage = (payload) => {
      let event = JSON.parse(payload.data)
      if (event.action === Event.ON_JOIN) {
        onJoinGame(dispatch, event.params)
      } else if (event.action === Event.ON_TAKE_SEAT) {
        onTakeSeat(dispatch, event.params)
      } else if (event.action === Event.NEW_MESSAGE) {
        newMessage(dispatch, event.params)
      } else if (event.action === Event.UPDATE_GAME) {
        updateGame(dispatch, event.params)
      } else {
        // TODO: Handle unknown message
      }
    }
    setClient(_client)
  }, [dispatch])

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