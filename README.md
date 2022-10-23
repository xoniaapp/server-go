# server

## Building

### 1. Make sure you have Go installed.

### 2. Git clone the repo.

```sh
git clone git@github.com:xoniaapp/server.git
```
### 3. Build it

```
cd server && go build .
```

## Pre-build binaries
[Latest Binary](https://github.com/xoniaapp/server/releases/latest) 

### Supported Platform

- `amd64` - linux
- `arm64`- linux

## Environment Variables

```
export PORT=""
export DATABASE_URL=""
export REDIS_URL=""
export CORS_ORIGIN=""
export SECRET=""
export AWS_ACCESS_KEY=""
export AWS_SECRET_ACCESS_KEY=""
export AWS_STORAGE_BUCKET_NAME=""
export AWS_S3_REGION=""
export HANDLER_TIMEOUT=""
export MAX_BODY_BYTES=""
export MAIL_USER=""
export MAIL_PASSWORD=""
```

## `systemd` file
```
[Unit]
Description=xoniaapp

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/usr/bin/xoniaapp

Environment=PORT=
Environment=DATABASE_URL=
Environment=REDIS_URL=
Environment=CORS_ORIGIN=
Environment=SECRET=
Environment=HANDLER_TIMEOUT=
Environment=MAX_BODY_BYTES=
Environment=AWS_ACCESS_KEY=
Environment=AWS_SECRET_ACCESS_KEY=
Environment=AWS_STORAGE_BUCKET_NAME=
Environment=AWS_S3_REGION=

[Install]
WantedBy=multi-user.target
```

## License
### [AGPL-3.0 LICENSE](./LICENSE)
