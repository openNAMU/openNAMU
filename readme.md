[(en-US)](./readme-en.md) | [(ko-KR)](./readme.md)

# 오픈나무
[![Go 1.25 이상](https://img.shields.io/badge/go-%3E%3D1.25-blue.svg)](https://go.dev/)
[![라이선스](https://img.shields.io/badge/license-BSD%203--Clause-lightgrey.svg)](./LICENSE)

![](https://raw.githubusercontent.com/openNAMU/openNAMU/beta/.github/logo.png)

오픈나무는 Go 기반의 위키 엔진입니다.

## 시작하기

릴리즈 페이지에서 사용하는 운영체제와 아키텍처에 맞는 바이너리를 내려받아 실행할 수 있습니다.

[릴리즈 페이지](https://github.com/openNAMU/openNAMU/releases)에서 다음 파일 중 하나를 내려받으세요.

| 운영체제 | 바이너리 |
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

기본 포트는 `3000`이며, 실행 시 첫 번째 인자로 변경할 수 있습니다.

```bash
./main.amd64.bin 3000
```

외부에서 접근하지 않고 같은 컴퓨터에서만 사용할 때는 `--localhost`를 추가합니다.

```bash
./main.amd64.bin --localhost
./main.amd64.bin 3000 --localhost
```

서버가 정상적으로 시작되지 않을 때는 이머전시 툴을 실행할 수 있습니다.

```bash
./main.amd64.bin em
```

## 개발자용

소스 코드를 수정하거나 직접 빌드하려는 경우에 참고하세요.

### 클론

아래 명령을 터미널(명령 프롬프트)에 입력하여 본 리포지토리를 클론할 수 있습니다.

- 일반: `git clone -b stable https://github.com/openNAMU/openNAMU.git`
- 베타: `git clone -b beta https://github.com/openNAMU/openNAMU.git`

### 빌드 및 테스트

소스에서 빌드하려면 Go 1.25 이상이 필요합니다.

```bash
./linux.sh linux_amd64
```

Windows에서는 다음 명령을 사용할 수 있습니다.

```powershell
.\windows.ps1 windows_amd64
```

직접 빌드하려면 다음 명령을 사용할 수 있습니다.

```bash
go build -o main .
./main
```

변경 사항을 확인할 때는 다음 명령을 사용할 수 있습니다.

```bash
go test ./...
go vet ./...
```

## 데이터베이스 설정

기본 설정은 `data/set.json`에서 지정합니다.

```json
{"db_type":"sqlite","db":"data"}
```

위 설정에서는 실행 디렉터리의 `data.db`를 SQLite 데이터베이스로 사용합니다. 설정이 없으면 SQLite의 `data.db`를 기본값으로 사용합니다.

MySQL을 사용하려면 `data/set.json`을 다음처럼 작성하고 `data/mysql.json`에 접속 정보를 넣습니다.

```json
{"db_type":"mysql","db":"openNAMU"}
```

```json
{"user":"username","password":"password","host":"127.0.0.1","port":"3306"}
```

기존 데이터베이스를 전환하기 전에는 반드시 데이터베이스와 설정 파일을 백업하세요.

## 현재 상태

오픈나무는 현재 beta 단계입니다.

- Go 단일 서버로 실행합니다.
- SQLite와 MySQL을 지원합니다.
- 기본 스킨과 언어 데이터는 바이너리에 포함됩니다.

## 기여

오픈나무에는 확인되지 않은 버그가 존재할 수 있습니다. 이를 보고해주시면 오픈나무의 발전을 도울 수 있습니다. [여기](https://github.com/openNAMU/openNAMU/issues/new)를 눌러 버그를 보고해주세요.

오픈나무는 오픈소스 프로젝트입니다. 원한다면 직접 코드를 수정하고 [Pull Request](https://github.com/openNAMU/openNAMU/compare)를 보낼 수 있습니다.

## 라이선스

오픈나무 프로젝트는 [BSD 3-Clause License](./LICENSE)를 따릅니다. 자세한 내용은 문서를 참고하세요.

### 포함된 외부 프로젝트

- [Quotes icon - Dave Gandy](http://www.flaticon.com/free-icon/quote-left_25672)
- [highlight.js](https://highlightjs.org/)
- [KaTeX](https://katex.org/)
- [Feather](https://feathericons.com/)
- [go.mod](./go.mod)

### 도움을 주신 분들

- [Team Croatia](https://github.com/TeamCroatia)
- Basix
- Efrit
- 기타 여러 사람들

## 지원 문법

- 나무마크 (NamuMark)
- 마크다운 (Markdown)
- MacroMark
- Raw

## 기타

- 첫 가입자에게 소유자 권한이 부여됩니다.
- 기본 포트는 `3000`입니다.
- [기여자 목록](https://github.com/openNAMU/openNAMU/graphs/contributors)
- [예전 히스토리 1](https://github.com/openNAMU/openNAMU-Backup)
