import { useEffect, useState } from 'react'
import {
  Button,
  Card,
  CardHeader,
  Field,
  Input,
  Mono,
  SignalDot,
  StatusBadge,
  Tabs,
  Eyebrow,
} from './ui.jsx'
import { ConfigForm } from './ConfigForm.jsx'
import { WebhookPanel } from './WebhookPanel.jsx'
import { TemplatePanel } from './TemplatePanel.jsx'
import { TestMessageModal } from './modals/TestMessageModal.jsx'
import { SessionWebhookModal } from './modals/SessionWebhookModal.jsx'
import { TemplateModal } from './modals/TemplateModal.jsx'
import * as api from '../api/endpoints.js'
import {
  configFormFromSnapshot,
  configFormToPayload,
  defaultConfigForm,
  defaultTemplateForm,
  defaultWebhookForm,
  templateFormFromItem,
  templateFormToPayload,
  webhookFormFromItem,
  webhookFormToPayload,
} from '../lib/forms.js'
import {
  digitsOnly,
  sessionStatus,
  statusLabel,
  statusTone,
} from '../lib/sessions.js'

const TABS = [
  { id: 'overview', label: 'Visão geral' },
  { id: 'config', label: 'Configuração' },
  { id: 'webhooks', label: 'Webhooks' },
  { id: 'templates', label: 'Templates' },
]

function DetailRow({ label, children }) {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-border py-2 last:border-0">
      <span className="text-xs text-muted">{label}</span>
      <span className="min-w-0 truncate text-right text-sm text-fg">{children}</span>
    </div>
  )
}

