import React, { useEffect, useRef, useState } from 'react'
import { createRoot } from 'react-dom/client'
import htm from 'htm'
import './styles.css'

const LOGO_LIGHT = 'https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_black.svg'
const LOGO_DARK = 'https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_white.svg'
const REFRESH_INTERVAL_MS = 7000
const QR_REFRESH_MS = 3000
const THEME_STORAGE_KEY = 'whatsmeow-dashboard-theme'

const html = htm.bind(React.createElement)

function enc(value) {
  return encodeURIComponent(value)
}

function getInitialTheme() {
  const stored = localStorage.getItem(THEME_STORAGE_KEY)
  if (stored === 'light' || stored === 'dark') {
    return stored
  }
  if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    return 'dark'
  }
  return 'light'
}

async function api(path, options = {}) {
  const request = { ...options }
  const headers = new Headers(request.headers || {})

  if (request.body && !(request.body instanceof FormData) && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }

  const response = await fetch(path, {
    ...request,
    headers,
  })

  if (!response.ok) {
    let message = `${response.status} ${response.statusText}`
    try {
      const contentType = response.headers.get('content-type') || ''
      if (contentType.includes('application/json')) {
        const payload = await response.json()
        message =
          payload?.error?.message ||
          payload?.message ||
          payload?.error ||
          JSON.stringify(payload)
      } else {
        const text = await response.text()
        if (text) {
          message = text
        }
      }
    } catch (_) {
      // ignored
    }
    throw new Error(message)
  }

  const contentType = response.headers.get('content-type') || ''
  return contentType.includes('application/json') ? response.json() : null
}

