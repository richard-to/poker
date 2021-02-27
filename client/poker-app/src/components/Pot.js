import classNames from 'classnames'
import PropTypes from 'prop-types'
import React from 'react'

const cssInfo = classNames(
  'bg-opacity-50',
  'bg-gray-700',
  'border-2',
  'border-gray-800',
  'border-opacity-10',
  'rounded-xl',

  // spacing
  'px-3',
  'py-1',

  // text
  'text-center',
  'text-gray-100',
  'text-sm',
)

const Pot = ({amount}) => (
  <div className="flex justify-center">
    <div className={cssInfo}>Main Pot ℝ{amount}</div>
  </div>
)

Pot.propTypes = {
  amount: PropTypes.number.isRequired,
}

export default Pot
