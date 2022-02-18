FROM node:16
LABEL org.opencontainers.image.source https://github.com/linuxsuren/goplay
WORKDIR /usr/src/app

COPY package*.json ./
COPY . .

ENV NODE_ENV=production
RUN npm install

EXPOSE 3000
CMD [ "npm", "run", "start1" ]
