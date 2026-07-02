export type UserRole = 'admin' | 'member'

export interface User {
  id: string
  email: string
  name: string
  role: UserRole
  active: boolean
  created_at: string
  updated_at: string
}

export type ProjectStatus = 'active' | 'archived' | 'completed'

export interface Project {
  id: string
  name: string
  description: string
  status: ProjectStatus
  owner_id: string
  owner?: User
  created_at: string
  updated_at: string
}

export type TicketStatus = 'open' | 'in_progress' | 'resolved' | 'closed'
export type TicketPriority = 'low' | 'medium' | 'high' | 'critical'

export interface Ticket {
  id: string
  project_id: string
  title: string
  description: string
  status: TicketStatus
  priority: TicketPriority
  assignee_id?: string
  assignee?: User
  reporter_id: string
  reporter?: User
  created_at: string
  updated_at: string
}

export type TaskStatus = 'todo' | 'in_progress' | 'done' | 'cancelled'

export interface Task {
  id: string
  ticket_id: string
  title: string
  description: string
  status: TaskStatus
  assignee_id?: string
  assignee?: User
  order: number
  created_at: string
  updated_at: string
}

export type EntityType = 'ticket' | 'task'

export interface Comment {
  id: string
  entity_type: EntityType
  entity_id: string
  user_id: string
  user?: User
  content: string
  created_at: string
  updated_at: string
}

export interface Attachment {
  id: string
  entity_type: EntityType
  entity_id: string
  filename: string
  original_name: string
  size: number
  mime_type: string
  uploaded_by_id: string
  uploaded_by?: User
  created_at: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface PaginationParams {
  page?: number
  page_size?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
  search?: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface ApiError {
  error: string
  details?: Record<string, string>
}

export interface DashboardStats {
  projects: number
  tickets: number
  tasks: number
  openTickets: number
}
