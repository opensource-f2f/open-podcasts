FROM node:16 as builder
LABEL org.opencontainers.image.source https://github.com/linuxsuren/goplay
WORKDIR /usr/src/app

COPY package*.json ./
COPY . .

ENV NODE_ENV=production
RUN npm run build

EXPOSE 5000
CMD [ "node", "index.js" ]