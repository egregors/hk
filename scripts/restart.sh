#!/bin/bash

# Use this script to quick restart t-hk server.

BIN="t-hk-srv"

echo "prev log:"
cat nohup.out
rm nohup.out
echo "---"

echo "kill prev srv"
pkill "$BIN" || { echo "Failed to kill $BIN"; exit 1; }

nohup ./"$BIN" &
sleep 1  # Give nohup some time to create the output file
cat nohup.out

echo "done"