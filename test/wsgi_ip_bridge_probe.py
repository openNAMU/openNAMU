#!/usr/bin/env python3

import importlib.util
import io
import json
import os
import sys
import tempfile
from pathlib import Path


def load_wsgi_module():
    repository_root = Path(__file__).resolve().parent.parent
    specification = importlib.util.spec_from_file_location(
        "opennamu_wsgi_app", repository_root / "app.py"
    )
    if specification is None or specification.loader is None:
        raise RuntimeError("could not load app.py")
    module = importlib.util.module_from_spec(specification)
    specification.loader.exec_module(module)
    return module


def run_case(module, case):
    environ = {
        "REQUEST_METHOD": "GET",
        "PATH_INFO": "/_test/wsgi_identity",
        "QUERY_STRING": "",
        "SERVER_NAME": "wiki.test",
        "SERVER_PORT": "80",
        "SERVER_PROTOCOL": "HTTP/1.1",
        "HTTP_HOST": "wiki.test",
        "REMOTE_ADDR": case.get("remote_addr", ""),
        "CONTENT_LENGTH": "0",
        "wsgi.url_scheme": "http",
        "wsgi.input": io.BytesIO(b""),
    }
    for header_name, header_value in case.get("headers", {}).items():
        environ["HTTP_" + header_name.upper().replace("-", "_")] = header_value

    response_info = {}

    def start_response(status, headers, exc_info=None):
        response_info["status"] = status
        response_info["headers"] = headers

    response_body = b"".join(module.wsgi_application(environ, start_response))
    return {
        "status": response_info.get("status", ""),
        "body": response_body.decode("utf-8", errors="replace"),
    }


def main():
    if len(sys.argv) != 3:
        raise SystemExit("usage: wsgi_ip_bridge_probe.py HOST PORT")

    module = load_wsgi_module()
    repository_root = Path(__file__).resolve().parent.parent
    temporary_root = repository_root / "tmp"
    temporary_root.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory(prefix="opennamu-wsgi-test-", dir=temporary_root) as runtime_dir:
        os.environ["NAMU_WSGI_HOST"] = sys.argv[1]
        os.environ["NAMU_WSGI_PORT"] = sys.argv[2]
        os.environ["NAMU_WSGI_RUNTIME_DIR"] = runtime_dir
        os.environ["NAMU_WSGI_PORT_FILE"] = str(Path(runtime_dir) / "missing.port")

        cases = json.load(sys.stdin)
        results = [run_case(module, case) for case in cases]
        json.dump(results, sys.stdout)
        sys.stdout.write("\n")


if __name__ == "__main__":
    main()
