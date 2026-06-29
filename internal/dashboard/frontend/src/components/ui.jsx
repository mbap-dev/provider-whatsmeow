import { useEffect } from 'react'

function cx(...parts) {
  return parts.filter(Boolean).join(' ')
}

// --- Button -----------------------------------------------------------------
const BUTTON_VARIANTS = {
  primary:
    'bg-brand text-brand-ink hover:bg-brand-strong border border-transparent shadow-sm',
  secondary:
    'bg-surface-2 text-fg border border-border-strong hover:border-brand hover:text-brand',
  ghost: 'bg-transparent text-muted border border-transparent hover:text-fg hover:bg-surface-2',
  danger:
    'bg-transparent text-down border border-[color-mix(in_oklab,var(--color-down)_40%,transparent)] hover:bg-down-soft',
}

const BUTTON_SIZES = {
  sm: 'h-8 px-3 text-xs gap-1.5',
  md: 'h-10 px-4 text-sm gap-2',
}

export function Button({
  variant = 'secondary',
  size = 'md',
  type = 'button',
  className,
  disabled,
  loading,
  children,
  ...rest
}) {
  return (
    <button
      type={type}
      disabled={disabled || loading}
      className={cx(
        'inline-flex items-center justify-center rounded-lg font-medium tracking-tight',
        'transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed',
        'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand',
        BUTTON_SIZES[size],
        BUTTON_VARIANTS[variant],
        className,
      )}
      {...rest}
    >
      {loading && <Spinner className="size-3.5" />}
      {children}
    </button>
  )
}

// --- Surfaces ---------------------------------------------------------------
export function Card({ className, children, ...rest }) {
  return (
    <section
      className={cx(
        'rounded-[var(--radius-card)] border border-border bg-surface',
        className,
      )}
      {...rest}
    >
      {children}
    </section>
  )
}

export function CardHeader({ eyebrow, title, description, actions }) {
  return (
    <header className="flex items-start justify-between gap-4 border-b border-border px-5 py-4">
      <div className="min-w-0">
        {eyebrow && <Eyebrow>{eyebrow}</Eyebrow>}
        <h2 className="font-display text-base font-semibold text-fg">{title}</h2>
        {description && <p className="mt-0.5 text-sm text-muted">{description}</p>}
      </div>
      {actions && <div className="flex shrink-0 items-center gap-2">{actions}</div>}
    </header>
  )
}

export function Eyebrow({ children, className }) {
  return (
    <p
      className={cx(
        'font-mono text-[0.68rem] font-medium uppercase tracking-[0.18em] text-faint',
        className,
      )}
    >
      {children}
    </p>
  )
}

// --- Signal language (signature) -------------------------------------------
const TONE_TEXT = {
  live: 'text-live',
  pending: 'text-pending',
  down: 'text-down',
}
const TONE_BG = {
  live: 'bg-live',
  pending: 'bg-pending',
  down: 'bg-down',
}

export function SignalDot({ tone = 'down', pulse = false, className }) {
  return (
    <span className={cx('relative inline-flex', TONE_TEXT[tone], className)}>
      <span
        className={cx(
          'size-2.5 rounded-full',
          TONE_BG[tone],
          pulse && 'animate-[var(--animate-pulse-ring)]',
        )}
      />
    </span>
  )
}

const TONE_BADGE = {
  live: 'text-live bg-live-soft',
  pending: 'text-pending bg-pending-soft',
  down: 'text-down bg-down-soft',
  neutral: 'text-muted bg-surface-2',
}

export function StatusBadge({ tone = 'neutral', pulse = false, children }) {
  return (
    <span
      className={cx(
        'inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium',
        TONE_BADGE[tone],
      )}
    >
      {tone !== 'neutral' && <SignalDot tone={tone} pulse={pulse} />}
      {children}
    </span>
  )
}

export function Tag({ children, muted = false }) {
  return (
    <span
      className={cx(
        'inline-flex items-center rounded-md px-2 py-0.5 text-[0.72rem] font-medium',
        muted
          ? 'bg-transparent text-faint ring-1 ring-inset ring-border'
          : 'bg-surface-2 text-muted',
      )}
    >
      {children}
    </span>
  )
}

export function Mono({ children, className }) {
  return <code className={cx('font-mono text-xs text-fg', className)}>{children}</code>
}

// --- Form controls ----------------------------------------------------------
export function Field({ label, hint, htmlFor, children }) {
  return (
    <label htmlFor={htmlFor} className="flex flex-col gap-1.5">
      {label && <span className="text-xs font-medium text-muted">{label}</span>}
      {children}
      {hint && <span className="text-[0.72rem] text-faint">{hint}</span>}
    </label>
  )
}

const INPUT_BASE =
  'w-full rounded-lg border border-border bg-surface-2 px-3 text-sm text-fg placeholder:text-faint ' +
  'transition-colors focus:border-brand focus:bg-surface focus:outline-none'

