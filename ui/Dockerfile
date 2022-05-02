FROM node:16.14 as builder
LABEL org.opencontainers.image.source=https://github.com/opensource-f2f/open-podcasts
WORKDIR /usr/src/app

COPY config/nginx.conf ./
COPY public ./public
COPY src ./src
COPY static ./static
COPY package.json compare.js ./

RUN npm install
RUN npm run build

FROM nginx:1.17.1-alpine
EXPOSE 80
COPY --from=builder /usr/src/app/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /usr/src/app/build /usr/share/nginx/html
