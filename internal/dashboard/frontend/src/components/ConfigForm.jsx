import { Button, Field, Input, Textarea, Select, Switch, Eyebrow } from './ui.jsx'

export function ConfigForm({ form, snapshot, saving, onChange, onSave, onClear }) {
  const set = (patch) => onChange({ ...form, ...patch })
  const hasOverride = Boolean(snapshot?.has_override)

  return (
    <form onSubmit={onSave} className="flex flex-col gap-6">
      <fieldset className="flex flex-col gap-3">
        <Eyebrow>Atendimento</Eyebrow>
        <Switch
          label="Recusar chamadas"
          description="Encerra chamadas recebidas automaticamente."
          checked={form.reject_calls}
          onChange={(v) => set({ reject_calls: v })}
        />
        {form.reject_calls && (
          <Field label="Mensagem ao recusar">
            <Textarea
              rows={3}
              value={form.reject_calls_message}
              placeholder="Não atendemos chamadas por aqui."
              onChange={(e) => set({ reject_calls_message: e.target.value })}
            />
          </Field>
        )}
        <Switch
          label="Marcar como lida ao receber"
          description="Confirma a leitura das mensagens recebidas."
          checked={form.mark_read_on_message}
          onChange={(v) => set({ mark_read_on_message: v })}
        />
        <Switch
          label="Ignorar status (stories)"
          checked={form.ignore_status_broadcast}
          onChange={(v) => set({ ignore_status_broadcast: v })}
        />
        <Switch
          label="Ignorar newsletters"
          checked={form.ignore_newsletters}
          onChange={(v) => set({ ignore_newsletters: v })}
        />
        <Switch
          label="Coexistência de mensagens enviadas"
          description="Repassa também mensagens enviadas por outros aparelhos."
          checked={form.outgoing_messages_coex}
          onChange={(v) => set({ outgoing_messages_coex: v })}
        />
      </fieldset>

      <fieldset className="flex flex-col gap-3">
        <Eyebrow>Presença</Eyebrow>
        <Switch
          label="Sempre online"
          description="Mantém a sessão marcada como online."
          checked={form.always_online}
          onChange={(v) => set({ always_online: v })}
        />
        {form.always_online && (
          <Field label="Intervalo (segundos)">
            <Input
              type="number"
              min="1"
              value={form.always_online_interval_seconds}
              onChange={(e) => set({ always_online_interval_seconds: e.target.value })}
            />
          </Field>
        )}
      </fieldset>

      <fieldset className="flex flex-col gap-3">
        <Eyebrow>Histórico</Eyebrow>
        <Switch
          label="Sincronizar histórico"
          description="Importa mensagens anteriores ao parear."
          checked={form.history_sync_enabled}
          onChange={(v) => set({ history_sync_enabled: v })}
        />
        {form.history_sync_enabled && (
          <Field label="Dias de histórico">
            <Input
              type="number"
              min="1"
              value={form.history_sync_days}
              onChange={(e) => set({ history_sync_days: e.target.value })}
            />
          </Field>
        )}
      </fieldset>

      <fieldset className="flex flex-col gap-3">
        <Eyebrow>Proxy</Eyebrow>
        <Switch
          label="Usar proxy"
          description="Roteia a conexão da sessão por um proxy."
          checked={form.proxy_enabled}
          onChange={(v) => set({ proxy_enabled: v })}
        />
        {form.proxy_enabled && (
          <div className="grid gap-3 sm:grid-cols-2">
            <Field label="Protocolo">
              <Select
                value={form.proxy_protocol}
                onChange={(e) => set({ proxy_protocol: e.target.value })}
              >
                <option value="">Selecione</option>
                <option value="http">http</option>
                <option value="https">https</option>
                <option value="socks5">socks5</option>
              </Select>
            </Field>
            <Field label="Host">
              <Input
                value={form.proxy_host}
                placeholder="proxy.exemplo.com"
                onChange={(e) => set({ proxy_host: e.target.value })}
              />
            </Field>
            <Field label="Porta">
              <Input
                value={form.proxy_port}
                placeholder="1080"
                onChange={(e) => set({ proxy_port: e.target.value })}
              />
            </Field>
            <Field label="Usuário">
              <Input
                value={form.proxy_username}
                onChange={(e) => set({ proxy_username: e.target.value })}
              />
            </Field>
            <Field label="Senha">
              <Input
                type="password"
                value={form.proxy_password}
                onChange={(e) => set({ proxy_password: e.target.value })}
              />
            </Field>
          </div>
        )}
      </fieldset>

      <div className="flex items-center gap-2 border-t border-border pt-4">
        <Button type="submit" variant="primary" loading={saving}>
          Salvar configuração
        </Button>
        <Button
          type="button"
          variant="ghost"
          onClick={onClear}
          disabled={saving || !hasOverride}
        >
          Restaurar padrão
        </Button>
      </div>
    </form>
  )
}