export function Input({ className, ...rest }) {
  return <input className={cx(INPUT_BASE, 'h-10', className)} {...rest} />
}

export function Textarea({ className, rows = 4, ...rest }) {
  return <textarea rows={rows} className={cx(INPUT_BASE, 'resize-y py-2 leading-relaxed', className)} {...rest} />
}

export function Select({ className, children, ...rest }) {
  return (
    <select className={cx(INPUT_BASE, 'h-10 appearance-none', className)} {...rest}>
      {children}
    </select>
  )
}

export function Switch({ label, description, checked, onChange }) {
  return (
    <label className="flex cursor-pointer items-start gap-3 rounded-lg p-1 -m-1 hover:bg-surface-2">
      <button
        type="button"
        role="switch"
        aria-checked={checked}
        onClick={() => onChange(!checked)}
        className={cx(
          'mt-0.5 inline-flex h-5 w-9 shrink-0 items-center rounded-full p-0.5 transition-colors',
          checked ? 'bg-brand' : 'bg-border-strong',
        )}
      >
        <span
          className={cx(
            'size-4 rounded-full bg-white shadow-sm transition-transform',
            checked && 'translate-x-4',
          )}
        />
      </button>
      <span className="min-w-0">
        <span className="block text-sm font-medium text-fg">{label}</span>
        {description && <span className="block text-xs text-muted">{description}</span>}
      </span>
    </label>
  )
}

// --- Tabs -------------------------------------------------------------------
export function Tabs({ tabs, active, onChange }) {
  return (
    <div role="tablist" className="flex gap-1 overflow-x-auto rounded-lg bg-surface-2 p-1">
      {tabs.map((tab) => {
        const selected = tab.id === active
        return (
          <button
            key={tab.id}
            role="tab"
            aria-selected={selected}
            onClick={() => onChange(tab.id)}
            className={cx(
              'flex-1 whitespace-nowrap rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
              selected
                ? 'bg-surface text-fg shadow-sm'
                : 'text-muted hover:text-fg',
            )}
          >
            {tab.label}
          </button>
        )
      })}
    </div>
  )
}

// --- Modal ------------------------------------------------------------------
export function Modal({ open, onClose, title, eyebrow, children, maxWidth = 'max-w-sm' }) {
  useEffect(() => {
    if (!open) return undefined
    const onKey = (event) => event.key === 'Escape' && onClose()
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open, onClose])

  if (!open) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
      onClick={(event) => event.target === event.currentTarget && onClose()}
    >
      <div
        className={cx(
          'flex max-h-[calc(100vh-2rem)] w-full flex-col overflow-hidden rounded-[var(--radius-card)] border border-border bg-surface shadow-2xl',
          maxWidth,
        )}
      >
        <header className="flex items-start justify-between gap-4 border-b border-border px-5 py-4">
          <div>
            {eyebrow && <Eyebrow>{eyebrow}</Eyebrow>}
            <h2 className="font-display text-base font-semibold text-fg">{title}</h2>
          </div>
          <Button variant="ghost" size="sm" onClick={onClose}>
            Fechar
          </Button>
        </header>
        <div className="overflow-y-auto p-5">{children}</div>
      </div>
    </div>
  )
}

// --- Feedback ---------------------------------------------------------------
export function EmptyState({ title, description, action }) {
  return (
    <div className="flex flex-col items-center gap-2 rounded-[var(--radius-card)] border border-dashed border-border px-6 py-12 text-center">
      <SignalDot tone="down" />
      <p className="font-display text-sm font-semibold text-fg">{title}</p>
      {description && <p className="max-w-xs text-sm text-muted">{description}</p>}
      {action && <div className="mt-2">{action}</div>}
    </div>
  )
}

export function Spinner({ className }) {
  return (
    <span
      className={cx(
        'inline-block animate-spin rounded-full border-2 border-current border-t-transparent',
        className || 'size-4',
      )}
      aria-hidden="true"
    />
  )
}

const TOAST_TONE = {
  success: 'border-l-live',
  error: 'border-l-down',
  info: 'border-l-brand',
}

export function ToastStack({ toasts, onDismiss }) {
  return (
    <div className="pointer-events-none fixed bottom-4 right-4 z-[60] flex w-[min(20rem,calc(100vw-2rem))] flex-col gap-2">
      {toasts.map((toast) => (
        <button
          key={toast.id}
          onClick={() => onDismiss(toast.id)}
          className={cx(
            'pointer-events-auto rounded-lg border border-border border-l-4 bg-surface px-4 py-3 text-left text-sm text-fg shadow-lg',
            TOAST_TONE[toast.tone] || TOAST_TONE.info,
          )}
        >
          {toast.message}
        </button>
      ))}
    </div>
  )
}

export { cx }
