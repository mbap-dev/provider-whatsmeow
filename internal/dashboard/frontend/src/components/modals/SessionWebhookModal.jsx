import { Field, Input, Switch } from '../ui.jsx'
import { FormModal } from './FormModal.jsx'

export function SessionWebhookModal({
  open,
  onClose,
  form,
  saving,
  onChange,
  onSave,
}) {
  const set = (patch) => onChange({ ...form, ...patch })
  const editing = Boolean(form.id)

  return (
    <FormModal
      open={open}
      onClose={onClose}
      eyebrow={editing ? 'Editar webhook' : 'Novo webhook'}
      title={editing ? form.name || form.id : 'Cadastrar webhook'}
      maxWidth="max-w-2xl"
      onSubmit={onSave}
      submitLabel={editing ? 'Salvar alterações' : 'Criar webhook'}
      submitting={saving}
    >
      <Field label="Nome">
        <Input value={form.name} onChange={(event) => set({ name: event.target.value })} />
      </Field>
      <Field label="URL base" hint="A sessão é anexada ao final quando habilitado.">
        <Input
          value={form.url}
          placeholder="https://api.exemplo.com/webhook"
          onChange={(event) => set({ url: event.target.value })}
        />
      </Field>
      <Field label="URL absoluta" hint="Usada exatamente como está, sem anexar a sessão.">
        <Input
          value={form.url_absolute}
          placeholder="https://api.exemplo.com/hooks/123"
          onChange={(event) => set({ url_absolute: event.target.value })}
        />
      </Field>
      <div className="grid gap-3 sm:grid-cols-2">
        <Field label="Header de auth">
          <Input value={form.header} onChange={(event) => set({ header: event.target.value })} />
        </Field>
        <Field label="Token">
          <Input value={form.token} onChange={(event) => set({ token: event.target.value })} />
        </Field>
        <Field label="Timeout (ms)">
          <Input
            type="number"
            min="1"
            value={form.timeout_ms}
            onChange={(event) => set({ timeout_ms: event.target.value })}
          />
        </Field>
      </div>
      <div className="flex flex-col gap-2 rounded-lg border border-border p-3">
        <Switch label="Ativo" checked={form.enabled} onChange={(value) => set({ enabled: value })} />
        <Switch
          label="Anexar sessão à URL"
          checked={form.append_session}
          onChange={(value) => set({ append_session: value })}
        />
        <Switch
          label="Mensagens recebidas"
          checked={form.send_incoming_messages}
          onChange={(value) => set({ send_incoming_messages: value })}
        />
        <Switch
          label="Mensagens enviadas"
          checked={form.send_outgoing_messages}
          onChange={(value) => set({ send_outgoing_messages: value })}
        />
        <Switch
          label="Atualizações de status"
          checked={form.send_update_messages}
          onChange={(value) => set({ send_update_messages: value })}
        />
        <Switch
          label="Mensagens de grupo"
          checked={form.send_group_messages}
          onChange={(value) => set({ send_group_messages: value })}
        />
        <Switch
          label="Newsletters"
          checked={form.send_newsletter_messages}
          onChange={(value) => set({ send_newsletter_messages: value })}
        />
        <Switch
          label="Sessão desconectada"
          checked={form.send_session_disconnected}
          onChange={(value) => set({ send_session_disconnected: value })}
        />
        <Switch
          label="Instância offline"
          description="Dispara quando a instância permanece offline por 1 minuto."
          checked={form.send_session_offline}
          onChange={(value) => set({ send_session_offline: value })}
        />
      </div>
    </FormModal>
  )
}
