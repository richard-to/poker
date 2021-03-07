import classNames from 'classnames'
import { noop } from 'lodash'
import PropTypes from 'prop-types'
import React, { useState } from 'react'
import Autosuggest from 'react-autosuggest'

const renderSuggestion = suggestion => <div>{suggestion.label}</div>

const RaiseInput = ({
  maxRaiseAmount,
  minRaiseAmount,
  onRaiseByAmount,
  placeholder,
  raiseSuggestions,
}) => {
  const [error, setError] = useState(false)
  const [value, setValue] = useState('')
  const [suggestions, setSuggestions] = useState([])

  const theme = {
    container: 'relative',
    input: classNames(
      {
        'focus:ring-blue-600': !error,
      },
      {
        'border-red-500': error,
        'focus:ring-red-400': error,
        'text-red-500': error,
      },

      'border',
      'w-full',

      // focus
      'focus:outline-none',
      'focus:ring-2',

      // spacing
      'p-3',
    ),
    suggestionsContainer: 'absolute bottom-14 left-0 right-0',
    suggestionsList: 'bg-white border-1 border-gray-300 rounded p-2 shadow-lg',
    suggestion: 'hover:bg-blue-100 cursor-pointer font-normal text-left p-1',
}

  const handleOnChange = (inputValue) => {
    const value = parseInt(inputValue)
    if (value > maxRaiseAmount) {
      onRaiseByAmount(maxRaiseAmount)
      setValue(maxRaiseAmount.toString())
      setError(false)
    } else if (value >= minRaiseAmount) {
      onRaiseByAmount(value)
      setValue(inputValue)
      setError(false)
    } else {
      setValue(inputValue)
      setError(true)
    }

  }

  const inputProps = {
    onChange: (event, { newValue }) => handleOnChange(newValue),
    placeholder,
    value,
  }

  return (
    <Autosuggest
      inputProps={inputProps}
      getSuggestionValue={suggestion => suggestion.value}
      onSuggestionsClearRequested={() => setSuggestions([])}
      onSuggestionsFetchRequested={() => setSuggestions(raiseSuggestions)}
      renderSuggestion={renderSuggestion}
      suggestions={suggestions}
      shouldRenderSuggestions={() => true}
      theme={theme}

    />
  )
}

RaiseInput.defaultProps = {
  onRaiseByAmount: noop,
  placeholder: 'Enter a raise by amount',
}

RaiseInput.propTypes = {
  maxRaiseAmount: PropTypes.number.isRequired,
  minRaiseAmount: PropTypes.number.isRequired,
  onRaiseByAmount: PropTypes.func.isRequired,
  placeholder: PropTypes.string.isRequired,
  raiseSuggestions: PropTypes.arrayOf(PropTypes.shape({
    label: PropTypes.string.isRequired,
    value: PropTypes.string.isRequired,
  })),
}

export default RaiseInput
