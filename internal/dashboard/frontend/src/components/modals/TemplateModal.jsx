import { Field, Input, Select, Textarea } from '../ui.jsx'
import { FormModal } from './FormModal.jsx'

export function TemplateModal({
  open,
  onClose,
  form,
  saving,
  onChange,
  onSave,
}) {
  const set = (patch) => onChange({ ...form, ...patch })
  const editing = Boolean(form.id)

  return (
    <FormModal
      open={open}
      onClose={onClose}
      eyebrow={editing ? 'Editar template' : 'Novo template'}
      title={editing ? form.name || form.id : 'Cadastrar template'}
      maxWidth="max-w-2xl"
      onSubmit={onSave}
      submitLabel={editing ? 'Salvar alterações' : 'Criar template'}
      submitting={saving}
    >
      <div className="grid gap-3 sm:grid-cols-2">
        <Field label="Nome">
          <Input value={form.name} onChange={(event) => set({ name: event.target.value })} />
        </Field>
        <Field label="Idioma">
          <Input value={form.language} onChange={(event) => set({ language: event.target.value })} />
        </Field>
        <Field label="Categoria">
          <Input value={form.category} onChange={(event) => set({ category: event.target.value })} />
        </Field>
        <Field label="Tipo de cabeçalho">
          <Select
            value={form.header_type}
            onChange={(event) => set({ header_type: event.target.value })}
          >
            <option value="">Sem cabeçalho</option>
            <option value="text">Texto</option>
            <option value="image">Imagem</option>
            <option value="video">Vídeo</option>
            <option value="document">Documento</option>
          </Select>
        </Field>
      </div>

      {form.header_type === 'text' && (
        <Field label="Texto do cabeçalho">
          <Input
            value={form.header_text}
            onChange={(event) => set({ header_text: event.target.value })}
          />
        </Field>
      )}
      {['image', 'video', 'document'].includes(form.header_type) && (
        <div className="grid gap-3 sm:grid-cols-2">
          <Field label="Link da mídia">
            <Input
              value={form.header_link}
              onChange={(event) => set({ header_link: event.target.value })}
            />
          </Field>
          <Field label="ID da mídia">
            <Input
              value={form.header_id}
              onChange={(event) => set({ header_id: event.target.value })}
            />
          </Field>
          <Field label="Nome do arquivo">
            <Input
              value={form.header_filename}
              onChange={(event) => set({ header_filename: event.target.value })}
            />
          </Field>
        </div>
      )}

      <Field label="Corpo">
        <Textarea
          rows={5}
          value={form.body}
          placeholder="Olá {{1}}, sua solicitação foi recebida."
          onChange={(event) => set({ body: event.target.value })}
        />
      </Field>
      <Field label="Rodapé">
        <Input value={form.footer} onChange={(event) => set({ footer: event.target.value })} />
      </Field>
    </FormModal>
  )
}