function makeToastId() {
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`
}

function createDefaultConfigForm() {
  return {
    reject_calls: false,
    reject_calls_message: '',
    mark_read_on_message: false,
    ignore_status_broadcast: false,
    ignore_newsletters: false,
    outgoing_messages_coex: false,
    history_sync_enabled: false,
    history_sync_days: '7',
    always_online: false,
    always_online_interval_seconds: '60',
    proxy_enabled: false,
    proxy_protocol: '',
    proxy_host: '',
    proxy_port: '',
    proxy_username: '',
    proxy_password: '',
  }
}

function createConfigFormFromSnapshot(snapshot) {
  const cfg = snapshot?.config || {}
  return {
    reject_calls: Boolean(cfg.reject_calls),
    reject_calls_message: cfg.reject_calls_message || '',
    mark_read_on_message: Boolean(cfg.mark_read_on_message),
    ignore_status_broadcast: Boolean(cfg.ignore_status_broadcast),
    ignore_newsletters: Boolean(cfg.ignore_newsletters),
    outgoing_messages_coex: Boolean(cfg.outgoing_messages_coex),
    history_sync_enabled: Boolean(cfg.history_sync_enabled),
    history_sync_days: String(
      Number.isFinite(Number(cfg.history_sync_days)) && Number(cfg.history_sync_days) > 0
        ? Number(cfg.history_sync_days)
        : 7
    ),
    always_online: Boolean(cfg.always_online),
    always_online_interval_seconds: String(
      Number.isFinite(Number(cfg.always_online_interval_seconds))
        ? Number(cfg.always_online_interval_seconds)
        : 60
    ),
    proxy_enabled: Boolean(cfg.proxy_enabled),
    proxy_protocol: cfg.proxy_protocol || '',
    proxy_host: cfg.proxy_host || '',
    proxy_port: cfg.proxy_port || '',
    proxy_username: cfg.proxy_username || '',
    proxy_password: cfg.proxy_password || '',
  }
}

function webhookFormToPayload(form) {
  const timeout = parseInt(form.timeout_ms, 10)
  return {
    id: form.id.trim(),
    name: form.name.trim(),
    url: form.url.trim(),
    url_absolute: form.url_absolute.trim(),
    header: form.header.trim(),
    token: form.token.trim(),
    timeout_ms: Number.isFinite(timeout) && timeout > 0 ? timeout : 5000,
    enabled: form.enabled,
    append_session: form.append_session,
    send_incoming_messages: form.send_incoming_messages,
    send_outgoing_messages: form.send_outgoing_messages,
    send_update_messages: form.send_update_messages,
    send_group_messages: form.send_group_messages,
    send_newsletter_messages: form.send_newsletter_messages,
  }
}

function historyConfigPayload(enabled, days) {
  const parsedDays = parseInt(days, 10)
  return {
    history_sync_enabled: Boolean(enabled),
    history_sync_days: Number.isFinite(parsedDays) && parsedDays > 0 ? parsedDays : 0,
  }
}

function createDefaultWebhookForm() {
  return {
    id: '',
    name: '',
    url: '',
    url_absolute: '',
    header: '',
    token: '',
    timeout_ms: '5000',
    enabled: true,
    append_session: true,
    send_incoming_messages: true,
    send_outgoing_messages: true,
    send_update_messages: true,
    send_group_messages: true,
    send_newsletter_messages: false,
  }
}

function createWebhookFormFromItem(item) {
  return {
    id: item?.id || '',
    name: item?.name || '',
    url: item?.url || '',
    url_absolute: item?.url_absolute || '',
    header: item?.header || '',
    token: item?.token || '',
    timeout_ms: String(
      Number.isFinite(Number(item?.timeout_ms)) ? Number(item.timeout_ms) : 5000
    ),
    enabled: Boolean(item?.enabled),
    append_session: Boolean(item?.append_session),
    send_incoming_messages: Boolean(item?.send_incoming_messages),
    send_outgoing_messages: Boolean(item?.send_outgoing_messages),
    send_update_messages: Boolean(item?.send_update_messages),
    send_group_messages: Boolean(item?.send_group_messages),
    send_newsletter_messages: Boolean(item?.send_newsletter_messages),
  }
}

function createDefaultTemplateForm() {
  return {
    id: '',
    name: '',
    language: 'pt_BR',
    category: '',
    header_type: '',
    header_text: '',
    header_link: '',
    header_id: '',
    header_filename: '',
    body: '',
    footer: '',
  }
}

function createTemplateFormFromItem(item) {
  const header = item?.header || {}
  return {
    id: item?.id || '',
    name: item?.name || '',
    language: item?.language || 'pt_BR',
    category: item?.category || '',
    header_type: header?.type || '',
    header_text: header?.text || '',
    header_link: header?.link || '',
    header_id: header?.id || '',
    header_filename: header?.filename || '',
    body: item?.body || '',
    footer: item?.footer || '',
  }
}

function templateFormToPayload(form) {
  return {
    id: form.id.trim(),
    name: form.name.trim(),
    language: form.language.trim(),
    category: form.category.trim(),
    header: {
      type: form.header_type.trim(),
      text: form.header_text.trim(),
      link: form.header_link.trim(),
      id: form.header_id.trim(),
      filename: form.header_filename.trim(),
    },
    body: form.body.trim(),
    footer: form.footer.trim(),
  }
}

function sessionStatus(session) {
  const raw = String(session?.status || '').trim().toLowerCase()
  if (raw === 'connected') {
    return 'online'
  }
  if (raw) {
    return raw
  }
  return session?.connected ? 'online' : 'offline'
}

function statusTone(status) {
  if (status === 'online') {
    return 'online'
  }
  if (status === 'connecting') {
    return 'connecting'
  }
  return 'offline'
}

function sortSessions(left, right) {
  const leftStatus = sessionStatus(left)
  const rightStatus = sessionStatus(right)
  const leftScore =
    (left?.assigned ? 0 : 4) +
    (leftStatus === 'online' ? 0 : leftStatus === 'connecting' ? 1 : 2)
  const rightScore =
    (right?.assigned ? 0 : 4) +
    (rightStatus === 'online' ? 0 : rightStatus === 'connecting' ? 1 : 2)

  if (leftScore !== rightScore) {
    return leftScore - rightScore
  }
  return String(left?.session_id || '').localeCompare(String(right?.session_id || ''), 'pt-BR')
}

function formatTime(date) {
  if (!date) {
    return 'Nunca'
  }
  return new Intl.DateTimeFormat('pt-BR', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(date)
}

function truncate(value, limit = 56) {
  if (!value) {
    return '—'
  }
  return value.length > limit ? `${value.slice(0, limit)}…` : value
}

function digitsOnlyValue(value) {
  return String(value || '').replace(/\D+/g, '')
}

function isEnterOrSpace(event) {
  return event.key === 'Enter' || event.key === ' '
}

function StatCard({ label, value, caption, tone = 'neutral' }) {
  return html`
    <article className=${`stat-card stat-card--${tone}`}>
      <p className="stat-card__label">${label}</p>
      <strong className="stat-card__value">${value}</strong>
      <span className="stat-card__caption">${caption}</span>
    </article>
  `
}

function InputField({
  label,
  value,
  onInput,
  placeholder = '',
  type = 'text',
  min,
  step,
}) {
  return html`
    <label className="field">
      <span>${label}</span>
      <input
        type=${type}
        value=${value}
        placeholder=${placeholder}
        min=${min}
        step=${step}
        onInput=${(event) => onInput(event.target.value)}
      />
    </label>
  `
}

function TextAreaField({ label, value, onInput, placeholder = '', rows = 4 }) {
  return html`
    <label className="field">
      <span>${label}</span>
      <textarea
        rows=${rows}
        value=${value}
        placeholder=${placeholder}
        onInput=${(event) => onInput(event.target.value)}
      ></textarea>
    </label>
  `
}

function CheckField({ label, description = '', checked, onChange }) {
  return html`
    <label className="check-field">
      <input
        type="checkbox"
        checked=${checked}
        onChange=${(event) => onChange(event.target.checked)}
      />
      <span>
        <strong>${label}</strong>
        ${description ? html`<small>${description}</small>` : null}
      </span>
    </label>
  `
}

function FilterChip({ active, label, count, onClick }) {
  return html`
    <button
      type="button"
      className=${`filter-chip ${active ? 'filter-chip--active' : ''}`}
      onClick=${onClick}
    >
      <span>${label}</span>
      <strong>${count}</strong>
    </button>
  `
}

function MiniButton({ label, tone = 'ghost', onClick, disabled = false }) {
  return html`
    <button
      type="button"
      className=${`mini-button mini-button--${tone}`}
      disabled=${disabled}
      onClick=${(event) => {
        event.stopPropagation()
        onClick()
      }}
    >
      ${label}
    </button>
  `
}

function SessionCard({
  session,
  selected,
  busyMap,
  onSelect,
  onSessionAction,
  onOpenTab,
  onShowQR,
}) {
  const status = sessionStatus(session)
  const tone = statusTone(status)
  const assignLabel = session.assigned ? 'Desatribuir' : 'Atribuir'
  const assignAction = session.assigned ? 'unassign' : 'assign'

  return html`
    <article
      className=${`session-card ${selected ? 'session-card--selected' : ''}`}
      role="button"
      tabIndex="0"
      onClick=${() => onSelect(session.session_id)}
      onKeyDown=${(event) => {
        if (isEnterOrSpace(event)) {
          event.preventDefault()
          onSelect(session.session_id)
        }
      }}
    >
      <div className="session-card__top">
        <div>
          <p className="session-card__eyebrow">Instância</p>
          <h3 className="session-card__title">${session.session_id}</h3>
        </div>
        <span className=${`status-pill status-pill--${tone}`}>${status}</span>
      </div>

      <div className="tag-row">
        <span className="tag">${session.connected ? 'Conectada' : 'Desconectada'}</span>
        <span className="tag">${session.assigned ? 'Atribuída' : 'Sem atribuição'}</span>
        <span className="tag">${session.has_qr ? 'QR disponível' : 'Sem QR'}</span>
      </div>

      <div className="session-card__meta">
        <span>Login: ${session.logged_in ? 'ativo' : 'pendente'}</span>
        <span>Device: ${session.device_jid ? truncate(session.device_jid, 30) : '—'}</span>
        <span>WA Web: ${session.whatsapp_web_version || '—'}</span>
      </div>

      <div className="session-card__actions">
        <${MiniButton}
          label="QR"
          onClick=${() => onShowQR(session.session_id)}
          disabled=${!session.has_qr}
        />
        <${MiniButton}
          label="Conectar"
          onClick=${() => onSessionAction('connect', session.session_id)}
          disabled=${Boolean(busyMap[`connect:${session.session_id}`])}
        />
        <${MiniButton}
          label="Recarregar"
          onClick=${() => onSessionAction('reload', session.session_id)}
          disabled=${Boolean(busyMap[`reload:${session.session_id}`])}
        />
        <${MiniButton}
          label="Desconectar"
          tone="danger"
          onClick=${() => onSessionAction('disconnect', session.session_id)}
          disabled=${Boolean(busyMap[`disconnect:${session.session_id}`])}
        />
        <${MiniButton}
          label="Excluir"
          tone="danger"
          onClick=${() => onSessionAction('delete', session.session_id)}
          disabled=${Boolean(busyMap[`delete:${session.session_id}`])}
        />
        <${MiniButton}
          label=${assignLabel}
          onClick=${() => onSessionAction(assignAction, session.session_id)}
          disabled=${Boolean(busyMap[`${assignAction}:${session.session_id}`])}
        />
        <${MiniButton}
          label="Config"
          onClick=${() => onOpenTab(session.session_id, 'config')}
        />
        <${MiniButton}
          label="Hooks"
          onClick=${() => onOpenTab(session.session_id, 'webhooks')}
        />
        <${MiniButton}
          label="Templates"
          onClick=${() => onOpenTab(session.session_id, 'templates')}
        />
      </div>
    </article>
  `
}

function DetailItem({ label, value, code = false }) {
  return html`
    <div className="detail-item">
      <span>${label}</span>
      ${code ? html`<code>${value}</code>` : html`<strong>${value}</strong>`}
    </div>
  `
}

function QRModal({ sessionId, onClose }) {
  const [tick, setTick] = useState(Date.now())

  useEffect(() => {
    if (!sessionId) {
      return undefined
    }

    setTick(Date.now())
    const timer = window.setInterval(() => setTick(Date.now()), QR_REFRESH_MS)
    const handleKeydown = (event) => {
      if (event.key === 'Escape') {
        onClose()
      }
    }
    window.addEventListener('keydown', handleKeydown)

    return () => {
      window.clearInterval(timer)
      window.removeEventListener('keydown', handleKeydown)
    }
  }, [sessionId])

  if (!sessionId) {
    return null
  }

  return html`
    <div
      className="modal-backdrop"
      onClick=${(event) => {
        if (event.target === event.currentTarget) {
          onClose()
        }
      }}
    >
      <div className="modal-card">
        <div className="modal-card__header">
          <div>
            <p className="eyebrow">Pareamento</p>
            <h3>QR da sessão ${sessionId}</h3>
          </div>
          <button type="button" className="ghost-button" onClick=${onClose}>Fechar</button>
        </div>
        <img
          className="modal-card__image"
          src=${`/sessions/${enc(sessionId)}/qr?ts=${tick}`}
          alt=${`QR da sessão ${sessionId}`}
        />
        <p className="hint-text">
          O QR atualiza automaticamente a cada poucos segundos enquanto a sessão estiver em pareamento.
        </p>
      </div>
    </div>
  `
}

function WebhookCard({ item, onEdit, onDelete }) {
  return html`
    <article className="webhook-card">
      <div className="webhook-card__top">
        <div>
          <strong>${item.id}</strong>
          <p>${item.name || 'Webhook sem nome'}</p>
        </div>
        <div className="inline-actions">
          <button type="button" className="text-button" onClick=${onEdit}>Editar</button>
          <button type="button" className="text-button text-button--danger" onClick=${onDelete}>
            Excluir
          </button>
        </div>
      </div>

      <code className="webhook-card__url">${item.url_absolute || item.url || '—'}</code>

      <div className="tag-row">
        <span className="tag">${item.enabled ? 'Ativo' : 'Desligado'}</span>
        <span className="tag">${item.append_session ? 'Anexa sessão' : 'URL fixa'}</span>
        <span className="tag">${item.send_incoming_messages ? 'Incoming' : 'Sem incoming'}</span>
        <span className="tag">${item.send_outgoing_messages ? 'Outgoing' : 'Sem outgoing'}</span>
        <span className="tag">${item.send_update_messages ? 'Update' : 'Sem update'}</span>
      </div>
    </article>
  `
}

function TemplateCard({ item, onEdit, onDelete }) {
  const header = item?.header || {}
  const headerType = header.type || (header.text ? 'text' : 'sem cabeçalho')
  return html`
    <article className="webhook-card">
      <div className="webhook-card__top">
        <div>
          <strong>${item.name}</strong>
          <p>${item.language || 'idioma padrão'}${item.category ? ` · ${item.category}` : ''}</p>
        </div>
        <div className="inline-actions">
          <button type="button" className="text-button" onClick=${onEdit}>Editar</button>
          <button type="button" className="text-button text-button--danger" onClick=${onDelete}>
            Excluir
          </button>
        </div>
      </div>

      <code className="webhook-card__url">${item.body || '—'}</code>

      <div className="tag-row">
        <span className="tag">Header: ${headerType}</span>
        <span className="tag">${item.footer ? 'Com footer' : 'Sem footer'}</span>
        <span className="tag">${item.id}</span>
      </div>
    </article>
  `
}

function App() {
  const [theme, setTheme] = useState(getInitialTheme)
  const [worker, setWorker] = useState(null)
  const [sessions, setSessions] = useState([])
  const [activePage, setActivePage] = useState('home')
  const [selectedSessionId, setSelectedSessionId] = useState('')
  const [selectedTab, setSelectedTab] = useState('quick')
  const [overviewLoading, setOverviewLoading] = useState(true)
  const [autoRefresh, setAutoRefresh] = useState(true)
  const [lastUpdated, setLastUpdated] = useState(null)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')
  const [toasts, setToasts] = useState([])
  const [busyMap, setBusyMap] = useState({})
  const [qrSessionId, setQrSessionId] = useState('')

  const [createSessionForm, setCreateSessionForm] = useState({
    id: '',
    assign: true,
    connect: true,
    history_sync_enabled: false,
    history_sync_days: '7',
    webhook: createDefaultWebhookForm(),
  })
  const [messageForm, setMessageForm] = useState({
    session: '',
    to: '',
    body: '',
  })
  const [configState, setConfigState] = useState({
    loading: false,
    saving: false,
    snapshot: null,
    form: createDefaultConfigForm(),
  })
  const [pairingState, setPairingState] = useState({
    phone: '',
    code: '',
    loading: false,
  })
  const [webhookState, setWebhookState] = useState({
    loading: false,
    saving: false,
    items: [],
    form: createDefaultWebhookForm(),
  })
  const [templateState, setTemplateState] = useState({
    loading: false,
    saving: false,
    items: [],
    form: createDefaultTemplateForm(),
  })

  const refreshOverviewRef = useRef(null)

  function pushToast(message, tone = 'info') {
    const id = makeToastId()
    setToasts((current) => [...current, { id, message, tone }])
    window.setTimeout(() => {
      setToasts((current) => current.filter((item) => item.id !== id))
    }, 3400)
  }

  function setBusy(key, value) {
    setBusyMap((current) => {
      const next = { ...current }
      if (value) {
        next[key] = true
      } else {
        delete next[key]
      }
      return next
    })
  }

  async function refreshOverview({ silent = false } = {}) {
    if (!silent) {
      setOverviewLoading(true)
    }
    try {
      const [workerResponse, sessionsResponse] = await Promise.all([
        api('/api/v1/worker'),
        api('/api/v1/sessions'),
      ])
      setWorker(workerResponse)
      setSessions(Array.isArray(sessionsResponse?.sessions) ? sessionsResponse.sessions : [])
      setLastUpdated(new Date())
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      if (!silent) {
        setOverviewLoading(false)
      }
    }
  }

  refreshOverviewRef.current = refreshOverview

  async function loadConfigForSession(sessionId) {
    const sid = String(sessionId || '').trim()
    if (!sid) {
      return
    }

    setConfigState((current) => ({ ...current, loading: true }))
    try {
      const snapshot = await api(`/api/v1/sessions/${enc(sid)}/config`)
      setConfigState((current) => ({
        ...current,
        loading: false,
        snapshot,
        form: createConfigFormFromSnapshot(snapshot),
      }))
    } catch (error) {
      setConfigState((current) => ({ ...current, loading: false }))
      pushToast(error.message, 'error')
    }
  }

  async function loadWebhooksForSession(sessionId, { resetForm = false } = {}) {
    const sid = String(sessionId || '').trim()
    if (!sid) {
      return
    }

    setWebhookState((current) => ({ ...current, loading: true }))
    try {
      const payload = await api(`/api/v1/sessions/${enc(sid)}/webhooks`)
      setWebhookState((current) => ({
        ...current,
        loading: false,
        items: Array.isArray(payload?.webhooks) ? payload.webhooks : [],
        form: resetForm ? createDefaultWebhookForm() : current.form,
      }))
    } catch (error) {
      setWebhookState((current) => ({ ...current, loading: false }))
      pushToast(error.message, 'error')
    }
  }

  async function loadTemplatesForSession(sessionId, { resetForm = false } = {}) {
    const sid = String(sessionId || '').trim()
    if (!sid) {
      return
    }

    setTemplateState((current) => ({ ...current, loading: true }))
    try {
      const payload = await api(`/api/v1/sessions/${enc(sid)}/templates`)
      setTemplateState((current) => ({
        ...current,
        loading: false,
        items: Array.isArray(payload?.templates) ? payload.templates : [],
        form: resetForm ? createDefaultTemplateForm() : current.form,
      }))
    } catch (error) {
      setTemplateState((current) => ({ ...current, loading: false }))
      pushToast(error.message, 'error')
    }
  }

  useEffect(() => {
    document.documentElement.dataset.theme = theme
    localStorage.setItem(THEME_STORAGE_KEY, theme)
  }, [theme])

  useEffect(() => {
    document.title = worker?.ui_title || 'WhatsMeow Control Center'
  }, [worker])

  useEffect(() => {
    refreshOverview()
  }, [])

  useEffect(() => {
    if (!autoRefresh) {
      return undefined
    }
    const timer = window.setInterval(() => {
      if (!document.hidden && refreshOverviewRef.current) {
        refreshOverviewRef.current({ silent: true })
      }
    }, REFRESH_INTERVAL_MS)
    return () => window.clearInterval(timer)
  }, [autoRefresh])

  useEffect(() => {
    if (sessions.length === 0) {
      setSelectedSessionId('')
      return
    }
    if (selectedSessionId && sessions.some((session) => session.session_id === selectedSessionId)) {
      return
    }
    const next = [...sessions].sort(sortSessions)[0]
    if (next) {
      setSelectedSessionId(next.session_id)
    }
  }, [sessions, selectedSessionId])

  useEffect(() => {
    if (!selectedSessionId) {
      setConfigState({
        loading: false,
        saving: false,
        snapshot: null,
        form: createDefaultConfigForm(),
      })
      setPairingState({
        phone: '',
        code: '',
        loading: false,
      })
      setWebhookState({
        loading: false,
        saving: false,
        items: [],
        form: createDefaultWebhookForm(),
      })
      setTemplateState({
        loading: false,
        saving: false,
        items: [],
        form: createDefaultTemplateForm(),
      })
      return
    }

    setMessageForm((current) => ({
      ...current,
      session: selectedSessionId,
    }))
    setPairingState({
      phone: digitsOnlyValue(selectedSessionId),
      code: '',
      loading: false,
    })

    let active = true
    setConfigState({
      loading: true,
      saving: false,
      snapshot: null,
      form: createDefaultConfigForm(),
    })
    setWebhookState({
      loading: true,
      saving: false,
      items: [],
      form: createDefaultWebhookForm(),
    })
    setTemplateState({
      loading: true,
      saving: false,
      items: [],
      form: createDefaultTemplateForm(),
    })

    Promise.all([
      api(`/api/v1/sessions/${enc(selectedSessionId)}/config`),
      api(`/api/v1/sessions/${enc(selectedSessionId)}/webhooks`),
      api(`/api/v1/sessions/${enc(selectedSessionId)}/templates`),
    ])
      .then(([snapshot, payload, templatesPayload]) => {
        if (!active) {
          return
        }
        setConfigState({
          loading: false,
          saving: false,
          snapshot,
          form: createConfigFormFromSnapshot(snapshot),
        })
        setWebhookState({
          loading: false,
          saving: false,
          items: Array.isArray(payload?.webhooks) ? payload.webhooks : [],
          form: createDefaultWebhookForm(),
        })
        setTemplateState({
          loading: false,
          saving: false,
          items: Array.isArray(templatesPayload?.templates) ? templatesPayload.templates : [],
          form: createDefaultTemplateForm(),
        })
      })
      .catch((error) => {
        if (!active) {
          return
        }
        setConfigState((current) => ({ ...current, loading: false }))
        setWebhookState((current) => ({ ...current, loading: false }))
        setTemplateState((current) => ({ ...current, loading: false }))
        pushToast(error.message, 'error')
      })

    return () => {
      active = false
    }
  }, [selectedSessionId])

  const title = worker?.ui_title || 'WhatsMeow Control Center'
  const selectedSession =
    sessions.find((session) => session.session_id === selectedSessionId) || null
  const onlineCount = sessions.filter(
    (session) => session.connected || sessionStatus(session) === 'online'
  ).length
  const withQRCount = sessions.filter((session) => session.has_qr).length
  const assignedCount = sessions.filter((session) => session.assigned).length
  const sortedSessions = [...sessions].sort(sortSessions)
  const normalizedSearch = search.trim().toLowerCase()
  const filteredSessions = sortedSessions.filter((session) => {
    if (statusFilter === 'assigned' && !session.assigned) {
      return false
    }
    if (
      statusFilter === 'online' &&
      !(session.connected || sessionStatus(session) === 'online')
    ) {
      return false
    }
    if (statusFilter === 'qr' && !session.has_qr) {
      return false
    }

    if (!normalizedSearch) {
      return true
    }

    const haystack = [
      session.session_id,
      session.device_jid,
      sessionStatus(session),
      session.assigned ? 'atribuida' : 'livre',
    ]
      .filter(Boolean)
      .join(' ')
      .toLowerCase()

    return haystack.includes(normalizedSearch)
  })
  const selectedSessionStatus = selectedSession ? sessionStatus(selectedSession) : ''
  const showInlineQR = Boolean(
    selectedSession && selectedSession.has_qr && selectedSessionStatus === 'connecting'
  )

  async function handleCreateSession(event) {
    event.preventDefault()
    const sid = createSessionForm.id.trim()
    if (!sid) {
      pushToast('Informe o ID da sessão.', 'error')
      return
    }
    const historyDays = parseInt(createSessionForm.history_sync_days, 10)
    if (createSessionForm.history_sync_enabled) {
      if (!Number.isFinite(historyDays) || historyDays <= 0) {
        pushToast('Informe a quantidade de dias do histórico.', 'error')
        return
      }
      if (
        !createSessionForm.webhook.url.trim() &&
        !createSessionForm.webhook.url_absolute.trim()
      ) {
        pushToast('Informe o webhook antes de ativar a sincronização de histórico.', 'error')
        return
      }
    }

    const busyKey = `create:${sid}`
    setBusy(busyKey, true)
    try {
      const config = historyConfigPayload(
        createSessionForm.history_sync_enabled,
        createSessionForm.history_sync_days
      )
      const body = {
        id: sid,
        assign: createSessionForm.assign,
        connect: createSessionForm.connect,
        config,
      }
      if (createSessionForm.history_sync_enabled) {
        body.webhook = webhookFormToPayload(createSessionForm.webhook)
      }
      await api('/api/v1/sessions', {
        method: 'POST',
        body: JSON.stringify(body),
      })
      pushToast(`Sessão ${sid} solicitada com sucesso.`, 'success')
      setCreateSessionForm({
        id: '',
        assign: true,
        connect: true,
        history_sync_enabled: false,
        history_sync_days: '7',
        webhook: createDefaultWebhookForm(),
      })
      setSelectedSessionId(sid)
      setSelectedTab('overview')
      await refreshOverview({ silent: true })
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      setBusy(busyKey, false)
    }
  }

  async function handleQuickSend(event) {
    event.preventDefault()
    const session = messageForm.session.trim()
    const to = messageForm.to.trim()
    const body = messageForm.body.trim()

    if (!session || !to || !body) {
      pushToast('Preencha sessão, destino e mensagem.', 'error')
      return
    }

    const busyKey = `send:${session}`
    setBusy(busyKey, true)
    try {
      await api(`/api/v1/sessions/${enc(session)}/messages`, {
        method: 'POST',
        body: JSON.stringify({
          type: 'text',
          to,
          text: { body },
        }),
      })
      pushToast(`Mensagem enviada pela sessão ${session}.`, 'success')
      setMessageForm((current) => ({
        ...current,
        body: '',
      }))
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      setBusy(busyKey, false)
    }
  }

  async function handleSessionAction(action, sessionId) {
    const sid = String(sessionId || '').trim()
    if (!sid) {
      return
    }

    if (action === 'disconnect' || action === 'delete') {
      const prompt =
        action === 'delete'
          ? `Excluir permanentemente a sessão "${sid}"? Isso remove o diretório da sessão, a atribuição, overrides e webhooks.`
          : `Desconectar a sessão "${sid}"?`
      const confirmed = window.confirm(prompt)
      if (!confirmed) {
        return
      }
    }

    const busyKey = `${action}:${sid}`
    setBusy(busyKey, true)
    try {
      if (action === 'assign') {
        await api('/api/v1/worker/assignments', {
          method: 'POST',
          body: JSON.stringify({ session_id: sid }),
        })
      } else if (action === 'unassign') {
        await api(`/api/v1/worker/assignments/${enc(sid)}`, {
          method: 'DELETE',
        })
      } else if (action === 'delete') {
        await api(`/api/v1/sessions/${enc(sid)}`, {
          method: 'DELETE',
        })
      } else if (action === 'qr') {
        setQrSessionId(sid)
      } else {
        await api(`/sessions/${enc(sid)}/${action}`, {
          method: 'POST',
        })
      }

      if (action !== 'qr') {
        pushToast(
          action === 'delete'
            ? `Sessão ${sid} excluída.`
            : `Ação "${action}" aplicada em ${sid}.`,
          'success'
        )
        await refreshOverview({ silent: true })
      }
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      setBusy(busyKey, false)
    }
  }

  function handleOpenTab(sessionId, tab) {
    setSelectedSessionId(sessionId)
    setActivePage('instances')
    setSelectedTab(tab)
  }

  async function handleSaveConfig(event) {
    event.preventDefault()
    if (!selectedSessionId) {
      return
    }

    const interval = parseInt(configState.form.always_online_interval_seconds, 10)
    if (!Number.isFinite(interval) || interval <= 0) {
      pushToast('O intervalo de always-online deve ser maior que zero.', 'error')
      return
    }
    const historyDays = parseInt(configState.form.history_sync_days, 10)
    if (configState.form.history_sync_enabled && (!Number.isFinite(historyDays) || historyDays <= 0)) {
      pushToast('A quantidade de dias do histórico deve ser maior que zero.', 'error')
      return
    }
    if (configState.form.proxy_enabled) {
      if (!configState.form.proxy_host.trim()) {
        pushToast('Informe o host do proxy.', 'error')
        return
      }
      if (!configState.form.proxy_port.trim()) {
        pushToast('Informe a porta do proxy.', 'error')
        return
      }
    }

    setConfigState((current) => ({ ...current, saving: true }))
    try {
      await api(`/api/v1/sessions/${enc(selectedSessionId)}/config`, {
        method: 'PUT',
        body: JSON.stringify({
          reject_calls: configState.form.reject_calls,
          reject_calls_message: configState.form.reject_calls_message.trim(),
          mark_read_on_message: configState.form.mark_read_on_message,
          ignore_status_broadcast: configState.form.ignore_status_broadcast,
          ignore_newsletters: configState.form.ignore_newsletters,
          outgoing_messages_coex: configState.form.outgoing_messages_coex,
          history_sync_enabled: configState.form.history_sync_enabled,
          history_sync_days: configState.form.history_sync_enabled ? historyDays : 0,
          always_online: configState.form.always_online,
          always_online_interval_seconds: interval,
          proxy_enabled: configState.form.proxy_enabled,
          proxy_protocol: configState.form.proxy_protocol.trim(),
          proxy_host: configState.form.proxy_host.trim(),
          proxy_port: configState.form.proxy_port.trim(),
          proxy_username: configState.form.proxy_username.trim(),
          proxy_password: configState.form.proxy_password,
        }),
      })
      pushToast(`Configuração da sessão ${selectedSessionId} salva.`, 'success')
      await loadConfigForSession(selectedSessionId)
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      setConfigState((current) => ({ ...current, saving: false }))
    }
  }

  async function handleRequestPairingCode(event) {
    event.preventDefault()
    if (!selectedSessionId) {
      return
    }

    const phone = digitsOnlyValue(pairingState.phone)
    if (!phone) {
      pushToast('Informe o telefone para gerar o código de pareamento.', 'error')
      return
    }

    setPairingState((current) => ({
      ...current,
      phone,
      code: '',
      loading: true,
    }))
    try {
      const payload = await api(`/api/v1/sessions/${enc(selectedSessionId)}/pairing-code`, {
        method: 'POST',
        body: JSON.stringify({ phone }),
      })
      const code = String(payload?.pairing_code || '').trim()
      if (!code) {
        throw new Error('pairing code not returned by server')
      }
      setPairingState((current) => ({
        ...current,
        phone,
        code,
        loading: false,
      }))
      pushToast(`Código de pareamento gerado para ${selectedSessionId}.`, 'success')
    } catch (error) {
      setPairingState((current) => ({
        ...current,
        loading: false,
      }))
      pushToast(error.message, 'error')
    }
  }

  async function handleClearConfig() {
    if (!selectedSessionId) {
      return
    }
    if (!window.confirm(`Remover o override da sessão "${selectedSessionId}"?`)) {
      return
    }

    setConfigState((current) => ({ ...current, saving: true }))
    try {
      await api(`/api/v1/sessions/${enc(selectedSessionId)}/config`, {
        method: 'DELETE',
      })
      pushToast(`Override da sessão ${selectedSessionId} removido.`, 'success')
      await loadConfigForSession(selectedSessionId)
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      setConfigState((current) => ({ ...current, saving: false }))
    }
  }

  async function handleSaveWebhook(event) {
    event.preventDefault()
    if (!selectedSessionId) {
      return
    }

    const timeout = parseInt(webhookState.form.timeout_ms, 10)
    if (!webhookState.form.url.trim() && !webhookState.form.url_absolute.trim()) {
      pushToast('Informe a URL base ou a URL absoluta do webhook.', 'error')
      return
    }
    if (!Number.isFinite(timeout) || timeout <= 0) {
      pushToast('O timeout do webhook deve ser maior que zero.', 'error')
      return
    }

    setWebhookState((current) => ({ ...current, saving: true }))
    try {
      await api(`/api/v1/sessions/${enc(selectedSessionId)}/webhooks`, {
        method: 'POST',
        body: JSON.stringify(webhookFormToPayload(webhookState.form)),
      })
      pushToast(`Webhook salvo para a sessão ${selectedSessionId}.`, 'success')
      await loadWebhooksForSession(selectedSessionId, { resetForm: true })
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      setWebhookState((current) => ({ ...current, saving: false }))
    }
  }

  async function handleDeleteWebhook(webhookId) {
    if (!selectedSessionId || !webhookId) {
      return
    }
    if (!window.confirm(`Excluir o webhook "${webhookId}" da sessão "${selectedSessionId}"?`)) {
      return
    }

    setWebhookState((current) => ({ ...current, saving: true }))
    try {
      await api(`/api/v1/sessions/${enc(selectedSessionId)}/webhooks/${enc(webhookId)}`, {
        method: 'DELETE',
      })
      pushToast(`Webhook ${webhookId} excluído.`, 'success')
      await loadWebhooksForSession(selectedSessionId, { resetForm: true })
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      setWebhookState((current) => ({ ...current, saving: false }))
    }
  }

  async function handleSaveTemplate(event) {
    event.preventDefault()
    if (!selectedSessionId) {
      return
    }
    if (!templateState.form.name.trim()) {
      pushToast('Informe o nome do template.', 'error')
      return
    }
    if (!templateState.form.body.trim()) {
      pushToast('Informe o corpo do template.', 'error')
      return
    }

    setTemplateState((current) => ({ ...current, saving: true }))
    try {
      await api(`/api/v1/sessions/${enc(selectedSessionId)}/templates`, {
        method: 'POST',
        body: JSON.stringify(templateFormToPayload(templateState.form)),
      })
      pushToast(`Template salvo para a sessão ${selectedSessionId}.`, 'success')
      await loadTemplatesForSession(selectedSessionId, { resetForm: true })
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      setTemplateState((current) => ({ ...current, saving: false }))
    }
  }

  async function handleDeleteTemplate(templateId) {
    if (!selectedSessionId || !templateId) {
      return
    }
    if (!window.confirm(`Excluir o template "${templateId}" da sessão "${selectedSessionId}"?`)) {
      return
    }

    setTemplateState((current) => ({ ...current, saving: true }))
    try {
      await api(`/api/v1/sessions/${enc(selectedSessionId)}/templates/${enc(templateId)}`, {
        method: 'DELETE',
      })
      pushToast(`Template ${templateId} excluído.`, 'success')
      await loadTemplatesForSession(selectedSessionId, { resetForm: true })
    } catch (error) {
      pushToast(error.message, 'error')
    } finally {
      setTemplateState((current) => ({ ...current, saving: false }))
    }
  }

  const filterOptions = [
    { id: 'all', label: 'Todas', count: sessions.length },
    { id: 'assigned', label: 'Atribuídas', count: assignedCount },
    { id: 'online', label: 'Online', count: onlineCount },
    { id: 'qr', label: 'Com QR', count: withQRCount },
  ]

  return html`
    <div className="page">
      <main className="dashboard-shell">
        <aside className="system-sidebar">
          <div className="system-sidebar__brand">
            <img
              className="system-sidebar__logo"
              src=${theme === 'dark' ? LOGO_DARK : LOGO_LIGHT}
              alt="Envolve"
            />
            <div>
              <strong>WhatsMeow</strong>
              <span>${worker?.worker_id || 'worker'}</span>
            </div>
          </div>

          <nav className="system-nav" aria-label="Navegação principal">
            <button
              type="button"
              className=${`system-nav__item ${activePage === 'home' ? 'system-nav__item--active' : ''}`}
              onClick=${() => setActivePage('home')}
            >
              Home
            </button>
            <button
              type="button"
              className=${`system-nav__item ${activePage === 'instances' ? 'system-nav__item--active' : ''}`}
              onClick=${() => setActivePage('instances')}
            >
              Instâncias
            </button>
          </nav>

          <div className="system-sidebar__meta">
            <span>${onlineCount} online</span>
            <span>${assignedCount} atribuída(s)</span>
          </div>
        </aside>

        <section className="provider-topbar panel">
          <div className="provider-topbar__brand">
            <div>
              <p className="eyebrow">Provider WhatsMeow</p>
              <h1>${activePage === 'home' ? 'Home' : 'Instâncias'}</h1>
            </div>
          </div>

          <div className="provider-topbar__actions">
            <span className="hero-panel__timestamp">
              ${lastUpdated ? `Atualizado às ${formatTime(lastUpdated)}` : 'Carregando dados'}
            </span>
            <button
              type="button"
              className=${`ghost-button ${autoRefresh ? 'ghost-button--active' : ''}`}
              onClick=${() => setAutoRefresh((current) => !current)}
            >
              ${autoRefresh ? 'Auto ligado' : 'Auto desligado'}
            </button>
            <button
              type="button"
              className="solid-button"
              disabled=${overviewLoading}
              onClick=${() => refreshOverview()}
            >
              ${overviewLoading ? 'Atualizando...' : 'Atualizar'}
            </button>
            <button
              type="button"
              className="ghost-button"
              onClick=${() => setTheme((current) => (current === 'dark' ? 'light' : 'dark'))}
            >
              ${theme === 'dark' ? 'Tema claro' : 'Tema escuro'}
            </button>
          </div>
        </section>

        ${activePage === 'home'
          ? html`
        <section className="summary-grid">
          <${StatCard}
            label="Worker ativo"
            value=${worker?.worker_id || '—'}
            caption=${worker?.worker_strict ? 'Modo strict habilitado' : 'Aceita sessões fora da lista'}
            tone="brand"
          />
          <${StatCard}
            label="Sessões atribuídas"
            value=${String(assignedCount)}
            caption="visíveis para este worker"
            tone="accent"
          />
          <${StatCard}
            label="Sessões online"
            value=${String(onlineCount)}
            caption="conexões ativas agora"
            tone="online"
          />
          <${StatCard}
            label="Fila AMQP"
            value=${worker?.amqp_queue || '—'}
            caption=${worker?.amqp_enabled ? 'consumidor habilitado' : 'AMQP desabilitado'}
            tone="secondary"
          />
        </section>

        <section className="home-grid">
          <section className="panel">
            <div className="panel__header">
              <div>
                <p className="eyebrow">Instâncias</p>
                <h2>Resumo operacional</h2>
              </div>
            </div>
            <div className="detail-grid">
              <${DetailItem} label="Total" value=${String(sessions.length)} />
              <${DetailItem} label="Online" value=${String(onlineCount)} />
              <${DetailItem} label="Com QR" value=${String(withQRCount)} />
              <${DetailItem} label="Atribuídas" value=${String(assignedCount)} />
            </div>
            <div className="inline-actions inline-actions--wrap">
              <button type="button" className="solid-button" onClick=${() => setActivePage('instances')}>
                Abrir instâncias
              </button>
              <button type="button" className="ghost-button" onClick=${() => refreshOverview()}>
                Atualizar dados
              </button>
            </div>
          </section>

          <section className="panel">
            <div className="panel__header">
              <div>
                <p className="eyebrow">Worker</p>
                <h2>Filas e bindings</h2>
              </div>
            </div>
            <div className="meta-stack">
              <div className="detail-item">
                <span>Exchange</span>
                <strong>${worker?.amqp_exchange || '—'}</strong>
              </div>
              <div className="detail-item">
                <span>Queue</span>
                <strong>${worker?.amqp_queue || '—'}</strong>
              </div>
              <div className="detail-item">
                <span>Auto-connect</span>
                <strong>${worker?.worker_auto_connect_assigned ? 'sim' : 'não'}</strong>
              </div>
            </div>
            <div className="tag-row tag-row--spaced">
              ${(worker?.amqp_bindings || []).length
                ? worker.amqp_bindings.map(
                    (binding) => html`<span key=${binding} className="tag">${binding}</span>`
                  )
                : html`<span className="empty-inline">Nenhum binding configurado.</span>`}
            </div>
          </section>
        </section>
          `
          : null}

        ${activePage === 'instances'
          ? html`
        <section className="workspace-grid">
          <div className="workspace-column workspace-column--sessions">
            <section className="panel">
              <div className="panel__header panel__header--space-between">
                <div>
                  <p className="eyebrow">Instâncias</p>
                  <h2>Instâncias</h2>
                </div>
                <span className="panel__hint">
                  ${filteredSessions.length} de ${sessions.length}
                </span>
              </div>

              <details className="create-instance">
                <summary>Nova instância</summary>
                <form className="stack-form" onSubmit=${handleCreateSession}>
                  <${InputField}
                    label="ID da sessão"
                    value=${createSessionForm.id}
                    placeholder="Ex: 5562999999999"
                    onInput=${(value) =>
                      setCreateSessionForm((current) => ({ ...current, id: value }))
                    }
                  />
                  <div className="check-grid check-grid--single">
                    <${CheckField}
                      label="Atribuir ao worker atual"
                      checked=${createSessionForm.assign}
                      onChange=${(checked) =>
                        setCreateSessionForm((current) => ({ ...current, assign: checked }))
                      }
                    />
                    <${CheckField}
                      label="Conectar imediatamente"
                      checked=${createSessionForm.connect}
                      onChange=${(checked) =>
                        setCreateSessionForm((current) => ({ ...current, connect: checked }))
                      }
                    />
                    <${CheckField}
                      label="Sincronizar histórico"
                      description="Configure o webhook antes da conexão para receber o HistorySync."
                      checked=${createSessionForm.history_sync_enabled}
                      onChange=${(checked) =>
                        setCreateSessionForm((current) => ({
                          ...current,
                          history_sync_enabled: checked,
                        }))
                      }
                    />
                  </div>

                  ${createSessionForm.history_sync_enabled
                    ? html`
                        <div className="form-grid">
                          <${InputField}
                            label="Dias de histórico"
                            value=${createSessionForm.history_sync_days}
                            type="number"
                            min="1"
                            step="1"
                            onInput=${(value) =>
                              setCreateSessionForm((current) => ({
                                ...current,
                                history_sync_days: value,
                              }))
                            }
                          />
                          <${InputField}
                            label="Webhook ID"
                            value=${createSessionForm.webhook.id}
                            placeholder="vazio = auto"
                            onInput=${(value) =>
                              setCreateSessionForm((current) => ({
                                ...current,
                                webhook: { ...current.webhook, id: value },
                              }))
                            }
                          />
                        </div>

                        <${InputField}
                          label="URL base do webhook"
                          value=${createSessionForm.webhook.url}
                          placeholder="https://seu-app.com/webhooks/whatsapp"
                          onInput=${(value) =>
                            setCreateSessionForm((current) => ({
                              ...current,
                              webhook: { ...current.webhook, url: value },
                            }))
                          }
                        />
                        <${InputField}
                          label="URL absoluta do webhook"
                          value=${createSessionForm.webhook.url_absolute}
                          placeholder="https://seu-app.com/webhooks/whatsapp/5562..."
                          onInput=${(value) =>
                            setCreateSessionForm((current) => ({
                              ...current,
                              webhook: { ...current.webhook, url_absolute: value },
                            }))
                          }
                        />

                        <div className="form-grid">
                          <${InputField}
                            label="Header auth"
                            value=${createSessionForm.webhook.header}
                            placeholder="Authorization"
                            onInput=${(value) =>
                              setCreateSessionForm((current) => ({
                                ...current,
                                webhook: { ...current.webhook, header: value },
                              }))
                            }
                          />
                          <${InputField}
                            label="Token"
                            value=${createSessionForm.webhook.token}
                            placeholder="Bearer xxxxx"
                            onInput=${(value) =>
                              setCreateSessionForm((current) => ({
                                ...current,
                                webhook: { ...current.webhook, token: value },
                              }))
                            }
                          />
                        </div>
                      `
                    : null}
                  <button
                    type="submit"
                    className="solid-button solid-button--wide"
                    disabled=${Boolean(busyMap[`create:${createSessionForm.id.trim()}`])}
                  >
                    Criar instância
                  </button>
                </form>
              </details>

              <div className="toolbar">
                <label className="search-field">
                  <span>Buscar</span>
                  <input
                    type="text"
                    value=${search}
                    placeholder="Filtrar por sessão, status ou device"
                    onInput=${(event) => setSearch(event.target.value)}
                  />
                </label>
                <div className="filter-row">
                  ${filterOptions.map(
                    (option) => html`
                      <${FilterChip}
                        key=${option.id}
                        active=${statusFilter === option.id}
                        label=${option.label}
                        count=${option.count}
                        onClick=${() => setStatusFilter(option.id)}
                      />
                    `
                  )}
                </div>
              </div>

              ${overviewLoading
                ? html`
                    <div className="empty-state">
                      <strong>Carregando sessões…</strong>
                      <p>Consultando o worker e o catálogo de sessões disponíveis.</p>
                    </div>
                  `
                : filteredSessions.length === 0
                  ? html`
                      <div className="empty-state">
                        <strong>Nenhuma sessão encontrada</strong>
                        <p>Ajuste o filtro ou crie uma nova instância.</p>
                      </div>
                    `
                  : html`
                      <div className="session-grid">
                        ${filteredSessions.map(
                          (session) => html`
                            <${SessionCard}
                              key=${session.session_id}
                              session=${session}
                              selected=${selectedSessionId === session.session_id}
                              busyMap=${busyMap}
                              onSelect=${(sessionId) => {
                                setSelectedSessionId(sessionId)
                                setActivePage('instances')
                              }}
                              onSessionAction=${handleSessionAction}
                              onOpenTab=${handleOpenTab}
                              onShowQR=${setQrSessionId}
                            />
                          `
                        )}
                      </div>
                    `}
            </section>
          </div>

          <div className="workspace-column workspace-column--inspector">
            <section className="panel panel--sticky">
              ${selectedSession
                ? html`
                    <div className="panel__header">
                      <div>
                        <p className="eyebrow">Instância</p>
                        <h2>${selectedSession.session_id}</h2>
                      </div>
                      <div className="session-header-actions">
                        <span className=${`status-pill status-pill--${statusTone(selectedSessionStatus)}`}>
                          ${selectedSessionStatus}
                        </span>
                        <button
                          type="button"
                          className="ghost-button"
                          onClick=${() => handleSessionAction('connect', selectedSession.session_id)}
                        >
                          Conectar
                        </button>
                        <button
                          type="button"
                          className="ghost-button"
                          onClick=${() => handleSessionAction('reload', selectedSession.session_id)}
                        >
                          Recarregar
                        </button>
                        <button
                          type="button"
                          className="ghost-button ghost-button--danger"
                          onClick=${() => handleSessionAction('disconnect', selectedSession.session_id)}
                        >
                          Desconectar
                        </button>
                        <button
                          type="button"
                          className="ghost-button ghost-button--danger"
                          onClick=${() => handleSessionAction('delete', selectedSession.session_id)}
                        >
                          Excluir
                        </button>
                      </div>
                    </div>

                    <div className="session-status-strip">
                      <${DetailItem} label="Worker" value=${selectedSession.assigned ? worker?.worker_id || 'atribuída' : 'livre'} />
                      <${DetailItem} label="Conexão" value=${selectedSession.connected ? 'conectada' : 'desconectada'} />
                      <${DetailItem} label="Login" value=${selectedSession.logged_in ? 'ativo' : 'pendente'} />
                      <${DetailItem} label="WA Web" value=${selectedSession.whatsapp_web_version || '—'} code=${Boolean(selectedSession.whatsapp_web_version)} />
                    </div>

                    ${showInlineQR
                      ? html`
                          <div className="inline-qr">
                            <img
                              src=${`/sessions/${enc(selectedSession.session_id)}/qr?ts=${Date.now()}`}
                              alt=${`QR da sessão ${selectedSession.session_id}`}
                            />
                            <div>
                              <strong>QR disponível</strong>
                              <span>Escaneie para concluir a conexão desta instância.</span>
                            </div>
                          </div>
                        `
                      : null}

                    <div className="tab-row">
                      ${['quick', 'config', 'webhooks', 'templates'].map(
                        (tab) => html`
                          <button
                            key=${tab}
                            type="button"
                            className=${`tab-button ${selectedTab === tab ? 'tab-button--active' : ''}`}
                            onClick=${() => setSelectedTab(tab)}
                          >
                            ${tab === 'quick'
                              ? 'Teste rápido'
                              : tab === 'config'
                                ? 'Config'
                                : tab === 'webhooks'
                                  ? 'Webhook'
                                  : 'Templates'}
                          </button>
                        `
                      )}
                    </div>

                    ${selectedTab === 'quick'
                      ? html`
                          <div className="detail-grid">
                            <${DetailItem}
                              label="Device JID"
                              value=${selectedSession.device_jid || '—'}
                              code=${Boolean(selectedSession.device_jid)}
                            />
                            <${DetailItem}
                              label="WA Web"
                              value=${selectedSession.whatsapp_web_version || '—'}
                              code=${Boolean(selectedSession.whatsapp_web_version)}
                            />
                            <${DetailItem}
                              label="Status"
                              value=${sessionStatus(selectedSession)}
                            />
                            <${DetailItem}
                              label="Worker"
                              value=${selectedSession.assigned ? worker?.worker_id || 'atribuída' : 'livre'}
                            />
                            <${DetailItem}
                              label="Pairing"
                              value=${selectedSession.has_qr ? 'QR disponível' : 'sem QR'}
                            />
                          </div>

                          <form className="stack-form" onSubmit=${handleQuickSend}>
                            <${InputField}
                              label="Sessão"
                              value=${messageForm.session}
                              placeholder="Selecione uma sessão ou informe manualmente"
                              onInput=${(value) =>
                                setMessageForm((current) => ({ ...current, session: value }))
                              }
                            />
                            <${InputField}
                              label="Destino"
                              value=${messageForm.to}
                              placeholder="5562999999999"
                              onInput=${(value) =>
                                setMessageForm((current) => ({ ...current, to: value }))
                              }
                            />
                            <${TextAreaField}
                              label="Mensagem"
                              value=${messageForm.body}
                              rows="4"
                              placeholder="Texto da mensagem"
                              onInput=${(value) =>
                                setMessageForm((current) => ({ ...current, body: value }))
                              }
                            />
                            <button
                              type="submit"
                              className="solid-button"
                              disabled=${Boolean(busyMap[`send:${messageForm.session.trim()}`])}
                            >
                              Enviar texto
                            </button>
                          </form>

                          <form className="stack-form" onSubmit=${handleRequestPairingCode}>
                            <div className="callout callout--muted">
                              <strong>Pareamento por código</strong>
                              <p>
                                Além do QR code, você pode gerar um código de pareamento para vincular
                                a sessão sem leitura de QR. O provider conecta a sessão, espera a fase
                                de QR internamente e então pede o código ao WhatsApp.
                              </p>
                            </div>

                            <${InputField}
                              label="Telefone para o pareamento"
                              value=${pairingState.phone}
                              placeholder="5562999999999"
                              onInput=${(value) =>
                                setPairingState((current) => ({
                                  ...current,
                                  phone: value,
                                  code: '',
                                }))
                              }
                            />

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${pairingState.loading}
                              >
                                ${pairingState.loading ? 'Gerando código...' : 'Gerar código de pareamento'}
                              </button>
                            </div>

                            ${pairingState.code
                              ? html`
                                  <div className="pairing-code-card">
                                    <span className="pairing-code-card__label">Código atual</span>
                                    <code className="pairing-code-card__value">${pairingState.code}</code>
                                    <p className="hint-text pairing-code-card__hint">
                                      Esse código expira rápido. Se o pareamento não concluir a tempo,
                                      gere um novo código.
                                    </p>
                                  </div>
                                `
                              : null}
                          </form>
                        `
                      : null}

                    ${selectedTab === 'config'
                      ? html`
                          <form className="stack-form" onSubmit=${handleSaveConfig}>
                            <div
                              className=${`callout ${configState.snapshot?.has_override ? 'callout--accent' : 'callout--muted'}`}
                            >
                              <strong>
                                ${configState.snapshot?.has_override
                                  ? 'Existe um override ativo para esta sessão.'
                                  : 'A sessão está usando os defaults do worker.'}
                              </strong>
                              <p>
                                ${configState.loading
                                  ? 'Carregando configuração efetiva…'
                                  : 'Edite os parâmetros operacionais e salve no Redis/arquivo da sessão.'}
                              </p>
                            </div>

                            <div className="check-grid">
                              <${CheckField}
                                label="Rejeitar ligações"
                                checked=${configState.form.reject_calls}
                                onChange=${(checked) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, reject_calls: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Marcar conversa como lida"
                                checked=${configState.form.mark_read_on_message}
                                onChange=${(checked) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, mark_read_on_message: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Ignorar status broadcast"
                                checked=${configState.form.ignore_status_broadcast}
                                onChange=${(checked) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, ignore_status_broadcast: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Ignorar newsletters"
                                checked=${configState.form.ignore_newsletters}
                                onChange=${(checked) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, ignore_newsletters: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Outgoing em smb_message_echoes"
                                checked=${configState.form.outgoing_messages_coex}
                                onChange=${(checked) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, outgoing_messages_coex: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Manter sempre online"
                                checked=${configState.form.always_online}
                                onChange=${(checked) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, always_online: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Sincronizar histórico"
                                description="Envia HistorySync recebido para webhooks cadastrados."
                                checked=${configState.form.history_sync_enabled}
                                onChange=${(checked) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, history_sync_enabled: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Usar proxy nesta sessão"
                                checked=${configState.form.proxy_enabled}
                                onChange=${(checked) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, proxy_enabled: checked },
                                  }))
                                }
                              />
                            </div>

                            <${TextAreaField}
                              label="Mensagem de rejeição"
                              value=${configState.form.reject_calls_message}
                              rows="3"
                              placeholder="Opcional"
                              onInput=${(value) =>
                                setConfigState((current) => ({
                                  ...current,
                                  form: { ...current.form, reject_calls_message: value },
                                }))
                              }
                            />
                            <${InputField}
                              label="Intervalo always-online (segundos)"
                              value=${configState.form.always_online_interval_seconds}
                              type="number"
                              min="1"
                              step="1"
                              onInput=${(value) =>
                                setConfigState((current) => ({
                                  ...current,
                                  form: { ...current.form, always_online_interval_seconds: value },
                                }))
                              }
                            />
                            <${InputField}
                              label="Dias de histórico"
                              value=${configState.form.history_sync_days}
                              type="number"
                              min="1"
                              step="1"
                              onInput=${(value) =>
                                setConfigState((current) => ({
                                  ...current,
                                  form: { ...current.form, history_sync_days: value },
                                }))
                              }
                            />

                            <div className="form-grid">
                              <${InputField}
                                label="Proxy protocolo"
                                value=${configState.form.proxy_protocol}
                                placeholder="http, https ou socks5"
                                onInput=${(value) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, proxy_protocol: value },
                                  }))
                                }
                              />
                              <${InputField}
                                label="Proxy host"
                                value=${configState.form.proxy_host}
                                placeholder="proxy.exemplo.com"
                                onInput=${(value) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, proxy_host: value },
                                  }))
                                }
                              />
                              <${InputField}
                                label="Proxy porta"
                                value=${configState.form.proxy_port}
                                placeholder="8080"
                                onInput=${(value) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, proxy_port: value },
                                  }))
                                }
                              />
                              <${InputField}
                                label="Proxy usuário"
                                value=${configState.form.proxy_username}
                                placeholder="Opcional"
                                onInput=${(value) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, proxy_username: value },
                                  }))
                                }
                              />
                              <${InputField}
                                label="Proxy senha"
                                value=${configState.form.proxy_password}
                                type="password"
                                placeholder="Opcional"
                                onInput=${(value) =>
                                  setConfigState((current) => ({
                                    ...current,
                                    form: { ...current.form, proxy_password: value },
                                  }))
                                }
                              />
                            </div>

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${configState.loading || configState.saving}
                              >
                                ${configState.saving ? 'Salvando...' : 'Salvar configuração'}
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                disabled=${configState.loading}
                                onClick=${() => loadConfigForSession(selectedSessionId)}
                              >
                                Recarregar
                              </button>
                              <button
                                type="button"
                                className="ghost-button ghost-button--danger"
                                disabled=${configState.saving}
                                onClick=${handleClearConfig}
                              >
                                Limpar override
                              </button>
                            </div>
                          </form>
                        `
                      : null}

                    ${selectedTab === 'webhooks'
                      ? html`
                          <form className="stack-form" onSubmit=${handleSaveWebhook}>
                            <div className="form-grid">
                              <${InputField}
                                label="Webhook ID"
                                value=${webhookState.form.id}
                                placeholder="vazio = auto"
                                onInput=${(value) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, id: value },
                                  }))
                                }
                              />
                              <${InputField}
                                label="Nome"
                                value=${webhookState.form.name}
                                placeholder="Chatwoot principal"
                                onInput=${(value) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, name: value },
                                  }))
                                }
                              />
                            </div>

                            <${InputField}
                              label="URL base"
                              value=${webhookState.form.url}
                              placeholder="https://seu-app.com/webhooks/whatsapp"
                              onInput=${(value) =>
                                setWebhookState((current) => ({
                                  ...current,
                                  form: { ...current.form, url: value },
                                }))
                              }
                            />
                            <${InputField}
                              label="URL absoluta"
                              value=${webhookState.form.url_absolute}
                              placeholder="https://seu-app.com/webhooks/whatsapp/5562..."
                              onInput=${(value) =>
                                setWebhookState((current) => ({
                                  ...current,
                                  form: { ...current.form, url_absolute: value },
                                }))
                              }
                            />

                            <div className="form-grid">
                              <${InputField}
                                label="Header auth"
                                value=${webhookState.form.header}
                                placeholder="Authorization"
                                onInput=${(value) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, header: value },
                                  }))
                                }
                              />
                              <${InputField}
                                label="Token"
                                value=${webhookState.form.token}
                                placeholder="Bearer xxxxx"
                                onInput=${(value) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, token: value },
                                  }))
                                }
                              />
                            </div>

                            <${InputField}
                              label="Timeout (ms)"
                              value=${webhookState.form.timeout_ms}
                              type="number"
                              min="1"
                              step="100"
                              onInput=${(value) =>
                                setWebhookState((current) => ({
                                  ...current,
                                  form: { ...current.form, timeout_ms: value },
                                }))
                              }
                            />

                            <div className="check-grid">
                              <${CheckField}
                                label="Webhook habilitado"
                                checked=${webhookState.form.enabled}
                                onChange=${(checked) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, enabled: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Adicionar sessão à URL"
                                checked=${webhookState.form.append_session}
                                onChange=${(checked) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, append_session: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Enviar incoming"
                                checked=${webhookState.form.send_incoming_messages}
                                onChange=${(checked) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, send_incoming_messages: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Enviar outgoing"
                                checked=${webhookState.form.send_outgoing_messages}
                                onChange=${(checked) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, send_outgoing_messages: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Enviar updates"
                                checked=${webhookState.form.send_update_messages}
                                onChange=${(checked) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, send_update_messages: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Enviar grupos"
                                checked=${webhookState.form.send_group_messages}
                                onChange=${(checked) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, send_group_messages: checked },
                                  }))
                                }
                              />
                              <${CheckField}
                                label="Enviar newsletters"
                                checked=${webhookState.form.send_newsletter_messages}
                                onChange=${(checked) =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: { ...current.form, send_newsletter_messages: checked },
                                  }))
                                }
                              />
                            </div>

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${webhookState.loading || webhookState.saving}
                              >
                                ${webhookState.saving ? 'Salvando...' : 'Salvar webhook'}
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                onClick=${() =>
                                  setWebhookState((current) => ({
                                    ...current,
                                    form: createDefaultWebhookForm(),
                                  }))
                                }
                              >
                                Limpar formulário
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                disabled=${webhookState.loading}
                                onClick=${() =>
                                  loadWebhooksForSession(selectedSessionId, { resetForm: false })
                                }
                              >
                                Recarregar lista
                              </button>
                            </div>
                          </form>

                          <div className="subsection">
                            <div className="panel__header panel__header--space-between">
                              <div>
                                <p className="eyebrow">Cadastrados</p>
                                <h3>Webhooks ativos</h3>
                              </div>
                              <span className="panel__hint">${webhookState.items.length} item(ns)</span>
                            </div>

                            ${webhookState.loading
                              ? html`
                                  <div className="empty-state empty-state--compact">
                                    <strong>Carregando webhooks…</strong>
                                  </div>
                                `
                              : webhookState.items.length === 0
                                ? html`
                                    <div className="empty-state empty-state--compact">
                                      <strong>Nenhum webhook cadastrado</strong>
                                      <p>Preencha o formulário acima para salvar o primeiro callback.</p>
                                    </div>
                                  `
                                : html`
                                    <div className="webhook-list">
                                      ${webhookState.items.map(
                                        (item) => html`
                                          <${WebhookCard}
                                            key=${item.id}
                                            item=${item}
                                            onEdit=${() =>
                                              setWebhookState((current) => ({
                                                ...current,
                                                form: createWebhookFormFromItem(item),
                                              }))
                                            }
                                            onDelete=${() => handleDeleteWebhook(item.id)}
                                          />
                                        `
                                      )}
                                    </div>
                                  `}
                          </div>
                        `
                      : null}

                    ${selectedTab === 'templates'
                      ? html`
                          <form className="stack-form" onSubmit=${handleSaveTemplate}>
                            <div className="form-grid">
                              <${InputField}
                                label="Template ID"
                                value=${templateState.form.id}
                                placeholder="vazio = nome__idioma"
                                onInput=${(value) =>
                                  setTemplateState((current) => ({
                                    ...current,
                                    form: { ...current.form, id: value },
                                  }))
                                }
                              />
                              <${InputField}
                                label="Nome Cloud API"
                                value=${templateState.form.name}
                                placeholder="boleto_disponivel"
                                onInput=${(value) =>
                                  setTemplateState((current) => ({
                                    ...current,
                                    form: { ...current.form, name: value },
                                  }))
                                }
                              />
                              <${InputField}
                                label="Idioma"
                                value=${templateState.form.language}
                                placeholder="pt_BR"
                                onInput=${(value) =>
                                  setTemplateState((current) => ({
                                    ...current,
                                    form: { ...current.form, language: value },
                                  }))
                                }
                              />
                              <${InputField}
                                label="Categoria"
                                value=${templateState.form.category}
                                placeholder="utility"
                                onInput=${(value) =>
                                  setTemplateState((current) => ({
                                    ...current,
                                    form: { ...current.form, category: value },
                                  }))
                                }
                              />
                            </div>

                            <div className="form-grid">
                              <label className="field">
                                <span>Tipo do cabeçalho</span>
                                <select
                                  value=${templateState.form.header_type}
                                  onChange=${(event) =>
                                    setTemplateState((current) => ({
                                      ...current,
                                      form: { ...current.form, header_type: event.target.value },
                                    }))
                                  }
                                >
                                  <option value="">Sem cabeçalho</option>
                                  <option value="text">Texto</option>
                                  <option value="image">Imagem</option>
                                  <option value="document">Documento</option>
                                </select>
                              </label>
                              <${InputField}
                                label="Texto do cabeçalho"
                                value=${templateState.form.header_text}
                                placeholder="Olá {{nome}}"
                                onInput=${(value) =>
                                  setTemplateState((current) => ({
                                    ...current,
                                    form: { ...current.form, header_text: value },
                                  }))
                                }
                              />
                            </div>

                            ${templateState.form.header_type === 'image' ||
                            templateState.form.header_type === 'document'
                              ? html`
                                  <div className="form-grid">
                                    <${InputField}
                                      label="Link padrão da mídia"
                                      value=${templateState.form.header_link}
                                      placeholder="https://..."
                                      onInput=${(value) =>
                                        setTemplateState((current) => ({
                                          ...current,
                                          form: { ...current.form, header_link: value },
                                        }))
                                      }
                                    />
                                    <${InputField}
                                      label="Media ID padrão"
                                      value=${templateState.form.header_id}
                                      placeholder="opcional"
                                      onInput=${(value) =>
                                        setTemplateState((current) => ({
                                          ...current,
                                          form: { ...current.form, header_id: value },
                                        }))
                                      }
                                    />
                                    <${InputField}
                                      label="Nome do arquivo"
                                      value=${templateState.form.header_filename}
                                      placeholder="documento.pdf"
                                      onInput=${(value) =>
                                        setTemplateState((current) => ({
                                          ...current,
                                          form: { ...current.form, header_filename: value },
                                        }))
                                      }
                                    />
                                  </div>
                                `
                              : null}

                            <${TextAreaField}
                              label="Corpo"
                              value=${templateState.form.body}
                              rows="5"
                              placeholder="Olá {{nome}}, seu boleto {{1}} está disponível."
                              onInput=${(value) =>
                                setTemplateState((current) => ({
                                  ...current,
                                  form: { ...current.form, body: value },
                                }))
                              }
                            />

                            <${TextAreaField}
                              label="Footer"
                              value=${templateState.form.footer}
                              rows="2"
                              placeholder="Opcional"
                              onInput=${(value) =>
                                setTemplateState((current) => ({
                                  ...current,
                                  form: { ...current.form, footer: value },
                                }))
                              }
                            />

                            <div className="callout callout--muted">
                              <strong>Variáveis aceitas</strong>
                              <p>
                                Use placeholders como {{'{{1}}'}} ou {{'{{nome}}'}}. No envio, eles podem
                                vir em components.parameters por ordem, com parameter_name, ou no objeto
                                template.variables.
                              </p>
                            </div>

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${templateState.loading || templateState.saving}
                              >
                                ${templateState.saving ? 'Salvando...' : 'Salvar template'}
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                onClick=${() =>
                                  setTemplateState((current) => ({
                                    ...current,
                                    form: createDefaultTemplateForm(),
                                  }))
                                }
                              >
                                Limpar formulário
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                disabled=${templateState.loading}
                                onClick=${() =>
                                  loadTemplatesForSession(selectedSessionId, { resetForm: false })
                                }
                              >
                                Recarregar lista
                              </button>
                            </div>
                          </form>

                          <div className="subsection">
                            <div className="panel__header panel__header--space-between">
                              <div>
                                <p className="eyebrow">Cadastrados</p>
                                <h3>Templates da sessão</h3>
                              </div>
                              <span className="panel__hint">${templateState.items.length} item(ns)</span>
                            </div>

                            ${templateState.loading
                              ? html`
                                  <div className="empty-state empty-state--compact">
                                    <strong>Carregando templates…</strong>
                                  </div>
                                `
                              : templateState.items.length === 0
                                ? html`
                                    <div className="empty-state empty-state--compact">
                                      <strong>Nenhum template cadastrado</strong>
                                      <p>Preencha o formulário acima para habilitar envios type=template.</p>
                                    </div>
                                  `
                                : html`
                                    <div className="webhook-list">
                                      ${templateState.items.map(
                                        (item) => html`
                                          <${TemplateCard}
                                            key=${item.id}
                                            item=${item}
                                            onEdit=${() =>
                                              setTemplateState((current) => ({
                                                ...current,
                                                form: createTemplateFormFromItem(item),
                                              }))
                                            }
                                            onDelete=${() => handleDeleteTemplate(item.id)}
                                          />
                                        `
                                      )}
                                    </div>
                                  `}
                          </div>
                        `
                      : null}
                  `
                : html`
                    <div className="empty-state">
                      <strong>Nenhuma sessão selecionada</strong>
                      <p>
                        Escolha uma sessão na coluna central para abrir o resumo, a configuração e os
                        webhooks dessa sessão.
                      </p>
                    </div>
                  `}
            </section>
          </div>
        </section>
          `
          : null}
      </main>

      <div className="toast-stack">
        ${toasts.map(
          (toast) => html`
            <div key=${toast.id} className=${`toast toast--${toast.tone}`}>
              <strong>${toast.tone === 'error' ? 'Erro' : toast.tone === 'success' ? 'OK' : 'Info'}</strong>
              <span>${toast.message}</span>
            </div>
          `
        )}
      </div>

      <${QRModal} sessionId=${qrSessionId} onClose=${() => setQrSessionId('')} />
    </div>
  `
}

const mountNode = document.getElementById('root')

if (mountNode) {
  createRoot(mountNode).render(html`<${App} />`)
}
