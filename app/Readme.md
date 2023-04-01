Creating a migration and applying it:
```shell
bin/console make:migration
bin/console doctrine:migrations:migrate
```

### Tailwind CSS:

Run these commands from inside the `app`.

Rebuild:
```shell
docker run --rm -v ${PWD}:/var/www -w /var/www node:19-alpine npx tailwindcss -i ./tailwind-input.css -o ./public/assets/output.css
```

If not yet initialized:
```shell
docker run --rm -v ${PWD}:/var/www -w /var/www node:19-alpine npm install -D tailwindcss @tailwindcss/forms
sudo chown -R <user>:<group> node_modules

# Below commands are just for reference of the initial setup and are no longer needed.
sudo chown <user>:<group> package*
docker run --rm -v ${PWD}:/var/www -w /var/www node:19-alpine npx tailwindcss init
sudo chown <user>:<group> tailwind.config.js
```
