#!/bin/bash

echo "Starting Bowling Score Calculator HTTPS server..."
echo "Open https://localhost:8443 in your browser"
echo "Note: You'll need to accept the self-signed certificate warning"
echo "Press Ctrl+C to stop the server"
echo ""

go run main.go
