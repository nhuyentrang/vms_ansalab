# VMS share you ip camera to world!

VMS converts your RTSP streams to formats consumable in a web browser
like MSE (Media Source Extensions), WebRTC, or HLS. It's fully native Golang
without the use of FFmpeg or GStreamer!

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Command-line](#command-line)
- [API documentation](#api-documentation)
- [Limitations](#limitations)
- [Performance](#performance)
- [Authors](#authors)
- [License](#license)

## Installation

### Installation from source

1. Download source
   ```bash
   $ git clone https://github.com/deepch/VMS
   ```
1. CD to Directory
   ```bash
    $ cd VMS/
   ```
1. Test Run
   ```bash
    $ GO111MODULE=on go run *.go
   ```
1. Open Browser
    ```bash
    open web browser http://127.0.0.1:8083 work chrome, safari, firefox
    ```

## Installation from docker

1. Run docker container
    ```bash
    $ docker run --name rtsp-to-web --network host ghcr.io/deepch/rtsptoweb:latest 
    ```
1. Open Browser
    ```bash
    open web browser http://127.0.0.1:8083 in chrome, safari, firefox
    ```

You may override the <a href="#example-configjson">configuration</a> `/PATH_TO_CONFIG/config.json` and mount as a docker volume:

```bash
$ docker run --name rtsp-to-web \
    -v /PATH_TO_CONFIG/config.json:/config/config.json \
    --network host \
    ghcr.io/deepch/rtsptoweb:latest 
```

## Configuration

### Server settings

```text
debug           - enable debug output
log_level       - log level (trace, debug, info, warning, error, fatal, or panic)

http_demo       - serve static files
http_debug      - debug http api server
http_login      - http auth login
http_password   - http auth password
http_port       - http server port
http_dir        - path to serve static files from
ice_servers     - array of servers to use for STUN/TURN
ice_username    - username to use for STUN/TURN
ice_credential  - credential to use for STUN/TURN
webrtc_port_min - minimum WebRTC port to use (UDP)
webrtc_port_max - maximum WebRTC port to use (UDP)

https
https_auto_tls
https_auto_tls_name
https_cert
https_key
https_port

rtsp_port       - rtsp server port
```

### Stream settings

```text
name            - stream name
```

### Channel settings

```text
name            - channel name
url             - channel rtsp url
on_demand       - stream mode static (run any time) or ondemand (run only has viewers)
debug           - enable debug output (RTSP client)
audio           - enable audio
status          - default stream status
```

#### Authorization play video

1 - enable config

```text
"token": {
"enable": true,
"backend": "http://127.0.0.1/file.php"
}
```

2 - try

```text
rtsp://127.0.0.1:5541/demo/0?token=you_key
```

file.php need response json

```text
   status: "1" or "0"
 ```

#### RTSP pull modes

  * **on demand** (on_demand=true) - only pull video from the source when there's a viewer
  * **static** (on_demand=false) - pull video from the source constantly

### Example config.json

```json
{
  "server": {
    "debug": true,
    "log_level": "info",
    "http_demo": true,
    "http_debug": false,
    "http_login": "demo",
    "http_password": "demo",
    "http_port": ":8083",
    "ice_servers": ["stun:stun.l.google.com:19302"],
    "rtsp_port": ":5541"
  },
  "streams": {
    "demo1": {
      "name": "test video stream 1",
      "channels": {
        "0": {
          "name": "ch1",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "audio": true,
          "status": 0
        },
        "1": {
          "name": "ch2",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "audio": true,
          "status": 0
        }
      }
    },
    "demo2": {
      "name": "test video stream 2",
      "channels": {
        "0": {
          "name": "ch1",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        },
        "1": {
          "name": "ch2",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        }
      }
    }
  },
  "channel_defaults": {
    "on_demand": true
  }
}
```

## Command-line

### Use help to show available args

```bash
./VMS --help
```

#### Response

```bash
Usage of ./VMS:
  -config string
        config patch (/etc/server/config.json or config.json) (default "config.json")
  -debug
        set debug mode (default true)
```

## API documentation

See the [API docs](/docs/api.md)

## Limitations

Video Codecs Supported: H264 all profiles

Audio Codecs Supported: no

