# Accountant

This project is in its very early stages meaning breaking changes are going to be introduced sooner rather than later.
You should be prepared for a complete data wipe due to the schema changes since a safe migration process is currently
not implemented.

## Usage

This project requires a Docker and Docker Compose installation.
A script for running it is present in the _bin_ folder.
You can familiarize yourself with it by running it with the _-h_ flag.

Please read the Notes before running the app.

### Production

Substitute the domains in the `nginx/config/app_prod.nginx` with your own. 

Run these commands:
```shell script
sudo chmod u+x bin/accountant.sh
# Creates the secrets folder with required credential files.
./bin/accountant.sh secrets --production
# Generates a Let's encrypt certificate.
./bin/accountant.sh cert -d <your.domain> -e <email@your.domain>
# Generates a .htpasswd file.
./bin/accountant.sh secure
./bin/accountant.sh up --production
```

### Development

Add the _accountant.net_ domain to your `/etc/hosts` file.

```shell script
sudo chmod u+x bin/accountant.sh
./bin/accountant.sh secrets
./bin/accountant.sh up
```

## Notes

This application doesn't handle errors gracefully, so be prepared to see some non-formatted error messages.

On that note, a category cannot be removed if it's linked to a transaction.
In order to remove a category no transaction can use it.

When a schema changes you have to remove existing data or manually change it.
The easiest way to remove it is by removing the Docker volume with `docker volume rm accountant_db_vol`.
You can manually alter the tables by connecting to the DB container `docker exec -it accountant_database /bin/sh`,
where you can log in to the MySQL DB and alter the schema.
You can find your DB credentials in the _secrets_ folder.
