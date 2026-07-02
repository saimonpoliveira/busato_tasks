import { useState, useMemo } from 'react'
import type { ColDef } from 'ag-grid-community'
import { DataGrid } from '@/components/common/DataGrid'
import { SearchInput } from '@/components/common/SearchInput'
import { useComments } from '@/hooks/useApi'
import type { Comment } from '@/types'

export function CommentsPage() {
  const [search, setSearch] = useState('')
  const { data, isLoading } = useComments({ search, page_size: 100 })

  const columnDefs = useMemo<ColDef<Comment>[]>(
    () => [
      {
        field: 'user.name',
        headerName: 'Autor',
        valueGetter: (params) => params.data?.user?.name || '-',
      },
      { field: 'entity_type', headerName: 'Tipo' },
      { field: 'content', headerName: 'Conteúdo', flex: 3 },
      {
        field: 'created_at',
        headerName: 'Criado em',
        valueFormatter: (params) =>
          params.value ? new Date(params.value as string).toLocaleString('pt-BR') : '',
      },
    ],
    []
  )

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Comentários</h2>
        <p className="text-muted-foreground">Visualize os comentários em chamados e tarefas</p>
      </div>

      <SearchInput value={search} onChange={setSearch} placeholder="Pesquisar comentários..." />

      <DataGrid rowData={data?.data ?? []} columnDefs={columnDefs} loading={isLoading} />
    </div>
  )
}
