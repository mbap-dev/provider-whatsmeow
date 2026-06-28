import { Card, CardHeader, SignalDot, Tag, Mono, Button, Eyebrow } from './ui.jsx'
import { CreateSession } from './CreateSession.jsx'

function StatTile({ tone, value, label }) {
  return (
    <div className="flex items-center gap-3 rounded-[var(--radius-card)] border border-border bg-surface px-4 py-3">
      <SignalDot tone={tone} pulse={tone === 'live' && value > 0} />
      <div>
        <p className="font-display text-2xl font-semibold leading-none text-fg">{value}</p>
        <p className="mt-1 text-xs text-muted">{label}</p>
      </div>
    </div>
  )
}

function Info({ label, children }) {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-border py-2 last:border-0">
      <span className="text-xs text-muted">{label}</span>
      <span className="min-w-0 truncate text-right text-sm text-fg">{children}</span>
    </div>
  )
}

export function Home({ worker, counts, notify, onCreated, onGoToInstances }) {
  const bindings = worker?.amqp_bindings || []
  return (
    <div className="flex flex-col gap-6">
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <StatTile tone="neutral" value={counts.total} label="Instâncias" />
        <StatTile tone="live" value={counts.online} label="Online" />
        <StatTile tone="pending" value={counts.withQR} label="Aguardando QR" />
        <StatTile tone="live" value={counts.assigned} label="Atribuídas" />
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader
            eyebrow="Worker"
            title={worker?.worker_id || '—'}
            description={
              worker?.worker_strict
                ? 'Modo strict: aceita apenas sessões da lista.'
                : 'Aceita sessões fora da lista atribuída.'
            }
            actions={
              <Button size="sm" onClick={onGoToInstances}>
                Ver instâncias
              </Button>
            }
          />
          <div className="px-5 py-4">
            <Info label="Fila AMQP">
              {worker?.amqp_enabled ? <Mono>{worker?.amqp_queue || '—'}</Mono> : 'desabilitada'}
            </Info>
            <Info label="Exchange">
              <Mono>{worker?.amqp_exchange || '—'}</Mono>
            </Info>
            <Info label="Auto-conectar atribuídas">
              {worker?.worker_auto_connect_assigned ? 'sim' : 'não'}
            </Info>
            <Info label="Sessões atribuídas">{(worker?.assigned_sessions || []).length}</Info>
            {bindings.length > 0 && (
              <div className="pt-3">
                <Eyebrow>Bindings</Eyebrow>
                <div className="mt-2 flex flex-wrap gap-1.5">
                  {bindings.map((binding) => (
                    <Tag key={binding} muted>
                      {binding}
                    </Tag>
                  ))}
                </div>
              </div>
            )}
          </div>
        </Card>

        <Card>
          <CardHeader
            eyebrow="Nova instância"
            title="Criar sessão"
            description="Provisione e conecte um novo número."
          />
          <CreateSession notify={notify} onCreated={onCreated} />
        </Card>
      </div>
    </div>
  )
}
