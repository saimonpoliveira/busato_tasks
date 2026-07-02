import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'
import type {
  AuthResponse,
  Comment,
  PaginatedResponse,
  PaginationParams,
  Project,
  Task,
  Ticket,
  User,
} from '@/types'

export const queryKeys = {
  me: ['me'] as const,
  users: (params?: PaginationParams & { role?: string }) => ['users', params] as const,
  projects: (params?: PaginationParams & { status?: string }) => ['projects', params] as const,
  tickets: (params?: PaginationParams & { project_id?: string; status?: string; priority?: string; assignee_id?: string }) => ['tickets', params] as const,
  tasks: (params?: PaginationParams & { ticket_id?: string; status?: string; assignee_id?: string }) => ['tasks', params] as const,
  comments: (params?: PaginationParams & { entity_type?: string; entity_id?: string; user_id?: string }) => ['comments', params] as const,
}

export function useLogin() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: { email: string; password: string }) => {
      const response = await api.post<AuthResponse>('/auth/login', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.me })
    },
  })
}

export function useRegister() {
  return useMutation({
    mutationFn: async (data: { email: string; password: string; name: string }) => {
      const response = await api.post<AuthResponse>('/auth/register', data)
      return response.data
    },
  })
}

export function useMe() {
  return useQuery({
    queryKey: queryKeys.me,
    queryFn: async () => {
      const response = await api.get<User>('/me')
      return response.data
    },
    retry: false,
  })
}

export function useUsers(params?: PaginationParams & { role?: string }) {
  return useQuery({
    queryKey: queryKeys.users(params),
    queryFn: async () => {
      const response = await api.get<PaginatedResponse<User>>('/users', { params })
      return response.data
    },
  })
}

export function useCreateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: { email: string; password: string; name: string; role?: string }) => {
      const response = await api.post<User>('/users', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
  })
}

export function useUpdateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: Partial<User> & { password?: string } }) => {
      const response = await api.put<User>(`/users/${id}`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
  })
}

export function useDeleteUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/users/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
  })
}

export function useProjects(params?: PaginationParams & { status?: string }) {
  return useQuery({
    queryKey: queryKeys.projects(params),
    queryFn: async () => {
      const response = await api.get<PaginatedResponse<Project>>('/projects', { params })
      return response.data
    },
  })
}

export function useCreateProject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: { name: string; description?: string; status?: string }) => {
      const response = await api.post<Project>('/projects', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}

export function useUpdateProject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: Partial<Project> }) => {
      const response = await api.put<Project>(`/projects/${id}`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}

export function useDeleteProject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/projects/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}

export function useTickets(params?: PaginationParams & { project_id?: string; status?: string; priority?: string; assignee_id?: string }) {
  return useQuery({
    queryKey: queryKeys.tickets(params),
    queryFn: async () => {
      const response = await api.get<PaginatedResponse<Ticket>>('/tickets', { params })
      return response.data
    },
  })
}

export function useCreateTicket() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: {
      project_id: string
      title: string
      description?: string
      status?: string
      priority?: string
      assignee_id?: string
    }) => {
      const response = await api.post<Ticket>('/tickets', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tickets'] })
    },
  })
}

export function useUpdateTicket() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: Partial<Ticket> }) => {
      const response = await api.put<Ticket>(`/tickets/${id}`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tickets'] })
    },
  })
}

export function useDeleteTicket() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/tickets/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tickets'] })
    },
  })
}

export function useTasks(params?: PaginationParams & { ticket_id?: string; status?: string; assignee_id?: string }) {
  return useQuery({
    queryKey: queryKeys.tasks(params),
    queryFn: async () => {
      const response = await api.get<PaginatedResponse<Task>>('/tasks', { params })
      return response.data
    },
  })
}

export function useCreateTask() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: {
      ticket_id: string
      title: string
      description?: string
      status?: string
      assignee_id?: string
    }) => {
      const response = await api.post<Task>('/tasks', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export function useUpdateTask() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: Partial<Task> }) => {
      const response = await api.put<Task>(`/tasks/${id}`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export function useDeleteTask() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/tasks/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export function useComments(params?: PaginationParams & { entity_type?: string; entity_id?: string; user_id?: string }) {
  return useQuery({
    queryKey: queryKeys.comments(params),
    queryFn: async () => {
      const response = await api.get<PaginatedResponse<Comment>>('/comments', { params })
      return response.data
    },
  })
}

export function useCreateComment() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: { entity_type: string; entity_id: string; content: string }) => {
      const response = await api.post<Comment>('/comments', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments'] })
    },
  })
}

export function useDashboardStats() {
  return useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: async () => {
      const [projects, tickets, tasks, openTickets] = await Promise.all([
        api.get<PaginatedResponse<Project>>('/projects', { params: { page_size: 1 } }),
        api.get<PaginatedResponse<Ticket>>('/tickets', { params: { page_size: 1 } }),
        api.get<PaginatedResponse<Task>>('/tasks', { params: { page_size: 1 } }),
        api.get<PaginatedResponse<Ticket>>('/tickets', { params: { page_size: 1, status: 'open' } }),
      ])

      return {
        projects: projects.data.total,
        tickets: tickets.data.total,
        tasks: tasks.data.total,
        openTickets: openTickets.data.total,
      }
    },
  })
}
