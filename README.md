![Group 1](https://grafana.com/static/assets/img/blog/grafana-labs-lgtm-graphic.png)
# Observability for simple Microservices written in Golang

Welcome!

## 🛠 Pre-requisites

Before diving in, ensure you have the following:

1. **`psql`**: Not installed? 🛑 Installation both for Mac and Ubuntu Follow the guide [here](https://www.timescale.com/blog/how-to-install-psql-on-mac-ubuntu-debian-windows/).

2.  Execute the following commands to set up your database:
```bash
docker-compose up -d postgres
psql -h localhost -p 5432 -U postgres -c "CREATE DATABASE users;"
psql -h localhost -p 5432 -U postgres -d users -a -f ./resources/users.sql
```

### How to run
```bash
make up_build
```

## 🌐 UI Access

To see the UI in action, navigate to:
```
127.0.0.1:80
```
#### Once you're there, click on "Test Auth" and "Test Broker" multiple times to generate events. Enjoy the interactivity!


## 📊 Access observability stack

To see logs, metrics, traces and profiles, open Grafana UI at:

```
http://localhost:3000

``` 