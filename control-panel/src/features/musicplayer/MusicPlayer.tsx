import React, { useEffect } from 'react'
import { useAppSelector, useAppDispatch } from '../../app/hooks'
import { setPlaying } from './musicPlayerSlice'

export function MusicPlayer() {
  const playerState = useAppSelector((state) => state.musicPlayer.isPlaying)
  const dispatch = useAppDispatch()

  useEffect(() => {
    dispatch(setPlaying(false))
  }, [])

  return (
    <div>
      <div>
      </div>
    </div>
  )
}