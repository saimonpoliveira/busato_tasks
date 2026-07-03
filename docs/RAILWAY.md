# Deploy no Railway

Guia passo a passo para publicar o **Busato Tasks** no [Railway](https://railway.com).

## Arquitetura no Railway

Você vai criar **2 serviços** no Railway (+ banco Neon externo):

| Serviço | Pasta raiz | Descrição |
|---------|------------|-----------|
| Backend | `backend` | API Go |
| Frontend | `frontend` | React (Nginx) |
| Banco | Neon (externo) | PostgreSQL — connection string em `DATABASE_URL` |

---

## ⚠️ Erro: "Railpack could not determine how to build" / "start.sh not found"

Esse erro acontece quando o Railway tenta buildar a **raiz do repositório** em vez da pasta `backend/` ou `frontend/`.

### Correção (faça em cada serviço)

1. Abra o serviço no Railway → **Settings**
2. Em **Source → Root Directory**, defina:
   - Backend: `backend`
   - Frontend: `frontend`
3. Em **Variables**, adicione:
   - `RAILWAY_DOCKERFILE_PATH` = `Dockerfile`
   - `NO_CACHE` = `1` (apenas no primeiro redeploy após corrigir)
4. Clique em **Redeploy**

O log de build deve mostrar: `Using Detected Dockerfile`

> Se ainda falhar, o projeto inclui `start.sh` e `nixpacks.toml` como fallback automático.

---

## Passo 1 — Criar projeto no Railway

1. Acesse [railway.com](https://railway.com) e faça login com GitHub
2. Clique em **New Project**
3. Escolha **Deploy from GitHub repo**
4. Selecione o repositório `busato_tasks`
5. Escolha a branch `cursor/railway-deploy-71ab` (ou `main` após merge)

---

## Passo 2 — Banco de dados (Neon)

Use sua connection string do Neon na variável `DATABASE_URL` do backend.

Não é necessário criar PostgreSQL no Railway se você já usa Neon.

---

## Passo 3 — Deploy do Backend

1. Clique em **+ New → GitHub Repo** (ou **Empty Service** e conecte o repo)
2. Em **Settings → Source**:
   - **Root Directory:** `backend` ← **obrigatório**
3. Em **Settings → Build**:
   - O builder deve usar **Dockerfile** (via `railway.json`)
4. Em **Variables**, configure:

| Variável | Valor |
|----------|-------|
| `DATABASE_URL` | Sua connection string do Neon |
| `JWT_SECRET` | Uma string longa e aleatória |
| `GIN_MODE` | `release` |
| `RAILWAY_DOCKERFILE_PATH` | `Dockerfile` |
| `CORS_ORIGINS` | URL do frontend (configure após o passo 4) |

5. Em **Settings → Networking**, clique em **Generate Domain**
6. Anote a URL pública, ex: `https://busato-api-production.up.railway.app`

### Testar o backend

```bash
curl https://SUA-URL-BACKEND.up.railway.app/health
# Resposta esperada: {"status":"ok"}
```

---

## Passo 4 — Deploy do Frontend

1. Clique em **+ New → GitHub Repo**
2. Em **Settings → Source**:
   - **Root Directory:** `frontend` ← **obrigatório**
3. Em **Variables**, configure:

| Variável | Valor |
|----------|-------|
| `VITE_API_URL` | `https://SUA-URL-BACKEND.up.railway.app/api/v1` |
| `RAILWAY_DOCKERFILE_PATH` | `Dockerfile` |

> **Importante:** `VITE_API_URL` é usada no **build**. Se mudar a URL do backend, faça **Redeploy** do frontend.

4. Em **Settings → Build → Docker Build Args** (se disponível), adicione:
   - `VITE_API_URL` = mesma URL acima

5. Em **Settings → Networking**, clique em **Generate Domain**
6. Anote a URL, ex: `https://busato-web-production.up.railway.app`

---

## Passo 5 — Ajustar CORS no Backend

Volte ao serviço **Backend** e atualize a variável:

```
CORS_ORIGINS=https://SUA-URL-FRONTEND.up.railway.app
```

Salve e aguarde o redeploy automático.

---

## Passo 6 — Acessar o app

Abra a URL do frontend no navegador (ou no celular Android):

```
https://SUA-URL-FRONTEND.up.railway.app
```

Crie uma conta em **Criar conta** e comece a usar.

---

## Variáveis de ambiente — referência completa

### Backend

```env
PORT=8080                    # Railway injeta automaticamente
DATABASE_URL=${{Postgres.DATABASE_URL}}
JWT_SECRET=sua-chave-secreta-longa
JWT_EXPIRATION_HOURS=24
GIN_MODE=release
CORS_ORIGINS=https://seu-frontend.up.railway.app
UPLOAD_DIR=/app/uploads
MAX_UPLOAD_SIZE_MB=10
```

### Frontend

```env
VITE_API_URL=https://seu-backend.up.railway.app/api/v1
```

---

## Acessar pelo celular Android

Depois do deploy, o app fica acessível de qualquer lugar:

1. Abra o Chrome no Android
2. Acesse `https://SUA-URL-FRONTEND.up.railway.app`
3. Opcional: **Menu → Adicionar à tela inicial** para usar como app

---

## Solução de problemas

### Backend não conecta ao banco
- Verifique se `DATABASE_URL` contém a connection string completa do Neon
- A string deve incluir `sslmode=require`

### Erro "Railpack could not determine how to build"
- Confirme **Root Directory** = `backend` ou `frontend` (sem barra no final)
- Adicione `RAILWAY_DOCKERFILE_PATH=Dockerfile`
- Faça redeploy com `NO_CACHE=1`
- Use a branch `cursor/railway-deploy-71ab` (contém os arquivos de deploy)

### Erro de CORS no navegador
- `CORS_ORIGINS` deve ser exatamente a URL do frontend (com `https://`, sem barra no final)

### Frontend não chama a API
- Confirme `VITE_API_URL` com a URL pública do backend + `/api/v1`
- Faça redeploy do frontend após alterar essa variável

### Health check falha
- O endpoint `/health` deve responder `{"status":"ok"}`
- Verifique os logs em **Deployments → View Logs**

### Anexos somem após redeploy
- O filesystem do Railway é efêmero. Para produção, use storage externo (S3, Cloudflare R2). Os anexos funcionam, mas são perdidos em redeploys.

---

## Custos

O Railway oferece créditos gratuitos mensais. Monitore em **Project → Usage**.

Para reduzir custos, use um único serviço Postgres compartilhado e evite múltiplos ambientes de teste.
