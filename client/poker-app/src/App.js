import React, { useContext } from 'react'

import { appStore } from './appStore'
import Game from './routes/Game'
import Join from './routes/Join'


function App() {
  const appContext = useContext(appStore)
  const { appState } = appContext

  if (appState.userID) {
    return <Game />
  }

  return <Join />
}

export default App
