import { useState, useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import type { ColDef } from 'ag-grid-community'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { DataGrid } from '@/components/common/DataGrid'
import { SearchInput } from '@/components/common/SearchInput'
import { useTasks, useCreateTask, useUpdateTask, useDeleteTask, useTickets } from '@/hooks/useApi'
import { getErrorMessage } from '@/lib/api'
import type { Task } from '@/types'

const taskSchema = z.object({
  ticket_id: z.string().min(1, 'Chamado é obrigatório'),
  title: z.string().min(2, 'Título deve ter no mínimo 2 caracteres'),
  description: z.string().optional(),
  status: z.enum(['todo', 'in_progress', 'done', 'cancelled']).optional(),
})

type TaskForm = z.infer<typeof taskSchema>

const statusLabels: Record<string, string> = {
  todo: 'A Fazer',
  in_progress: 'Em Progresso',
  done: 'Concluída',
  cancelled: 'Cancelada',
}

export function TasksPage() {
  const [search, setSearch] = useState('')
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingTask, setEditingTask] = useState<Task | null>(null)

  const { data, isLoading } = useTasks({ search, page_size: 100 })
  const { data: ticketsData } = useTickets({ page_size: 100 })
  const createMutation = useCreateTask()
  const updateMutation = useUpdateTask()
  const deleteMutation = useDeleteTask()

  const form = useForm<TaskForm>({
    resolver: zodResolver(taskSchema),
    defaultValues: { ticket_id: '', title: '', description: '', status: 'todo' },
  })

  const columnDefs = useMemo<ColDef<Task>[]>(
    () => [
      { field: 'title', headerName: 'Título', flex: 2 },
      {
        field: 'status',
        headerName: 'Status',
        valueFormatter: (params) => statusLabels[params.value as string] || params.value,
      },
      {
        field: 'assignee.name',
        headerName: 'Responsável',
        valueGetter: (params) => params.data?.assignee?.name || '-',
      },
      { field: 'order', headerName: 'Ordem', width: 100 },
      {
        field: 'created_at',
        headerName: 'Criado em',
        valueFormatter: (params) =>
          params.value ? new Date(params.value as string).toLocaleDateString('pt-BR') : '',
      },
      {
        headerName: 'Ações',
        sortable: false,
        filter: false,
        width: 160,
        cellRenderer: (params: { data: Task }) => (
          <div className="flex gap-2 py-2">
            <Button size="sm" variant="outline" onClick={() => openEdit(params.data)}>
              Editar
            </Button>
            <Button
              size="sm"
              variant="destructive"
              onClick={() => deleteMutation.mutate(params.data.id)}
            >
              Excluir
            </Button>
          </div>
        ),
      },
    ],
    [deleteMutation]
  )

  const openCreate = () => {
    setEditingTask(null)
    form.reset({ ticket_id: '', title: '', description: '', status: 'todo' })
    setDialogOpen(true)
  }

  const openEdit = (task: Task) => {
    setEditingTask(task)
    form.reset({
      ticket_id: task.ticket_id,
      title: task.title,
      description: task.description,
      status: task.status,
    })
    setDialogOpen(true)
  }

  const onSubmit = async (formData: TaskForm) => {
    try {
      if (editingTask) {
        await updateMutation.mutateAsync({ id: editingTask.id, data: formData })
      } else {
        await createMutation.mutateAsync(formData)
      }
      setDialogOpen(false)
      form.reset()
    } catch (err) {
      alert(getErrorMessage(err))
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Tarefas</h2>
          <p className="text-muted-foreground">Gerencie as tarefas vinculadas aos chamados</p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={openCreate}>
              <Plus className="mr-2 h-4 w-4" />
              Nova Tarefa
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingTask ? 'Editar Tarefa' : 'Nova Tarefa'}</DialogTitle>
            </DialogHeader>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <div className="space-y-2">
                <Label>Chamado</Label>
                <Select
                  value={form.watch('ticket_id')}
                  onValueChange={(value) => form.setValue('ticket_id', value)}
                  disabled={!!editingTask}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione um chamado" />
                  </SelectTrigger>
                  <SelectContent>
                    {ticketsData?.data.map((ticket) => (
                      <SelectItem key={ticket.id} value={ticket.id}>
                        {ticket.title}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {form.formState.errors.ticket_id && (
                  <p className="text-sm text-destructive">{form.formState.errors.ticket_id.message}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label>Título</Label>
                <Input {...form.register('title')} />
                {form.formState.errors.title && (
                  <p className="text-sm text-destructive">{form.formState.errors.title.message}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label>Descrição</Label>
                <Textarea {...form.register('description')} />
              </div>
              <div className="space-y-2">
                <Label>Status</Label>
                <Select
                  value={form.watch('status')}
                  onValueChange={(value) => form.setValue('status', value as TaskForm['status'])}
                >
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    {Object.entries(statusLabels).map(([value, label]) => (
                      <SelectItem key={value} value={value}>{label}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <Button type="submit" className="w-full" disabled={createMutation.isPending || updateMutation.isPending}>
                {editingTask ? 'Salvar' : 'Criar'}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <SearchInput value={search} onChange={setSearch} placeholder="Pesquisar tarefas..." />

      <DataGrid rowData={data?.data ?? []} columnDefs={columnDefs} loading={isLoading} />
    </div>
  )
}
