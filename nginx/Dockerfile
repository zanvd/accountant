FROM nginx:alpine

RUN rm /etc/nginx/conf.d/default.conf

COPY ./config/app_dev.nginx /etc/nginx/conf.d/app_dev.conf
