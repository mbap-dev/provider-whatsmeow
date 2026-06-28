// Thin fetch wrapper. Auth is cookie-based (wm_session) for browser sessions
// and Bearer wm_* for API key usage. We normalize JSON and errors here.

const enc = encodeURIComponent

export class ApiError extends Error {
  constructor(message, status) {
    super(message)
    this.status = status
    this.name = 'ApiError'
  }
}

export async function request(path, options = {}) {
  const init = { ...options }
  const headers = new Headers(init.headers || {})

  if (init.body && !(init.body instanceof FormData) && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }

  const response = await fetch(path, { credentials: 'include', ...init, headers })

  if (!response.ok) {
    throw new ApiError(await extractError(response), response.status)
  }

  const contentType = response.headers.get('content-type') || ''
  return contentType.includes('application/json') ? response.json() : null
}

async function extractError(response) {
  const fallback = `${response.status} ${response.statusText}`.trim()
  try {
    const contentType = response.headers.get('content-type') || ''
    if (contentType.includes('application/json')) {
      const payload = await response.json()
      return (
        payload?.error?.message ||
        payload?.message ||
        payload?.error ||
        JSON.stringify(payload)
      )
    }
    const text = (await response.text()).trim()
    return text || fallback
  } catch {
    return fallback
  }
}

export const get = (path) => request(path)
export const post = (path, body) =>
  request(path, { method: 'POST', body: body == null ? undefined : JSON.stringify(body) })
export const put = (path, body) =>
  request(path, { method: 'PUT', body: JSON.stringify(body) })
export const del = (path) => request(path, { method: 'DELETE' })

export { enc }
