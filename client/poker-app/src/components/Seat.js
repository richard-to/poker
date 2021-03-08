import classNames from 'classnames'
import { motion, useAnimation } from "framer-motion"
import { noop } from 'lodash'
import PropTypes from 'prop-types'
import React, { useEffect, useRef } from 'react'

import AppPropTypes from '../AppPropTypes'
import { PlayerLocation, PlayerStatus, Stage } from '../enums'
import { getCardImage } from '../helpers'

const CARD_ANIM_VARIANTS = {
  deal: {
    opacity: [0, 1, 1],
    scale: [1, 1, 1],
    x: [500, 500, 0],
    y: [1000, 1000, 0],
  },
  initial: {
    opacity: 0,
    x: 500,
    y: 1000,
  },
  visible: {
    opacity: 1,
  },
}

const actionCss = classNames(
  'bg-purple-700',
  'border',
  'border-gray-900',
  'rounded-full',
  'shadow-md',

  // spacing
  'm-1',
  'px-2',
  'py-1',

  // text
  'font-medium',
  'text-gray-50',
  'text-xs',
)

const chipsInfoCss = classNames(
  'bg-gray-700',
  'bg-opacity-50',
  'border-1',
  'border-gray-800',
  'border-opacity-10',
  'rounded',
  'shadow-md',

  // spacing
  'm-1',
  'p-1',
  'lg:p-2',

  // text
  'text-center',
  'text-gray-100',
  'text-xs',
)

const getWrapCss = (location) => (
  classNames(
    {
      'flex-col': location === PlayerLocation.BOTTOM,
      'flex-col-reverse': location === PlayerLocation.TOP,
    },

    'flex',
    'flex-1',

    // spacing
    'sm:mx-2',
    'lg:mx-4',
    'my-2',
    'p-1',

    // text
    'text-xs',
  )
)

const getCardWrapCss = (location) => (
  classNames(
    {
      'flex-col': location === PlayerLocation.TOP,
      'flex-col-reverse': location === PlayerLocation.BOTTOM,
    },
    'flex',
    'flex-col',

    'w-5/12',
  )
)

const getNameOverlayCss = (player) => (
  classNames(
    {
      'bg-opacity-60': player.isActive,
      'bg-yellow-400': player.isActive,
      'border-yellow-500': player.isActive,
      'border-opacity-60': player.isActive,
      'text-gray-900': player.isActive,

      'bg-opacity-50': !player.isActive,
      'bg-gray-700': !player.isActive,
      'border-gray-800': !player.isActive,
      'border-opacity-10': !player.isActive,
      'text-gray-100': !player.isActive,
    },
    'border-1',
    'shadow-md',

    // position
    'absolute',
    'bottom-0',
    'left-0',
    'right-0',

    // spacing
    'p-1',

    // text
    'text-xs',
  )
)



const getButtonCss = (seatID) => (
  classNames(
    {
      'cursor-default': seatID !== null,
      'hover:bg-gray-800': seatID === null,
      'hover:bg-opacity-50': seatID === null,
      'hover:text-gray-200': seatID === null,
    },

    'flex-1',

    'bg-gray-700',
    'bg-opacity-50',
    'border-1',
    'border-gray-800',
    'border-opacity-10',
    'rounded',
    'shadow-md',

    // spacing
    'mb-2',
    'ml-1',
    'p-2',

    // text
    'text-gray-100',
    'text-xs',
  )
)


const getDealerCss = (player) => (
  classNames(
    {
      "invisible": !player.isDealer,
    },

    'bg-white',
    'border',
    'border-gray-500',
    'rounded-full',
    'shadow-md',

    // spacing
    'm-1',
    'px-2',
    'py-1',

    // text
    'font-medium',
    'text-sm',
  )
)

const getCardCss = (player) => (
  classNames(
    {
      'invisible': player.status !== PlayerStatus.ACTIVE || player.hasFolded,
    },
    'm-1',
  )
)

const getPlayerStatusMessage = (player) => {
  if (player.status === PlayerStatus.VACATED) {
    return 'Open Seat'
  } else if (player.status === PlayerStatus.SITTING_OUT) {
    return 'Waiting...'
  } else if (player.hasFolded) {
    return 'Folded'
  } else if (player.chips === 0) {
    return 'All In'
  }
  return `ℝ${player.chips}`
}

const Seat = ({
  dealDelay,
  location,
  onTakeSeat,
  player,
  seatID,
  stage,
  stream,
}) => {
  const videoRef = useRef(null)

  const card1Anim = useAnimation()
  const card2Anim = useAnimation()

  const playerActive = player.status === PlayerStatus.ACTIVE

  useEffect(() => {
    const dealSeq = async () => {
      if (!playerActive) {
        return
      }
      if (stage === Stage.PREFLOP) {
        card1Anim.set('initial')
        card2Anim.set('initial')
        card1Anim.start('deal', { delay: dealDelay[0], duration: 0.75, times: [0, 1]})
        card2Anim.start('deal', { delay: dealDelay[1], duration: 0.75, times: [0, 1]})
      } else if (stage !== Stage.PREFLOP) {
        card1Anim.set('deal')
        card2Anim.set('deal')
      }
    }
    dealSeq()
  }, [card1Anim, card2Anim, dealDelay, playerActive, stage])

  useEffect(() => {
    if (videoRef.current) {
      videoRef.current.srcObject = stream
    }
  }, [videoRef, stream])

  if (player.status === PlayerStatus.VACATED) {
    return (
      <div className={getWrapCss(location)}>
        <div className={getDealerCss(player)}>D</div>
        <button
          className={getButtonCss(seatID)}
          disabled={seatID !== null}
          onClick={() => onTakeSeat(player.id)}
        >
          {seatID === null ? 'Take Seat' : 'Open Seat'}
        </button>
      </div>
    )
  }


  return (
    <div className={getWrapCss(location)}>
      <div className="flex justify-center mb-2">
        <div className={getDealerCss(player)}>D</div>
        {player.chipsInPot > 0 && <div className={actionCss}>ℝ{player.chipsInPot}</div>}
      </div>
      <div className="flex">
          <div className="w-7/12">
            <div className="relative bg-black">
              <video className="rounded shadow-lg" ref={videoRef} autoPlay />
              <div className={getNameOverlayCss(player)}>
                <p>{player.name}</p>
              </div>
            </div>
          </div>
          <div className={getCardWrapCss(location)}>
            <p className={chipsInfoCss}>{getPlayerStatusMessage(player)}</p>
            <div className="flex justify-center items-end">
              <motion.div animate={card1Anim} className={getCardCss(player)} variants={CARD_ANIM_VARIANTS}>
                <img alt="Card" className="max-h-20" src={getCardImage(player.holeCards[0])} />
              </motion.div>
              <motion.div animate={card2Anim} className={getCardCss(player)} variants={CARD_ANIM_VARIANTS}>
                <img alt="Card" className="max-h-20" src={getCardImage(player.holeCards[1])} />
              </motion.div>
            </div>
          </div>
        </div>
    </div>
  )
}

Seat.defaultProps = {
  location: PlayerLocation.BOTTOM,
  onTakeSeat: noop,
}

Seat.propTypes = {
  dealDelay: PropTypes.arrayOf(PropTypes.number),
  location: PropTypes.oneOf([PlayerLocation.BOTTOM, PlayerLocation.TOP]),
  onTakeSeat: PropTypes.func,
  player: AppPropTypes.player.isRequired,
  seatID: PropTypes.string,
  stage: AppPropTypes.stage.isRequired,
  stream: PropTypes.any,
}

export default Seat
