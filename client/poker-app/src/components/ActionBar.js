import classNames from 'classnames'
import { noop } from 'lodash'
import PropTypes from 'prop-types'
import React, { useState, useEffect } from 'react'

import AppPropTypes from '../AppPropTypes'
import { Event } from '../enums'

const cssButton = classNames(
  'flex-1',

  'bg-blue-600',
  'hover:bg-blue-700',

  // spacing
  'm-1',
  'p-2',

  // text
  'font-medium',
  'text-center',
  'text-sm',
  'text-white',
)

const cssRaiseWrap = classNames(
  'flex',
  'flex-col',
  'flex-1',

  'bg-gray-50',

  // spacing
  'm-1',
  'p-2',

  // text
  'font-medium',
  'text-center',
  'text-sm',
  'text-gray-900',
)

const ActionBar = ({ actions, callAmount, chipsInPot, maxRaiseAmount, minRaiseAmount, onAction, totalChips }) => {
  const [raiseByAmount, setRaiseByAmount] = useState(minRaiseAmount)

  useEffect(() => {
    setRaiseByAmount(minRaiseAmount)
  }, [minRaiseAmount])

  const onRaiseByAmount = (e) => {
    const value = parseInt(e.target.value)
    if (value) {
      setRaiseByAmount(value)
    }
  }

  const raiseToAmount = callAmount + raiseByAmount
  let callRemaining = callAmount - chipsInPot

  let callRemainingText = `ℝ${callRemaining}`
  if (callRemaining >= totalChips) {
    callRemainingText = 'ALL IN'
  }

  let raiseToAmountText = `ℝ${raiseToAmount}`

  if (raiseByAmount === maxRaiseAmount) {
    raiseToAmountText = 'ALL IN'
  }
  const actionButtons = actions.map(action => {
    if (action === Event.FOLD) {
      return (
        <button key={action} className={cssButton} onClick={() => onAction(action)}>
          {action.toUpperCase()}
        </button>
      )
    }

    if (action === Event.CHECK) {
      return (
        <button key={action} className={cssButton} onClick={() => onAction(action)}>
          {action.toUpperCase()}
        </button>
      )
    }

    if (action === Event.CALL) {
      return (
        <button key={action} className={cssButton} onClick={() => onAction(action)}>
          {action.toUpperCase()}<br />{callRemainingText}
        </button>
      )
    }

    return (
      <button key={action} className={cssButton} onClick={() => onAction(action, {value: raiseToAmount})}>
        RAISE TO<br />{raiseToAmountText}
      </button>
    )
  })

  return (
    <div className="flex bg-gray-800">
      {actionButtons}
      <div className={cssRaiseWrap}>
        <input
          type="range"
          min={minRaiseAmount}
          max={maxRaiseAmount}
          onChange={(e) => onRaiseByAmount(e)}
          value={raiseByAmount}
        />
        <input
          type="number"
          className="text-center border"
          min={minRaiseAmount}
          max={maxRaiseAmount}
          placeholder={minRaiseAmount}
          onChange={(e) => onRaiseByAmount(e)}
          value={raiseByAmount}
        />
      </div>
    </div>
  )
}

ActionBar.defaultProps = {
  onAction: noop,
}

ActionBar.propTypes = {
  actions: AppPropTypes.actions.isRequired,
  callAmount: PropTypes.number.isRequired,
  chipsInPot: PropTypes.number.isRequired,
  maxRaiseAmount: PropTypes.number.isRequired,
  minRaiseAmount: PropTypes.number.isRequired,
  onAction: PropTypes.func,
  totalChips: PropTypes.number.isRequired,
}

export default ActionBar
