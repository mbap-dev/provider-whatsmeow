import { Button, Field, Input, Mono, Switch, Tag, Eyebrow, EmptyState, Spinner } from './ui.jsx'

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
      </div>
    </article>
  )
}

export function WebhookPanel({
  items,
  loading,
  saving,
  form,
  onChange,
  onSave,
  onDelete,
  onEdit,
  onResetForm,
}) {
  const set = (patch) => onChange({ ...form, ...patch })
  const editing = Boolean(form.id)

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      <div className="flex flex-col gap-3">
        <Eyebrow>Webhooks configurados</Eyebrow>
        {loading ? (
          <div className="flex justify-center py-8 text-muted">
            <Spinner />
          </div>
        ) : items.length === 0 ? (
          <EmptyState
            title="Nenhum webhook"
            description="Crie um endpoint para receber os eventos desta sessão."
          />
        ) : (
          items.map((item) => (
            <WebhookItem key={item.id} item={item} onEdit={onEdit} onDelete={onDelete} />
          ))
        )}
      </div>

      <form onSubmit={onSave} className="flex flex-col gap-3">
        <div className="flex items-center justify-between">
          <Eyebrow>{editing ? `Editando ${form.id}` : 'Novo webhook'}</Eyebrow>
          {editing && (
            <Button size="sm" variant="ghost" onClick={onResetForm}>
              Cancelar edição
            </Button>
          )}
        </div>
        <Field label="Nome">
          <Input value={form.name} onChange={(e) => set({ name: e.target.value })} />
        </Field>
        <Field label="URL base" hint="A sessão é anexada ao final quando habilitado.">
          <Input
            value={form.url}
            placeholder="https://api.exemplo.com/webhook"
            onChange={(e) => set({ url: e.target.value })}
          />
        </Field>
        <Field label="URL absoluta" hint="Usada exatamente como está, sem anexar a sessão.">
          <Input
            value={form.url_absolute}
            placeholder="https://api.exemplo.com/hooks/123"
            onChange={(e) => set({ url_absolute: e.target.value })}
          />
        </Field>
        <div className="grid gap-3 sm:grid-cols-2">
          <Field label="Header de auth">
            <Input value={form.header} onChange={(e) => set({ header: e.target.value })} />
          </Field>
          <Field label="Token">
            <Input value={form.token} onChange={(e) => set({ token: e.target.value })} />
          </Field>
          <Field label="Timeout (ms)">
            <Input
              type="number"
              min="1"
              value={form.timeout_ms}
              onChange={(e) => set({ timeout_ms: e.target.value })}
            />
          </Field>
        </div>
        <div className="flex flex-col gap-2 rounded-lg border border-border p-3">
          <Switch label="Ativo" checked={form.enabled} onChange={(v) => set({ enabled: v })} />
          <Switch
            label="Anexar sessão à URL"
            checked={form.append_session}
            onChange={(v) => set({ append_session: v })}
          />
          <Switch
            label="Mensagens recebidas"
            checked={form.send_incoming_messages}
            onChange={(v) => set({ send_incoming_messages: v })}
          />
          <Switch
            label="Mensagens enviadas"
            checked={form.send_outgoing_messages}
            onChange={(v) => set({ send_outgoing_messages: v })}
          />
          <Switch
            label="Atualizações de status"
            checked={form.send_update_messages}
            onChange={(v) => set({ send_update_messages: v })}
          />
          <Switch
            label="Mensagens de grupo"
            checked={form.send_group_messages}
            onChange={(v) => set({ send_group_messages: v })}
          />
          <Switch
            label="Newsletters"
            checked={form.send_newsletter_messages}
            onChange={(v) => set({ send_newsletter_messages: v })}
          />
        </div>
        <Button type="submit" variant="primary" loading={saving}>
          {editing ? 'Salvar alterações' : 'Criar webhook'}
        </Button>
      </form>
    </div>
  )
}
