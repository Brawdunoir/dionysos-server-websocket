[![Documentation Status](https://readthedocs.org/projects/dionysos-server/badge/?version=latest)](https://dionysos-server.readthedocs.io/en/latest/)

# What is it ?
dionysos-server is a server instance for the [dionysos-client](https://github.com/Brawdunoir/dionysos-client) project, enabling users to **share cinematic experiences**.

# I want my own instance
If you have your own server setup and you don’t want to use the default instance address to connect with your friends, feel free to host your own dionysos server instance.

Note that your server and its ports need to be publicly accessible from your friends. We recommend using a reverse proxy such as traefik.

## Using docker 🐳
The preferred method is to run the default dionysos docker image.
```
docker run -p 8080:8080 -e DIONYSOS_ENVIRONMENT="PROD" -v /path/to/logs:/logs brawdunoir/dionysos-server:master
```

This command will:
- Pull the `brawdunoir/dionysos-server:master` image from Docker Hub
- Forward port 8080 from the docker container to the host
- Set to "PROD" the environment
- Write logs to the `/path/to/logs` folder

Below is a simple `docker-compose.yml` file:
```
version: '3.9'

services:
  dionysos:
    container_name: dionysos
    image: brawdunoir/dionysos-server:master
    ports:
      - 8080:8080
    environment:
      - DIONYSOS_ENVIRONMENT="PROD"
    volumes:
      - /path/to/logs:/logs
```

# I want to participate 🍵
You can:
- Fill issues, we will try to fix them asap
- Fork this repository, make code changes and create a pull request



# Just looking for the API
If you are developing a client for dionysos and are looking for the API, head over the [docs](https://dionysos-server.readthedocs.io/en/latest/) !
