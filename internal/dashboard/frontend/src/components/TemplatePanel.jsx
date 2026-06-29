import {
  Button,
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
  onDelete,
  onEdit,
  onCreate,
}) {
  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-between gap-3">
        <Eyebrow>Templates salvos</Eyebrow>
        <Button size="sm" variant="primary" onClick={onCreate}>
          Novo template
        </Button>
      </div>
      {loading ? (
        <div className="flex justify-center py-8 text-muted">
          <Spinner />
        </div>
      ) : items.length === 0 ? (
        <EmptyState
          title="Nenhum template"
          description="Crie modelos reutilizáveis de mensagem para esta sessão."
          action={
            <Button size="sm" variant="primary" onClick={onCreate}>
              Novo template
            </Button>
          }
        />
      ) : (
        items.map((item) => (
          <TemplateItem key={item.id} item={item} onEdit={onEdit} onDelete={onDelete} />
        ))
      )}
    </div>
  )
}
