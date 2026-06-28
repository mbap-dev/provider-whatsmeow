import { Button, Mono, SignalDot, StatusBadge, Tag, cx } from './ui.jsx'
import {
  isOnline,
  sessionStatus,
  statusLabel,
  statusTone,
  truncate,
} from '../lib/sessions.js'

export function SessionCard({ session, selected, busy, onSelect, onAction }) {
  const status = sessionStatus(session)
  const tone = statusTone(status)
  const id = session.session_id
  const assigned = Boolean(session.assigned)
  const busyFor = (action) => Boolean(busy?.[`${action}:${id}`])

  return (
    <article
      role="button"
      tabIndex={0}
      onClick={() => onSelect(id)}
      onKeyDown={(event) => {
        if (event.key === 'Enter' || event.key === ' ') {
          event.preventDefault()
          onSelect(id)
        }
      }}
      className={cx(
        'group cursor-pointer rounded-[var(--radius-card)] border bg-surface p-4 transition-colors',
        selected
          ? 'border-brand ring-1 ring-brand'
          : 'border-border hover:border-border-strong',
      )}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex min-w-0 items-center gap-2.5">
          <SignalDot tone={tone} pulse={tone === 'live'} className="mt-0.5" />
          <div className="min-w-0">
            <h3 className="truncate font-display text-sm font-semibold text-fg">{id}</h3>
            <p className="text-xs text-muted">
              {statusLabel(status)} · {session.logged_in ? 'login ativo' : 'login pendente'}
            </p>
          </div>
        </div>
        {session.has_qr && <StatusBadge tone="pending">QR pronto</StatusBadge>}
      </div>

      <div className="mt-3 flex flex-wrap gap-1.5">
        <Tag>{assigned ? 'Atribuída' : 'Sem atribuição'}</Tag>
        {session.whatsapp_web_version && <Tag muted>WA {session.whatsapp_web_version}</Tag>}
        {session.device_jid && <Tag muted>{truncate(session.device_jid, 22)}</Tag>}
      </div>

      <div
        className="mt-3 flex flex-wrap gap-1.5"
        onClick={(event) => event.stopPropagation()}
      >
        <Button
          size="sm"
          variant="primary"
          onClick={() => onAction('connect', id)}
          loading={busyFor('connect')}
        >
          Conectar
        </Button>
        <Button size="sm" onClick={() => onAction('reload', id)} loading={busyFor('reload')}>
          Recarregar
        </Button>
        <Button
          size="sm"
          onClick={() => onAction(assigned ? 'unassign' : 'assign', id)}
          loading={busyFor(assigned ? 'unassign' : 'assign')}
        >
          {assigned ? 'Desatribuir' : 'Atribuir'}
        </Button>
        <Button
          size="sm"
          variant="danger"
          onClick={() => onAction('disconnect', id)}
          loading={busyFor('disconnect')}
        >
          Desconectar
        </Button>
        <Button
          size="sm"
          variant="danger"
          onClick={() => onAction('delete', id)}
          loading={busyFor('delete')}
        >
          Excluir
        </Button>
      </div>
    </article>
  )
}
