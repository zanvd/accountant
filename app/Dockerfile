# Prepared with help from: https://github.com/dunglas/symfony-docker/blob/main/Dockerfile

# Base image is also production.
FROM php:8.2-fpm-alpine AS base

ENV COMPOSER_ALLOW_SUPERUSER=1

COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

RUN mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"

RUN apk --no-cache update \
    && apk --no-cache upgrade \
    && apk add --update icu \
    && apk --no-cache add --virtual .build-deps icu-dev ${PHPIZE_DEPS} \
    && docker-php-ext-configure intl \
    && docker-php-ext-install intl mysqli pdo pdo_mysql \
    && pecl install redis-5.3.7 \
    && docker-php-ext-enable intl pdo_mysql redis \
    && { find /usr/local/lib -type f -print0 | xargs -0r strip --strip-all -p 2>/dev/null || true; } \
    && apk del .build-deps

WORKDIR /app

COPY composer.* symfony.* ./

RUN composer install --prefer-dist --no-autoloader --no-scripts --no-progress --no-dev \
    && composer clear-cache

COPY . .

RUN mkdir -p var/cache var/log \
    && composer dump-autoload --classmap-authoritative --no-dev \
    && composer dump-env prod \
    && composer run-script --no-dev post-install-cmd \
    && chmod +x bin/console \
    && sync

# Development image
FROM base AS dev

RUN mv "$PHP_INI_DIR/php.ini" "$PHP_INI_DIR/php.ini-production" \
    && mv "$PHP_INI_DIR/php.ini-development" "$PHP_INI_DIR/php.ini"

RUN apk --no-cache update \
    && apk --no-cache upgrade \
    && apk --no-cache add --virtual .build-deps linux-headers ${PHPIZE_DEPS} \
    && pecl install xdebug-3.2.0 \
    && docker-php-ext-enable xdebug \
    && apk del .build-deps

RUN echo "xdebug.client_host='host.docker.internal'" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && echo "xdebug.idekey=PHPSTORM" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && echo "xdebug.mode=debug" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && composer dump-env dev \
    && composer install --no-progress \
    && composer clear-cache \
    && sync
