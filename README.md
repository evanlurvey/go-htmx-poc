## Objective

I just wanted to learn go templating and htmx. Ended up exploring some other ideas including:

- automatic browser reload
- layouts / templates
- components which are just functions
- automatic csrf
- automatic forms
- validations

## TODO

- actually use tailwind generated css
- menu bar based on login state
- header with name
- edit name menu

##### Infrastructure

k3s provisioned [with](flux https://fluxcd.io/flux/get-started/)

##### Running

```console
# Auto reloading local dev
go install github.com/cosmtrek/air@latest  # install if needed
air

# Build with rancher desktop / nerdctl
nerdctl -n k8s.io build -t ghcr.io/evanlurvey/htmx-poc .
```

##### Tailwind

```console
# watch
tailwindcss -o app/web/static/style.css -w
# prod / ci
tailwindcss -o app/web/static/style.css -m
```
