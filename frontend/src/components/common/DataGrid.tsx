import { useMemo } from 'react'
import { AgGridReact } from 'ag-grid-react'
import type { ColDef, GridOptions } from 'ag-grid-community'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-quartz.css'
import { cn } from '@/lib/utils'

interface DataGridProps<T> {
  rowData: T[]
  columnDefs: ColDef<T>[]
  loading?: boolean
  className?: string
  onRowClicked?: (data: T) => void
  gridOptions?: GridOptions<T>
}

export function DataGrid<T>({
  rowData,
  columnDefs,
  loading = false,
  className,
  onRowClicked,
  gridOptions,
}: DataGridProps<T>) {
  const defaultColDef = useMemo<ColDef>(
    () => ({
      sortable: true,
      filter: true,
      resizable: true,
      flex: 1,
      minWidth: 120,
    }),
    []
  )

  return (
    <div className={cn('ag-theme-quartz ag-grid-custom w-full', className)} style={{ height: 'calc(100vh - 280px)' }}>
      <AgGridReact<T>
        rowData={rowData}
        columnDefs={columnDefs}
        defaultColDef={defaultColDef}
        loading={loading}
        animateRows
        rowHeight={56}
        headerHeight={48}
        pagination
        paginationPageSize={20}
        paginationPageSizeSelector={[10, 20, 50, 100]}
        onRowClicked={(event) => {
          if (event.data && onRowClicked) {
            onRowClicked(event.data)
          }
        }}
        gridOptions={gridOptions}
      />
    </div>
  )
}
