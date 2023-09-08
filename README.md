I just updated the csrf lib to use its own middleware and cookies so
we can run sessions on their own separately (like logging in and out)

## TODO

- create sessions for user
- menu bar based on login state
- header with name
- edit name menu

##### Tailwind

```console
# watch
tailwindcss -o app/static/style.css -w
# prod / ci
tailwindcss -o app/static/style.css -m
```
