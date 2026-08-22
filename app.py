#!/usr/bin/env python3

import contextlib
import errno
import http.client
import json
import os
import platform
import socket
import subprocess
import sys
import time
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


def wsgi_log(message):
    print("[openNAMU WSGI] " + message, file=sys.stderr, flush=True)






def wsgi_autostart_enabled(target_host):
    enabled = os.environ.get("NAMU_WSGI_AUTOSTART", "1").lower()
    local_host = target_host.lower() in {"127.0.0.1", "localhost", "::1"}
    return local_host and enabled not in {"0", "false", "no", "off"}


def wsgi_connection_refused(error):
    return isinstance(error, ConnectionRefusedError) or getattr(error, "errno", None) == errno.ECONNREFUSED


def wsgi_backend_is_listening(target_host, target_port):
    try:
        with socket.create_connection((target_host, target_port), timeout=0.25):
            return True
    except OSError:
        return False


def wsgi_runtime_dir(base_dir):
    runtime_dir = Path(os.environ.get("NAMU_WSGI_RUNTIME_DIR", str(base_dir))).expanduser()
    runtime_dir.mkdir(parents=True, exist_ok=True)
    return runtime_dir


@contextlib.contextmanager
def wsgi_backend_lock(lock_path):
    lock_file = lock_path.open("a+")
    try:
        if os.name != "nt":
            import fcntl

            fcntl.flock(lock_file.fileno(), fcntl.LOCK_EX)
        yield
    finally:
        if os.name != "nt":
            fcntl.flock(lock_file.fileno(), fcntl.LOCK_UN)
        lock_file.close()


def wsgi_backend_log_path(base_dir):
    log_path = Path(os.environ.get("NAMU_WSGI_LOG", str(base_dir / "backend.log"))).expanduser()
    if not log_path.is_absolute():
        log_path = base_dir / log_path
    log_path.parent.mkdir(parents=True, exist_ok=True)
    return log_path


def wsgi_backend_wait(target_host, target_port, process, log_path):
    try:
        timeout = float(os.environ.get("NAMU_WSGI_START_TIMEOUT", "15"))
    except ValueError:
        timeout = 15
    timeout = max(1, min(timeout, 60))
    deadline = time.monotonic() + timeout

    while time.monotonic() < deadline:
        if wsgi_backend_is_listening(target_host, target_port):
            return
        return_code = process.poll()
        if return_code is not None:
            raise RuntimeError(
                "backend exited with code " + str(return_code) + "; see " + str(log_path)
            )
        time.sleep(0.1)

    raise RuntimeError(
        "backend did not start listening on "
        + target_host
        + ":"
        + str(target_port)
        + " within "
        + str(timeout)
        + " seconds; see "
        + str(log_path)
    )


def wsgi_start_backend(target_host, target_port):
    base_dir = Path(__file__).resolve().parent
    runtime_dir = wsgi_runtime_dir(base_dir)
    lock_path = Path(
        os.environ.get("NAMU_WSGI_LOCK", str(runtime_dir / ".opennamu-wsgi.lock"))
    ).expanduser()
    if not lock_path.is_absolute():
        lock_path = runtime_dir / lock_path
    lock_path.parent.mkdir(parents=True, exist_ok=True)
    log_path = wsgi_backend_log_path(base_dir)

    with wsgi_backend_lock(lock_path):
        if wsgi_backend_is_listening(target_host, target_port):
            return

        binary_path = base_dir / get_binary_name()
        if not binary_path.is_file():
            wsgi_log("backend binary is missing; downloading " + binary_path.name)
            download_binary(binary_path)

        if os.name != "nt":
            binary_path.chmod(0o755)

        command = [str(binary_path), str(target_port)]
        environment = os.environ.copy()
        environment.pop("NAMU_START_DELAY_MS", None)
        with log_path.open("ab", buffering=0) as log_file:
            log_file.write(
                ("\n[openNAMU WSGI] starting: " + " ".join(command) + "\n").encode("utf-8")
            )
            process_options = {
                "cwd": str(base_dir),
                "env": environment,
                "stdin": subprocess.DEVNULL,
                "stdout": log_file,
                "stderr": subprocess.STDOUT,
            }
            if os.name != "nt":
                process_options["start_new_session"] = True
            process = subprocess.Popen(command, **process_options)

        wsgi_log("started backend process " + str(process.pid) + "; see " + str(log_path))
        wsgi_backend_wait(target_host, target_port, process, log_path)


def wsgi_backend_request(environ, target_host, target_port, request_body):
    connection = http.client.HTTPConnection(target_host, target_port, timeout=download_timeout)
    try:
        connection.request(
            environ.get("REQUEST_METHOD", "GET"),
            wsgi_path(environ),
            body=request_body,
            headers=wsgi_headers(environ, target_host, target_port),
        )
        return connection, connection.getresponse()
    except Exception:
        connection.close()
        raise


def wsgi_backend_error(start_response, target_host, target_port, error, start_error=None):
    message = (
        "openNAMU backend is not running: "
        + str(error)
        + " (target "
        + target_host
        + ":"
        + str(target_port)
        + ")"
    )
    if start_error is not None:
        message += "\nbackend auto-start failed: " + str(start_error)
    wsgi_log(message.replace("\n", " | "))
    body = message.encode("utf-8")
    start_response(
        "503 Service Unavailable",
        [
            ("Content-Type", "text/plain; charset=utf-8"),
            ("Content-Length", str(len(body))),
        ],
    )
    return [body]


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

    request_body = wsgi_request_body(environ)
    try:
        connection, response = wsgi_backend_request(
            environ, target_host, target_port, request_body
        )
    except (OSError, http.client.HTTPException) as error:
        if not wsgi_connection_refused(error) or not wsgi_autostart_enabled(target_host):
            return wsgi_backend_error(start_response, target_host, target_port, error)

        try:
            wsgi_start_backend(target_host, target_port)
        except Exception as start_error:
            return wsgi_backend_error(
                start_response, target_host, target_port, error, start_error
            )

        try:
            connection, response = wsgi_backend_request(
                environ, target_host, target_port, request_body
            )
        except Exception as retry_error:
            return wsgi_backend_error(start_response, target_host, target_port, retry_error)

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
