const workerIdEl = document.getElementById('worker-id')
const assignedCountEl = document.getElementById('assigned-count')
const onlineCountEl = document.getElementById('online-count')
const queueNameEl = document.getElementById('queue-name')
const pageTitleEl = document.getElementById('page-title')
const sessionsBodyEl = document.getElementById('sessions-body')
const lastUpdateEl = document.getElementById('last-update')
const toastEl = document.getElementById('toast')

const sessionForm = document.getElementById('session-form')
const messageForm = document.getElementById('message-form')
const sessionConfigForm = document.getElementById('session-config-form')
const configSessionInput = document.getElementById('config-session')
const configLoadBtn = document.getElementById('config-load')
const configClearBtn = document.getElementById('config-clear')

const cfgRejectCalls = document.getElementById('cfg-reject-calls')
const cfgRejectMessage = document.getElementById('cfg-reject-message')
const cfgMarkRead = document.getElementById('cfg-mark-read')
const cfgIgnoreStatus = document.getElementById('cfg-ignore-status')
const cfgIgnoreNewsletters = document.getElementById('cfg-ignore-newsletters')
const cfgOutgoingCoex = document.getElementById('cfg-outgoing-coex')
const cfgAlwaysOnline = document.getElementById('cfg-always-online')
const cfgAlwaysOnlineInterval = document.getElementById('cfg-always-online-interval')

const webhookForm = document.getElementById('webhook-form')
const webhookSessionInput = document.getElementById('webhook-session')
const webhookListBody = document.getElementById('webhook-list-body')
const webhookLoadBtn = document.getElementById('webhook-load')
const webhookClearFormBtn = document.getElementById('webhook-clear-form')
const webhookDeleteCurrentBtn = document.getElementById('webhook-delete-current')

const webhookIdInput = document.getElementById('webhook-id')
const webhookNameInput = document.getElementById('webhook-name')
const webhookUrlInput = document.getElementById('webhook-url')
const webhookUrlAbsoluteInput = document.getElementById('webhook-url-absolute')
const webhookHeaderInput = document.getElementById('webhook-header')
const webhookTokenInput = document.getElementById('webhook-token')
const webhookTimeoutInput = document.getElementById('webhook-timeout')

const webhookEnabledInput = document.getElementById('webhook-enabled')
const webhookAppendSessionInput = document.getElementById('webhook-append-session')
const webhookSendIncomingInput = document.getElementById('webhook-send-incoming')
const webhookSendOutgoingInput = document.getElementById('webhook-send-outgoing')
const webhookSendUpdateInput = document.getElementById('webhook-send-update')
const webhookSendGroupInput = document.getElementById('webhook-send-group')
const webhookSendNewsletterInput = document.getElementById('webhook-send-newsletter')

const qrDialog = document.getElementById('qr-dialog')
const qrImage = document.getElementById('qr-image')
const qrTitle = document.getElementById('qr-title')
const qrClose = document.getElementById('qr-close')

const refreshBtn = document.getElementById('refresh-btn')
const toggleRefreshBtn = document.getElementById('toggle-auto-refresh')

let refreshTimer = null
let autoRefresh = true

async function api(path, opts = {}) {
  const res = await fetch(path, {
    headers: { 'Content-Type': 'application/json' },
    ...opts,
  })
  if (!res.ok) {
    let msg = `${res.status} ${res.statusText}`
    try {
      const text = await res.text()
      if (text) {
        msg = text
      }
    } catch (err) {
      // ignore parse issues
    }
    throw new Error(msg)
  }

  const contentType = res.headers.get('content-type') || ''
  if (contentType.includes('application/json')) {
    return res.json()
  }
  return null
}

function toast(message) {
  toastEl.textContent = message
  toastEl.classList.add('show')
  setTimeout(() => toastEl.classList.remove('show'), 2200)
}

function statusBadge(status) {
  const value = (status || 'offline').toLowerCase()
  let cls = 'badge-offline'
  if (value === 'online' || value === 'connected') {
    cls = 'badge-online'
  } else if (value === 'connecting') {
    cls = 'badge-connecting'
  }
  return `<span class="badge ${cls}">${value}</span>`
}

