# Port Scanner

Lightweight port scanner CLI tool written in Go


## Installation

### Docker

Pull the latest image:
```bash
docker pull omercengel/port-scanner:latest
```

Run the container:
```bash
docker run --rm omercengel/port-scanner -h
```

### Build Docker Image

Clone the repository:
```bash
git clone https://github.com/omercengel/port-scanner.git
```

Change directory:
```bash
cd port-scanner
```

Build image from Dockerfile:
```bash
docker build -t port-scanner .
```

Run the container:
```bash
docker run --rm port-scanner -h
```

### Build From Source

Clone the repository:
```bash
git clone https://github.com/omercengel/port-scanner.git
```

Change directory:
```bash
cd port-scanner
```

Build the binary:
```bash
go build -o port-scanner ./cmd/main.go
```

Run the binary:
```bash
./port-scanner -h
```

## Usage

### Docker

```bash
docker run --rm -v /path/to/your/output:/output port-scanner -a 192.168.1.134
```

```bash
docker run --rm -v /path/to/your/output:/output port-scanner -a 192.168.1.134 -p 1-1024 -m stealth
```

```bash
docker run --rm -v /path/to/your/output:/output port-scanner -a 192.168.1.134 -p 80,443 -o results -f json
```

### Source

```bash
./port-scanner -a 192.168.1.134
```

```bash
./port-scanner -a 192.168.1.134 -p 1-1024 -m stealth
```

```bash
./port-scanner -a 192.168.1.134 -p 80,443 -o results -f json
```


## Flags

| Flag        | Short | Type   | Required | Default               | Description                      |
| :---------- | :---- | :----- | :------- | :-------------------- | :------------------------------- |
| `--address` | `-a`  | string | `true`   | -                   | domain or ip address             |
| `--ports`   | `-p`  | string | `false`  | 1-65535             | range: 1-1024 or list: 80,443    |
| `--mode`    | `-m`  | string | `false`  | default             | stealth, default, rapid          |
| `--output`  | `-o`  | string | `false`  | YYYY-MM-DD_HH:MM:SS | output file name                 |
| `--format`  | `-f`  | string | `false`  | txt                 | txt, json, csv                   |
| `--timeout` | `-t`  | int    | `false`  | mode's timeout      | timeout per port in milliseconds |

## Contributing

Contributions are welcome! Whether you want to fix bugs, add new features, improve documentation, you can contribute to this project by following these steps:

#### Fork the repository

Click the "Fork" button at the top-right of the repo page to create your own copy.

#### Clone your fork
```bash
git clone https://github.com/your-username/port-scanner.git
```

#### Change directory:
```bash
cd port-scanner
```

#### Create a new branch:
```bash
git checkout -b feature-or-bugfix-name
```

#### Make your changes:
- Fix bugs, add features, or update documentation.
- Keep your code clean and follow existing formatting conventions.

#### Stage your changes:
```bash
git add .
```

#### Commit your changes:
```bash
git commit -m "Brief description of changes"
```

#### Push to your fork:
```bash
git push origin feature-or-bugfix-name
```

#### Open a Pull Request:
- Go to the original repository on GitHub and click "Compare & pull request".
- Provide a clear description of your changes and submit the PR.
## License
This project is licensed under the MIT License.  
For details, see the [LICENSE](https://github.com/omer-cengel/port-scanner/blob/master/LICENSE) file.
