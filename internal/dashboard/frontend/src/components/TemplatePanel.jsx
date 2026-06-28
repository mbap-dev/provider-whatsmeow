import {
  Button,
  Field,
  Input,
  Textarea,
  Select,
  Tag,
  Mono,
  Eyebrow,
  EmptyState,
  Spinner,
} from './ui.jsx'

function TemplateItem({ item, onEdit, onDelete }) {
  const header = item?.header || {}
  const headerType = header.type || (header.text ? 'text' : 'sem cabeçalho')
  return (
    <article className="rounded-lg border border-border bg-surface-2 p-3">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="truncate text-sm font-medium text-fg">{item.name}</p>
          <p className="text-xs text-muted">
            {item.language || 'idioma padrão'}
            {item.category ? ` · ${item.category}` : ''}
          </p>
        </div>
        <div className="flex shrink-0 gap-1">
          <Button size="sm" variant="ghost" onClick={() => onEdit(item)}>
            Editar
          </Button>
          <Button size="sm" variant="danger" onClick={() => onDelete(item.id)}>
            Excluir
          </Button>
        </div>
      </div>
      <p className="mt-2 line-clamp-3 whitespace-pre-wrap text-sm text-muted">
        {item.body || '—'}
      </p>
      <div className="mt-2 flex flex-wrap gap-1.5">
        <Tag muted>Header: {headerType}</Tag>
        {item.footer && <Tag muted>Com rodapé</Tag>}
        <Mono className="text-faint">{item.id}</Mono>
      </div>
    </article>
  )
}

export function TemplatePanel({
  items,
  loading,
  saving,
  form,
  onChange,
  onSave,
  onDelete,
  onEdit,
  onResetForm,
}) {
  const set = (patch) => onChange({ ...form, ...patch })
  const editing = Boolean(form.id)

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      <div className="flex flex-col gap-3">
        <Eyebrow>Templates salvos</Eyebrow>
        {loading ? (
          <div className="flex justify-center py-8 text-muted">
            <Spinner />
          </div>
        ) : items.length === 0 ? (
          <EmptyState
            title="Nenhum template"
            description="Crie modelos reutilizáveis de mensagem para esta sessão."
          />
        ) : (
          items.map((item) => (
            <TemplateItem key={item.id} item={item} onEdit={onEdit} onDelete={onDelete} />
          ))
        )}
      </div>

      <form onSubmit={onSave} className="flex flex-col gap-3">
        <div className="flex items-center justify-between">
          <Eyebrow>{editing ? `Editando ${form.id}` : 'Novo template'}</Eyebrow>
          {editing && (
            <Button size="sm" variant="ghost" onClick={onResetForm}>
              Cancelar edição
            </Button>
          )}
        </div>
        <div className="grid gap-3 sm:grid-cols-2">
          <Field label="Nome">
            <Input value={form.name} onChange={(e) => set({ name: e.target.value })} />
          </Field>
          <Field label="Idioma">
            <Input value={form.language} onChange={(e) => set({ language: e.target.value })} />
          </Field>
          <Field label="Categoria">
            <Input value={form.category} onChange={(e) => set({ category: e.target.value })} />
          </Field>
          <Field label="Tipo de cabeçalho">
            <Select
              value={form.header_type}
              onChange={(e) => set({ header_type: e.target.value })}
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
              onChange={(e) => set({ header_text: e.target.value })}
            />
          </Field>
        )}
        {['image', 'video', 'document'].includes(form.header_type) && (
          <div className="grid gap-3 sm:grid-cols-2">
            <Field label="Link da mídia">
              <Input
                value={form.header_link}
                onChange={(e) => set({ header_link: e.target.value })}
              />
            </Field>
            <Field label="ID da mídia">
              <Input
                value={form.header_id}
                onChange={(e) => set({ header_id: e.target.value })}
              />
            </Field>
            <Field label="Nome do arquivo">
              <Input
                value={form.header_filename}
                onChange={(e) => set({ header_filename: e.target.value })}
              />
            </Field>
          </div>
        )}

        <Field label="Corpo">
          <Textarea
            rows={5}
            value={form.body}
            placeholder="Olá {{1}}, sua solicitação foi recebida."
            onChange={(e) => set({ body: e.target.value })}
          />
        </Field>
        <Field label="Rodapé">
          <Input value={form.footer} onChange={(e) => set({ footer: e.target.value })} />
        </Field>
        <Button type="submit" variant="primary" loading={saving}>
          {editing ? 'Salvar alterações' : 'Criar template'}
        </Button>
      </form>
    </div>
  )
}
