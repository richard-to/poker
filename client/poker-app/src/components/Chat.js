import classNames from 'classnames'
import { noop } from 'lodash'
import PropTypes from 'prop-types'
import React, { useState, useEffect, useRef } from 'react'

const cssChatLog = classNames(
  'overflow-y-auto',

  // dimensions
  'h-screen',
  'w-full',

  // spacing
  'p-3',

  // text
  'text-gray-700',
)

const Chat = ({messages, onSend}) => {
  const [message, setMessage] = useState('')
  const scrollRef = useRef()

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollIntoView(
        {
          behavior: 'smooth',
          block: 'end',
          inline: 'nearest',
        })
    }
  }, [messages])


  const handleKeyPress = (e) => {
    if (e.code === 'Enter') {
      onSend(e.target.value)
      setMessage('')
    }
  }

  return (
    <>
      <div className={cssChatLog}>
        {messages.map(message => (
          <div key={message.id} className="pb-1">
            <strong>{message.username}</strong>: {message.message}
          </div>
        ))}
        <div ref={scrollRef}></div>
      </div>
      <div className="flex">
        <input
          className="border p-3 w-full"
          name="name"
          onChange={(e) => setMessage(e.target.value)}
          onKeyUp={handleKeyPress}
          placeholder="Send message"
          type="text"
          value={message}
        />
      </div>
    </>
  )
}

Chat.defaultProps = {
  onSend: noop,
}

Chat.propTypes = {
  messages: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      message: PropTypes.string.isRequired,
      username: PropTypes.string.isRequired,
    }).isRequired,
  ),
  onSend: PropTypes.func,
}

export default Chat
