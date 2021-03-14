import { faMicrophone, faMicrophoneSlash } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import classNames from 'classnames'
import { noop } from 'lodash'
import PropTypes from 'prop-types'
import React from 'react'

const getMicButtonCss = (muted) => (
  classNames(
    {
      'text-gray-50': !muted,
      'text-red-500': muted,
    },

    'flex-1',

    'bg-gray-800',
    'hover:bg-gray-900',

    // Spacing
    'p-2',
  )
)


const OptionsBar = ({muted, onMuteVideo}) => {
  const icon = muted ? faMicrophoneSlash : faMicrophone
  const label = muted ? 'Unmute' : 'Mute'
  return (
    <div className="flex">
      <button className={getMicButtonCss(muted)} onClick={() => onMuteVideo(!muted)}>
        <FontAwesomeIcon icon={icon} /> {label}
      </button>
    </div>
  )
}

OptionsBar.defaultProps = {
  onMuteVideo: noop,
}

OptionsBar.propTypes = {
  muted: PropTypes.bool.isRequired,
  onMuteVideo: PropTypes.func,
}

export default OptionsBar
