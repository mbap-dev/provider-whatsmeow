import { useEffect, useRef, useState } from 'react'
import { cx } from './ui.jsx'

const ROLE_STYLE = {
  owner: 'text-brand',
  admin: 'text-pending',
  user: 'text-muted',
}

export function UserMenu({ user, onLogout, onOpenAPIKeys }) {
  const [open, setOpen] = useState(false)
  const ref = useRef(null)
  const canManageAPIKeys = user?.role === 'owner' || user?.role === 'admin'

  useEffect(() => {
    if (!open) return undefined
    const handler = (e) => {
      if (ref.current && !ref.current.contains(e.target)) setOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [open])

  const initial = (user?.email ?? '?')[0].toUpperCase()

  return (
    <div ref={ref} className="relative">
      <button
        onClick={() => setOpen((v) => !v)}
        className="flex size-8 items-center justify-center rounded-full bg-brand text-[0.7rem] font-bold text-brand-ink transition-opacity hover:opacity-85"
        aria-label="Menu do usuário"
        aria-expanded={open}
      >
        {initial}
      </button>

      {open && (
        <div className="absolute right-0 top-10 z-50 w-52 rounded-[var(--radius-card)] border border-border bg-surface shadow-xl">
          <div className="border-b border-border px-4 py-3">
            <p className="truncate text-sm font-medium text-fg">{user?.email}</p>
            <p className={cx('mt-0.5 text-xs font-medium capitalize', ROLE_STYLE[user?.role] ?? 'text-muted')}>
              {user?.role}
            </p>
          </div>
          <div className="flex flex-col gap-0.5 p-1.5">
            {canManageAPIKeys && (
              <button
                onClick={() => {
                  setOpen(false)
                  onOpenAPIKeys()
                }}
                className="w-full rounded-md px-3 py-2 text-left text-sm text-fg transition-colors hover:bg-surface-2"
              >
                API Keys
              </button>
            )}
            <button
              onClick={() => {
                setOpen(false)
                onLogout()
              }}
              className="w-full rounded-md px-3 py-2 text-left text-sm text-down transition-colors hover:bg-down-soft"
            >
              Sair
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
