FROM php:8.1-fpm-alpine AS build

RUN apk --no-cache update \
    && apk --no-cache upgrade \
    && apk --no-cache add ${PHPIZE_DEPS} \
    && docker-php-ext-install mysqli pdo pdo_mysql \
    && pecl install redis-5.3.7 \
    && docker-php-ext-enable pdo_mysql redis \
    && apk del ${PHPIZE_DEPS}

WORKDIR /app

# Development image
FROM build AS dev

RUN apk --no-cache update \
    && apk --no-cache upgrade \
    && apk --no-cache add ${PHPIZE_DEPS} \
    && pecl install xdebug-3.1.5 \
    && docker-php-ext-enable xdebug \
    && apk del ${PHPIZE_DEPS}

RUN echo "xdebug.client_host='host.docker.internal'" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && echo "xdebug.mode=debug" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && echo "xdebug.start_with_request=yes" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini

# Production image
FROM build as prod

COPY . .
