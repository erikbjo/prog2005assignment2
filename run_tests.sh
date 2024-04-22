#!/bin/bash

# Start Firestore emulator on port 8888
echo "Starting Firestore emulator..."
nohup firebase emulators:start --only firestore > firestore_emulator.log 2>&1 &

# Wait for Firestore emulator to start
echo "Waiting for Firestore emulator to start..."
sleep 10

# Run Go tests
echo "Running Go tests..."
go test ./... -coverprofile=coverage.out

# Kill Firestore emulator process
echo "Stopping Firestore emulator..."
EMULATOR_PORT=8080
EMULATOR_PID=$(lsof -t -i tcp:$EMULATOR_PORT)

if [ -z "$EMULATOR_PID" ]; then
    echo "No Firestore emulator process found on port $EMULATOR_PORT."
else
    echo "Stopping Firestore emulator (PID: $EMULATOR_PID)..."
    kill -9 $EMULATOR_PID
fi