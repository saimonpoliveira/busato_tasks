import { useState, useMemo } from 'react'
import type { ColDef } from 'ag-grid-community'
import { DataGrid } from '@/components/common/DataGrid'
import { SearchInput } from '@/components/common/SearchInput'
import { useUsers } from '@/hooks/useApi'
import type { User } from '@/types'

export function UsersPage() {
  const [search, setSearch] = useState('')
  const { data, isLoading } = useUsers({ search, page_size: 100 })

  const columnDefs = useMemo<ColDef<User>[]>(
    () => [
      { field: 'name', headerName: 'Nome' },
      { field: 'email', headerName: 'E-mail', flex: 2 },
      { field: 'role', headerName: 'Papel' },
      {
        field: 'active',
        headerName: 'Ativo',
        valueFormatter: (params) => (params.value ? 'Sim' : 'Não'),
      },
      {
        field: 'created_at',
        headerName: 'Criado em',
        valueFormatter: (params) =>
          params.value ? new Date(params.value as string).toLocaleDateString('pt-BR') : '',
      },
    ],
    []
  )

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Usuários</h2>
        <p className="text-muted-foreground">Visualize os usuários do sistema</p>
      </div>

      <SearchInput value={search} onChange={setSearch} placeholder="Pesquisar usuários..." />

      <DataGrid rowData={data?.data ?? []} columnDefs={columnDefs} loading={isLoading} />
    </div>
  )
}
