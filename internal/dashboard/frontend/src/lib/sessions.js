// Session status derivation + sorting shared across the session views.

export function sessionStatus(session) {
  const raw = String(session?.status || '').trim().toLowerCase()
  if (raw === 'connected') return 'online'
  if (raw) return raw
  return session?.connected ? 'online' : 'offline'
}

// Maps any derived status onto one of three signal tones.
export function statusTone(status) {
  if (status === 'online') return 'live'
  if (status === 'connecting') return 'pending'
  return 'down'
}

const STATUS_LABEL = {
  online: 'Online',
  connecting: 'Conectando',
  offline: 'Offline',
}

export function statusLabel(status) {
  return STATUS_LABEL[status] || status || 'Desconhecido'
}

export function isOnline(session) {
  return Boolean(session?.connected) || sessionStatus(session) === 'online'
}

export function sortSessions(left, right) {
  const score = (s) =>
    (s?.assigned ? 0 : 4) +
    (sessionStatus(s) === 'online' ? 0 : sessionStatus(s) === 'connecting' ? 1 : 2)
  const diff = score(left) - score(right)
  if (diff !== 0) return diff
  return String(left?.session_id || '').localeCompare(
    String(right?.session_id || ''),
    'pt-BR',
  )
}

export function digitsOnly(value) {
  return String(value || '').replace(/\D+/g, '')
}

export function truncate(value, limit = 56) {
  if (!value) return '—'
  return value.length > limit ? `${value.slice(0, limit)}…` : value
}

export function formatTime(date) {
  if (!date) return '—'
  return new Intl.DateTimeFormat('pt-BR', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(date)
}
