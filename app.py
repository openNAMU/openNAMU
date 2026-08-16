#!/usr/bin/env python3

import json
import os
import platform
import subprocess
import sys
from pathlib import Path
from urllib.error import HTTPError, URLError
from urllib.parse import quote
from urllib.request import Request, urlopen

binary_name_map = {
    "linux": {
        "amd64": "main.amd64.bin",
        "arm64": "main.arm64.bin",
    },
    "windows": {
        "amd64": "main.amd64.exe",
        "arm64": "main.arm64.exe",
    },
    "darwin": {
        "arm64": "main.mac.arm64.bin",
    },
}

download_timeout = 60
download_chunk_size = 1024 * 1024
remote_version_url = "https://raw.githubusercontent.com/openNAMU/openNAMU/refs/heads/beta/version.json"


def get_architecture():
    machine = platform.machine().lower()

    if machine in ("amd64", "x86_64"):
        return "amd64"

    if machine in ("arm64", "aarch64"):
        return "arm64"

    raise RuntimeError("지원하지 않는 CPU 아키텍처입니다: " + machine)


def get_binary_name():
    system = platform.system().lower()
    architecture = get_architecture()

    if system not in binary_name_map:
        raise RuntimeError("지원하지 않는 운영체제입니다: " + system)

    if architecture not in binary_name_map[system]:
        raise RuntimeError("해당 운영체제의 바이너리가 없습니다: " + system + "/" + architecture)

    return binary_name_map[system][architecture]


def get_release_tag():
    request = Request(
        remote_version_url,
        headers={"User-Agent": "openNAMU-launcher"},
    )

    try:
        with urlopen(request, timeout=download_timeout) as response:
            if response.status != 200:
                raise RuntimeError("version.json 요청에 실패했습니다: " + str(response.status))
            version_data = json.load(response)
    except (HTTPError, URLError, OSError, json.JSONDecodeError) as error:
        raise RuntimeError("최신 version.json을 읽을 수 없습니다: " + str(error))

    release_tag = version_data.get("r_ver")
    if not isinstance(release_tag, str) or not release_tag:
        raise RuntimeError("최신 version.json에 r_ver이 없습니다.")
    return release_tag


def download_binary(binary_path):
    release_tag = get_release_tag()
    binary_url = "https://github.com/openNAMU/openNAMU/releases/download/" + quote(release_tag, safe="") + "/" + binary_path.name
    temporary_path = binary_path.with_name(binary_path.name + ".download")
    request = Request(
        binary_url,
        headers={"User-Agent": "openNAMU-launcher"},
    )

    try:
        with urlopen(request, timeout=download_timeout) as response:
            if response.status != 200:
                raise RuntimeError("바이너리 요청에 실패했습니다: " + str(response.status))
            with temporary_path.open("wb") as file_data:
                while True:
                    data = response.read(download_chunk_size)
                    if not data:
                        break
                    file_data.write(data)

        if not temporary_path.exists() or temporary_path.stat().st_size == 0:
            raise RuntimeError("다운로드된 바이너리가 비어 있습니다.")

        os.replace(temporary_path, binary_path)
        if os.name != "nt":
            binary_path.chmod(0o755)
    except Exception:
        if temporary_path.exists():
            temporary_path.unlink()
        raise


def run_binary(binary_path, arguments, base_dir):
    command = [str(binary_path)] + arguments

    if os.name != "nt":
        os.chdir(base_dir)
        os.execv(str(binary_path), command)

    return subprocess.call(command, cwd=str(base_dir))


def main():
    base_dir = Path(__file__).resolve().parent
    try:
        binary_path = base_dir / get_binary_name()
        if not binary_path.is_file():
            print("바이너리가 없어 최신 버전을 다운로드합니다.")
            download_binary(binary_path)
    except (HTTPError, URLError, OSError, RuntimeError) as error:
        print(str(error), file=sys.stderr)
        return 1

    try:
        return run_binary(binary_path, sys.argv[1:], base_dir)
    except OSError as error:
        print("바이너리 실행에 실패했습니다: " + str(error), file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
