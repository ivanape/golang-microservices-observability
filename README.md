# Observability for simple Microservices written in Golang

Welcome to our awesome microservices project written in Go!

## 🛠 Pre-requisites

Before diving in, ensure you have the following:

1. **`psql`**: Not installed? 🛑 Installation both for Mac and Ubuntu Follow the guide [here](https://www.timescale.com/blog/how-to-install-psql-on-mac-ubuntu-debian-windows/).

2.  Execute the following commands to set up your database:
```bash
mkdir -p ./data
psql -h 5432:5432 -U poostgres -d users -c "CREATE DATABASE users;"
psql -h hostname -U username -d users -a -f ./resources/users.sql
```

### How to run
```bash
docker-compose up --build -d
```