function button(label, action, sessionId) {
  return `<button class="btn" data-action="${action}" data-session="${sessionId}">${label}</button>`
}

function renderSessions(sessions) {
  sessionsBodyEl.innerHTML = ''

  if (!Array.isArray(sessions) || sessions.length === 0) {
    sessionsBodyEl.innerHTML = '<tr><td colspan="6">Nenhuma sessão encontrada.</td></tr>'
    return
  }

  const rows = sessions
    .map((s) => {
      const assigned = s.assigned ? 'sim' : 'não'
      const connected = s.connected ? 'sim' : 'não'
      const device = s.device_jid || '-'
      const actions = [
        button('Conectar', 'connect', s.session_id),
        button('Recarregar', 'reload', s.session_id),
        button('Desconectar', 'disconnect', s.session_id),
        button('QR', 'qr', s.session_id),
        button('Config', 'config', s.session_id),
        button('Webhooks', 'webhooks', s.session_id),
      ]
      if (s.assigned) {
        actions.push(button('Desatribuir', 'unassign', s.session_id))
      } else {
        actions.push(button('Atribuir', 'assign', s.session_id))
      }
      const actionsHtml = actions.join('')

      return `
        <tr>
          <td>${s.session_id}</td>
          <td>${statusBadge(s.status)}</td>
          <td>${connected}</td>
          <td>${device}</td>
          <td>${assigned}</td>
          <td><div class="actions">${actionsHtml}</div></td>
        </tr>
      `
    })
    .join('')

  sessionsBodyEl.innerHTML = rows
}

function bindTableActions() {
  sessionsBodyEl.querySelectorAll('button[data-action]').forEach((btn) => {
    btn.addEventListener('click', async () => {
      const action = btn.dataset.action
      const sessionId = btn.dataset.session
      if (!action || !sessionId) {
        return
      }

      try {
        if (action === 'qr') {
          openQR(sessionId)
          return
        }
        if (action === 'config') {
          await openSessionConfig(sessionId)
          return
        }
        if (action === 'webhooks') {
          await openSessionWebhooks(sessionId)
          return
        }

        if (action === 'assign') {
          await api('/api/v1/worker/assignments', {
            method: 'POST',
            body: JSON.stringify({ session_id: sessionId }),
          })
        } else if (action === 'unassign') {
          await api(`/api/v1/worker/assignments/${encodeURIComponent(sessionId)}`, {
            method: 'DELETE',
          })
        } else {
          await api(`/sessions/${encodeURIComponent(sessionId)}/${action}`, {
            method: 'POST',
          })
        }

        toast(`Ação ${action} enviada para ${sessionId}`)
        await loadAll()
      } catch (err) {
        toast(`Erro: ${err.message}`)
      }
    })
  })
}

function openQR(sessionId) {
  qrTitle.textContent = `QR da sessão ${sessionId}`
  qrImage.src = `/sessions/${encodeURIComponent(sessionId)}/qr?ts=${Date.now()}`
  if (typeof qrDialog.showModal === 'function') {
    qrDialog.showModal()
  }
}

function applySessionConfig(cfg = {}) {
  cfgRejectCalls.checked = Boolean(cfg.reject_calls)
  cfgRejectMessage.value = cfg.reject_calls_message || ''
  cfgMarkRead.checked = Boolean(cfg.mark_read_on_message)
  cfgIgnoreStatus.checked = Boolean(cfg.ignore_status_broadcast)
  cfgIgnoreNewsletters.checked = Boolean(cfg.ignore_newsletters)
  cfgOutgoingCoex.checked = Boolean(cfg.outgoing_messages_coex)
  cfgAlwaysOnline.checked = Boolean(cfg.always_online)
  const interval = Number(cfg.always_online_interval_seconds || 60)
  cfgAlwaysOnlineInterval.value = String(Number.isFinite(interval) && interval > 0 ? interval : 60)
}

async function loadSessionConfig(sessionId) {
  const sid = (sessionId || '').trim()
  if (!sid) {
    toast('Informe a sessão para carregar configuração')
    return
  }
  const payload = await api(`/api/v1/sessions/${encodeURIComponent(sid)}/config`)
  configSessionInput.value = payload.session_id || sid
  applySessionConfig(payload.config || {})
}

