import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Button,
  Card,
  CardHeader,
  EmptyState,
  Field,
  Input,
  Modal,
  Mono,
  Select,
  Spinner,
  StatusBadge,
  Switch,
  Tag,
  cx,
} from './ui.jsx'
import * as api from '../api/endpoints.js'
import {
  defaultGlobalWebhookForm,
  webhookFormFromItem,
  webhookFormToPayload,
} from '../lib/forms.js'

const EMPTY_FORM = {
  id: '',
  email: '',
  password: '',
  role: 'user',
  assigned_sessions: [],
}

const ROLE_LABEL = {
  owner: 'Owner',
  admin: 'Admin',
  user: 'Usuário',
}

const SECTIONS = [
  { id: 'users', label: 'Usuários' },
  { id: 'webhooks', label: 'Webhooks' },
]

function fmtDate(value) {
  return value ? new Date(value).toLocaleString('pt-BR') : '—'
}

function toForm(user) {
  return {
    id: user.id,
    email: user.email || '',
    password: '',
    role: user.role || 'user',
    assigned_sessions: Array.isArray(user.assigned_sessions) ? user.assigned_sessions : [],
  }
}

function sessionLabel(session) {
  return session?.session_id || session?.id || ''
}

export function Settings({ currentUser, notify, activeSection = 'users', onSectionChange }) {
  const [users, setUsers] = useState([])
  const [sessionOptions, setSessionOptions] = useState([])
  const [globalWebhooks, setGlobalWebhooks] = useState([])
  const [loading, setLoading] = useState(true)
  const [webhooksLoading, setWebhooksLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [webhookSaving, setWebhookSaving] = useState(false)
  const [deleting, setDeleting] = useState('')
  const [webhookDeleting, setWebhookDeleting] = useState('')
  const [form, setForm] = useState(EMPTY_FORM)
  const [webhookForm, setWebhookForm] = useState(defaultGlobalWebhookForm)
  const [userModalOpen, setUserModalOpen] = useState(false)
  const [webhookModalOpen, setWebhookModalOpen] = useState(false)
  const [error, setError] = useState('')
  const [webhookError, setWebhookError] = useState('')

  const canManage = currentUser?.role === 'owner'
  const editing = Boolean(form.id)
  const editingWebhook = Boolean(webhookForm.id)

  const reload = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const [data, sessionData] = await Promise.all([api.listUsers(), api.listSessions()])
      setUsers(Array.isArray(data) ? data : [])
      setSessionOptions(Array.isArray(sessionData?.sessions) ? sessionData.sessions : [])
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }, [])

  const reloadGlobalWebhooks = useCallback(async () => {
    setWebhooksLoading(true)
    setWebhookError('')
    try {
      const data = await api.listGlobalWebhooks()
      setGlobalWebhooks(Array.isArray(data?.webhooks) ? data.webhooks : [])
    } catch (e) {
      setWebhookError(e.message)
    } finally {
      setWebhooksLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!canManage) return
    reload()
    reloadGlobalWebhooks()
  }, [canManage, reload, reloadGlobalWebhooks])

  const sortedUsers = useMemo(
    () => [...users].sort((a, b) => String(a.email).localeCompare(String(b.email))),
    [users],
  )

  const sortedSessionOptions = useMemo(
    () => [...sessionOptions].sort((a, b) => sessionLabel(a).localeCompare(sessionLabel(b))),
    [sessionOptions],
  )

  const resetForm = () => {
    setForm(EMPTY_FORM)
    setError('')
  }

  const closeUserModal = () => {
    setUserModalOpen(false)
    resetForm()
  }

  const openCreateUserModal = () => {
    resetForm()
    setUserModalOpen(true)
  }

  const openEditUserModal = (user) => {
    setForm(toForm(user))
    setError('')
    setUserModalOpen(true)
  }

  const resetWebhookForm = () => {
    setWebhookForm(defaultGlobalWebhookForm())
    setWebhookError('')
  }

  const closeWebhookModal = () => {
    setWebhookModalOpen(false)
    resetWebhookForm()
  }

  const openCreateWebhookModal = () => {
    resetWebhookForm()
    setWebhookModalOpen(true)
  }

  const openEditWebhookModal = (webhook) => {
    setWebhookForm(webhookFormFromItem(webhook))
    setWebhookError('')
    setWebhookModalOpen(true)
  }

  const handleSave = async (event) => {
    event.preventDefault()
    setSaving(true)
    setError('')
    try {
      const body = {
        email: form.email.trim(),
        role: form.role,
        password: form.password,
        assigned_sessions: form.role === 'owner' ? [] : form.assigned_sessions,
      }
      if (editing) await api.updateUser(form.id, body)
      else await api.createUser(body)
      notify(editing ? 'Usuário atualizado.' : 'Usuário criado.', 'success')
      closeUserModal()
      await reload()
    } catch (e) {
      setError(e.message)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (user) => {
    if (user.id === currentUser?.id) {
      setError('Você não pode apagar o usuário atual.')
      return
    }
    if (!window.confirm(`Apagar o usuário "${user.email}"?`)) return
    setDeleting(user.id)
    setError('')
    try {
      await api.deleteUser(user.id)
      notify('Usuário apagado.', 'success')
      if (form.id === user.id) closeUserModal()
      await reload()
    } catch (e) {
      setError(e.message)
    } finally {
      setDeleting('')
    }
  }

  const handleSaveWebhook = async (event) => {
    event.preventDefault()
    const { url, url_absolute: urlAbsolute, timeout_ms: timeoutMS } = webhookForm
    if (!url.trim() && !urlAbsolute.trim()) {
      setWebhookError('Informe a URL base ou a URL absoluta.')
      return
    }
    if (!(parseInt(timeoutMS, 10) > 0)) {
      setWebhookError('O timeout deve ser maior que zero.')
      return
    }
    setWebhookSaving(true)
    setWebhookError('')
    try {
      await api.saveGlobalWebhook(webhookFormToPayload(webhookForm))
      notify(editingWebhook ? 'Webhook global atualizado.' : 'Webhook global criado.', 'success')
      closeWebhookModal()
      await reloadGlobalWebhooks()
    } catch (e) {
      setWebhookError(e.message)
    } finally {
      setWebhookSaving(false)
    }
  }

  const handleDeleteWebhook = async (webhook) => {
    if (!window.confirm(`Excluir o webhook global "${webhook.name || webhook.id}"?`)) return
    setWebhookDeleting(webhook.id)
    setWebhookError('')
    try {
      await api.deleteGlobalWebhook(webhook.id)
      notify('Webhook global excluído.', 'success')
      if (webhookForm.id === webhook.id) closeWebhookModal()
      await reloadGlobalWebhooks()
    } catch (e) {
      setWebhookError(e.message)
    } finally {
      setWebhookDeleting('')
    }
  }

  const toggleAssignedSession = (sessionId) => {
    setForm((current) => {
      const assigned = new Set(current.assigned_sessions)
      if (assigned.has(sessionId)) assigned.delete(sessionId)
      else assigned.add(sessionId)
      return { ...current, assigned_sessions: [...assigned].sort() }
    })
  }

  if (!canManage) {
    return (
      <EmptyState
        title="Acesso restrito"
        description="Somente usuários owner podem gerenciar contas."
      />
    )
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[13rem_minmax(0,1fr)]">
      <aside className="lg:sticky lg:top-24 lg:self-start">
        <nav className="flex gap-1 overflow-x-auto border-b border-border pb-2 lg:flex-col lg:overflow-visible lg:border-b-0 lg:pb-0">
          {SECTIONS.map((section) => (
            <button
              key={section.id}
              type="button"
              onClick={() => onSectionChange?.(section.id)}
              className={cx(
                'whitespace-nowrap rounded-lg px-3 py-2 text-left text-sm font-medium transition-colors',
                activeSection === section.id
                  ? 'bg-surface-2 text-fg'
                  : 'text-muted hover:bg-surface-2 hover:text-fg',
              )}
            >
              {section.label}
            </button>
          ))}
        </nav>
      </aside>

      <section className="min-w-0">
        {activeSection === 'users' ? (
          <Card>
            <CardHeader
              eyebrow="Configurações"
              title="Usuários"
              description="Contas com acesso ao dashboard e às rotas autenticadas."
              actions={
                <>
                  <Button size="sm" variant="ghost" onClick={reload} loading={loading}>
                    Atualizar
                  </Button>
                  <Button size="sm" variant="primary" onClick={openCreateUserModal}>
                    Novo usuário
                  </Button>
                </>
              }
            />
            <div className="p-5">
              {error && (
                <p className="mb-4 rounded-lg bg-down-soft px-3 py-2 text-sm text-down">
                  {error}
                </p>
              )}
              {loading ? (
                <div className="flex justify-center py-10">
                  <Spinner />
                </div>
              ) : sortedUsers.length === 0 ? (
                <EmptyState
                  title="Nenhum usuário"
                  description="Cadastre o primeiro usuário."
                  action={
                    <Button size="sm" variant="primary" onClick={openCreateUserModal}>
                      Novo usuário
                    </Button>
                  }
                />
              ) : (
                <div className="overflow-x-auto rounded-[var(--radius-card)] border border-border">
                  <table className="w-full min-w-[42rem] border-collapse text-left text-sm">
                    <thead className="bg-surface-2 text-xs uppercase text-faint">
                      <tr>
                        <th className="px-4 py-3 font-medium">E-mail</th>
                        <th className="px-4 py-3 font-medium">Papel</th>
                        <th className="px-4 py-3 font-medium">Instâncias</th>
                        <th className="px-4 py-3 font-medium">Criado em</th>
                        <th className="px-4 py-3 text-right font-medium">Ações</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border">
                      {sortedUsers.map((user) => (
                        <tr key={user.id}>
                          <td className="min-w-0 px-4 py-3">
                            <p className="truncate font-medium text-fg">{user.email}</p>
                            {user.id === currentUser?.id && (
                              <p className="mt-0.5 text-xs text-faint">usuário atual</p>
                            )}
                          </td>
                          <td className="px-4 py-3">
                            <StatusBadge tone={user.role === 'owner' ? 'live' : user.role === 'admin' ? 'pending' : 'neutral'}>
                              {ROLE_LABEL[user.role] || user.role}
                            </StatusBadge>
                          </td>
                          <td className="px-4 py-3 text-muted">
                            {user.role === 'owner'
                              ? 'Todas'
                              : `${Array.isArray(user.assigned_sessions) ? user.assigned_sessions.length : 0} vinculada(s)`}
                          </td>
                          <td className="whitespace-nowrap px-4 py-3 text-muted">
                            {fmtDate(user.created_at)}
                          </td>
                          <td className="px-4 py-3">
                            <div className="flex justify-end gap-2">
                              <Button size="sm" onClick={() => openEditUserModal(user)}>
                                Editar
                              </Button>
                              <Button
                                size="sm"
                                variant="danger"
                                loading={deleting === user.id}
                                disabled={user.id === currentUser?.id}
                                onClick={() => handleDelete(user)}
                              >
                                Apagar
                              </Button>
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </Card>
        ) : (
          <Card>
            <CardHeader
              eyebrow="Configurações"
              title="Webhooks"
              description="Endpoints globais que recebem eventos administrativos de todas as sessões deste worker."
              actions={
                <>
                  <Button size="sm" variant="ghost" onClick={reloadGlobalWebhooks} loading={webhooksLoading}>
                    Atualizar
                  </Button>
                  <Button size="sm" variant="primary" onClick={openCreateWebhookModal}>
                    Novo webhook
                  </Button>
                </>
              }
            />
            <div className="p-5">
              {webhookError && (
                <p className="mb-4 rounded-lg bg-down-soft px-3 py-2 text-sm text-down">
                  {webhookError}
                </p>
              )}
              {webhooksLoading ? (
                <div className="flex justify-center py-10">
                  <Spinner />
                </div>
              ) : globalWebhooks.length === 0 ? (
                <EmptyState
                  title="Nenhum webhook global"
                  description="Cadastre o primeiro webhook."
                  action={
                    <Button size="sm" variant="primary" onClick={openCreateWebhookModal}>
                      Novo webhook
                    </Button>
                  }
                />
              ) : (
                <div className="grid gap-3">
                  {globalWebhooks.map((webhook) => (
                    <article key={webhook.id} className="rounded-lg border border-border bg-surface-2 p-3">
                      <div className="flex items-start justify-between gap-3">
                        <div className="min-w-0">
                          <p className="truncate text-sm font-medium text-fg">
                            {webhook.name || 'Sem nome'}
                          </p>
                          <Mono className="text-faint">{webhook.id}</Mono>
                        </div>
                        <div className="flex shrink-0 gap-1">
                          <Button size="sm" variant="ghost" onClick={() => openEditWebhookModal(webhook)}>
                            Editar
                          </Button>
                          <Button
                            size="sm"
                            variant="danger"
                            loading={webhookDeleting === webhook.id}
                            onClick={() => handleDeleteWebhook(webhook)}
                          >
                            Excluir
                          </Button>
                        </div>
                      </div>
                      <Mono className="mt-2 block truncate text-muted">
                        {webhook.url_absolute || webhook.url || '—'}
                      </Mono>
                      <div className="mt-2 flex flex-wrap gap-1.5">
                        <Tag>{webhook.enabled ? 'Ativo' : 'Desligado'}</Tag>
                        {webhook.send_session_disconnected && <Tag muted>Sessão desconectada</Tag>}
                        {webhook.send_session_offline && <Tag muted>Instância offline</Tag>}
                      </div>
                    </article>
                  ))}
                </div>
              )}
            </div>
          </Card>
        )}
      </section>

      <Modal
        open={userModalOpen}
        onClose={closeUserModal}
        eyebrow={editing ? 'Editar' : 'Novo'}
        title={editing ? 'Editar usuário' : 'Cadastrar usuário'}
        maxWidth="max-w-2xl"
      >
        <form onSubmit={handleSave} className="flex flex-col gap-4">
          {error && (
            <p className="rounded-lg bg-down-soft px-3 py-2 text-sm text-down">
              {error}
            </p>
          )}
          <Field label="E-mail" htmlFor="settings-user-email">
            <Input
              id="settings-user-email"
              type="email"
              autoComplete="email"
              required
              value={form.email}
              onChange={(e) => setForm((current) => ({ ...current, email: e.target.value }))}
            />
          </Field>
          <Field label="Papel" htmlFor="settings-user-role">
            <Select
              id="settings-user-role"
              value={form.role}
              onChange={(e) => setForm((current) => ({ ...current, role: e.target.value }))}
            >
              <option value="user">Usuário</option>
              <option value="admin">Admin</option>
              <option value="owner">Owner</option>
            </Select>
          </Field>
          {form.role !== 'owner' && (
            <Field label="Instâncias vinculadas">
              <div className="max-h-56 overflow-y-auto rounded-lg border border-border bg-surface-2 p-2">
                {sortedSessionOptions.length === 0 ? (
                  <p className="px-2 py-3 text-sm text-muted">Nenhuma instância encontrada.</p>
                ) : (
                  <div className="grid gap-1">
                    {sortedSessionOptions.map((session) => {
                      const id = sessionLabel(session)
                      if (!id) return null
                      return (
                        <label
                          key={id}
                          className="flex items-center gap-3 rounded-md px-2 py-2 text-sm text-fg transition-colors hover:bg-surface"
                        >
                          <input
                            type="checkbox"
                            className="size-4 accent-brand"
                            checked={form.assigned_sessions.includes(id)}
                            onChange={() => toggleAssignedSession(id)}
                          />
                          <span className="min-w-0 flex-1 truncate font-mono text-xs">{id}</span>
                          {session.assigned && <Tag muted>worker</Tag>}
                        </label>
                      )
                    })}
                  </div>
                )}
              </div>
            </Field>
          )}
          <Field label={editing ? 'Nova senha' : 'Senha'} htmlFor="settings-user-password">
            <Input
              id="settings-user-password"
              type="password"
              autoComplete="new-password"
              required={!editing}
              value={form.password}
              onChange={(e) => setForm((current) => ({ ...current, password: e.target.value }))}
            />
          </Field>
          <div className="flex justify-end gap-2 border-t border-border pt-4">
            <Button variant="ghost" onClick={closeUserModal}>
              Cancelar
            </Button>
            <Button type="submit" variant="primary" loading={saving}>
              {editing ? 'Salvar' : 'Cadastrar'}
            </Button>
          </div>
        </form>
      </Modal>

      <Modal
        open={webhookModalOpen}
        onClose={closeWebhookModal}
        eyebrow={editingWebhook ? 'Editar' : 'Novo'}
        title={editingWebhook ? 'Editar webhook' : 'Cadastrar webhook'}
        maxWidth="max-w-lg"
      >
        <form onSubmit={handleSaveWebhook} className="flex flex-col gap-4">
          {webhookError && (
            <p className="rounded-lg bg-down-soft px-3 py-2 text-sm text-down">
              {webhookError}
            </p>
          )}
          <Field label="Nome">
            <Input
              value={webhookForm.name}
              onChange={(e) => setWebhookForm((current) => ({ ...current, name: e.target.value }))}
            />
          </Field>
          <Field label="URL base">
            <Input
              value={webhookForm.url}
              placeholder="https://api.exemplo.com/webhook"
              onChange={(e) => setWebhookForm((current) => ({ ...current, url: e.target.value }))}
            />
          </Field>
          <Field label="URL absoluta" hint="Usada exatamente como está, sem anexar a sessão.">
            <Input
              value={webhookForm.url_absolute}
              placeholder="https://api.exemplo.com/hooks/global"
              onChange={(e) => setWebhookForm((current) => ({ ...current, url_absolute: e.target.value }))}
            />
          </Field>
          <div className="grid gap-3 sm:grid-cols-2">
            <Field label="Header de auth">
              <Input
                value={webhookForm.header}
                onChange={(e) => setWebhookForm((current) => ({ ...current, header: e.target.value }))}
              />
            </Field>
            <Field label="Token">
              <Input
                value={webhookForm.token}
                onChange={(e) => setWebhookForm((current) => ({ ...current, token: e.target.value }))}
              />
            </Field>
            <Field label="Timeout (ms)">
              <Input
                type="number"
                min="1"
                value={webhookForm.timeout_ms}
                onChange={(e) => setWebhookForm((current) => ({ ...current, timeout_ms: e.target.value }))}
              />
            </Field>
          </div>
          <div className="flex flex-col gap-2 rounded-lg border border-border p-3">
            <Switch
              label="Ativo"
              checked={webhookForm.enabled}
              onChange={(value) => setWebhookForm((current) => ({ ...current, enabled: value }))}
            />
            <Switch
              label="Sessão desconectada"
              description="Dispara em logout real, remoção do dispositivo, restrição ou banimento."
              checked={webhookForm.send_session_disconnected}
              onChange={(value) => setWebhookForm((current) => ({ ...current, send_session_disconnected: value }))}
            />
            <Switch
              label="Instância offline"
              description="Dispara quando a instância permanece offline por 1 minuto."
              checked={webhookForm.send_session_offline}
              onChange={(value) => setWebhookForm((current) => ({ ...current, send_session_offline: value }))}
            />
          </div>
          <div className="flex justify-end gap-2 border-t border-border pt-4">
            <Button variant="ghost" onClick={closeWebhookModal}>
              Cancelar
            </Button>
            <Button type="submit" variant="primary" loading={webhookSaving}>
              {editingWebhook ? 'Salvar' : 'Cadastrar'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
