import classNames from 'classnames'
import { noop } from 'lodash'
import PropTypes from 'prop-types'
import React, { useState, useEffect } from 'react'

import AppPropTypes from '../AppPropTypes'
import { Event, Stage } from '../enums'

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

const getRaiseInputCss = (error) => (
  classNames(
    {
      'border-black': !error,
      'focus:ring-blue-600': !error,
      'text-black': !error,
    },
    {
      'border-red-500': error,
      'focus:ring-red-400': error,
      'text-red-500': error,
    },

    'border',

    // focus
    'focus:outline-none',
    'focus:ring-2',

    // text
    'text-center',
  )
)

const calcBetSizes = (minBetAmount, minRaiseAmount, maxRaiseAmount, stage, totalChips, totalPot) => {
  const bet3x = minBetAmount * 3
  if (stage === Stage.PREFLOP && minBetAmount === minRaiseAmount && bet3x < totalChips) {
    return [minRaiseAmount, bet3x, maxRaiseAmount]
  }
  const betSizes = [
    minRaiseAmount,
    Math.ceil(totalPot * 0.333),
    Math.ceil(totalPot * 0.667),
    totalPot,
    maxRaiseAmount,
  ]
  return betSizes.filter(bet => bet > minBetAmount && bet < totalChips)
}

const ActionBar = ({
    actions,
    callAmount,
    chipsInPot,
    maxRaiseAmount,
    minBetAmount,
    minRaiseAmount,
    onAction,
    stage,
    totalChips,
    totalPot,
}) => {
  const [raiseInputError, setRaiseInputError] = useState(false)
  const [raiseInput, setRaiseInput] = useState(minRaiseAmount)
  const [raiseByAmount, setRaiseByAmount] = useState(minRaiseAmount)

  useEffect(() => {
    setRaiseByAmount(minRaiseAmount)
    setRaiseInput(minRaiseAmount)
    setRaiseInputError(false)
  }, [minRaiseAmount])

  const onRaiseInputEntered = (e) => {
    if ((e.type === 'keyup' && e.code === 'Enter') || e.type === 'blur') {
      onRaiseByAmount(e)
    }
  }
  const onRaiseByAmount = (e) => {
    const value = parseInt(e.target.value)
    if (value > maxRaiseAmount) {
      setRaiseByAmount(maxRaiseAmount)
      setRaiseInput(maxRaiseAmount)
      setRaiseInputError(false)
    } else if (value >= minRaiseAmount) {
      setRaiseByAmount(value)
      setRaiseInput(value)
      setRaiseInputError(false)
    } else {
      setRaiseInputError(true)
    }
  }

  const raiseToAmount = callAmount + raiseByAmount
  const callRemaining = callAmount - chipsInPot

  let callRemainingLabel = `ℝ${callRemaining}`
  if (callRemaining >= totalChips) {
    callRemainingLabel = 'ALL IN'
  }

  let raiseToAmountLabel = `ℝ${raiseToAmount}`
  if (raiseByAmount === maxRaiseAmount || callRemaining + raiseByAmount > totalChips) {
    raiseToAmountLabel = 'ALL IN'
  }

  const raiseLabel = callAmount === 0 ? 'BET' : 'RAISE TO'

  const showRaiseSlider = actions.includes(Event.RAISE) && callRemaining + minRaiseAmount < totalChips

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
          {action.toUpperCase()}<br />{callRemainingLabel}
        </button>
      )
    }

    return (
      <button key={action} className={buttonCss} onClick={() => onAction(action, {value: raiseToAmount})}>
        {raiseLabel}<br />{raiseToAmountLabel}
      </button>
    )
  })


  const betSizes = calcBetSizes(minBetAmount, minRaiseAmount, maxRaiseAmount, stage, totalChips, totalPot)
  const betSizeOptions = betSizes.map(betSize => <option key={betSize} value={betSize} />)

  return (
    <div className="flex bg-gray-800">
      {actionButtons}
      {showRaiseSlider &&
        <div className={raiseWrapCss}>
          <input
            className="mb-2"
            type="range"
            list="bet-options"
            min={minRaiseAmount}
            max={maxRaiseAmount}
            onChange={(e) => onRaiseByAmount(e)}
            value={raiseByAmount}
          />
          <input
            type="number"
            className={getRaiseInputCss(raiseInputError)}
            min={minRaiseAmount}
            max={maxRaiseAmount}
            onBlur={(e) => onRaiseInputEntered(e)}
            onChange={(e) => setRaiseInput(e.target.value)}
            onKeyUp={(e) => onRaiseInputEntered(e)}
            placeholder={minRaiseAmount}
            value={raiseInput}
          />
          <datalist id="bet-options">
            {betSizeOptions}
          </datalist>
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
  minBetAmount: PropTypes.number.isRequired,
  minRaiseAmount: PropTypes.number.isRequired,
  onAction: PropTypes.func,
  stage: AppPropTypes.stage.isRequired,
  totalChips: PropTypes.number.isRequired,
  totalPot: PropTypes.number.isRequired,
}

export default ActionBar
