import { type ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth } from '@/contexts/AuthContext'
import { LoadingSpinner } from '@/components/common/LoadingSpinner'

interface ProtectedRouteProps {
  children: ReactNode
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated } = useAuth()

  if (!localStorage.getItem('token') && !isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

export function PublicRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated } = useAuth()

  if (isAuthenticated || localStorage.getItem('token')) {
    return <Navigate to="/dashboard" replace />
  }

  return <>{children}</>
}

export function PageLoader() {
  return (
    <div className="flex h-screen items-center justify-center">
      <LoadingSpinner size="lg" />
    </div>
  )
}
