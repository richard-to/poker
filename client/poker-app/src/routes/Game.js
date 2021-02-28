import { zip } from 'lodash'
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
const Game = () => {
  const ws = useContext(WebSocketContext)

  const appContext = useContext(appStore)
  const { appState } = appContext
  const { chat, gameState, seatID } = appState

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
    return (
      <div className="container-fluid">
        <div className="flex h-screen bg-green-600">
          <div className="sm:w-3/4 h-screen flex flex-col overflow-y-auto">LOADING</div>
        </div>
      </div>
    )
  }

  return (
    <div className="container-fluid">
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
              />
              <Seat
                dealDelay={cardDelay[1]}
                onTakeSeat={ws.takeSeat}
                player={gameState.players[1]}
                location={PlayerLocation.TOP}
                seatID={seatID}
                stage={gameState.stage}
              />
              <Seat
                dealDelay={cardDelay[2]}
                onTakeSeat={ws.takeSeat}
                player={gameState.players[2]}
                location={PlayerLocation.TOP}
                seatID={seatID}
                stage={gameState.stage}
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
                />
                <Seat
                  dealDelay={cardDelay[4]}
                  onTakeSeat={ws.takeSeat}
                  player={gameState.players[4]}
                  seatID={seatID}
                  stage={gameState.stage}
                />
                <Seat
                  dealDelay={cardDelay[3]}
                  onTakeSeat={ws.takeSeat}
                  player={gameState.players[3]}
                  seatID={seatID}
                  stage={gameState.stage}
                />
              </div>
            </div>

            <div className="flex-1 flex flex-col-reverse">
            {stage !== Stage.WAITING && <ActionBar
                actions={gameState.actionBar.actions}
                callAmount={gameState.actionBar.callAmount}
                chipsInPot={gameState.actionBar.chipsInPot}
                maxRaiseAmount={gameState.actionBar.maxRaiseAmount}
                minRaiseAmount={gameState.actionBar.minRaiseAmount}
                onAction={ws.sendPlayerAction}
                totalChips={gameState.actionBar.totalChips}
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