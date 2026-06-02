# ---- build stage ----
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm install
RUN npm ls @fontsource-variable/geist
RUN ls -la node_modules/@fontsource-variable
COPY . .
RUN npm run build

FROM docker.io/library/nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

# Route all paths to index.html for SPA routing
RUN printf 'server {\n\
    listen 80;\n\
    root /usr/share/nginx/html;\n\
    index index.html;\n\
    location / {\n\
        try_files $uri $uri/ /index.html;\n\
    }\n\
}\n' > /etc/nginx/conf.d/default.conf

EXPOSE 80
