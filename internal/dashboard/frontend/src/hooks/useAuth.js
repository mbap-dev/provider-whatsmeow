import { useCallback, useEffect, useState } from 'react'
import * as api from '../api/endpoints.js'

/**
 * Tracks Postgres-backed authentication state.
 *
 * authAvailable:
 *   null  → still checking (first render)
 *   false → backend has no auth system (no POSTGRES_URL) — render app as-is
 *   true  → auth enabled; user may be null (unauthenticated) or an object
 */
export function useAuth() {
  const [authAvailable, setAuthAvailable] = useState(null)
  const [user, setUser] = useState(null)

  const check = useCallback(async () => {
    try {
      const u = await api.getMe()
      setAuthAvailable(true)
      setUser(u)
    } catch (err) {
      // 404/405 → auth routes not registered (no Postgres configured)
      if (err.status === 404 || err.status === 405) {
        setAuthAvailable(false)
        setUser(null)
      } else {
        // 401 or anything else → auth enabled but not logged in
        setAuthAvailable(true)
        setUser(null)
      }
    }
  }, [])

  useEffect(() => {
    check()
  }, [check])

  const login = useCallback(
    async (email, password) => {
      await api.login(email, password)
      await check()
    },
    [check],
  )

  const logout = useCallback(async () => {
    try {
      await api.logout()
    } catch {
      // ignore — clear local state regardless
    }
    setUser(null)
  }, [])

  return { user, authAvailable, login, logout }
}
