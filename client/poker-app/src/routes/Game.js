import { has, zip } from 'lodash'
import React, { useContext, useEffect, useState } from 'react'

import { appStore } from '../appStore'
import ActionBar from '../components/ActionBar'
import Chat from '../components/Chat'
import CommunityCards from '../components/CommunityCards'
import Seat from '../components/Seat'
import Pot from '../components/Pot'
import { PlayerLocation, Stage } from '../enums'
import { WebSocketContext } from '../WebSocket'

const DELAY_INCREMENT = .15
const DEFAULT_CARD_DELAY = [
  [null, null],
  [null, null],
  [null, null],
  [null, null],
  [null, null],
  [null, null],
]

const getStream = (userID, player, clientSeatMap, peers, userStream) => {
  const peerID = clientSeatMap[player.id]
  if (userID === peerID) {
    return userStream
  } else if (has(peers, peerID)) {
    return peers[peerID].stream
  }
  return null
}

const Game = () => {
  const ws = useContext(WebSocketContext)

  const appContext = useContext(appStore)
  const { appState } = appContext
  const { chat, error, gameState, seatID, peers, userID, userStream } = appState
  console.log(peers)
  const [cardDelay, setCardDelay] = useState(DEFAULT_CARD_DELAY)
  const [newDeal, setNewDeal] = useState(true)

  const stage = (gameState) ? gameState.stage : null
  const players = (gameState) ? gameState.players : null

  useEffect(() => {
    if (newDeal === false && stage !== Stage.PREFLOP) {
      setNewDeal(true)
    }
  }, [newDeal, stage])

  useEffect(() => {
    if (newDeal && players && stage === Stage.PREFLOP) {
      let delay = 0

      const card1Delay = players.map((player) => {
        if (player.active) {
          const animDelay = delay
          delay += DELAY_INCREMENT
          return animDelay
        }
        return null
      })

      const card2Delay = players.map((player) => {
        if (player.active) {
          const animDelay = delay
          delay += DELAY_INCREMENT
          return animDelay
        }
        return null
      })
      setCardDelay(zip(card1Delay, card2Delay))
      setNewDeal(false)
    }
  }, [newDeal, players, stage])

  if (!gameState) {
    return <div className="container-fluid"></div>
  }

  const showActionBar = stage !== Stage.WAITING &&  seatID === gameState.actionBar.seatID

  return (
    <div className="container-fluid">
      {error &&
        <div className="flex w-full">
          <div className="flex-1 bg-red-500 shadow-md text-center p-2 text-gray-50">{error}</div>
        </div>
      }

      <div className="flex h-screen bg-green-600">
        <div className="sm:w-3/4 h-screen flex flex-col overflow-y-auto">
            <div className="flex">
              <Seat
                dealDelay={cardDelay[0]}
                onTakeSeat={ws.takeSeat}
                player={gameState.players[0]}
                location={PlayerLocation.TOP}
                seatID={seatID}
                stage={gameState.stage}
                stream={getStream(
                  userID,
                  gameState.players[0],
                  gameState.clientSeatMap,
                  peers,
                  userStream)}
              />
              <Seat
                dealDelay={cardDelay[1]}
                onTakeSeat={ws.takeSeat}
                player={gameState.players[1]}
                location={PlayerLocation.TOP}
                seatID={seatID}
                stage={gameState.stage}
                stream={getStream(
                  userID,
                  gameState.players[1],
                  gameState.clientSeatMap,
                  peers,
                  userStream)}
              />
              <Seat
                dealDelay={cardDelay[2]}
                onTakeSeat={ws.takeSeat}
                player={gameState.players[2]}
                location={PlayerLocation.TOP}
                seatID={seatID}
                stage={gameState.stage}
                stream={getStream(
                  userID,
                  gameState.players[2],
                  gameState.clientSeatMap,
                  peers,
                  userStream)}
              />
            </div>
            <div>
              {stage !== Stage.WAITING && <Pot amount={gameState.table.pot} />}

              <CommunityCards
                flop={gameState.table.flop}
                turn={gameState.table.turn}
                river={gameState.table.river}
                stage={gameState.stage}
              />
            </div>

            <div>
              <div className="flex">
                <Seat
                  dealDelay={cardDelay[5]}
                  onTakeSeat={ws.takeSeat}
                  player={gameState.players[5]}
                  seatID={seatID}
                  stage={gameState.stage}
                  stream={getStream(
                    userID,
                    gameState.players[5],
                    gameState.clientSeatMap,
                    peers,
                    userStream)}
                />
                <Seat
                  dealDelay={cardDelay[4]}
                  onTakeSeat={ws.takeSeat}
                  player={gameState.players[4]}
                  seatID={seatID}
                  stage={gameState.stage}
                  stream={getStream(
                    userID,
                    gameState.players[4],
                    gameState.clientSeatMap,
                    peers,
                    userStream)}
                />
                <Seat
                  dealDelay={cardDelay[3]}
                  onTakeSeat={ws.takeSeat}
                  player={gameState.players[3]}
                  seatID={seatID}
                  stage={gameState.stage}
                  stream={getStream(
                    userID,
                    gameState.players[3],
                    gameState.clientSeatMap,
                    peers,
                    userStream)}
                />
              </div>
            </div>

            <div className="flex-1 flex flex-col-reverse">
            {showActionBar && <ActionBar
                actions={gameState.actionBar.actions}
                callAmount={gameState.actionBar.callAmount}
                chipsInPot={gameState.actionBar.chipsInPot}
                maxRaiseAmount={gameState.actionBar.maxRaiseAmount}
                minBetAmount={gameState.actionBar.minBetAmount}
                minRaiseAmount={gameState.actionBar.minRaiseAmount}
                onAction={ws.sendPlayerAction}
                stage={gameState.stage}
                totalChips={gameState.actionBar.totalChips}
                totalPot={gameState.table.pot}
              />
            }
            </div>
          </div>

        <div className="hidden sm:flex flex-col w-1/4 bg-gray-50">
          <Chat messages={chat.messages} onSend={ws.sendMessage} />
        </div>
      </div>
    </div>
  )
}

export default Game
