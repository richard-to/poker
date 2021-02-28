import classNames from 'classnames'
import { motion, useAnimation } from "framer-motion"
import { noop } from 'lodash'
import PropTypes from 'prop-types'
import React, { useEffect } from 'react'

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
  'text-sm',
)


const takeSeatCss = classNames(
  'bg-opacity-50',
  'bg-gray-700',
  'border-gray-800',
  'border-opacity-10',
  'text-gray-100',

  'border-2',
  'rounded-xl',
  'shadow-md',

  // spacing
  'p-1',

  // text
  'text-center',

  // hover
  'hover:bg-opacity-50',
  'hover:bg-gray-800',
  'hover:text-gray-200',
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
    'lg:mx-10',
    'my-2',
    'p-1',

    // text
    'text-sm',
  )
)

const getCardWrapCss = (player) => (
  classNames(
    {
      'invisible': player.hasFolded || player.status !== PlayerStatus.ACTIVE,
    },

    'flex',
    'justify-center',
  )
)

const getInfoCss = (player) => (
  classNames(
    {
      'bg-opacity-80': player.isActive,
      'bg-yellow-400': player.isActive,
      'border-yellow-500': player.isActive,
      'border-opacity-80': player.isActive,
      'text-gray-700': player.isActive,

      'bg-opacity-50': !player.isActive,
      'bg-gray-700': !player.isActive,
      'border-gray-800': !player.isActive,
      'border-opacity-10': !player.isActive,
      'text-gray-100': !player.isActive,
    },

    'border-1',
    'rounded-xl',
    'shadow-md',

    // spacing
    'p-1',

    // text
    'text-center',
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
      'invisible': player.status !== PlayerStatus.ACTIVE,
    },
    'm-1',
  )
)

const getPlayerStatusMessage = (player) => {
  if (player.status === PlayerStatus.VACATED) {
    return 'OPEN SEAT'
  } else if (player.status === PlayerStatus.SITTING_OUT) {
    return 'SITTING OUT'
  } else if (player.hasFolded) {
    return 'FOLDED'
  } else if (player.chips === 0) {
    return 'ALL IN'
  }
  return `ℝ${player.chips}`
}

const Seat = ({dealDelay, location, onTakeSeat, player, seatID, stage}) => {
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

  return (
    <div className={getWrapCss(location)}>
      <div className="flex justify-center mb-2">
        <div className={getDealerCss(player)}>D</div>
        {player.chipsInPot > 0 && <div className={actionCss}>ℝ{player.chipsInPot}</div>}
      </div>
      <div className={getCardWrapCss(player)}>
        <motion.div animate={card1Anim} className={getCardCss(player)} variants={CARD_ANIM_VARIANTS}>
          <img
            alt="Card"
            className="max-h-32"
            src={getCardImage(player.holeCards[0])}
            variants={CARD_ANIM_VARIANTS}
          />
        </motion.div>
        <motion.div animate={card2Anim} className={getCardCss(player)} variants={CARD_ANIM_VARIANTS}>
          <img alt="Card" className="max-h-32" src={getCardImage(player.holeCards[1])} />
        </motion.div>
      </div>

      {player.status === PlayerStatus.VACATED && seatID === null
        ? <button className={takeSeatCss} onClick={() => onTakeSeat(player.id)}>Take Seat</button>
        : <div className={getInfoCss(player)}>
            <p>{player.name}</p>
            <p>{getPlayerStatusMessage(player)}</p>
          </div>
      }
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
}

export default Seat
