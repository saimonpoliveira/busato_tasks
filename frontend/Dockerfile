# Single-stage build (mais estável no Railway que multi-stage)
FROM node:20-bookworm-slim

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

RUN npm install -g serve

ENV PORT=3000

EXPOSE 3000

CMD ["sh", "-c", "serve -s dist -l ${PORT}"]
