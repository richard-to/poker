import classNames from 'classnames'
import { noop } from 'lodash'
import PropTypes from 'prop-types'
import React, { useState } from 'react'

import AppPropTypes from '../AppPropTypes'
import { Event, Stage } from '../enums'
import RaiseInput from './RaiseInput'

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
/*
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
*/
const makeBetSizesSuggestions = (
  minBetAmount,
  minRaiseAmount,
  maxRaiseAmount,
  stage,
  totalChips,
  totalPot
) => {
  const bet3x = minBetAmount * 3
  if (stage === Stage.PREFLOP && minBetAmount === minRaiseAmount && bet3x < totalChips) {
    return [
      {label: 'Min bet', value: minRaiseAmount.toString()},
      {label: '3BB', value: bet3x.toString()},
      {label: 'All In', value: maxRaiseAmount.toString()},
    ]
  }
  const betSizes = [
    {label: 'Min raise', value: minRaiseAmount},
    {label: '1/4 Pot', value: Math.ceil(totalPot * 0.25)},
    {label: '1/3 Pot', value: Math.ceil(totalPot * 0.333)},
    {label: '1/2 Pot', value: Math.ceil(totalPot * 0.5)},
    {label: '2/3 Pot', value: Math.ceil(totalPot * 0.667)},
    {label: '3/4 Pot', value: Math.ceil(totalPot * 0.75)},
    {label: 'Pot', value: totalPot},
    {label: 'All In', value: maxRaiseAmount},
  ]
  return betSizes
    .filter(bet => bet.value >= minRaiseAmount && bet.value < totalChips)
    .map(bet => ({label: bet.label, value: bet.value.toString()}))
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
  const [raiseByAmount, setRaiseByAmount] = useState(minRaiseAmount)

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
  const placeholder = callAmount === 0 ? 'Enter a bet' : 'Enter a raise'

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


  const betSizeSuggestions = makeBetSizesSuggestions(
    minBetAmount,
    minRaiseAmount,
    maxRaiseAmount,
    stage,
    totalChips,
    totalPot,
  )

  return (
    <div className="flex bg-gray-800">
      {actionButtons}
      {showRaiseSlider &&
        <div className={raiseWrapCss}>
          <RaiseInput
            maxRaiseAmount={maxRaiseAmount}
            minRaiseAmount={minRaiseAmount}
            onRaiseByAmount={setRaiseByAmount}
            placeholder={placeholder}
            raiseSuggestions={betSizeSuggestions}
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
  minBetAmount: PropTypes.number.isRequired,
  minRaiseAmount: PropTypes.number.isRequired,
  onAction: PropTypes.func,
  stage: AppPropTypes.stage.isRequired,
  totalChips: PropTypes.number.isRequired,
  totalPot: PropTypes.number.isRequired,
}

export default ActionBar
