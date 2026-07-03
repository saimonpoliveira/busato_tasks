# Build
FROM node:20-bookworm-slim AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

# Run
FROM node:20-bookworm-slim

WORKDIR /app

RUN npm install -g serve

COPY --from=builder /app/dist ./dist
COPY start.sh ./start.sh
RUN chmod +x ./start.sh

ENV PORT=3000

EXPOSE 3000

CMD ["./start.sh"]
