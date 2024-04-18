#!/bin/bash

# Start Firestore emulator on port 8888
echo "Starting Firestore emulator..."
nohup firebase emulators:start --only firestore > firestore_emulator.log 2>&1 &

# Store the process ID of the emulator
EMULATOR_PID=$!

# Wait for Firestore emulator to start
echo "Waiting for Firestore emulator to start..."
sleep 10

# Run Go tests
echo "Running Go tests..."
go test ./...

# Kill Firestore emulator process
echo "Stopping Firestore emulator..."
kill $EMULATOR_PID
