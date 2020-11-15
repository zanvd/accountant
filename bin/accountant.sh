#!/bin/bash
#### Output colors.
BACKGROUND_RED='\033[0;41m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[0;33m'
OUTPUT_RESET='\033[0m'

#### Settings.
LIVE_PATH=/etc/letsencrypt/live

#### Variables.
compose_cmd='docker-compose -f docker-compose.yml -f docker-compose-dev.yml'

#### Functions.
function build_images() {
  echo -e "$COLOR_YELLOW>> Building images.$OUTPUT_RESET"
  ${compose_cmd} build
}

function check_production_flag() {
  ! [[ $1 == "--production" ]]
  return $?
}

function create_secrets_folder() {
  if [[ -d "secrets" ]]; then
    echo -e "$COLOR_YELLOW>> Warning! This is going to overwrite your current database credentials!$OUTPUT_RESET"
    echo '>> You already have a secrets folder. Are you sure you want to proceed?'

    select answer in "Yes" "No"; do
      case ${answer} in
      Yes) break ;;
      No) return ;;
      esac
    done
  else
    echo -e "$COLOR_YELLOW>> Creating the secrets folder.$OUTPUT_RESET"
    mkdir secrets
  fi

  db_database='accountant'
  db_password="$(random_string)"
  db_user='accountant'
  db_root_password="$(random_string)"

  # Production credentials ought to be set on the fly!
  if (($1)); then
    answer=''
    read -p "Enter a name for the API database (default '${db_database}'):" answer
    if [[ ${answer} != "" ]]; then
      db_database=${answer}
    else
      echo -e "$BACKGROUND_RED>> You are using the default value at your own risk!$OUTPUT_RESET"
    fi

    read -p "Enter a username for the API database (default '${db_user}'):" answer
    if [[ ${answer} != "" ]]; then
      db_user=${answer}
    else
      echo -e "$BACKGROUND_RED>> You are using the default value at your own risk!$OUTPUT_RESET"
    fi
  fi

  write_secret_files \
    db_database "$db_database" \
    db_password "$db_password" \
    db_user "$db_user" \
    db_root_password "$db_root_password"
}

function down() {
  echo -e "$COLOR_YELLOW>> Stopping services, removing containers and networks.$OUTPUT_RESET"
  ${compose_cmd} down
}

# Argument passed to this function is the result of the check_production_flag.
function include_production_compose() {
  compose_cmd='docker-compose -f docker-compose.yml -f docker-compose-prod.yml'
}

function initialize_nginx_and_certbot() {
  domain=$1
  email=$2
  staging=$3

  CERT_PATH="$LIVE_PATH/$domain"

  echo -e "$COLOR_YELLOW>> Generating dummy certificates.$OUTPUT_RESET"
  ${compose_cmd} run --rm --entrypoint "sh -c '\
    mkdir -p $CERT_PATH \
    && openssl req -x509 -nodes -newkey rsa:1024 -days 1 \
    -keyout $CERT_PATH/privkey.pem \
    -out $CERT_PATH/fullchain.pem \
    -subj /CN=localhost'" cert
  echo -e "$COLOR_GREEN>> Dummy certificate generated.$OUTPUT_RESET"

  echo -e "$COLOR_YELLOW>> Starting nginx.$OUTPUT_RESET"
  ${compose_cmd} up --force-recreate -d nginx

  echo -e "$COLOR_YELLOW>> Deleting dummy certificates.$OUTPUT_RESET"
  ${compose_cmd} run --rm --entrypoint "\
    rm -rf $CERT_PATH \
    && rm -rf /etc/letsencrypt/archive/$domain \
    && rm -rf /etc/letsencrypt/renewal/$domain.conf" cert

  echo -e "$COLOR_YELLOW>> Generating live certificates.$OUTPUT_RESET"
  ${compose_cmd} run --rm --entrypoint "\
    certbot certonly --webroot -w /var/www/cert \
    $staging \
    --email=$email \
    -d $domain \
    -d www.$domain \
    --agree-tos \
    --non-interactive \
    --force-renewal" cert

  stop
  ${compose_cmd} rm -f
}

function random_string() {
  cat /dev/urandom | tr -dc 'a-zA-Z0-9!@#$%&\*()+=_\-' | fold -w ${1:-32} | head -n 1
}

