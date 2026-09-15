import http.server
import socketserver
import urllib.request
import subprocess
import os

PORT = 8085
TARGET = "https://rag-comparison-demo-4gwhjlgy7a-ew.a.run.app"

class ProxyHandler(http.server.BaseHTTPRequestHandler):
    def get_token(self):
        try:
            return subprocess.check_output(
                ["gcloud", "auth", "print-identity-token", "--account=william@hoffmannw.altostrat.com"]
            ).decode().strip()
        except Exception:
            return ""

    def do_GET(self):
        url = TARGET + self.path
        req = urllib.request.Request(url)
        token = self.get_token()
        if token:
            req.add_header("Authorization", f"Bearer {token}")
        
        try:
            with urllib.request.urlopen(req) as response:
                self.send_response(response.status)
                for header, val in response.getheaders():
                    if header.lower() not in ["content-length", "transfer-encoding", "connection"]:
                        self.send_header(header, val)
                self.end_headers()
                
                # Stream SSE or binary/html content
                while True:
                    chunk = response.read(1024)
                    if not chunk:
                        break
                    self.wfile.write(chunk)
                    self.wfile.flush()
        except urllib.error.HTTPError as e:
            self.send_response(e.code)
            self.end_headers()
            self.wfile.write(e.read())
        except Exception as e:
            self.send_response(500)
            self.end_headers()
            self.wfile.write(str(e).encode())

    def do_POST(self):
        url = TARGET + self.path
        length = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(length) if length > 0 else None
        
        req = urllib.request.Request(url, data=body, method="POST")
        token = self.get_token()
        if token:
            req.add_header("Authorization", f"Bearer {token}")
        if self.headers.get("Content-Type"):
            req.add_header("Content-Type", self.headers.get("Content-Type"))

        try:
            with urllib.request.urlopen(req) as response:
                self.send_response(response.status)
                for header, val in response.getheaders():
                    if header.lower() not in ["content-length", "transfer-encoding", "connection"]:
                        self.send_header(header, val)
                self.end_headers()
                self.wfile.write(response.read())
        except urllib.error.HTTPError as e:
            self.send_response(e.code)
            self.end_headers()
            self.wfile.write(e.read())

if __name__ == "__main__":
    socketserver.TCPServer.allow_reuse_address = True
    with socketserver.TCPServer(("", PORT), ProxyHandler) as httpd:
        print(f"Proxy local démarré sur http://0.0.0.0:{PORT}")
        httpd.serve_forever()
