import { useMemo, useState } from 'react'
import { Input, EmptyState, cx } from './ui.jsx'
import { SessionCard } from './SessionCard.jsx'
import { SessionDetail } from './SessionDetail.jsx'
import { isOnline, sessionStatus, sortSessions } from '../lib/sessions.js'

function FilterChip({ active, label, count, onClick }) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cx(
        'inline-flex items-center gap-1.5 rounded-full border px-3 py-1 text-xs font-medium transition-colors',
        active
          ? 'border-brand bg-brand-soft text-brand'
          : 'border-border text-muted hover:text-fg',
      )}
    >
      {label}
      <span className={cx('font-mono', active ? 'text-brand' : 'text-faint')}>{count}</span>
    </button>
  )
}

export function Instances({
  sessions,
  worker,
  selectedId,
  detailTab,
  busy,
  onSelect,
  onAction,
  onShowQR,
  notify,
}) {
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState('all')

  const counts = useMemo(
    () => ({
      all: sessions.length,
      assigned: sessions.filter((s) => s.assigned).length,
      online: sessions.filter(isOnline).length,
      qr: sessions.filter((s) => s.has_qr).length,
    }),
    [sessions],
  )

  const visible = useMemo(() => {
    const needle = search.trim().toLowerCase()
    return [...sessions]
      .sort(sortSessions)
      .filter((s) => {
        if (filter === 'assigned' && !s.assigned) return false
        if (filter === 'online' && !isOnline(s)) return false
        if (filter === 'qr' && !s.has_qr) return false
        if (!needle) return true
        return [s.session_id, s.device_jid, sessionStatus(s)]
          .filter(Boolean)
          .join(' ')
          .toLowerCase()
          .includes(needle)
      })
  }, [sessions, search, filter])

  const selected = sessions.find((s) => s.session_id === selectedId) || null

  return (
    <div className="grid gap-6 lg:grid-cols-[minmax(0,22rem)_minmax(0,1fr)]">
      <div className="flex flex-col gap-3">
        <Input
          value={search}
          placeholder="Buscar por ID, JID ou status…"
          onChange={(e) => setSearch(e.target.value)}
        />
        <div className="flex flex-wrap gap-1.5">
          <FilterChip label="Todas" count={counts.all} active={filter === 'all'} onClick={() => setFilter('all')} />
          <FilterChip label="Atribuídas" count={counts.assigned} active={filter === 'assigned'} onClick={() => setFilter('assigned')} />
          <FilterChip label="Online" count={counts.online} active={filter === 'online'} onClick={() => setFilter('online')} />
          <FilterChip label="Com QR" count={counts.qr} active={filter === 'qr'} onClick={() => setFilter('qr')} />
        </div>
        <div className="flex flex-col gap-3">
          {visible.length === 0 ? (
            <EmptyState
              title="Nenhuma instância"
              description="Ajuste os filtros ou crie uma sessão na Home."
            />
          ) : (
            visible.map((session) => (
              <SessionCard
                key={session.session_id}
                session={session}
                selected={session.session_id === selectedId}
                busy={busy}
                onSelect={onSelect}
                onAction={onAction}
              />
            ))
          )}
        </div>
      </div>

      <div>
        {selected ? (
          <SessionDetail
            key={selected.session_id}
            session={selected}
            worker={worker}
            initialTab={detailTab}
            notify={notify}
            onAction={onAction}
            onShowQR={onShowQR}
          />
        ) : (
          <EmptyState
            title="Selecione uma instância"
            description="Escolha uma sessão à esquerda para ver detalhes, configuração e webhooks."
          />
        )}
      </div>
    </div>
  )
}
