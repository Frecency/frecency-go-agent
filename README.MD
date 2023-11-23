
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

The makefile compiles the client for all platforms, removes symbols from the executables and subsequently packs them using upx.
As of writing Mac executables are not stripped.
NOTE: Makefile assumes commands upx and strip are available

### Debug builds
If anything doesn't work as anticipated, debugging builds are available for both server and client. These builds include debugging symbols, numerous printouts to track the activities of the program, and omit compressing the resulting binaries. Debugging builds are easy to reverse-engineer and therefore should not be used in production.

To compile a debugging build, use the command:
```
make debug
```
For both client and the server.

#### Configuration
The following parameters can be customized as per user requirement in /client/config/parameters.go:

```
Sleeptime (int) - The time between beacons in seconds
Jitter (int) - Random extra delay to be added to the Sleeptime
Retries (int) - Number of tries in case C2 is unreachable
UserAgent (string) - The User Agent that is sent on beacons
CCHost (string) - The base server URL
Endpoints (string[]) - The server endpoints to which beacon. They are chosen randomly
```

## Usage

### Server

It is recommended to get an SSL certificate and corresponding private key from Lets encrypt for best results. The agent will not trust self-signed certificates unless the root certificate has been installed first.

```
Usage of ./server:
  -cert string
    File with https cert (default "server.crt")
  -forwarded-timeout duration
    Forwarded-tcpip timeout (default 3h0m0s)
  -handlerport int
    Port of HTTP* listener (default 443)
  -http
    Use HTTP instead of HTTPS
  -main-timeout duration
    Client socket timeout (default 3h0m0s)
  -privkey string
    File with private key for cert (default "server.key")
  -sshport int
    Port of SSH listener (default 22)
```

#### examples
```
./server -handlerport 8080 -sshport 2222
./server -cert othercert.crt -privkey otherkey.key
./server -forwarded-timeout 0h10m0s
```


## Reverse shell and file access

Reverse shell access operates as follows: The client initiates an sshd listening on localhost and then establishes a reverse port forward to the server. The sshd on the client can then be accessed from the server using the forwarded port. There are commands quickssh and serveSSH on the server that start the reverse shell. Quickssh uses the built-in sshd on the Go Agent server to catch the port forward, serveSSH allows using any normal sshd for the purpose.

NOTE: For the client to forward port to the host, it needs to have valid SSH credentials. It must be ensured that the credentials cannot be used to obtain a shell on the server when using normal sshd. In the built-in SSH, only reverse port forwards are allowed by clients.

Filesystem access is executed using sftp.

### Unix (linux or mac):
After port forwarding from the client to the server, simply SSH in using the generated credentials.
```
ssh <username>@localhost -p <port>
```

### Windows:
After port forwarding from the client to the server, forward another port from the client machine using the client's sshd to localhost, then connect to the forwarded port using netcat
```
ssh -N -L <localport>:0.0.0.0:<destport> localhost -p <port>
nc localhost <localport>
```

NOTE: in order for the client to be able to forward port to the host, it needs to have valid SSH credentials. One must make sure that the credentials cannot be used to obtain a shell to the server when using normal sshd. In the inbuilt SSH, only reverse port forwards are allowed by clients.

Filesystem access is done using sftp.

### Unix (linux or mac):
After forwarding port from the client to server, just SSH in using the generated credentials.
```
ssh <username>@localhost -p <port>
```

### Windows:
After forwarding port from the client to server, forward another port from the client machine using the client's sshd to localhost, then connect to the forwarded port using netcat
```
ssh -N -L <localport>:0.0.0.0:<destport> localhost -p <port>
nc localhost <localport>
```

NOTE: right now gthe destport is not used! The remote host 0.0.0.0 is a magic string that tells the agent to serve a cmd shell to the client.

### Uploading / Downloading files with sftp (both windows and unix)
After forwarding port from the client to the server, you can access the filesystem using sftp.
```
sftp -P <port> <username>@localhost
```

## Port forwarding

Local ssh listener in the server allows a maximum of 5 reverse port forwards from the victim.

### Windows:
On windows, ports can be forwarded but no shell can be launched at the same time, therefore flag -N must be used. Omitting the flag causes the client to serve sftp instead
```
ssh -N -L <localport>:<desthost>:<destport> localhost -p <port>
```

### Unix (linux or mac):
On linux, ports can be forwarded similarly with the familiar sshd
```
ssh [-N] -L <localport>:<desthost>:<destport> localhost -p <port>
```

## Features for future releases
- Add x11 forwarding to Unix ssh handler
- Encrypt protocol strings

## Contributing
Contributions are highly appreciated! Please submit your thoroughly tested code via a pull request. Documentation writers are also encouraged to contribute.