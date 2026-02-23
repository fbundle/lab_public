from typing import Any, Callable
from html.parser import HTMLParser
import os
import urllib.parse
import urllib.request
import urllib.error
from collections import deque


Node = Any
Memo = Any

def bfs(
    memo: Memo,
    root_node: Node,
    visit_node: Callable[[Memo, Node], tuple[Memo, list[Node]]],
) -> Memo:
    visited: set[Node] = set()
    queue: deque[tuple[Node, int]] = deque([(root_node, 0)])
    while queue:
        node, depth = queue.popleft()
        print(f"queue size: {len(queue)}, visiting (depth {depth}) {node} ...")
        memo, children = visit_node(memo, node)
        visited.add(node)
        for child in children:
            if child in visited:
                continue
            if any(c == child for (c, _) in queue):
                continue
            queue.append((child, depth + 1))
    return memo


class HrefParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.hrefs: list[str] = []
    
    def handle_starttag(self, tag, attrs):
        tag = tag.lower()

        # <a href="...">
        if tag == "a":
            for k, v in attrs:
                if (k or "").lower() == "href" and v:
                    self.hrefs.append(v)
            return

        # <link rel="stylesheet" href="...">
        if tag == "link":
            rel = ""
            href = None
            for k, v in attrs:
                lk = (k or "").lower()
                if lk == "rel" and v:
                    rel = v.lower()
                elif lk == "href" and v:
                    href = v
            if href and "stylesheet" in rel:
                self.hrefs.append(href)
            return

        # <script src="..."></script>
        if tag == "script":
            for k, v in attrs:
                if (k or "").lower() == "src" and v:
                    self.hrefs.append(v)
            return

def parse_hrefs(html: str) -> list[str]:
    parser = HrefParser()
    parser.feed(html)
    hrefs = parser.hrefs
    return hrefs


class Memo:
    def __init__(self, root_dir: str, whitelist: set[str] | None = None):
        if whitelist is None:
            whitelist = set[str]()
        self.whitelist = whitelist
        self.root_dir = root_dir
    
    def parse_url(self, node: str) -> tuple[str, str, str, str]:
        url = urllib.parse.urlparse(node)
        protocol, website, path = url.scheme, url.netloc, url.path
        
        path_with_index = path
        if path_with_index == "" or path_with_index.endswith("/"):
            path_with_index = path_with_index + "index.html"
        if path_with_index.startswith("/"):
            path_with_index = path_with_index.lstrip("/")

        local_path = os.path.join(self.root_dir, protocol, website, path_with_index)

        return protocol, website, path, local_path
    
    def get_html(self, node: Node) -> str | None:
        protocol, website, path, local_path = self.parse_url(node)
        if protocol != "http" and protocol != "https":
            return None # skip non-http/https URLs
        
        if not os.path.exists(local_path):
            os.makedirs(os.path.dirname(local_path), exist_ok=True)
            try:
                node = urllib.parse.quote(node, safe=":/?&=%#")
                with urllib.request.urlopen(node) as response:
                    body = response.read()
                with open(local_path, "wb") as f:
                    f.write(body)

            except urllib.error.HTTPError as e:
                print("HTTP error", e.code, "for", node)
                return None
            except urllib.error.URLError as e:
                print("URL error", e.reason, "for", node)
                return None

        html: str | None = None
        with open(local_path, "rb") as f:
            body = f.read()
        try:
            html = body.decode("utf-8", errors="ignore")
        except Exception:
            return None

        return html

    def visit_node(self, node: Node) -> tuple[Memo, list[Node]]:
        html = self.get_html(node)
        
        children: list[Node] = []
        if html is not None:
            hrefs = parse_hrefs(html)

            absolute_hrefs = []
            for href in hrefs:
                # make it absolute relative to the current page
                abs_url = urllib.parse.urljoin(node, href)
                # drop fragment (#section)
                abs_url, _ = urllib.parse.urldefrag(abs_url)
                absolute_hrefs.append(abs_url)
            
            for href in absolute_hrefs:
                protocol, website, path, local_path = self.parse_url(href)
                if website not in self.whitelist:
                    # remove non-whitelist
                    continue
                children.append(href)
        
        return self, children


def fetch(root_url: str, root_dir: str) -> None:
    website = urllib.parse.urlparse(root_url).netloc
    memo = Memo(root_dir=root_dir, whitelist={website})

    memo = bfs(
        memo=memo,
        root_node=root_url,
        visit_node=Memo.visit_node,
    )


if __name__ == "__main__":
    fetch("https://eldenring.wiki.fextralife.com", "tmp/fextralife_eldenring")

