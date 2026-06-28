import { useEffect, useState } from 'react'
import { Modal } from './ui.jsx'
import { qrSrc } from '../api/endpoints.js'

const REFRESH_MS = 3000

export function QRModal({ sessionId, onClose }) {
  const [tick, setTick] = useState(Date.now())

  useEffect(() => {
    if (!sessionId) return undefined
    setTick(Date.now())
    const timer = window.setInterval(() => setTick(Date.now()), REFRESH_MS)
    return () => window.clearInterval(timer)
  }, [sessionId])

  return (
    <Modal
      open={Boolean(sessionId)}
      onClose={onClose}
      eyebrow="Pareamento"
      title={`Conectar ${sessionId || ''}`}
    >
      <div className="flex flex-col items-center gap-4">
        <div className="rounded-xl bg-white p-3">
          <img
            className="size-56"
            src={sessionId ? qrSrc(sessionId, tick) : ''}
            alt={`QR da sessão ${sessionId || ''}`}
          />
        </div>
        <p className="text-center text-sm text-muted">
          Abra o WhatsApp no celular, vá em Aparelhos conectados e aponte para o código. Ele
          se renova sozinho a cada poucos segundos.
        </p>
      </div>
    </Modal>
  )
}
