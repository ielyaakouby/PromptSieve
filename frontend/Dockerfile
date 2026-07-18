FROM node:20-alpine
LABEL org.opencontainers.image.title="PromptSieve Web"
LABEL org.opencontainers.image.description="Open-source sensitive data scanner for AI prompts"
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