async function openSessionConfig(sessionId) {
  configSessionInput.value = sessionId
  await loadSessionConfig(sessionId)
  sessionConfigForm.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

function resetWebhookForm(keepSession = true) {
  const currentSession = keepSession ? webhookSessionInput.value.trim() : ''
  webhookForm.reset()
  webhookSessionInput.value = currentSession
  webhookEnabledInput.checked = true
  webhookAppendSessionInput.checked = true
  webhookSendIncomingInput.checked = true
  webhookSendOutgoingInput.checked = true
  webhookSendUpdateInput.checked = true
  webhookSendGroupInput.checked = true
  webhookSendNewsletterInput.checked = false
  webhookTimeoutInput.value = '5000'
}

function webhookPayloadFromForm() {
  const timeoutMs = parseInt(webhookTimeoutInput.value, 10)
  return {
    id: webhookIdInput.value.trim(),
    name: webhookNameInput.value.trim(),
    url: webhookUrlInput.value.trim(),
    url_absolute: webhookUrlAbsoluteInput.value.trim(),
    header: webhookHeaderInput.value.trim(),
    token: webhookTokenInput.value.trim(),
    timeout_ms: Number.isFinite(timeoutMs) ? timeoutMs : 5000,
    enabled: webhookEnabledInput.checked,
    append_session: webhookAppendSessionInput.checked,
    send_incoming_messages: webhookSendIncomingInput.checked,
    send_outgoing_messages: webhookSendOutgoingInput.checked,
    send_update_messages: webhookSendUpdateInput.checked,
    send_group_messages: webhookSendGroupInput.checked,
    send_newsletter_messages: webhookSendNewsletterInput.checked,
  }
}

function fillWebhookForm(webhook) {
  webhookIdInput.value = webhook.id || ''
  webhookNameInput.value = webhook.name || ''
  webhookUrlInput.value = webhook.url || ''
  webhookUrlAbsoluteInput.value = webhook.url_absolute || ''
  webhookHeaderInput.value = webhook.header || ''
  webhookTokenInput.value = webhook.token || ''
  webhookTimeoutInput.value = String(webhook.timeout_ms || 5000)
  webhookEnabledInput.checked = Boolean(webhook.enabled)
  webhookAppendSessionInput.checked = Boolean(webhook.append_session)
  webhookSendIncomingInput.checked = Boolean(webhook.send_incoming_messages)
  webhookSendOutgoingInput.checked = Boolean(webhook.send_outgoing_messages)
  webhookSendUpdateInput.checked = Boolean(webhook.send_update_messages)
  webhookSendGroupInput.checked = Boolean(webhook.send_group_messages)
  webhookSendNewsletterInput.checked = Boolean(webhook.send_newsletter_messages)
}

function renderSessionWebhooks(sessionId, webhooks) {
  webhookListBody.innerHTML = ''
  if (!Array.isArray(webhooks) || webhooks.length === 0) {
    webhookListBody.innerHTML = '<tr><td colspan="6">Nenhum webhook cadastrado.</td></tr>'
    return
  }
  const rows = webhooks
    .map((w) => {
      const displayUrl = w.url_absolute || w.url || '-'
      const incoming = w.send_incoming_messages ? 'sim' : 'não'
      const outgoing = w.send_outgoing_messages ? 'sim' : 'não'
      const update = w.send_update_messages ? 'sim' : 'não'
      return `
        <tr>
          <td>${w.id}</td>
          <td>${displayUrl}</td>
          <td>${incoming}</td>
          <td>${outgoing}</td>
          <td>${update}</td>
          <td>
            <div class="actions">
              ${button('Editar', 'webhook-edit', w.id)}
              ${button('Excluir', 'webhook-delete', w.id)}
            </div>
          </td>
        </tr>
      `
    })
    .join('')
  webhookListBody.innerHTML = rows

  webhookListBody.querySelectorAll('button[data-action]').forEach((btn) => {
    btn.addEventListener('click', async () => {
      const action = btn.dataset.action
      const webhookId = btn.dataset.session
      if (!action || !webhookId) {
        return
      }
      const target = webhooks.find((w) => w.id === webhookId)
      try {
        if (action === 'webhook-edit') {
          if (target) {
            fillWebhookForm(target)
          }
          return
        }
        if (action === 'webhook-delete') {
          await api(`/api/v1/sessions/${encodeURIComponent(sessionId)}/webhooks/${encodeURIComponent(webhookId)}`, {
            method: 'DELETE',
          })
          if (webhookIdInput.value.trim() === webhookId) {
            resetWebhookForm(true)
          }
          await loadSessionWebhooks(sessionId)
          toast(`Webhook ${webhookId} removido`)
        }
      } catch (err) {
        toast(`Erro: ${err.message}`)
      }
    })
  })
}

async function loadSessionWebhooks(sessionId) {
  const sid = (sessionId || '').trim()
  if (!sid) {
    toast('Informe a sessão para carregar webhooks')
    return
  }
  const payload = await api(`/api/v1/sessions/${encodeURIComponent(sid)}/webhooks`)
  webhookSessionInput.value = payload.session_id || sid
  renderSessionWebhooks(sid, payload.webhooks || [])
}

async function openSessionWebhooks(sessionId) {
  webhookSessionInput.value = sessionId
  await loadSessionWebhooks(sessionId)
  webhookForm.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

async function loadWorker() {
  const worker = await api('/api/v1/worker')
  pageTitleEl.textContent = worker.ui_title || 'WhatsMeow Control Center'
  workerIdEl.textContent = worker.worker_id || '-'
  queueNameEl.textContent = worker.amqp_queue || '-'
  assignedCountEl.textContent = (worker.assigned_sessions || []).length
}

async function loadSessions() {
  const payload = await api('/api/v1/sessions')
  const sessions = payload.sessions || []
  const online = sessions.filter((s) => s.connected || (s.status || '').toLowerCase() === 'online').length
  onlineCountEl.textContent = online
  renderSessions(sessions)
  bindTableActions()
  lastUpdateEl.textContent = `Atualizado às ${new Date().toLocaleTimeString('pt-BR')}`
}

async function loadAll() {
  await Promise.all([loadWorker(), loadSessions()])
}

function setAutoRefresh(enabled) {
  autoRefresh = enabled
  toggleRefreshBtn.textContent = `Auto-refresh: ${enabled ? 'on' : 'off'}`

  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }

  if (enabled) {
    refreshTimer = setInterval(() => {
      loadAll().catch((err) => toast(`Erro: ${err.message}`))
    }, 7000)
  }
}

sessionForm.addEventListener('submit', async (ev) => {
  ev.preventDefault()
  const input = document.getElementById('session-id')
  const assign = document.getElementById('session-assign').checked
  const connect = document.getElementById('session-connect').checked
  const sessionId = input.value.trim()
  if (!sessionId) {
    return
  }

  try {
    await api('/api/v1/sessions', {
      method: 'POST',
      body: JSON.stringify({ id: sessionId, connect, assign }),
    })
    toast(`Sessão ${sessionId} solicitada`)
    input.value = ''
    await loadAll()
  } catch (err) {
    toast(`Erro: ${err.message}`)
  }
})

messageForm.addEventListener('submit', async (ev) => {
  ev.preventDefault()
  const session = document.getElementById('message-session').value.trim()
  const to = document.getElementById('message-to').value.trim()
  const body = document.getElementById('message-body').value.trim()

  if (!session || !to || !body) {
    toast('Preencha sessão, destino e texto')
    return
  }

  try {
    await api(`/api/v1/sessions/${encodeURIComponent(session)}/messages`, {
      method: 'POST',
      body: JSON.stringify({
        type: 'text',
        to,
        text: { body },
      }),
    })
    toast('Mensagem enviada para fila de envio')
    document.getElementById('message-body').value = ''
  } catch (err) {
    toast(`Erro: ${err.message}`)
  }
})

sessionConfigForm.addEventListener('submit', async (ev) => {
  ev.preventDefault()
  const sessionId = configSessionInput.value.trim()
  if (!sessionId) {
    toast('Informe a sessão para salvar configuração')
    return
  }
  const interval = parseInt(cfgAlwaysOnlineInterval.value, 10)
  if (!Number.isFinite(interval) || interval <= 0) {
    toast('Intervalo de always-online deve ser maior que zero')
    return
  }

  try {
    await api(`/api/v1/sessions/${encodeURIComponent(sessionId)}/config`, {
      method: 'PUT',
      body: JSON.stringify({
        reject_calls: cfgRejectCalls.checked,
        reject_calls_message: cfgRejectMessage.value.trim(),
        mark_read_on_message: cfgMarkRead.checked,
        ignore_status_broadcast: cfgIgnoreStatus.checked,
        ignore_newsletters: cfgIgnoreNewsletters.checked,
        outgoing_messages_coex: cfgOutgoingCoex.checked,
        always_online: cfgAlwaysOnline.checked,
        always_online_interval_seconds: interval,
      }),
    })
    toast(`Configuração salva para ${sessionId}`)
    await loadAll()
  } catch (err) {
    toast(`Erro: ${err.message}`)
  }
})

configLoadBtn.addEventListener('click', async () => {
  try {
    await loadSessionConfig(configSessionInput.value)
    toast('Configuração carregada')
  } catch (err) {
    toast(`Erro: ${err.message}`)
  }
})

configClearBtn.addEventListener('click', async () => {
  const sessionId = configSessionInput.value.trim()
  if (!sessionId) {
    toast('Informe a sessão para limpar override')
    return
  }
  try {
    await api(`/api/v1/sessions/${encodeURIComponent(sessionId)}/config`, {
      method: 'DELETE',
    })
    await loadSessionConfig(sessionId)
    toast(`Override removido de ${sessionId}`)
  } catch (err) {
    toast(`Erro: ${err.message}`)
  }
})

webhookForm.addEventListener('submit', async (ev) => {
  ev.preventDefault()
  const sessionId = webhookSessionInput.value.trim()
  if (!sessionId) {
    toast('Informe a sessão para salvar webhook')
    return
  }
  const payload = webhookPayloadFromForm()
  if (!payload.url && !payload.url_absolute) {
    toast('Informe URL base ou URL absoluta para o webhook')
    return
  }

  try {
    const created = await api(`/api/v1/sessions/${encodeURIComponent(sessionId)}/webhooks`, {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    if (created && created.webhook && created.webhook.id) {
      webhookIdInput.value = created.webhook.id
    }
    await loadSessionWebhooks(sessionId)
    toast(`Webhook salvo para ${sessionId}`)
  } catch (err) {
    toast(`Erro: ${err.message}`)
  }
})

webhookLoadBtn.addEventListener('click', async () => {
  try {
    await loadSessionWebhooks(webhookSessionInput.value)
    toast('Lista de webhooks carregada')
  } catch (err) {
    toast(`Erro: ${err.message}`)
  }
})

webhookClearFormBtn.addEventListener('click', () => {
  resetWebhookForm(true)
})

webhookDeleteCurrentBtn.addEventListener('click', async () => {
  const sessionId = webhookSessionInput.value.trim()
  const webhookId = webhookIdInput.value.trim()
  if (!sessionId || !webhookId) {
    toast('Informe sessão e webhook ID para excluir')
    return
  }
  try {
    await api(`/api/v1/sessions/${encodeURIComponent(sessionId)}/webhooks/${encodeURIComponent(webhookId)}`, {
      method: 'DELETE',
    })
    resetWebhookForm(true)
    await loadSessionWebhooks(sessionId)
    toast(`Webhook ${webhookId} removido`)
  } catch (err) {
    toast(`Erro: ${err.message}`)
  }
})

qrClose.addEventListener('click', () => qrDialog.close())

refreshBtn.addEventListener('click', () => {
  loadAll().catch((err) => toast(`Erro: ${err.message}`))
})

toggleRefreshBtn.addEventListener('click', () => {
  setAutoRefresh(!autoRefresh)
})

window.addEventListener('DOMContentLoaded', async () => {
  try {
    resetWebhookForm(false)
    await loadAll()
    setAutoRefresh(true)
  } catch (err) {
    toast(`Erro: ${err.message}`)
  }
})
