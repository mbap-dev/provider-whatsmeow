// Form factories + payload mappers. These mirror the API contract field names
// exactly (snake_case) so payloads can be sent without remapping.

export function defaultConfigForm() {
  return {
    reject_calls: false,
    reject_calls_message: '',
    mark_read_on_message: false,
    ignore_status_broadcast: false,
    ignore_newsletters: false,
    outgoing_messages_coex: false,
    history_sync_enabled: false,
    history_sync_days: '7',
    always_online: false,
    always_online_interval_seconds: '60',
    proxy_enabled: false,
    proxy_protocol: '',
    proxy_host: '',
    proxy_port: '',
    proxy_username: '',
    proxy_password: '',
  }
}

export function configFormFromSnapshot(snapshot) {
  const cfg = snapshot?.config || {}
  const positiveOr = (value, fallback) =>
    Number.isFinite(Number(value)) && Number(value) > 0 ? Number(value) : fallback
  return {
    reject_calls: Boolean(cfg.reject_calls),
    reject_calls_message: cfg.reject_calls_message || '',
    mark_read_on_message: Boolean(cfg.mark_read_on_message),
    ignore_status_broadcast: Boolean(cfg.ignore_status_broadcast),
    ignore_newsletters: Boolean(cfg.ignore_newsletters),
    outgoing_messages_coex: Boolean(cfg.outgoing_messages_coex),
    history_sync_enabled: Boolean(cfg.history_sync_enabled),
    history_sync_days: String(positiveOr(cfg.history_sync_days, 7)),
    always_online: Boolean(cfg.always_online),
    always_online_interval_seconds: String(
      Number.isFinite(Number(cfg.always_online_interval_seconds))
        ? Number(cfg.always_online_interval_seconds)
        : 60,
    ),
    proxy_enabled: Boolean(cfg.proxy_enabled),
    proxy_protocol: cfg.proxy_protocol || '',
    proxy_host: cfg.proxy_host || '',
    proxy_port: cfg.proxy_port || '',
    proxy_username: cfg.proxy_username || '',
    proxy_password: cfg.proxy_password || '',
  }
}

export function configFormToPayload(form) {
  const interval = parseInt(form.always_online_interval_seconds, 10)
  const days = parseInt(form.history_sync_days, 10)
  return {
    reject_calls: form.reject_calls,
    reject_calls_message: form.reject_calls_message.trim(),
    mark_read_on_message: form.mark_read_on_message,
    ignore_status_broadcast: form.ignore_status_broadcast,
    ignore_newsletters: form.ignore_newsletters,
    outgoing_messages_coex: form.outgoing_messages_coex,
    history_sync_enabled: form.history_sync_enabled,
    history_sync_days: form.history_sync_enabled && Number.isFinite(days) ? days : 0,
    always_online: form.always_online,
    always_online_interval_seconds: Number.isFinite(interval) ? interval : 60,
    proxy_enabled: form.proxy_enabled,
    proxy_protocol: form.proxy_protocol.trim(),
    proxy_host: form.proxy_host.trim(),
    proxy_port: form.proxy_port.trim(),
    proxy_username: form.proxy_username.trim(),
    proxy_password: form.proxy_password,
  }
}

export function defaultWebhookForm() {
  return {
    id: '',
    name: '',
    url: '',
    url_absolute: '',
    header: '',
    token: '',
    timeout_ms: '5000',
    enabled: true,
    append_session: true,
    send_incoming_messages: true,
    send_outgoing_messages: true,
    send_update_messages: true,
    send_group_messages: true,
    send_newsletter_messages: false,
    send_session_disconnected: false,
  }
}

export function webhookFormFromItem(item) {
  const num = (value, fallback) =>
    String(Number.isFinite(Number(value)) ? Number(value) : fallback)
  return {
    id: item?.id || '',
    name: item?.name || '',
    url: item?.url || '',
    url_absolute: item?.url_absolute || '',
    header: item?.header || '',
    token: item?.token || '',
    timeout_ms: num(item?.timeout_ms, 5000),
    enabled: Boolean(item?.enabled),
    append_session: Boolean(item?.append_session),
    send_incoming_messages: Boolean(item?.send_incoming_messages),
    send_outgoing_messages: Boolean(item?.send_outgoing_messages),
    send_update_messages: Boolean(item?.send_update_messages),
    send_group_messages: Boolean(item?.send_group_messages),
    send_newsletter_messages: Boolean(item?.send_newsletter_messages),
    send_session_disconnected: Boolean(item?.send_session_disconnected),
  }
}

export function webhookFormToPayload(form) {
  const timeout = parseInt(form.timeout_ms, 10)
  return {
    id: form.id.trim(),
    name: form.name.trim(),
    url: form.url.trim(),
    url_absolute: form.url_absolute.trim(),
    header: form.header.trim(),
    token: form.token.trim(),
    timeout_ms: Number.isFinite(timeout) && timeout > 0 ? timeout : 5000,
    enabled: form.enabled,
    append_session: form.append_session,
    send_incoming_messages: form.send_incoming_messages,
    send_outgoing_messages: form.send_outgoing_messages,
    send_update_messages: form.send_update_messages,
    send_group_messages: form.send_group_messages,
    send_newsletter_messages: form.send_newsletter_messages,
    send_session_disconnected: form.send_session_disconnected,
  }
}

export function defaultGlobalWebhookForm() {
  return {
    ...defaultWebhookForm(),
    append_session: false,
    send_incoming_messages: false,
    send_outgoing_messages: false,
    send_update_messages: false,
    send_group_messages: false,
    send_newsletter_messages: false,
    send_session_disconnected: true,
  }
}

export function defaultTemplateForm() {
  return {
    id: '',
    name: '',
    language: 'pt_BR',
    category: '',
    header_type: '',
    header_text: '',
    header_link: '',
    header_id: '',
    header_filename: '',
    body: '',
    footer: '',
  }
}

export function templateFormFromItem(item) {
  const header = item?.header || {}
  return {
    id: item?.id || '',
    name: item?.name || '',
    language: item?.language || 'pt_BR',
    category: item?.category || '',
    header_type: header?.type || '',
    header_text: header?.text || '',
    header_link: header?.link || '',
    header_id: header?.id || '',
    header_filename: header?.filename || '',
    body: item?.body || '',
    footer: item?.footer || '',
  }
}

export function templateFormToPayload(form) {
  return {
    id: form.id.trim(),
    name: form.name.trim(),
    language: form.language.trim(),
    category: form.category.trim(),
    header: {
      type: form.header_type.trim(),
      text: form.header_text.trim(),
      link: form.header_link.trim(),
      id: form.header_id.trim(),
      filename: form.header_filename.trim(),
    },
    body: form.body.trim(),
    footer: form.footer.trim(),
  }
}
