# serialcaster

Reads bytes from a serial port and broadcasts them to all connected TCP clients.

## Usage

```
serialcaster [flags]

Flags:
  -p, --port      string   serial device path (default "/dev/ttyUSB0")
  -b, --baud      int      baud rate (default 9600)
      --databits  int      data bits: 5, 6, 7, 8 (default 8)
      --stopbits  string   stop bits: 1, 1.5, 2 (default "1")
      --parity    string   parity: none, odd, even, mark, space (default "none")
  -l, --listen    string   TCP listen address (default ":8888")
```

Every flag can also be set via environment variable with the `SD_` prefix:

| Flag        | Environment variable |
|-------------|----------------------|
| `--port`     | `SD_PORT`            |
| `--baud`     | `SD_BAUD`            |
| `--databits` | `SD_DATABITS`        |
| `--stopbits` | `SD_STOPBITS`        |
| `--parity`   | `SD_PARITY`          |
| `--listen`   | `SD_LISTEN`          |

## Build

```sh
go mod tidy
go build -o serialcaster .
```

## Docker

```sh
# Build — override DIALOUT_GID if dialout is not gid 20 on your host
docker build \
  --build-arg DIALOUT_GID=$(getent group dialout | cut -d: -f3) \
  -t serialcaster .

# Run
docker run \
  --device /dev/ttyUSB0 \
  -p 8888:8888 \
  serialcaster --port /dev/ttyUSB0 --baud 115200
```

The runtime image is based on `scratch` and contains only the binary.
Serial device access requires the container's gid to match the `dialout` group
on the host. Use `--build-arg DIALOUT_GID` or `--user` at runtime to set it.

## Connecting

Any TCP client can subscribe to the serial stream:

```sh
nc <host> 8888
```

Multiple clients can connect simultaneously. Frames are broadcast to all of
them; if a client's receive buffer is full its frames are dropped rather than
blocking the other clients.
