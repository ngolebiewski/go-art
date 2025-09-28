#!/usr/bin/env bash
set -o errexit

# --- Build the frontend ---
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Copy frontend build into backend
rm -rf backend/dist
cp -r frontend/dist backend/

# --- Build the Go backend ---
echo "Building Go backend..."
cd backend
go build -tags netgo -ldflags "-s -w" -o app
cd ..
