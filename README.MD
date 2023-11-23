
# Frecency Go Agent
Frecency Go Agent is a multi-platform post-exploitation agent developed in Go. It comprises a server and a client (or agent). The agent connects to the server periodically using http POST requests. As a response to the connection, the server returns a list of commands to execute on the client. Built-in commands include port forwards, reverse shell, and configuration changes.

To harden detection and forensics:
- All communications in Go Agent are encrypted and create a 'normal' impression in traffic analysis. By default, Beacons use HTTPS, while reverse shells and port forwards use SSH tunnels.
- The http server seems like a default installation of nginx to those other than the agent.
- All sensitive strings within the agent binary, like configuration, are AES encrypted using a random key.
- Agent configuration is customizable, users can choose for example beacon interval, user agent of beacons, and beacon endpoints.

### Note
Go Agent is not designed for illegal activities. It has no inbuilt functionality for illicit activities, nor will it be incorporated. This includes no persistence or modules for stealing credentials from systems. The owner is not liable for any misuse of go-agent.

## Why?
As there was no long-deployable, cross-platform, flexible reverse shell with only basic functionality and standard protocols found. Go Agent is an attempt to create one.

## Features
- Encrypted communications
- Reverse shell access
- Upload / download files
- Reverse port forwarding

## Installation
```
apt install -y git make upx strip
git clone <repo name>
cd <repo name>
make deps
```

## Compilation
Both client and server of Go Agent are standalone executables that need to be compiled. Server has only been tested to compile and work on linux. Client has been tested to compile on Linux, and work on Linux and Windows.

### Server
1. Compile the server:
```
make prod
```

### Client
1. Set configuration in /client/config/parameters.go
2. Remove //build +ignore tag from client/config/parameters.go
3. Run `go generate` in /client/config directory
4. Add //build +ignore tag back to client/config/parameters.go
3. Compile the client for all platforms (in cmd/client): `make prod`
