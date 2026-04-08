import http.server
import socketserver

PORT = 8080

RELOAD_SCRIPT = b"""
<script>
    let lastModWasm = null;
    let lastModHtml = null;
    setInterval(() => {
        fetch('main.wasm', { method: 'HEAD' })
            .then(r => r.headers.get('Last-Modified'))
            .then(mod => {
                if (lastModWasm && lastModWasm !== mod) window.location.reload();
                lastModWasm = mod;
            }).catch(() => {});
        fetch('index.html', { method: 'HEAD' })
            .then(r => r.headers.get('Last-Modified'))
            .then(mod => {
                if (lastModHtml && lastModHtml !== mod) window.location.reload();
                lastModHtml = mod;
            }).catch(() => {});
    }, 1000);
</script>
</body>
"""

class Handler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory="dist", **kwargs)

    def do_GET(self):
        if self.path in ('/', '/index.html'):
            try:
                with open('dist/index.html', 'rb') as f:
                    content = f.read()
                self.send_response(200)
                self.send_header('Content-type', 'text/html')
                self.end_headers()
                # Inject the script right before the closing body tag
                content = content.replace(b'</body>', RELOAD_SCRIPT)
                self.wfile.write(content)
            except FileNotFoundError:
                self.send_error(404, "File not found")
        else:
            super().do_GET()

if __name__ == '__main__':
    socketserver.TCPServer.allow_reuse_address = True
    with socketserver.TCPServer(("", PORT), Handler) as httpd:
        print(f"Serving dev server at: http://localhost:{PORT}")
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\nShutting down dev server...")
