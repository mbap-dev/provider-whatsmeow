import { Button, Modal } from '../ui.jsx'

export function FormModal({
  open,
  onClose,
  eyebrow,
  title,
  maxWidth = 'max-w-lg',
  children,
  onSubmit,
  submitLabel = 'Salvar',
  submitting = false,
  cancelLabel = 'Cancelar',
}) {
  return (
    <Modal open={open} onClose={onClose} eyebrow={eyebrow} title={title} maxWidth={maxWidth}>
      <form onSubmit={onSubmit} className="flex flex-col gap-4">
        {children}
        <div className="flex justify-end gap-2 border-t border-border pt-4">
          <Button variant="ghost" onClick={onClose}>
            {cancelLabel}
          </Button>
          <Button type="submit" variant="primary" loading={submitting}>
            {submitLabel}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
