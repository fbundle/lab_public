import html
import os
from typing import Any, Callable
import urllib

Node = Any
Memo = Any

def dfs(
    memo: Memo,
    root_node: Node,
    visit_node: Callable[[Memo, Node], tuple[Memo, list[Node]]],
    is_visited: Callable[[Memo, Node], bool],
) -> Memo:
    stack: list[Node] = [root_node]
    while stack:
        node = stack.pop()
        memo, children = visit_node(memo, node)
        for child in children:
            if is_visited(memo, child):
                continue
            stack.append(child)
    return memo


class HrefParser(html.parser.HTMLParser):
    def __init__(self):
        super().__init__()
        self.hrefs: list[str] = []
    
    def handle_starttag(self, tag, attrs):
        if tag != "a":
            return
        for k, v in attrs:
            if k == "href" and v:
                self.hrefs.append(v)

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
        if path == "" or path.endswith("/"):
            path_with_index = path + "index.html"
        else:
            path_with_index = path
        local_path = os.path.join(self.root_dir, protocol, website, path_with_index)
        return protocol, website, path, local_path
    
    def is_visited(self, node: str) -> bool:
        protocol, website, path, local_path = self.parse_url(node)
        if protocol != "http" and protocol != "https":
            return True # skip non-http/https URLs
        return os.path.exists(local_path)
    
    def visit_node(self, node: Node) -> tuple[Memo, list[Node]]:
        protocol, website, path, local_path = self.parse_url(node)
        if protocol != "http" and protocol != "https":
            return self, [] # skip non-http/https URLs
        
        os.makedirs(os.path.dirname(local_path), exist_ok=True)
        with urllib.request.urlopen(node) as response:
            content_type = response.headers.get("Content-Type") or ""
            is_html = "text/html" in content_type.lower()
            body = response.read()
        with open(local_path, "wb") as f:
            f.write(body)
        
        children: list[Node] = []
        if is_html:
            html = body.decode("utf-8", errors="ignore")
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
    """Placeholder for your HTTP/HTML mirror logic."""
    raise NotImplementedError("implement me: fetch(root_url, root_dir)")


if __name__ == "__main__":
    fetch("https://eldenring.wiki.fextralife.com", "tmp/fextralife_eldenring")

