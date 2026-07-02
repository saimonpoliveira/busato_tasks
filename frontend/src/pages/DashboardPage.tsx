import { FolderKanban, Ticket, ListTodo, AlertCircle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { LoadingSpinner } from '@/components/common/LoadingSpinner'
import { useDashboardStats } from '@/hooks/useApi'

const statCards = [
  { key: 'projects' as const, label: 'Projetos', icon: FolderKanban, color: 'text-blue-500' },
  { key: 'tickets' as const, label: 'Chamados', icon: Ticket, color: 'text-green-500' },
  { key: 'tasks' as const, label: 'Tarefas', icon: ListTodo, color: 'text-purple-500' },
  { key: 'openTickets' as const, label: 'Chamados Abertos', icon: AlertCircle, color: 'text-orange-500' },
]

export function DashboardPage() {
  const { data: stats, isLoading } = useDashboardStats()

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
        <p className="text-muted-foreground">Visão geral do sistema de gestão de chamados e tarefas</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {statCards.map((card) => (
          <Card key={card.key}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{card.label}</CardTitle>
              <card.icon className={`h-4 w-4 ${card.color}`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats?.[card.key] ?? 0}</div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
