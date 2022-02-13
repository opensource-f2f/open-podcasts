FROM node:16
LABEL org.opencontainers.image.source https://github.com/linuxsuren/goplay
WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install

COPY . .

EXPOSE 5000
CMD [ "node", "index.js" ]
