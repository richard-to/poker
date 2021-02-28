import classNames from 'classnames'
import { noop } from 'lodash'
import PropTypes from 'prop-types'
import React, { useState, useEffect } from 'react'

import AppPropTypes from '../AppPropTypes'
import { Event } from '../enums'

const buttonCss = classNames(
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

const raiseWrapCss = classNames(
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
    console.log(value)
    if (callAmount + value < minRaiseAmount &&  callAmount + value <= totalChips) {
      setRaiseByAmount(totalChips - callAmount)
    }
    if (value < minRaiseAmount) {
      return
    }

    if (value && callAmount + value <= totalChips) {
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

  if (raiseByAmount === maxRaiseAmount || raiseToAmount > totalChips) {
    raiseToAmountText = 'ALL IN'
  }

  const showRaiseSlider = actions.includes(Event.RAISE) && callAmount + minRaiseAmount < totalChips

  const actionButtons = actions.map(action => {
    if (action === Event.FOLD) {
      return (
        <button key={action} className={buttonCss} onClick={() => onAction(action)}>
          {action.toUpperCase()}
        </button>
      )
    }

    if (action === Event.CHECK) {
      return (
        <button key={action} className={buttonCss} onClick={() => onAction(action)}>
          {action.toUpperCase()}
        </button>
      )
    }

    if (action === Event.CALL) {
      return (
        <button key={action} className={buttonCss} onClick={() => onAction(action)}>
          {action.toUpperCase()}<br />{callRemainingText}
        </button>
      )
    }

    return (
      <button key={action} className={buttonCss} onClick={() => onAction(action, {value: raiseToAmount})}>
        RAISE TO<br />{raiseToAmountText}
      </button>
    )
  })

  return (
    <div className="flex bg-gray-800">
      {actionButtons}
      {showRaiseSlider &&
        <div className={raiseWrapCss}>
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
      }
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
