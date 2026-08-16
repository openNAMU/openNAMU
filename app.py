#!/usr/bin/env python3

import os
import platform
import subprocess
import sys
from pathlib import Path

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
    except RuntimeError as error:
        print(str(error), file=sys.stderr)
        return 1

    if not binary_path.exists():
        print("바이너리를 찾을 수 없습니다: " + str(binary_path), file=sys.stderr)
        return 1

    try:
        return run_binary(binary_path, sys.argv[1:], base_dir)
    except OSError as error:
        print("바이너리 실행에 실패했습니다: " + str(error), file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
