FROM node:16.14 as builder
LABEL org.opencontainers.image.source=https://github.com/linuxsuren/open-podcasts
WORKDIR /usr/src/app

COPY crds ./crds
COPY public ./public
COPY src ./src
COPY static ./static
COPY package.json index.js ./

RUN npm install
RUN npm run build

EXPOSE 5000
CMD [ "node", "index.js" ]