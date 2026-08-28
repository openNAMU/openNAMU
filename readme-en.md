[(en-US)](./readme-en.md) | [(ko-KR)](./readme.md)

# openNAMU
[![Go 1.25+](https://img.shields.io/badge/go-%3E%3D1.25-blue.svg)](https://go.dev/)
[![LICENSE](https://img.shields.io/badge/license-BSD%203--Clause-lightgrey.svg)](./LICENSE)

![](https://raw.githubusercontent.com/openNAMU/openNAMU/beta/.github/logo.png)

openNAMU is a Go-based wiki engine.

## Getting Started

Download the binary for your operating system and architecture from the [Releases page](https://github.com/openNAMU/openNAMU/releases).

| Operating system | Binary |
| --- | --- |
| Linux amd64 | main.amd64.bin |
| Linux arm64 | main.arm64.bin |
| Windows amd64 | main.amd64.exe |
| Windows arm64 | main.arm64.exe |
| macOS arm64 | main.mac.arm64.bin |

### Linux

```bash
./main.amd64.bin
```

### Windows

```powershell
.\main.amd64.exe
```

The default port is `3000`. It can be changed with the first argument.

```bash
./main.amd64.bin 3000
```

Add `--localhost` when the server should only be accessible from the same computer.

```bash
./main.amd64.bin --localhost
./main.amd64.bin 3000 --localhost
```

If the server cannot start normally, run the emergency tool mode.

```bash
./main.amd64.bin em
```

## For developers

Use this section when modifying the source code or building the project yourself.

### Clone

You can clone this repository by entering the following command at the terminal or command prompt:

- Stable: `git clone -b stable https://github.com/openNAMU/openNAMU.git`
- Beta: `git clone -b beta https://github.com/openNAMU/openNAMU.git`

### Build and test

Go 1.25 or later is required to build from source.

```bash
./linux.sh linux_amd64
```

On Windows, use the following commands:

```powershell
.\windows.ps1 windows_amd64
```

You can also build the server directly:

```bash
go build -o main .
./main
```

Run the following checks before submitting changes:

```bash
go test ./...
go vet ./...
```

## Database configuration

The default database configuration is stored in `data/set.json`.

```json
{"db_type":"sqlite","db":"data"}
```

With this setting, SQLite uses `data.db` in the working directory. If no configuration exists, SQLite `data.db` is used by default.

To use MySQL, set `data/set.json` as follows and put the connection details in `data/mysql.json`.

```json
{"db_type":"mysql","db":"openNAMU"}
```

```json
{"user":"username","password":"password","host":"127.0.0.1","port":"3306"}
```

Always back up the database and configuration files before migrating an existing installation.

## Current status

openNAMU is currently in beta.

- Runs as a single Go server.
- Supports SQLite and MySQL.
- The default skins and language data are embedded in the binary.
- The NamuMark compatibility renderer is still being improved.
- Edit requests and section editing are currently unavailable.

## Contribute

openNAMU may contain untested bugs. Reporting them helps improve the project. [Create an Issue](https://github.com/openNAMU/openNAMU/issues/new).

openNAMU is an open source project. You can modify the code and submit a [Pull Request](https://github.com/openNAMU/openNAMU/compare).

## License

openNAMU is licensed under the [BSD 3-Clause License](./LICENSE). Please refer to the documentation for details.

### External Projects

- [Quotes icon - Dave Gandy](http://www.flaticon.com/free-icon/quote-left_25672)
- [highlight.js](https://highlightjs.org/)
- [KaTeX](https://katex.org/)
- [Feather](https://feathericons.com/)
- [go.mod](./go.mod)

### Contributors

- [Team Croatia](https://github.com/TeamCroatia)
- Basix
- Efrit
- Others

## Supported markup

- NamuMark
- Markdown
- MacroMark
- Raw

## Etc.

- Owner rights are granted to the first registrant.
- The default port is `3000`.
- [Contributors](https://github.com/openNAMU/openNAMU/graphs/contributors)
- [Old History 1](https://github.com/openNAMU/openNAMU-Backup)
