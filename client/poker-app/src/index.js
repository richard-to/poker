import React from 'react'
import ReactDOM from 'react-dom'

import './index.css'
import App from './App'
import { AppStateProvider } from './appStore'
import { WebSocketProvider } from './WebSocket'
import reportWebVitals from './reportWebVitals'

ReactDOM.render(
  <React.StrictMode>
    <AppStateProvider>
      <WebSocketProvider>
        <App />
      </WebSocketProvider>
    </AppStateProvider>
  </React.StrictMode>,
  document.getElementById('root')
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals(console.log)
