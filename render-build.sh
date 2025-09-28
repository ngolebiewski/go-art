#!/usr/bin/env bash
set -o errexit

# --- Build the frontend ---
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Copy build output into backend root (next to main.go)
rm -rf dist
cp -r frontend/dist .

# --- Build the Go app ---
echo "Building Go backend..."
go build -tags netgo -ldflags "-s -w" -o app
