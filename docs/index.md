# Dionysos – Share cinematic experiences

*Note: this documentation is for the server part of Dionysos.*

## What is dionysos ?

Dionysos is an open source client to connect with friends and watch movies synchronously and locally.
Equipped with a peer-to-peer system and a server for the social and connection part, our goal is to connect people.

Although the software was designed to share movies, any kind of multimedia material can be shared via the client.

## I want to enhance the server

You can:

- Fill issues, we will try to fix them asap
- Fork this repository, make code changes and create a pull request.

See also [How to contribute to open source](https://github.com/FreeCodeCamp/how-to-contribute-to-open-source).

## I want to improve the official client or develop a new one

Use this documentation to be aware of messages between the client and the server, called requests (client->server) and responses (server->client).

You can either develop in full local (more control) or use the public test server (easier).

### Using the test server

Simply connect with a websocket to `wss://dionysos-test.yannlacroix.fr` and start your journey.

### Local
Either use docker or compile server from source using Go.

Using docker:
```s
docker run -p 8080:8080 dionysos-server:master
```

Compile from source:
```s
git clone https://github.com/Brawdunoir/dionysos-server.git
cd dionysos-server
go run .
```

Then you should connect to `ws://localhost:8080`.
