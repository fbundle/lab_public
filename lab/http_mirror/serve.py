#!/usr/bin/env python3

import argparse
import http.server
import os
import socketserver


class MirrorHandler(http.server.SimpleHTTPRequestHandler):
    """
    Like `python -m http.server`, but:
    - If a request path has no extension and the exact file doesn't exist,
      try `<path>.html`, then `<path>/index.html`.
    This helps URLs like `/Armor` render instead of downloading.
    """

    def translate_path(self, path: str) -> str:
        base = super().translate_path(path)

        # If the path already exists, keep default behavior.
        if os.path.exists(base):
            return base

        _root, ext = os.path.splitext(base)
        if ext:  # has an extension already
            return base

        html_path = base + ".html"
        if os.path.exists(html_path):
            return html_path

        index_path = os.path.join(base, "index.html")
        if os.path.exists(index_path):
            return index_path

        return base

    def guess_type(self, path: str) -> str:
        """
        Ensure HTML is served as UTF‑8 so special characters render correctly.
        - Known .html/.htm files: force 'text/html; charset=utf-8'
        - Extensionless files that *look* like HTML: same.
        Otherwise, fall back to the default mapping.
        """
        _root, ext = os.path.splitext(path)
        if ext.lower() in {".html", ".htm"}:
            return "text/html; charset=utf-8"

        if ext:
            return super().guess_type(path)

        try:
            with open(path, "rb") as f:
                head = f.read(2048).lstrip().lower()
        except OSError:
            return super().guess_type(path)

        if head.startswith(b"<!doctype html") or head.startswith(b"<html") or b"<head" in head[:200]:
            return "text/html; charset=utf-8"

        return super().guess_type(path)


def main() -> None:
    p = argparse.ArgumentParser()
    p.add_argument(
        "dir",
        nargs="?",
        default=os.path.join("tmp", "fextralife_eldenring", "https", "eldenring.wiki.fextralife.com"),
        help="Directory to serve",
    )
    p.add_argument("-p", "--port", type=int, default=3000)
    args = p.parse_args()

    os.chdir(args.dir)

    with socketserver.TCPServer(("", args.port), MirrorHandler) as httpd:
        print(f"Serving {os.getcwd()} at http://localhost:{args.port}/")
        httpd.serve_forever()


if __name__ == "__main__":
    main()
