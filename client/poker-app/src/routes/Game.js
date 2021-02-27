import { zip } from 'lodash'
import React, { useContext, useEffect, useState } from 'react'

import { appStore } from '../appStore'
import ActionBar from '../components/ActionBar'
import Chat from '../components/Chat'
import CommunityCards from '../components/CommunityCards'
import Player from '../components/Player'
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
  const { chat, gameState } = appState

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
  console.log(cardDelay)
  return (
    <div className="container-fluid">
      <div className="flex h-screen bg-green-600">
        <div className="sm:w-3/4 h-screen flex flex-col overflow-y-auto">
            <div className="flex">
              <Player
                dealDelay={cardDelay[0]}
                player={gameState.players[0]}
                location={PlayerLocation.TOP}
                stage={gameState.stage}
              />
              <Player
                dealDelay={cardDelay[1]}
                player={gameState.players[1]}
                location={PlayerLocation.TOP}
                stage={gameState.stage}
              />
              <Player
                dealDelay={cardDelay[2]}
                player={gameState.players[2]}
                location={PlayerLocation.TOP}
                stage={gameState.stage}
              />
            </div>

            <div>
              <Pot amount={gameState.table.pot} />
              <CommunityCards
                flop={gameState.table.flop}
                turn={gameState.table.turn}
                river={gameState.table.river}
                stage={gameState.stage}
              />
            </div>

            <div>
              <div className="flex">
                <Player
                  dealDelay={cardDelay[5]}
                  player={gameState.players[5]}
                  stage={gameState.stage}
                />
                <Player
                  dealDelay={cardDelay[4]}
                  player={gameState.players[4]}
                  stage={gameState.stage}
                />
                <Player
                  dealDelay={cardDelay[3]}
                  player={gameState.players[3]}
                  stage={gameState.stage}
                />
              </div>
            </div>

            <div className="flex-1 flex flex-col-reverse">
              <ActionBar
                actions={gameState.actionBar.actions}
                callAmount={gameState.actionBar.callAmount}
                chipsInPot={gameState.actionBar.chipsInPot}
                maxRaiseAmount={gameState.actionBar.maxRaiseAmount}
                minRaiseAmount={gameState.actionBar.minRaiseAmount}
                onAction={ws.sendPlayerAction}
                totalChips={gameState.actionBar.totalChips}
              />
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
