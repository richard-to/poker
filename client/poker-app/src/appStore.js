/**
 * The appStore use the createContext and useReducer hooks to create
 * a global state for the applications.
 *
 * The context and reducer hooks are used in place of redux reducers.
 *
 * See this article for more implementations details:
 *   - https://blog.logrocket.com/use-hooks-and-context-not-react-and-redux/
 */
import update from 'immutability-helper'
import PropTypes from 'prop-types'
import React, { createContext, useReducer } from 'react'

import actionTypes from './actionTypes'

const initialState = {
  chat: {
    messages: [],
  },
  error: null,
  gameState: null,
  seatID: null,
  streams: {},
  streamSeatMap: {},
  userID: null,
  username: null,
  userStream: null,
}

const appStore = createContext(initialState)
const { Provider } = appStore

const AppStateProvider = ({ children }) => {
  const [appState, dispatch] = useReducer((state, action) => {
    switch (action.type) {
      case actionTypes.CHAT.NEW_MESSAGE:
        return {
          ...state,
          chat: {
            messages: update(state.chat.messages, {$push: [action.message]})
          },
        }
      case actionTypes.GAME.UPDATE:
        return {
          ...state,
          gameState: action.gameState,
        }
      case actionTypes.SERVER.ERROR:
        return {
          ...state,
          error: action.error,
        }
      case actionTypes.SERVER.ON_JOIN:
        return {
          ...state,
          userID: action.userID,
          username: action.username,
        }
      case actionTypes.SERVER.ON_TAKE_SEAT:
        return {
          ...state,
          seatID: action.seatID,
          userStream: action.userStream,
        }
      case actionTypes.WEBRTC.REMOVE_STREAM:
        return {
          ...state,
          streams: update(state.streams, {$unset: [action.streamID]}),
        }
      case actionTypes.WEBRTC.SET_STREAM:
        return {
          ...state,
          streams: update(
            state.streams, {
              [action.streamID]: {$set: {
                peer: action.peer,
                peerID: action.peerID,
                stream: action.stream,
                streamID: action.streamID,
              }},
            },
          ),
        }
      case actionTypes.WEBRTC.SET_STREAM_SEAT_MAP:
        return {
          ...state,
          streamSeatMap: action.streamSeatMap,
        }
      default:
        throw new Error()
    }
  }, initialState)

  return <Provider value={{ appState, dispatch }}>{children}</Provider>
}

AppStateProvider.propTypes = {
  children: PropTypes.node.isRequired,
}

export { appStore, AppStateProvider }
