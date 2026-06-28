import { useCallback, useRef, useState } from 'react'

const TTL_MS = 3600

export function useToasts() {
  const [toasts, setToasts] = useState([])
  const timers = useRef(new Map())

  const dismiss = useCallback((id) => {
    setToasts((current) => current.filter((toast) => toast.id !== id))
    const timer = timers.current.get(id)
    if (timer) {
      window.clearTimeout(timer)
      timers.current.delete(id)
    }
  }, [])

  const notify = useCallback(
    (message, tone = 'info') => {
      const id = `${Date.now()}-${Math.random().toString(16).slice(2)}`
      setToasts((current) => [...current, { id, message, tone }])
      timers.current.set(
        id,
        window.setTimeout(() => dismiss(id), TTL_MS),
      )
      return id
    },
    [dismiss],
  )

  return { toasts, notify, dismiss }
}
