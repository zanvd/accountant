# Accountant

This project is in its very early stages meaning breaking changes are going to be introduced sooner rather than later.

## Usage

This project requires a Docker and Docker Compose installation.
A script for running it is present in the _bin_ folder.
You can familiarize yourself with it by running it with the _-h_ flag.

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
