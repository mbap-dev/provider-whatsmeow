import { useEffect, useRef } from 'react'

// Calls `callback` on an interval while `enabled`, skipping ticks when the tab
// is hidden. The latest callback is always used without resetting the timer.
export function usePolling(callback, intervalMs, enabled = true) {
  const saved = useRef(callback)
  saved.current = callback

  useEffect(() => {
    if (!enabled) return undefined
    const timer = window.setInterval(() => {
      if (!document.hidden) saved.current()
    }, intervalMs)
    return () => window.clearInterval(timer)
  }, [intervalMs, enabled])
}
