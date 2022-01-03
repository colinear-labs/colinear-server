# Colinear Payments Server

## Installation (SUBJECT TO CHANGE)

Download from github releases:
```shell
<curl link once initial release is out>
cd <directory>
```

Run the server:
```shell
./xserver -mode=single -port=80
```

## Development

Clone the repository:
```shell
git clone --recursive https://github.com/colinear-labs/colinear-server.git
```

Build & run in dev mode:
```shell
make dev
go run main.go
```

Build a full local release:
```shell
make
```