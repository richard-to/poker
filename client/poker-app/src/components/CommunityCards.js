import { motion, useAnimation } from "framer-motion"
import { concat } from 'lodash'
import PropTypes from 'prop-types'
import React, { useEffect } from 'react'

import AppPropTypes from '../AppPropTypes'
import { Stage } from '../enums'
import { getCardImage } from '../helpers'

const CARD_ANIM_VARIANTS = {
  deal: {
    opacity: 1,
    scale: 1,
    transition: { duration: 0.75 },
    x: 0,
    y: 0,
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

const CommunityCards = ({flop, turn, river, stage}) => {
  const flopAnim1 = useAnimation()
  const flopAnim2 = useAnimation()
  const flopAnim3 = useAnimation()
  const turnAnim = useAnimation()
  const riverAnim = useAnimation()

  const cardControls = [
    flopAnim1,
    flopAnim2,
    flopAnim3,
    turnAnim,
    riverAnim,
  ]

  useEffect(() => {
    const dealSeq = async () => {
      if (stage === Stage.FLOP) {
        flopAnim1.set('visible')
        await flopAnim1.start('deal')
        flopAnim2.set('visible')
        await flopAnim2.start('deal')
        flopAnim3.set('visible')
        await flopAnim3.start('deal')
      } else if (stage === Stage.TURN) {
        flopAnim1.set('deal')
        flopAnim2.set('deal')
        flopAnim3.set('deal')
        turnAnim.set('visible')
        await turnAnim.start('deal')
      } else if (stage === Stage.RIVER) {
        flopAnim1.set('deal')
        flopAnim2.set('deal')
        flopAnim3.set('deal')
        turnAnim.set('deal')
        riverAnim.set('visible')
        await riverAnim.start('deal')
      }
    }
    dealSeq()
  }, [stage, flopAnim1, flopAnim2, flopAnim3, turnAnim, riverAnim])

  const cards = concat([], flop, turn, river).map((card, index) => {
    if (card) {
      return (
        <div key={index} className="relative m-2">
          <div className="invisible">
            <img alt="Card" className="max-h-32" src="images/cards/1B.svg" />
          </div>
          <motion.div
            animate={cardControls[index]}
            className="absolute top-0 left-0"
            initial="initial"
            key={index}
            variants={CARD_ANIM_VARIANTS}
          >
            <img alt="Card" className="max-h-32" src={getCardImage(card)} />
          </motion.div>
        </div>
      )
    }
    return (
      <div className="m-2 invisible" key={index}>
        <img alt="Card" className="max-h-32" src="images/cards/1B.svg" />
      </div>
    )
  })
  return (
    <div className="flex justify-center">
      {cards}
    </div>
  )
}

CommunityCards.propTypes = {
  flop: PropTypes.arrayOf(AppPropTypes.card),
  turn: AppPropTypes.card,
  river: AppPropTypes.card,
  stage: AppPropTypes.stage.isRequired,
}

export default CommunityCards
