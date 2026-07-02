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
import {
  useTickets,
  useCreateTicket,
  useUpdateTicket,
  useDeleteTicket,
  useProjects,
} from '@/hooks/useApi'
import { getErrorMessage } from '@/lib/api'
import type { Ticket } from '@/types'

const ticketSchema = z.object({
  project_id: z.string().min(1, 'Projeto é obrigatório'),
  title: z.string().min(2, 'Título deve ter no mínimo 2 caracteres'),
  description: z.string().optional(),
  status: z.enum(['open', 'in_progress', 'resolved', 'closed']).optional(),
  priority: z.enum(['low', 'medium', 'high', 'critical']).optional(),
})

type TicketForm = z.infer<typeof ticketSchema>

const statusLabels: Record<string, string> = {
  open: 'Aberto',
  in_progress: 'Em Progresso',
  resolved: 'Resolvido',
  closed: 'Fechado',
}

const priorityLabels: Record<string, string> = {
  low: 'Baixa',
  medium: 'Média',
  high: 'Alta',
  critical: 'Crítica',
}

export function TicketsPage() {
  const [search, setSearch] = useState('')
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingTicket, setEditingTicket] = useState<Ticket | null>(null)

  const { data, isLoading } = useTickets({ search, page_size: 100 })
  const { data: projectsData } = useProjects({ page_size: 100 })
  const createMutation = useCreateTicket()
  const updateMutation = useUpdateTicket()
  const deleteMutation = useDeleteTicket()

  const form = useForm<TicketForm>({
    resolver: zodResolver(ticketSchema),
    defaultValues: { project_id: '', title: '', description: '', status: 'open', priority: 'medium' },
  })

  const columnDefs = useMemo<ColDef<Ticket>[]>(
    () => [
      { field: 'title', headerName: 'Título', flex: 2 },
      {
        field: 'status',
        headerName: 'Status',
        valueFormatter: (params) => statusLabels[params.value as string] || params.value,
      },
      {
        field: 'priority',
        headerName: 'Prioridade',
        valueFormatter: (params) => priorityLabels[params.value as string] || params.value,
      },
      {
        field: 'assignee.name',
        headerName: 'Responsável',
        valueGetter: (params) => params.data?.assignee?.name || '-',
      },
      {
        field: 'reporter.name',
        headerName: 'Reportado por',
        valueGetter: (params) => params.data?.reporter?.name || '-',
      },
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
        cellRenderer: (params: { data: Ticket }) => (
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
    setEditingTicket(null)
    form.reset({ project_id: '', title: '', description: '', status: 'open', priority: 'medium' })
    setDialogOpen(true)
  }

  const openEdit = (ticket: Ticket) => {
    setEditingTicket(ticket)
    form.reset({
      project_id: ticket.project_id,
      title: ticket.title,
      description: ticket.description,
      status: ticket.status,
      priority: ticket.priority,
    })
    setDialogOpen(true)
  }

  const onSubmit = async (formData: TicketForm) => {
    try {
      if (editingTicket) {
        await updateMutation.mutateAsync({ id: editingTicket.id, data: formData })
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
          <h2 className="text-3xl font-bold tracking-tight">Chamados</h2>
          <p className="text-muted-foreground">Gerencie os chamados de suporte e desenvolvimento</p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={openCreate}>
              <Plus className="mr-2 h-4 w-4" />
              Novo Chamado
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingTicket ? 'Editar Chamado' : 'Novo Chamado'}</DialogTitle>
            </DialogHeader>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <div className="space-y-2">
                <Label>Projeto</Label>
                <Select
                  value={form.watch('project_id')}
                  onValueChange={(value) => form.setValue('project_id', value)}
                  disabled={!!editingTicket}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione um projeto" />
                  </SelectTrigger>
                  <SelectContent>
                    {projectsData?.data.map((project) => (
                      <SelectItem key={project.id} value={project.id}>
                        {project.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {form.formState.errors.project_id && (
                  <p className="text-sm text-destructive">{form.formState.errors.project_id.message}</p>
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
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Status</Label>
                  <Select
                    value={form.watch('status')}
                    onValueChange={(value) => form.setValue('status', value as TicketForm['status'])}
                  >
                    <SelectTrigger><SelectValue /></SelectTrigger>
                    <SelectContent>
                      {Object.entries(statusLabels).map(([value, label]) => (
                        <SelectItem key={value} value={value}>{label}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Prioridade</Label>
                  <Select
                    value={form.watch('priority')}
                    onValueChange={(value) => form.setValue('priority', value as TicketForm['priority'])}
                  >
                    <SelectTrigger><SelectValue /></SelectTrigger>
                    <SelectContent>
                      {Object.entries(priorityLabels).map(([value, label]) => (
                        <SelectItem key={value} value={value}>{label}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <Button type="submit" className="w-full" disabled={createMutation.isPending || updateMutation.isPending}>
                {editingTicket ? 'Salvar' : 'Criar'}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <SearchInput value={search} onChange={setSearch} placeholder="Pesquisar chamados..." />

      <DataGrid rowData={data?.data ?? []} columnDefs={columnDefs} loading={isLoading} />
    </div>
  )
}
