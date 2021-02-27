import React, { useContext, useState } from 'react'

import { WebSocketContext } from '../WebSocket'

const Join = () => {
  const ws = useContext(WebSocketContext)

  const [isJoining, setJoining] = useState(false)
  const [username, setUsername] = useState('')

  const handleJoinGame = (username) => {
    const trimmedUsername = username.trim()
    if (trimmedUsername) {
      setJoining(true)
      ws.joinGame(username)
    }
  }

  const handleKeyPress = (e) => {
    if (e.code === 'Enter') {
      handleJoinGame(e.target.value)
    }
  }

  return (
    <div className="container mx-auto">
      <div className="flex h-screen">
        <div className="border shadow-lg m-auto p-10 font-bold text-black">
          <label className="sr-only" htmlFor="name">Name</label>
          <input
            autoFocus
            className="border rounded-sm mr-2 p-3"
            name="name"
            onChange={(e) => setUsername(e.target.value)}
            onKeyUp={handleKeyPress}
            placeholder="Jane Doe"
            type="text"
            value={username}
          />
          <button
            className="bg-blue-700 px-6 py-3 font-bold text-white"
            disabled={isJoining}
            onClick={() => handleJoinGame(username)}
          >
            Join
          </button>
	      </div>
      </div>
    </div>
  )
}

export default Join
