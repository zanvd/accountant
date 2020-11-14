FROM nginx:alpine AS base

RUN rm /etc/nginx/conf.d/default.conf

### Development image. ###
FROM base AS dev

COPY ./config/app_dev.nginx /etc/nginx/conf.d/accountant.conf

### Production image. ###
FROM base AS prod

RUN apk add --update --no-cache apache2-utils \
    && mkdir /etc/nginx/auth

COPY ./config/app_prod.nginx /etc/nginx/conf.d/accountant.conf
