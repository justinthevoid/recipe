#!/usr/bin/env python3
"""
Simple HTTP server for testing Recipe WASM locally.
Serves files with proper MIME types for WASM.
"""

import http.server
import socketserver
import os

PORT = 8080

class WAMSHandler(http.server.SimpleHTTPRequestHandler):
    """HTTP handler with WASM MIME type support"""

    extensions_map = {
        **http.server.SimpleHTTPRequestHandler.extensions_map,
        '.wasm': 'application/wasm',
        '.js': 'application/javascript',
    }

    def end_headers(self):
        # Add CORS headers for local development
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        http.server.SimpleHTTPRequestHandler.end_headers(self)

if __name__ == '__main__':
    # Change to web directory
    os.chdir(os.path.dirname(os.path.abspath(__file__)))

    with socketserver.TCPServer(("", PORT), WAMSHandler) as httpd:
        print(f"✅ Recipe WASM test server running")
        print(f"   URL: http://localhost:{PORT}")
        print(f"   Press Ctrl+C to stop")
        print()
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\n\nServer stopped.")
