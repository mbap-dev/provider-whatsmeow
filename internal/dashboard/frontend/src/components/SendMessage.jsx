import { useState } from 'react'
import { Button, Field, Input, Textarea, Eyebrow } from './ui.jsx'
import { sendText } from '../api/endpoints.js'

export function SendMessage({ sessionId, notify }) {
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
    } catch (error) {
      notify(error.message, 'error')
    } finally {
      setSending(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-3">
      <Eyebrow>Enviar mensagem de teste</Eyebrow>
      <Field label="Destino" hint="Número com DDI ou JID completo.">
        <Input
          value={to}
          placeholder="5511999999999"
          onChange={(e) => setTo(e.target.value)}
        />
      </Field>
      <Field label="Mensagem">
        <Textarea
          rows={3}
          value={body}
          placeholder="Escreva a mensagem…"
          onChange={(e) => setBody(e.target.value)}
        />
      </Field>
      <Button type="submit" variant="primary" loading={sending}>
        Enviar
      </Button>
    </form>
  )
}
