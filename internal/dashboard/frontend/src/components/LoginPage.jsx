import { useState } from 'react'
import { Button, Field, Input } from './ui.jsx'
import { useTheme } from '../hooks/useTheme.js'

const LOGO_LIGHT = 'https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_black.svg'
const LOGO_DARK = 'https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_white.svg'

export function LoginPage({ onLogin }) {
  const { theme } = useTheme()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await onLogin(email, password)
    } catch (err) {
      setError(err.message || 'Credenciais inválidas')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-surface px-4">
      <div className="w-full max-w-sm">
        <div className="mb-8 flex flex-col items-center gap-3">
          <img
            src={theme === 'dark' ? LOGO_DARK : LOGO_LIGHT}
            alt="Envolve"
            className="h-8 w-auto"
          />
          <div className="text-center">
            <h1 className="font-display text-lg font-semibold text-fg">
              WhatsMeow Control Center
            </h1>
            <p className="mt-1 text-sm text-muted">Faça login para continuar</p>
          </div>
        </div>

        <form
          onSubmit={handleSubmit}
          className="flex flex-col gap-4 rounded-[var(--radius-card)] border border-border bg-surface-2 p-6 shadow-sm"
        >
          <Field label="E-mail" htmlFor="login-email">
            <Input
              id="login-email"
              type="email"
              placeholder="admin@example.com"
              autoComplete="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </Field>

          <Field label="Senha" htmlFor="login-password">
            <Input
              id="login-password"
              type="password"
              placeholder="••••••••"
              autoComplete="current-password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </Field>

          {error && (
            <p className="rounded-lg bg-down-soft px-3 py-2 text-sm text-down">{error}</p>
          )}

          <Button type="submit" variant="primary" loading={loading} className="mt-1 w-full">
            Entrar
          </Button>
        </form>
      </div>
    </div>
  )
}
