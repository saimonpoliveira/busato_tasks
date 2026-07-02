import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  FolderKanban,
  Ticket,
  ListTodo,
  Users,
  MessageSquare,
} from 'lucide-react'
import { cn } from '@/lib/utils'

const navItems = [
  { to: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/projects', label: 'Projetos', icon: FolderKanban },
  { to: '/tickets', label: 'Chamados', icon: Ticket },
  { to: '/tasks', label: 'Tarefas', icon: ListTodo },
  { to: '/users', label: 'Usuários', icon: Users },
  { to: '/comments', label: 'Comentários', icon: MessageSquare },
]

export function Sidebar() {
  return (
    <aside className="flex h-full w-64 flex-col border-r border-sidebar-border bg-sidebar">
      <div className="flex h-16 items-center border-b border-sidebar-border px-6">
        <h1 className="text-lg font-bold text-sidebar-foreground">Busato Tasks</h1>
      </div>
      <nav className="flex-1 space-y-1 p-4">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
                isActive
                  ? 'bg-sidebar-accent text-sidebar-foreground'
                  : 'text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-foreground'
              )
            }
          >
            <item.icon className="h-4 w-4" />
            {item.label}
          </NavLink>
        ))}
      </nav>
    </aside>
  )
}