export function SessionDetail({
  session,
  worker,
  initialTab,
  notify,
  onAction,
  onShowQR,
  onTabChange,
}) {
  const sessionId = session.session_id
  const status = sessionStatus(session)
  const tone = statusTone(status)

  const [tab, setTab] = useState(initialTab || 'overview')

  const [config, setConfig] = useState({
    loading: true,
    saving: false,
    snapshot: null,
    form: defaultConfigForm(),
  })
  const [webhooks, setWebhooks] = useState({
    loading: true,
    saving: false,
    items: [],
    form: defaultWebhookForm(),
  })
  const [templates, setTemplates] = useState({
    loading: true,
    saving: false,
    items: [],
    form: defaultTemplateForm(),
  })
  const [pairing, setPairing] = useState({ phone: '', code: '', loading: false })
  const [messageModalOpen, setMessageModalOpen] = useState(false)
  const [webhookModalOpen, setWebhookModalOpen] = useState(false)
  const [templateModalOpen, setTemplateModalOpen] = useState(false)

  useEffect(() => {
    if (initialTab) setTab(initialTab)
  }, [initialTab, sessionId])

  // Load all per-session data when the selected session changes.
  useEffect(() => {
    let active = true
    setConfig((c) => ({ ...c, loading: true }))
    setWebhooks((w) => ({ ...w, loading: true }))
    setTemplates((t) => ({ ...t, loading: true }))
    setPairing({ phone: digitsOnly(sessionId), code: '', loading: false })
    setMessageModalOpen(false)
    setWebhookModalOpen(false)
    setTemplateModalOpen(false)

    Promise.all([
      api.getConfig(sessionId),
      api.listWebhooks(sessionId),
      api.listTemplates(sessionId),
    ])
      .then(([snapshot, webhookPayload, templatePayload]) => {
        if (!active) return
        setConfig({
          loading: false,
          saving: false,
          snapshot,
          form: configFormFromSnapshot(snapshot),
        })
        setWebhooks({
          loading: false,
          saving: false,
          items: webhookPayload?.webhooks || [],
          form: defaultWebhookForm(),
        })
        setTemplates({
          loading: false,
          saving: false,
          items: templatePayload?.templates || [],
          form: defaultTemplateForm(),
        })
      })
      .catch((error) => {
        if (!active) return
        setConfig((c) => ({ ...c, loading: false }))
        setWebhooks((w) => ({ ...w, loading: false }))
        setTemplates((t) => ({ ...t, loading: false }))
        notify(error.message, 'error')
      })

    return () => {
      active = false
    }
  }, [sessionId, notify])

  async function reloadConfig() {
    const snapshot = await api.getConfig(sessionId)
    setConfig((c) => ({ ...c, snapshot, form: configFormFromSnapshot(snapshot) }))
  }
  async function reloadWebhooks() {
    const payload = await api.listWebhooks(sessionId)
    setWebhooks((w) => ({ ...w, items: payload?.webhooks || [], form: defaultWebhookForm() }))
  }
  async function reloadTemplates() {
    const payload = await api.listTemplates(sessionId)
    setTemplates((t) => ({ ...t, items: payload?.templates || [], form: defaultTemplateForm() }))
  }

  async function handleSaveConfig(event) {
    event.preventDefault()
    setConfig((c) => ({ ...c, saving: true }))
    try {
      await api.updateConfig(sessionId, configFormToPayload(config.form))
      notify('Configuração salva.', 'success')
      await reloadConfig()
    } catch (error) {
      notify(error.message, 'error')
    } finally {
      setConfig((c) => ({ ...c, saving: false }))
    }
  }

  async function handleClearConfig() {
    if (!window.confirm(`Restaurar a configuração padrão de "${sessionId}"?`)) return
    setConfig((c) => ({ ...c, saving: true }))
    try {
      await api.clearConfig(sessionId)
      notify('Configuração restaurada para o padrão.', 'success')
      await reloadConfig()
    } catch (error) {
      notify(error.message, 'error')
    } finally {
      setConfig((c) => ({ ...c, saving: false }))
    }
  }

  async function handleSaveWebhook(event) {
    event.preventDefault()
    const { url, url_absolute, timeout_ms } = webhooks.form
    if (!url.trim() && !url_absolute.trim()) {
      notify('Informe a URL base ou a URL absoluta.', 'error')
      return
    }
    if (!(parseInt(timeout_ms, 10) > 0)) {
      notify('O timeout deve ser maior que zero.', 'error')
      return
    }
    setWebhooks((w) => ({ ...w, saving: true }))
    try {
      await api.saveWebhook(sessionId, webhookFormToPayload(webhooks.form))
      notify('Webhook salvo.', 'success')
      await reloadWebhooks()
      setWebhookModalOpen(false)
    } catch (error) {
      notify(error.message, 'error')
    } finally {
      setWebhooks((w) => ({ ...w, saving: false }))
    }
  }

  async function handleDeleteWebhook(webhookId) {
    if (!window.confirm(`Excluir o webhook "${webhookId}"?`)) return
    try {
      await api.deleteWebhook(sessionId, webhookId)
      notify('Webhook excluído.', 'success')
      await reloadWebhooks()
    } catch (error) {
      notify(error.message, 'error')
    }
  }

  async function handleSaveTemplate(event) {
    event.preventDefault()
    if (!templates.form.name.trim() || !templates.form.body.trim()) {
      notify('Nome e corpo do template são obrigatórios.', 'error')
      return
    }
    setTemplates((t) => ({ ...t, saving: true }))
    try {
      await api.saveTemplate(sessionId, templateFormToPayload(templates.form))
      notify('Template salvo.', 'success')
      await reloadTemplates()
      setTemplateModalOpen(false)
    } catch (error) {
      notify(error.message, 'error')
    } finally {
      setTemplates((t) => ({ ...t, saving: false }))
    }
  }

  async function handleDeleteTemplate(templateId) {
    if (!window.confirm(`Excluir o template "${templateId}"?`)) return
    try {
      await api.deleteTemplate(sessionId, templateId)
      notify('Template excluído.', 'success')
      await reloadTemplates()
    } catch (error) {
      notify(error.message, 'error')
    }
  }

  async function handlePairing(event) {
    event.preventDefault()
    const phone = digitsOnly(pairing.phone)
    if (!phone) {
      notify('Informe o telefone para gerar o código.', 'error')
      return
    }
    setPairing((p) => ({ ...p, phone, code: '', loading: true }))
    try {
      const payload = await api.requestPairingCode(sessionId, phone)
      const code = String(payload?.pairing_code || '').trim()
      if (!code) throw new Error('O servidor não retornou um código.')
      setPairing((p) => ({ ...p, code, loading: false }))
      notify('Código de pareamento gerado.', 'success')
    } catch (error) {
      setPairing((p) => ({ ...p, loading: false }))
      notify(error.message, 'error')
    }
  }

  function openCreateWebhookModal() {
    setWebhooks((w) => ({ ...w, form: defaultWebhookForm() }))
    setWebhookModalOpen(true)
  }

  function openEditWebhookModal(item) {
    setWebhooks((w) => ({ ...w, form: webhookFormFromItem(item) }))
    setWebhookModalOpen(true)
  }

  function openCreateTemplateModal() {
    setTemplates((t) => ({ ...t, form: defaultTemplateForm() }))
    setTemplateModalOpen(true)
  }

  function openEditTemplateModal(item) {
    setTemplates((t) => ({ ...t, form: templateFormFromItem(item) }))
    setTemplateModalOpen(true)
  }

  return (
    <Card className="flex flex-col">
      <CardHeader
        eyebrow="Instância"
        title={sessionId}
        description={`${statusLabel(status)} · ${session.assigned ? worker?.worker_id || 'atribuída' : 'sem atribuição'}`}
        actions={
          <>
            <StatusBadge tone={tone} pulse={tone === 'live'}>
              {statusLabel(status)}
            </StatusBadge>
            {session.has_qr && (
              <Button size="sm" onClick={() => onShowQR(sessionId)}>
                Ver QR
              </Button>
            )}
          </>
        }
      />

      <div className="px-5 pt-4">
        <Tabs
          tabs={TABS}
          active={tab}
          onChange={(nextTab) => {
            setTab(nextTab)
            onTabChange?.(nextTab)
          }}
        />
      </div>

      <div className="p-5">
        {tab === 'overview' && (
          <div className="grid gap-6 lg:grid-cols-2">
            <div className="flex flex-col gap-5">
              <div>
                <Eyebrow>Conexão</Eyebrow>
                <div className="mt-2">
                  <DetailRow label="Status">
                    <span className="inline-flex items-center gap-2">
                      <SignalDot tone={tone} pulse={tone === 'live'} />
                      {statusLabel(status)}
                    </span>
                  </DetailRow>
                  <DetailRow label="Login">
                    {session.logged_in ? 'ativo' : 'pendente'}
                  </DetailRow>
                  <DetailRow label="Device JID">
                    <Mono>{session.device_jid || '—'}</Mono>
                  </DetailRow>
                  <DetailRow label="WhatsApp Web">{session.whatsapp_web_version || '—'}</DetailRow>
                  <DetailRow label="Worker">
                    {session.assigned ? worker?.worker_id || 'atribuída' : 'livre'}
                  </DetailRow>
                </div>
                <div className="mt-3 flex flex-wrap gap-2">
                  <Button
                    size="sm"
                    variant="primary"
                    onClick={() => onAction('connect', sessionId)}
                  >
                    Conectar
                  </Button>
                  <Button size="sm" onClick={() => onAction('reload', sessionId)}>
                    Recarregar
                  </Button>
                  <Button size="sm" variant="danger" onClick={() => onAction('disconnect', sessionId)}>
                    Desconectar
                  </Button>
                </div>
              </div>

              <form onSubmit={handlePairing} className="flex flex-col gap-3 border-t border-border pt-5">
                <Eyebrow>Parear por código</Eyebrow>
                <Field label="Telefone" hint="Somente dígitos, com DDI.">
                  <Input
                    value={pairing.phone}
                    placeholder="5511999999999"
                    onChange={(e) => setPairing((p) => ({ ...p, phone: e.target.value }))}
                  />
                </Field>
                <Button type="submit" loading={pairing.loading}>
                  Gerar código
                </Button>
                {pairing.code && (
                  <div className="rounded-lg border border-border bg-surface-2 p-3 text-center">
                    <Eyebrow>Código</Eyebrow>
                    <p className="mt-1 font-mono text-2xl font-medium tracking-[0.3em] text-fg">
                      {pairing.code}
                    </p>
                  </div>
                )}
              </form>
            </div>

            <div className="border-t border-border pt-5 lg:border-0 lg:border-l lg:pl-6 lg:pt-0">
              <Eyebrow>Teste</Eyebrow>
              <div className="mt-2 rounded-lg border border-border bg-surface-2 p-4">
                <p className="text-sm text-muted">
                  Envie uma mensagem rápida por esta instância para validar conexão e entrega.
                </p>
                <Button
                  className="mt-3"
                  variant="primary"
                  onClick={() => setMessageModalOpen(true)}
                >
                  Enviar mensagem de teste
                </Button>
              </div>
            </div>
          </div>
        )}

        {tab === 'config' &&
          (config.loading ? (
            <p className="py-8 text-center text-muted">Carregando configuração…</p>
          ) : (
            <ConfigForm
              form={config.form}
              snapshot={config.snapshot}
              saving={config.saving}
              onChange={(form) => setConfig((c) => ({ ...c, form }))}
              onSave={handleSaveConfig}
              onClear={handleClearConfig}
            />
          ))}

        {tab === 'webhooks' && (
          <WebhookPanel
            items={webhooks.items}
            loading={webhooks.loading}
            onDelete={handleDeleteWebhook}
            onEdit={openEditWebhookModal}
            onCreate={openCreateWebhookModal}
          />
        )}

        {tab === 'templates' && (
          <TemplatePanel
            items={templates.items}
            loading={templates.loading}
            onDelete={handleDeleteTemplate}
            onEdit={openEditTemplateModal}
            onCreate={openCreateTemplateModal}
          />
        )}
      </div>

      <TestMessageModal
        open={messageModalOpen}
        onClose={() => setMessageModalOpen(false)}
        sessionId={sessionId}
        notify={notify}
      />
      <SessionWebhookModal
        open={webhookModalOpen}
        onClose={() => setWebhookModalOpen(false)}
        form={webhooks.form}
        saving={webhooks.saving}
        onChange={(form) => setWebhooks((w) => ({ ...w, form }))}
        onSave={handleSaveWebhook}
      />
      <TemplateModal
        open={templateModalOpen}
        onClose={() => setTemplateModalOpen(false)}
        form={templates.form}
        saving={templates.saving}
        onChange={(form) => setTemplates((t) => ({ ...t, form }))}
        onSave={handleSaveTemplate}
      />
    </Card>
  )
}
