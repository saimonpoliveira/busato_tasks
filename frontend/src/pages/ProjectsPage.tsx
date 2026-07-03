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
import { useProjects, useCreateProject, useUpdateProject, useDeleteProject } from '@/hooks/useApi'
import { getErrorMessage } from '@/lib/api'
import type { Project } from '@/types'

const projectSchema = z.object({
  name: z.string().min(2, 'Nome deve ter no mínimo 2 caracteres'),
  description: z.string().optional(),
  status: z.enum(['active', 'archived', 'completed']).optional(),
})

type ProjectForm = z.infer<typeof projectSchema>

export function ProjectsPage() {
  const [search, setSearch] = useState('')
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingProject, setEditingProject] = useState<Project | null>(null)

  const { data, isLoading } = useProjects({ search, page_size: 100 })
  const createMutation = useCreateProject()
  const updateMutation = useUpdateProject()
  const deleteMutation = useDeleteProject()

  const form = useForm<ProjectForm>({
    resolver: zodResolver(projectSchema),
    defaultValues: { name: '', description: '', status: 'active' },
  })

  const columnDefs = useMemo<ColDef<Project>[]>(
    () => [
      { field: 'name', headerName: 'Nome' },
      { field: 'description', headerName: 'Descrição', flex: 2 },
      { field: 'status', headerName: 'Status' },
      {
        field: 'owner.name',
        headerName: 'Proprietário',
        valueGetter: (params) => params.data?.owner?.name || '-',
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
        width: 180,
        cellClass: 'ag-cell-actions',
        cellRenderer: (params: { data: Project }) => (
          <div className="flex h-full items-center gap-2">
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
    setEditingProject(null)
    form.reset({ name: '', description: '', status: 'active' })
    setDialogOpen(true)
  }

  const openEdit = (project: Project) => {
    setEditingProject(project)
    form.reset({
      name: project.name,
      description: project.description,
      status: project.status,
    })
    setDialogOpen(true)
  }

  const onSubmit = async (formData: ProjectForm) => {
    try {
      if (editingProject) {
        await updateMutation.mutateAsync({ id: editingProject.id, data: formData })
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
          <h2 className="text-3xl font-bold tracking-tight">Projetos</h2>
          <p className="text-muted-foreground">Gerencie os projetos de desenvolvimento</p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={openCreate}>
              <Plus className="mr-2 h-4 w-4" />
              Novo Projeto
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingProject ? 'Editar Projeto' : 'Novo Projeto'}</DialogTitle>
            </DialogHeader>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <div className="space-y-2">
                <Label>Nome</Label>
                <Input {...form.register('name')} />
                {form.formState.errors.name && (
                  <p className="text-sm text-destructive">{form.formState.errors.name.message}</p>
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
                  onValueChange={(value) => form.setValue('status', value as ProjectForm['status'])}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="active">Ativo</SelectItem>
                    <SelectItem value="archived">Arquivado</SelectItem>
                    <SelectItem value="completed">Concluído</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <Button type="submit" className="w-full" disabled={createMutation.isPending || updateMutation.isPending}>
                {editingProject ? 'Salvar' : 'Criar'}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <SearchInput value={search} onChange={setSearch} placeholder="Pesquisar projetos..." />

      <DataGrid rowData={data?.data ?? []} columnDefs={columnDefs} loading={isLoading} />
    </div>
  )
}
