# Busato Tasks

Sistema completo de gestão de chamados e tarefas para desenvolvimento de software.

## Stack

### Backend
- Go + Gin
- GORM + PostgreSQL (Neon)
- JWT Authentication
- Clean Architecture (controllers, services, repositories, models, dto, validators, routes, middlewares, utils)

### Frontend
- React + TypeScript + Vite
- TanStack Query + React Router
- Tailwind CSS + shadcn/ui
- AG Grid Community
- React Hook Form + Zod

## Estrutura do Projeto

```
backend/          # API REST em Go
frontend/         # SPA em React
```

## Configuração

### Backend

1. Copie o arquivo de ambiente:
```bash
cp backend/.env.example backend/.env
```

2. Configure as variáveis no `.env`:
- `DATABASE_URL` - Connection string do PostgreSQL (Neon)
- `JWT_SECRET` - Chave secreta para tokens JWT

3. Execute o servidor:
```bash
cd backend
go run cmd/server/main.go
```

O servidor inicia em `http://localhost:8080` com migrations automáticas.

### Frontend

1. Copie o arquivo de ambiente:
```bash
cp frontend/.env.example frontend/.env
```

2. Instale dependências e execute:
```bash
cd frontend
npm install
npm run dev
```

O frontend inicia em `http://localhost:5173`.

## API Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/register` | Registro |
| GET | `/api/v1/me` | Usuário autenticado |
| CRUD | `/api/v1/users` | Usuários |
| CRUD | `/api/v1/projects` | Projetos |
| CRUD | `/api/v1/tickets` | Chamados |
| CRUD | `/api/v1/tasks` | Tarefas |
| CRUD | `/api/v1/comments` | Comentários |
| CRUD | `/api/v1/attachments` | Anexos |

Todos os endpoints de listagem suportam paginação, filtros, ordenação e pesquisa via query params:
- `page`, `page_size`, `sort_by`, `sort_order`, `search`

## Deploy no Railway

Guia completo em [`docs/RAILWAY.md`](docs/RAILWAY.md).

Resumo rápido:
1. Crie um projeto Railway com 3 serviços: **PostgreSQL**, **backend** (`backend/`), **frontend** (`frontend/`)
2. Backend: `DATABASE_URL`, `JWT_SECRET`, `CORS_ORIGINS`
3. Frontend: `VITE_API_URL=https://seu-backend.up.railway.app/api/v1`
4. Gere domínios públicos para backend e frontend

Após o deploy, acesse pelo navegador ou celular Android pela URL pública do frontend.

## Funcionalidades

- Autenticação JWT com registro e login
- CRUD completo de usuários, projetos, chamados, tarefas e comentários
- Upload e download de anexos
- Dashboard com estatísticas
- Listagens com AG Grid (filtros, ordenação, paginação)
- Formulários validados com React Hook Form + Zod
- Tema claro e escuro
