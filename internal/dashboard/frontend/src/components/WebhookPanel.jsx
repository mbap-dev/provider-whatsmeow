import { Button, Mono, Tag, Eyebrow, EmptyState, Spinner } from './ui.jsx'

function WebhookItem({ item, onEdit, onDelete }) {
  return (
    <article className="rounded-lg border border-border bg-surface-2 p-3">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="truncate text-sm font-medium text-fg">{item.name || 'Sem nome'}</p>
          <Mono className="text-faint">{item.id}</Mono>
        </div>
        <div className="flex shrink-0 gap-1">
          <Button size="sm" variant="ghost" onClick={() => onEdit(item)}>
            Editar
          </Button>
          <Button size="sm" variant="danger" onClick={() => onDelete(item.id)}>
            Excluir
          </Button>
        </div>
      </div>
      <Mono className="mt-2 block truncate text-muted">
        {item.url_absolute || item.url || '—'}
      </Mono>
      <div className="mt-2 flex flex-wrap gap-1.5">
        <Tag>{item.enabled ? 'Ativo' : 'Desligado'}</Tag>
        {item.send_incoming_messages && <Tag muted>Recebidas</Tag>}
        {item.send_outgoing_messages && <Tag muted>Enviadas</Tag>}
        {item.send_update_messages && <Tag muted>Status</Tag>}
        {item.send_group_messages && <Tag muted>Grupos</Tag>}
        {item.send_newsletter_messages && <Tag muted>Newsletters</Tag>}
        {item.send_session_disconnected && <Tag muted>Sessão desconectada</Tag>}
        {item.send_session_offline && <Tag muted>Instância offline</Tag>}
      </div>
    </article>
  )
}

export function WebhookPanel({
  items,
  loading,
  onDelete,
  onEdit,
  onCreate,
}) {
  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-between gap-3">
        <Eyebrow>Webhooks configurados</Eyebrow>
        <Button size="sm" variant="primary" onClick={onCreate}>
          Novo webhook
        </Button>
      </div>
      {loading ? (
        <div className="flex justify-center py-8 text-muted">
          <Spinner />
        </div>
      ) : items.length === 0 ? (
        <EmptyState
          title="Nenhum webhook"
          description="Crie um endpoint para receber os eventos desta sessão."
          action={
            <Button size="sm" variant="primary" onClick={onCreate}>
              Novo webhook
            </Button>
          }
        />
      ) : (
        items.map((item) => (
          <WebhookItem key={item.id} item={item} onEdit={onEdit} onDelete={onDelete} />
        ))
      )}
    </div>
  )
}
