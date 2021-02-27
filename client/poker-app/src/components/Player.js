import classNames from 'classnames'
import { motion, useAnimation } from "framer-motion"
import PropTypes from 'prop-types'
import React, { useEffect } from 'react'

import AppPropTypes from '../AppPropTypes'
import { PlayerLocation, Stage } from '../enums'
import { getCardImage } from '../helpers'

const cssAction = classNames(
  'bg-purple-700',
  'border',
  'border-gray-900',
  'rounded-full',

  // spacing
  'm-1',
  'px-2',
  'py-1',

  // text
  'font-medium',
  'text-gray-50',
  'text-sm',
)

const Player = ({dealDelay, location, player, stage}) => {
  const variants = {
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

  const card1Anim = useAnimation()
  const card2Anim = useAnimation()

  useEffect(() => {
    const dealSeq = async () => {
      if (stage === Stage.PREFLOP) {
        // console.log(dealDelay[0])
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
  }, [card1Anim, card2Anim, dealDelay, stage])

  const cssWrap = classNames(
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
  const cssCardsWrap = classNames(
    {
      'invisible': player.hasFolded,
    },

    'flex',
    'justify-center',
  )
  const cssInfo = classNames(
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

    'border-2',
    'rounded-xl',

    // spacing
    'p-1',

    // text
    'text-center',
  )

  const cssDealer = classNames(
    {
      "invisible": !player.isDealer,
    },

    'bg-white',
    'border',
    'border-gray-500',
    'rounded-full',

    // spacing
    'm-1',
    'px-2',
    'py-1',

    // text
    'font-medium',
    'text-sm',
  )

  const cssCard = classNames(
    {
      'invisible': !player.active,
    },
    'm-1',
  )

  let status = `ℝ${player.chips}`
  if (!player.active) {
    status = 'SITTING OUT'
  } else if (player.hasFolded) {
    status = 'FOLDED'
  } else if (player.chips === 0) {
    status = 'ALL IN'
  }

  return (
    <div className={cssWrap}>
      <div className="flex justify-center mb-2">
        <div className={cssDealer}>D</div>
        {player.chipsInPot > 0 && <div className={cssAction}>ℝ{player.chipsInPot}</div>}
      </div>
      <div className={cssCardsWrap}>
        <motion.div animate={card1Anim} className={cssCard} variants={variants}>
          <img
            alt="Card"
            className="max-h-32"
            src={getCardImage(player.holeCards[0])}
            variants={variants}
          />
        </motion.div>
        <motion.div animate={card2Anim} className={cssCard} variants={variants}>
          <img alt="Card" className="max-h-32" src={getCardImage(player.holeCards[1])} />
        </motion.div>
      </div>
      <div className={cssInfo}>
        <p>{player.name}</p>
        <p>{status}</p>
      </div>
    </div>
  )
}

Player.defaultProps = {
  location: PlayerLocation.BOTTOM,
}

Player.propTypes = {
  dealDelay: PropTypes.arrayOf(PropTypes.number),
  location: PropTypes.oneOf([PlayerLocation.BOTTOM, PlayerLocation.TOP]),
  player: AppPropTypes.player.isRequired,
  stage: AppPropTypes.stage.isRequired,
}

export default Player
