import { useCallback, useEffect, useMemo, useState } from 'react'
import { Button, SignalDot, ToastStack, cx } from './components/ui.jsx'
import { Home } from './components/Home.jsx'
import { Instances } from './components/Instances.jsx'
import { QRModal } from './components/QRModal.jsx'
import { LoginPage } from './components/LoginPage.jsx'
import { UserMenu } from './components/UserMenu.jsx'
import { APIKeysModal } from './components/APIKeysModal.jsx'
import { useTheme } from './hooks/useTheme.js'
import { useToasts } from './hooks/useToasts.js'
import { usePolling } from './hooks/usePolling.js'
import { useAuth } from './hooks/useAuth.js'
import * as api from './api/endpoints.js'
import { formatTime, isOnline, sortSessions } from './lib/sessions.js'

const LOGO_LIGHT = 'https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_black.svg'
const LOGO_DARK = 'https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_white.svg'
const REFRESH_MS = 7000

const CONFIRM = {
  disconnect: (id) => `Desconectar a sessão "${id}"?`,
  delete: (id) =>
    `Excluir permanentemente "${id}"? Isso remove o diretório da sessão, a atribuição, overrides e webhooks.`,
}

export function App() {
  const { theme, toggleTheme } = useTheme()
  const { toasts, notify, dismiss } = useToasts()
  const { user, authAvailable, login, logout } = useAuth()

  const [worker, setWorker] = useState(null)
  const [sessions, setSessions] = useState([])
  const [loading, setLoading] = useState(true)
  const [lastUpdated, setLastUpdated] = useState(null)
  const [page, setPage] = useState('home')
  const [selectedId, setSelectedId] = useState('')
  const [detailTab, setDetailTab] = useState('overview')
  const [busy, setBusy] = useState({})
  const [qrSessionId, setQrSessionId] = useState('')
  const [apiKeysOpen, setApiKeysOpen] = useState(false)

  const refresh = useCallback(
    async ({ silent = false } = {}) => {
      if (!silent) setLoading(true)
      try {
        const [workerResp, sessionsResp] = await Promise.all([
          api.getWorker(),
          api.listSessions(),
        ])
        setWorker(workerResp)
        setSessions(Array.isArray(sessionsResp?.sessions) ? sessionsResp.sessions : [])
        setLastUpdated(new Date())
      } catch (error) {
        notify(error.message, 'error')
      } finally {
        if (!silent) setLoading(false)
      }
    },
    [notify],
  )

  useEffect(() => {
    refresh()
  }, [refresh])

  usePolling(() => refresh({ silent: true }), REFRESH_MS, true)

  useEffect(() => {
    document.title = worker?.ui_title || 'WhatsMeow Control Center'
  }, [worker])

  // Keep a valid selection as the session list changes.
  useEffect(() => {
    if (sessions.length === 0) {
      setSelectedId('')
      return
    }
    if (selectedId && sessions.some((s) => s.session_id === selectedId)) return
    setSelectedId([...sessions].sort(sortSessions)[0]?.session_id || '')
  }, [sessions, selectedId])

  const counts = useMemo(
    () => ({
      total: sessions.length,
      online: sessions.filter(isOnline).length,
      withQR: sessions.filter((s) => s.has_qr).length,
      assigned: sessions.filter((s) => s.assigned).length,
    }),
    [sessions],
  )

  const setBusyFor = (key, value) =>
    setBusy((current) => {
      const next = { ...current }
      if (value) next[key] = true
      else delete next[key]
      return next
    })

  const handleAction = useCallback(
    async (action, sessionId) => {
      const id = String(sessionId || '').trim()
      if (!id) return
      if (CONFIRM[action] && !window.confirm(CONFIRM[action](id))) return

      const key = `${action}:${id}`
      setBusyFor(key, true)
      try {
        if (action === 'assign') await api.assignSession(id)
        else if (action === 'unassign') await api.unassignSession(id)
        else if (action === 'delete') await api.deleteSession(id)
        else await api.sessionAction(id, action)

        notify(
          action === 'delete' ? `Sessão ${id} excluída.` : `Ação "${action}" aplicada a ${id}.`,
          'success',
        )
        await refresh({ silent: true })
      } catch (error) {
        notify(error.message, 'error')
      } finally {
        setBusyFor(key, false)
      }
    },
    [notify, refresh],
  )

  const handleSelect = useCallback((id) => {
    setSelectedId(id)
    setDetailTab('overview')
    setPage('instances')
  }, [])

  const handleCreated = useCallback(
    async (id) => {
      setSelectedId(id)
      setDetailTab('overview')
      setPage('instances')
      await refresh({ silent: true })
    },
    [refresh],
  )

  // Auth gate: null = still checking, show blank; true+no user = show login
  if (authAvailable === null) return <div className="min-h-screen bg-surface" />
  if (authAvailable === true && user === null) {
    return <LoginPage onLogin={login} />
  }

  const navItems = [
    { id: 'home', label: 'Home' },
    { id: 'instances', label: 'Instâncias' },
  ]

  return (
    <div className="min-h-screen">
      <header className="sticky top-0 z-40 border-b border-border bg-surface/80 backdrop-blur">
        <div className="mx-auto flex h-16 max-w-7xl items-center gap-4 px-4 sm:px-6">
          <div className="flex items-center gap-2.5">
            <img
              src={theme === 'dark' ? LOGO_DARK : LOGO_LIGHT}
              alt="Envolve"
              className="h-6 w-auto"
            />
            <span className="hidden font-display text-sm font-semibold text-fg sm:inline">
              Control Center
            </span>
          </div>

          <nav className="ml-2 flex items-center gap-1">
            {navItems.map((item) => (
              <button
                key={item.id}
                onClick={() => setPage(item.id)}
                className={cx(
                  'rounded-lg px-3 py-1.5 text-sm font-medium transition-colors',
                  page === item.id ? 'bg-surface-2 text-fg' : 'text-muted hover:text-fg',
                )}
              >
                {item.label}
              </button>
            ))}
          </nav>

          <div className="ml-auto flex items-center gap-3">
            <span className="hidden items-center gap-2 text-xs text-muted md:flex">
              <SignalDot tone={counts.online > 0 ? 'live' : 'down'} pulse={counts.online > 0} />
              {counts.online}/{counts.total} online
            </span>
            <span className="hidden font-mono text-[0.72rem] text-faint lg:inline">
              ↻ {formatTime(lastUpdated)}
            </span>
            <Button size="sm" variant="ghost" onClick={() => refresh()} loading={loading}>
              Atualizar
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={toggleTheme}
              aria-label="Alternar tema"
            >
              {theme === 'dark' ? 'Claro' : 'Escuro'}
            </Button>
            {authAvailable && user && (
              <UserMenu
                user={user}
                onLogout={logout}
                onOpenAPIKeys={() => setApiKeysOpen(true)}
              />
            )}
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-4 py-6 sm:px-6">
        {page === 'home' ? (
          <Home
            worker={worker}
            counts={counts}
            notify={notify}
            onCreated={handleCreated}
            onGoToInstances={() => setPage('instances')}
          />
        ) : (
          <Instances
            sessions={sessions}
            worker={worker}
            selectedId={selectedId}
            detailTab={detailTab}
            busy={busy}
            onSelect={handleSelect}
            onAction={handleAction}
            onShowQR={setQrSessionId}
            notify={notify}
          />
        )}
      </main>

      <QRModal sessionId={qrSessionId} onClose={() => setQrSessionId('')} />
      <APIKeysModal open={apiKeysOpen} onClose={() => setApiKeysOpen(false)} />
      <ToastStack toasts={toasts} onDismiss={dismiss} />
    </div>
  )
}
