# Build
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

# Run
FROM nginx:1.27-alpine

RUN apk add --no-cache nodejs npm

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY start.sh /start.sh
RUN chmod +x /start.sh

COPY --from=builder /app/dist /usr/share/nginx/html

ENV PORT=80

EXPOSE 80

CMD ["/start.sh"]
