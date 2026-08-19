#!/usr/bin/env python3

import http.client
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


def wsgi_request_body(environ):
    content_length = environ.get("CONTENT_LENGTH", "")
    if content_length:
        try:
            content_length = int(content_length)
        except ValueError:
            return b""

        if content_length > 0:
            return environ["wsgi.input"].read(content_length)

    if environ.get("wsgi.input_terminated"):
        return environ["wsgi.input"].read()

    return b""


def wsgi_path(environ):
    raw_uri = environ.get("RAW_URI", "")
    if raw_uri.startswith("/"):
        return raw_uri

    path = environ.get("PATH_INFO", "/") or "/"
    path = quote(path, safe="/%:@-._~!$&'()*+,;=")
    query = environ.get("QUERY_STRING", "")
    if query:
        path += "?" + query
    return path


def wsgi_headers(environ, target_host, target_port):
    headers = {}
    excluded = {
        "connection",
        "content-length",
        "host",
        "keep-alive",
        "proxy-authenticate",
        "proxy-authorization",
        "te",
        "trailer",
        "transfer-encoding",
        "upgrade",
    }

    for key, value in environ.items():
        if key.startswith("HTTP_"):
            header_name = key[5:].replace("_", "-")
        elif key == "CONTENT_TYPE":
            header_name = "Content-Type"
        else:
            continue

        if header_name.lower() not in excluded:
            headers[header_name] = value

    headers["Host"] = environ.get("HTTP_HOST", target_host + ":" + str(target_port))
    headers["X-Real-IP"] = environ.get("HTTP_X_REAL_IP", environ.get("REMOTE_ADDR", ""))
    headers["X-Forwarded-Proto"] = environ.get(
        "HTTP_X_FORWARDED_PROTO",
        environ.get("wsgi.url_scheme", "http"),
    )
    return headers


def wsgi_application(environ, start_response):
    target_host = os.environ.get("NAMU_WSGI_HOST", "127.0.0.1")
    target_port = os.environ.get("NAMU_WSGI_PORT", "3000")
    try:
        target_port = int(target_port)
    except ValueError:
        target_port = 3000

    connection = http.client.HTTPConnection(target_host, target_port, timeout=download_timeout)
    try:
        connection.request(
            environ.get("REQUEST_METHOD", "GET"),
            wsgi_path(environ),
            body=wsgi_request_body(environ),
            headers=wsgi_headers(environ, target_host, target_port),
        )
        response = connection.getresponse()
    except (OSError, http.client.HTTPException) as error:
        connection.close()
        body = ("openNAMU backend is not running: " + str(error)).encode("utf-8")
        start_response(
            "503 Service Unavailable",
            [("Content-Type", "text/plain; charset=utf-8"), ("Content-Length", str(len(body)))],
        )
        return [body]

    hop_by_hop_headers = {
        "connection",
        "keep-alive",
        "proxy-authenticate",
        "proxy-authorization",
        "te",
        "trailer",
        "transfer-encoding",
        "upgrade",
    }
    response_headers = [
        (name, value)
        for name, value in response.getheaders()
        if name.lower() not in hop_by_hop_headers
    ]
    start_response(str(response.status) + " " + (response.reason or ""), response_headers)

    def response_body():
        try:
            while True:
                data = response.read(download_chunk_size)
                if not data:
                    break
                yield data
        finally:
            response.close()
            connection.close()

    return response_body()


app = wsgi_application
application = app


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
