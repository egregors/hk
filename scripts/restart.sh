#!/bin/bash

# Use this script to quick restart t-hk server.

BIN="t-hk-srv"

echo "prev log:"
cat nohup.out 2>/dev/null || echo "(no previous log found)"
rm nohup.out 2>/dev/null || true
echo "---"

echo "kill prev srv"
pkill "$BIN" || { echo "Failed to kill $BIN"; }

nohup ./"$BIN" &
SERVER_PID=$!
sleep 3  # Give the server time to initialize

# Wait a bit more if nohup.out doesn't exist yet
if [ ! -f nohup.out ]; then
    sleep 2
fi

cat nohup.out 2>/dev/null || echo "(no log output yet)"

# Check if the server process is still running
if kill -0 "$SERVER_PID" 2>/dev/null; then
    echo "done"
else
    echo "Server failed to start properly. Check the log above for error details."
    exit 1
fi