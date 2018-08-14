# Closing requests from bots

If request from bot wil return status 403 forbidden

## Running container

```
docker run -p8000:80 smotrovalilit/traefik-auth-botforbidden
```

### In Docker-compose
```
  version: '3.4'
  services:
    traefik:
      image: "traefik"
      container_name: "traefik"
      ports:
        - "80:80"
        - "443:443"
        - "8083:8083"
      restart: always
      volumes:
        - "/var/run/docker.sock:/var/run/docker.sock:ro"
        - "/opt/traefik/traefik.toml:/traefik.toml"
        - "/opt/traefik/acme.json:/acme.json"
      environment:
          DO_AUTH_TOKEN: afb2affe6272e6c5e5505282823ee3cd8fee480c31743c653b526355a276a7
      links:
        - "auth:auth"

    auth:
      image: "smotrovalilit/traefik-auth-botforbidden"
```

## Configuration in traefik
- http
```
[entryPoints]
  [entryPoints.http]
  address = ":80"
  [entryPoints.http.auth.forward]
  address = "http://auth:80/"
```
- https
```
  [entryPoints.https]
  address = ":443"
  [entryPoints.https.auth.forward]
  address = "http://auth:80/"
```

