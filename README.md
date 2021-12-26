# X Server

## Installation (SUBJECT TO CHANGE)

### Prerequisites:
- `python3`

Download from github releases:
```shell
<curl link once release is out>
cd <directory>
```

Generate a config & walk through prompt:
```shell
./genconfig.py
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

Build in dev mode:
```shell
make dev
```

Build a full local release:
```shell
make
```