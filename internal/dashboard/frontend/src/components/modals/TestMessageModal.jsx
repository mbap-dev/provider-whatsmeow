import { useState } from 'react'
import { Field, Input, Textarea } from '../ui.jsx'
import { sendText } from '../../api/endpoints.js'
import { FormModal } from './FormModal.jsx'

export function TestMessageModal({ open, onClose, sessionId, notify }) {
  const [to, setTo] = useState('')
  const [body, setBody] = useState('')
  const [sending, setSending] = useState(false)

  async function handleSubmit(event) {
    event.preventDefault()
    const target = to.trim()
    const text = body.trim()
    if (!target || !text) {
      notify('Informe o destino e a mensagem.', 'error')
      return
    }
    setSending(true)
    try {
      await sendText(sessionId, target, text)
      notify(`Mensagem enviada por ${sessionId}.`, 'success')
      setBody('')
      onClose()
    } catch (error) {
      notify(error.message, 'error')
    } finally {
      setSending(false)
    }
  }

  return (
    <FormModal
      open={open}
      onClose={onClose}
      eyebrow="Mensagem"
      title="Enviar mensagem de teste"
      onSubmit={handleSubmit}
      submitLabel="Enviar"
      submitting={sending}
    >
      <Field label="Destino" hint="Número com DDI ou JID completo.">
        <Input
          value={to}
          placeholder="5511999999999"
          onChange={(event) => setTo(event.target.value)}
        />
      </Field>
      <Field label="Mensagem">
        <Textarea
          rows={4}
          value={body}
          placeholder="Escreva a mensagem..."
          onChange={(event) => setBody(event.target.value)}
        />
      </Field>
    </FormModal>
  )
}
