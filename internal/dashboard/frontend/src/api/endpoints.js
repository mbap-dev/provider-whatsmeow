import { get, post, put, del, enc } from './client.js'

const V1 = '/api/v1'

// --- Overview ---------------------------------------------------------------
export const getWorker = () => get(`${V1}/worker`)
export const listSessions = () => get(`${V1}/sessions`)

// --- Session lifecycle ------------------------------------------------------
export const createSession = (body) => post(`${V1}/sessions`, body)
export const deleteSession = (id) => del(`${V1}/sessions/${enc(id)}`)

// connect / disconnect / reload live on the unversioned /sessions routes
export const sessionAction = (id, action) => post(`/sessions/${enc(id)}/${action}`)

export const qrSrc = (id, ts) => `/sessions/${enc(id)}/qr?ts=${ts}`
export const requestPairingCode = (id, phone) =>
  post(`${V1}/sessions/${enc(id)}/pairing-code`, { phone })

// --- Worker assignments -----------------------------------------------------
export const assignSession = (id) => post(`${V1}/worker/assignments`, { session_id: id })
export const unassignSession = (id) => del(`${V1}/worker/assignments/${enc(id)}`)

// --- Session config ---------------------------------------------------------
export const getConfig = (id) => get(`${V1}/sessions/${enc(id)}/config`)
export const updateConfig = (id, body) => put(`${V1}/sessions/${enc(id)}/config`, body)
export const clearConfig = (id) => del(`${V1}/sessions/${enc(id)}/config`)

// --- Webhooks ---------------------------------------------------------------
export const listWebhooks = (id) => get(`${V1}/sessions/${enc(id)}/webhooks`)
export const saveWebhook = (id, body) => post(`${V1}/sessions/${enc(id)}/webhooks`, body)
export const deleteWebhook = (id, webhookId) =>
  del(`${V1}/sessions/${enc(id)}/webhooks/${enc(webhookId)}`)

// --- Templates --------------------------------------------------------------
export const listTemplates = (id) => get(`${V1}/sessions/${enc(id)}/templates`)
export const saveTemplate = (id, body) => post(`${V1}/sessions/${enc(id)}/templates`, body)
export const deleteTemplate = (id, templateId) =>
  del(`${V1}/sessions/${enc(id)}/templates/${enc(templateId)}`)

// --- Messaging --------------------------------------------------------------
export const sendText = (id, to, body) =>
  post(`${V1}/sessions/${enc(id)}/messages`, { type: 'text', to, text: { body } })

export const checkUser = (id, numbers) =>
  post(`${V1}/sessions/${enc(id)}/user/check`, { number: numbers })

// --- Auth (Postgres-backed, may return 404 when unavailable) ----------------
export const getMe = () => get('/auth/me')
export const login = (email, password) => post('/auth/login', { email, password })
export const logout = () => post('/auth/logout')
export const listAPIKeys = () => get('/auth/api-keys')
export const createAPIKey = (body) => post('/auth/api-keys', body)
export const deleteAPIKey = (id) => del(`/auth/api-keys/${enc(id)}`)
