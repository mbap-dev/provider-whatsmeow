import { useEffect, useState } from 'react'

const SESSION_TABS = new Set(['overview', 'config', 'webhooks', 'templates'])
const SETTINGS_SECTIONS = new Set(['users', 'webhooks'])

function cleanPath(pathname) {
  const path = pathname || '/'
  if (path === '/') return path
  return path.replace(/\/+$/, '') || '/'
}

function decodePathPart(value) {
  try {
    return decodeURIComponent(value)
  } catch {
    return value
  }
}

export function pathForHome() {
  return '/'
}

export function pathForInstances() {
  return '/instances'
}

export function pathForInstance(sessionId, tab = 'overview') {
  const id = encodeURIComponent(String(sessionId || '').trim())
  if (!id) return pathForInstances()
  return tab && tab !== 'overview' ? `/instances/${id}/${tab}` : `/instances/${id}`
}

export function pathForSettings(section = 'users') {
  const active = SETTINGS_SECTIONS.has(section) ? section : 'users'
  return `/settings/${active}`
}

export function parseRoute(pathname = window.location.pathname) {
  const path = cleanPath(pathname)
  if (path === '/') return { page: 'home' }
  if (path === '/instances') return { page: 'instances', sessionId: '', tab: 'overview' }
  if (path.startsWith('/instances/')) {
    const parts = path.slice('/instances/'.length).split('/')
    const sessionId = decodePathPart(parts[0] || '')
    const tab = SESSION_TABS.has(parts[1]) ? parts[1] : 'overview'
    return { page: 'instances', sessionId, tab }
  }
  if (path === '/settings') return { page: 'settings', section: 'users' }
  if (path.startsWith('/settings/')) {
    const section = path.slice('/settings/'.length).split('/')[0]
    return {
      page: 'settings',
      section: SETTINGS_SECTIONS.has(section) ? section : 'users',
    }
  }
  return { page: 'not-found' }
}

export function navigate(path, { replace = false } = {}) {
  const next = path || '/'
  if (window.location.pathname === next) return
  const method = replace ? 'replaceState' : 'pushState'
  window.history[method](null, '', next)
  window.dispatchEvent(new PopStateEvent('popstate'))
}

export function useRoute() {
  const [route, setRoute] = useState(() => parseRoute())

  useEffect(() => {
    const handleChange = () => setRoute(parseRoute())
    window.addEventListener('popstate', handleChange)
    return () => window.removeEventListener('popstate', handleChange)
  }, [])

  return route
}