function secure_nginx_with_htpasswd() {
  user=''
  password=''

  read -p 'Enter a username for the HTTP Basic Authentication:' answer
  if [[ ${answer} != '' ]]; then
    user=${answer}
  else
    echo -e "$BACKGROUND_RED>> Empty username is not allowed! Rerun the command.$OUTPUT_RESET"
    exit 1
  fi

  read -p 'Enter a password for the HTTP Basic Authentication:' answer
  if [[ ${answer} != '' ]]; then
    password=${answer}
  else
    echo -e "$BACKGROUND_RED>> Empty password is not allowed! Rerun the command.$OUTPUT_RESET"
    exit 1
  fi

  ${compose_cmd} run --rm --entrypoint "htpasswd -Bbc /etc/nginx/auth/.htpasswd $user $password" nginx

  stop
  ${compose_cmd} rm -f
}

function stop() {
  echo -e "$COLOR_YELLOW>> Stopping services.$OUTPUT_RESET"
  ${compose_cmd} stop
}

function up() {
  echo -e "$COLOR_YELLOW>> Creating container and networks and starting services.$OUTPUT_RESET"
  ${compose_cmd} up -d
}

function usage() {
  echo "A script for managing project's Docker services.

Usage:
  accountant.sh (down|secrets-folder|stop) [--production]
  accountant.sh up [--production]
  accountant.sh cert (-d|--domain) (-e|--email) [--staging]
  accountant.sh -h|--help

Options:
  -h, --help  Displays this help message.

Commands:;
  down    Brings the services to a halt and remove container and networks.
            --production Runs the command in the production mode.
  cert    Creates a valid Let's encrypt certificate. Always runs in production mode.
          This command has two required arguments (domain and email) and one option (staging).
            -d, --domain  Fully qualified domain name used for the certificate. The certificate is going to be built
                          for the 'www' subdomain, as well, so you should leave it out.
            -e, --email   Email for the certificate.
            --staging     Use this when testing certificate generation in order to avoid hitting the Let's encrypt's limit.
  secrets Sets up the secrets folder with the files required by the database.
            --production  When generating the secrets it prompts you to input database name and username.
  secure  Creates the .htpasswd file for the HTTP Basic Authentication. Always runs in production mode.
  stop    Halts the services.
            --production  Runs the command in the production mode.
  up      Builds and runs the services. Make sure to run the 'cert' command in production environment before the first up.
            --production  Runs the command in the production mode.

Examples:
  accountant.sh cert --domain example.org --email info@example.org --staging
  accountant.sh up
  accountant.sh stop --production
  accountant.sh down
  "
}

function write_secret_files() {
  while [[ "$1" != "" ]]; do
    file_name=$1
    shift
    value=$1
    shift

    echo -n "${value}" >secrets/"${file_name}".txt
  done
}

#### Main.
if [[ $# -eq 0 ]]; then
  usage
  exit 1
fi

if ! [[ -x "$(command -v docker-compose)" ]]; then
  echo -e "$BACKGROUND_RED>> Command not found: docker-compose$OUTPUT_RESET"
  exit 1
fi

if ! [[ -f "docker-compose.yml" ]]; then
  echo -e "$BACKGROUND_RED>> This script has to be executed from where the docker-compose.yml resides.$OUTPUT_RESET"
  echo '>> Execute it with ./bin/accountant.sh'
  exit 1
fi

case $1 in
cert)
  domain=''
  email=''
  staging=''

  shift
  while [[ "$1" != "" ]]; do
    case $1 in
    -d | --domain)
      shift
      domain="$1"
      ;;
    -e | --email)
      shift
      email="$1"
      ;;
    --staging)
      staging='--staging'
      ;;
    esac
    shift
  done

  if [[ "$domain" == "" ]] || [[ "$email" == "" ]]; then
    echo -e "$BACKGROUND_RED>> Options '--domain' and '--email' are required for the 'cert' command.$OUTPUT_RESET"
    exit 1
  fi

  include_production_compose
  build_images
  initialize_nginx_and_certbot "$domain" "$email" "$staging"
  ;;
down)
  shift
  check_production_flag "$1"
  (($?)) && include_production_compose
  echo "${compose_cmd}"
  down
  ;;
-h | --help)
  usage
  exit 0
  ;;
secrets)
  shift
  check_production_flag "$1"
  create_secrets_folder $?
  ;;
secure)
  include_production_compose
  build_images
  secure_nginx_with_htpasswd
  ;;
stop)
  shift
  check_production_flag "$1"
  (($?)) && include_production_compose
  stop
  ;;
up)
  shift
  check_production_flag "$1"
  (($?)) && include_production_compose
  build_images
  up
  ;;
*)
  usage
  exit 1
  ;;
esac

echo -e "$COLOR_GREEN>> Done. <<$OUTPUT_RESET"

exit 0
