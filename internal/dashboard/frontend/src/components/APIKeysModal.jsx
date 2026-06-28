import { useCallback, useEffect, useState } from 'react'
import { Button, Input, Modal, Spinner } from './ui.jsx'
import * as api from '../api/endpoints.js'

function fmtDate(s) {
  return s ? new Date(s).toLocaleDateString('pt-BR') : null
}

export function APIKeysModal({ open, onClose }) {
  const [keys, setKeys] = useState([])
  const [loading, setLoading] = useState(false)
  const [newName, setNewName] = useState('')
  const [creating, setCreating] = useState(false)
  const [newKey, setNewKey] = useState(null) // revealed once on creation
  const [copied, setCopied] = useState(false)
  const [error, setError] = useState('')

  const reload = useCallback(async () => {
    setLoading(true)
    try {
      const data = await api.listAPIKeys()
      setKeys(Array.isArray(data) ? data : [])
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!open) return
    setNewKey(null)
    setNewName('')
    setError('')
    reload()
  }, [open, reload])

  const handleCreate = async (e) => {
    e.preventDefault()
    if (!newName.trim()) return
    setCreating(true)
    setError('')
    try {
      const k = await api.createAPIKey({ name: newName.trim() })
      setNewKey(k)
      setNewName('')
      await reload()
    } catch (e) {
      setError(e.message)
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async (id) => {
    if (!window.confirm('Revogar esta API Key?')) return
    try {
      await api.deleteAPIKey(id)
      await reload()
    } catch (e) {
      setError(e.message)
    }
  }

  const handleCopy = async () => {
    if (!newKey?.key) return
    await navigator.clipboard.writeText(newKey.key)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <Modal open={open} onClose={onClose} eyebrow="Autenticação" title="API Keys" maxWidth="max-w-lg">
      <div className="flex flex-col gap-5">
        {/* Revealed key — shown once after creation */}
        {newKey && (
          <div className="flex flex-col gap-2 rounded-lg border border-live bg-live-soft p-4">
            <p className="text-xs font-medium text-live">
              Chave criada — copie agora, ela não será exibida novamente.
            </p>
            <div className="flex items-center gap-2">
              <code className="flex-1 truncate rounded bg-surface px-2 py-1.5 font-mono text-xs text-fg">
                {newKey.key}
              </code>
              <Button size="sm" variant="secondary" onClick={handleCopy}>
                {copied ? 'Copiado!' : 'Copiar'}
              </Button>
            </div>
          </div>
        )}

        {error && (
          <p className="rounded-lg bg-down-soft px-3 py-2 text-sm text-down">{error}</p>
        )}

        {/* Key list */}
        {loading ? (
          <div className="flex justify-center py-6">
            <Spinner />
          </div>
        ) : keys.length === 0 ? (
          <p className="py-4 text-center text-sm text-muted">
            Nenhuma API Key. Crie uma abaixo.
          </p>
        ) : (
          <div className="flex flex-col gap-2">
            {keys.map((k) => (
              <div
                key={k.id}
                className="flex items-center justify-between gap-3 rounded-lg border border-border px-3 py-2.5"
              >
                <div className="min-w-0">
                  <p className="truncate text-sm font-medium text-fg">{k.name}</p>
                  <p className="font-mono text-[0.68rem] text-faint">
                    Criada {fmtDate(k.created_at)}
                    {k.expires_at ? ` · Expira ${fmtDate(k.expires_at)}` : ''}
                    {k.last_used_at ? ` · Último uso ${fmtDate(k.last_used_at)}` : ''}
                  </p>
                </div>
                <Button size="sm" variant="danger" onClick={() => handleDelete(k.id)}>
                  Revogar
                </Button>
              </div>
            ))}
          </div>
        )}

        {/* Create form */}
        <form onSubmit={handleCreate} className="flex gap-2 border-t border-border pt-4">
          <Input
            placeholder="Nome da chave"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            className="flex-1"
            required
          />
          <Button type="submit" variant="primary" loading={creating}>
            Criar
          </Button>
        </form>
      </div>
    </Modal>
  )
}
