FROM certbot/certbot

COPY cron-conf ./conf

RUN cat ./conf/cert-renew >> /var/spool/cron/crontabs/root \
    && rm -r ./conf

# Clear the entrypoint since we want to run a cronjob.
ENTRYPOINT []
CMD ["crond", "-f", "-l", "0", "-d", "0"]
