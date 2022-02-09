# Accountant

This project is in its very early stages meaning breaking changes are going to be introduced sooner rather than later.
You should be prepared for a complete data wipe due to the schema changes since a safe migration process is currently
not implemented.

## Usage

This project requires a Docker and Docker Compose installation.
A script for running it is present in the _bin_ folder.
You can familiarize yourself with it by running it with the _-h_ flag.

Please read the Notes bellow before running the app.

### Production

Substitute the domains in the `nginx/config/app_prod.nginx` with your own. 

Run these commands:
```shell
sudo chmod u+x bin/accountant.sh
# Creates the secrets folder with required credential files.
./bin/accountant.sh secrets --production
# Generates a Let's encrypt certificate.
./bin/accountant.sh cert -d <your.domain> -e <email@your.domain>
./bin/accountant.sh up --production
# When setting up for the first time:
docker exec -it accountant_app ./accountant-cmd createTables
```

If you require HTTP Basic Auth, you can set it up with:
```shell
# Generates a .htpasswd file.
./bin/accountant.sh secure
```
You also have to uncomment:
* `docker-compose-prod.yml`: build arguments and auth volume
* `nginx/Dockerfile`: the `RUN` command

### Development

Add the _accountant.test_ domain to your `/etc/hosts` file.

```shell
sudo chmod u+x bin/accountant.sh
./bin/accountant.sh secrets
./bin/accountant.sh up
# When setting up for the first time:
docker exec -it accountant_app ./accountant-cmd createTables
```

## Notes

This application doesn't handle errors gracefully, so be prepared to see some non-formatted error messages.

On that note, a category cannot be removed if it's linked to a transaction or transaction template.
In order to remove a category no transaction or transaction template can use it.

When a schema changes you have to remove existing data or manually change it.
The easiest way to remove it is by removing the Docker volume with `docker volume rm accountant_db_vol`.
You can manually alter the tables by connecting to the DB container `docker exec -it accountant_database /bin/sh`,
where you can log in to the MySQL DB and alter the schema.
You can find your DB credentials in the _secrets_ folder.
