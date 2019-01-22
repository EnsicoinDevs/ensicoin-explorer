FROM alpine:latest

RUN apk add --update nginx nodejs npm

RUN mkdir -p /tmp/nginx/app
RUN mkdir -p /var/log/nginx
RUN mkdir -p /var/www/html

COPY nginx_config/nginx.conf /etc/nginx/nginx.conf
COPY nginx_config/default.conf /etc/nginx/conf.d/default.conf

WORKDIR /tmp/nginx/app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

RUN cp -r dist/* /var/www/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
