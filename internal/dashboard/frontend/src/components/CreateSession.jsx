import { useState } from 'react'
import { Button, Field, Input, Switch, Eyebrow } from './ui.jsx'
import * as api from '../api/endpoints.js'
import { defaultWebhookForm, webhookFormToPayload } from '../lib/forms.js'

function initialForm() {
  return {
    id: '',
    assign: true,
    connect: true,
    history_sync_enabled: false,
    history_sync_days: '7',
    webhook: defaultWebhookForm(),
  }
}

export function CreateSession({ notify, onCreated }) {
  const [form, setForm] = useState(initialForm)
  const [busy, setBusy] = useState(false)
  const set = (patch) => setForm((current) => ({ ...current, ...patch }))
  const setWebhook = (patch) =>
    setForm((current) => ({ ...current, webhook: { ...current.webhook, ...patch } }))

  async function handleSubmit(event) {
    event.preventDefault()
    const id = form.id.trim()
    if (!id) {
      notify('Informe o ID da sessão.', 'error')
      return
    }
    const days = parseInt(form.history_sync_days, 10)
    if (form.history_sync_enabled) {
      if (!(days > 0)) {
        notify('Informe a quantidade de dias do histórico.', 'error')
        return
      }
      if (!form.webhook.url.trim() && !form.webhook.url_absolute.trim()) {
        notify('Informe o webhook antes de ativar a sincronização de histórico.', 'error')
        return
      }
    }

    setBusy(true)
    try {
      const body = {
        id,
        assign: form.assign,
        connect: form.connect,
        config: {
          history_sync_enabled: form.history_sync_enabled,
          history_sync_days: form.history_sync_enabled && days > 0 ? days : 0,
        },
      }
      if (form.history_sync_enabled) {
        body.webhook = webhookFormToPayload(form.webhook)
      }
      await api.createSession(body)
      notify(`Sessão ${id} criada.`, 'success')
      setForm(initialForm())
      onCreated?.(id)
    } catch (error) {
      notify(error.message, 'error')
    } finally {
      setBusy(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4 p-5">
      <Field label="ID da sessão" hint="Geralmente o número do telefone, sem símbolos.">
        <Input
          value={form.id}
          placeholder="5511999999999"
          onChange={(e) => set({ id: e.target.value })}
        />
      </Field>
      <div className="flex flex-col gap-2 rounded-lg border border-border p-3">
        <Switch
          label="Atribuir a este worker"
          checked={form.assign}
          onChange={(v) => set({ assign: v })}
        />
        <Switch
          label="Conectar imediatamente"
          checked={form.connect}
          onChange={(v) => set({ connect: v })}
        />
        <Switch
          label="Sincronizar histórico ao parear"
          checked={form.history_sync_enabled}
          onChange={(v) => set({ history_sync_enabled: v })}
        />
      </div>

      {form.history_sync_enabled && (
        <div className="flex flex-col gap-3 rounded-lg border border-border bg-surface-2 p-3">
          <Eyebrow>Histórico exige um webhook de destino</Eyebrow>
          <Field label="Dias de histórico">
            <Input
              type="number"
              min="1"
              value={form.history_sync_days}
              onChange={(e) => set({ history_sync_days: e.target.value })}
            />
          </Field>
          <Field label="URL do webhook">
            <Input
              value={form.webhook.url}
              placeholder="https://api.exemplo.com/webhook"
              onChange={(e) => setWebhook({ url: e.target.value })}
            />
          </Field>
          <div className="grid gap-3 sm:grid-cols-2">
            <Field label="Header de auth" hint="Nome do header enviado ao webhook.">
              <Input
                value={form.webhook.header}
                placeholder="Authorization"
                onChange={(e) => setWebhook({ header: e.target.value })}
              />
            </Field>
            <Field label="Token" hint="Valor enviado no header.">
              <Input
                value={form.webhook.token}
                onChange={(e) => setWebhook({ token: e.target.value })}
              />
            </Field>
          </div>
        </div>
      )}

      <Button type="submit" variant="primary" loading={busy}>
        Criar sessão
      </Button>
    </form>
  )
}
