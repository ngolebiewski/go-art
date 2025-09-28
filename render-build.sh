#!/usr/bin/env bash
set -o errexit

echo "=== STARTING BUILD PROCESS ==="

echo "Building frontend..."
cd frontend
npm ci
npm run build

echo "=== FRONTEND BUILD CONTENTS ==="
ls -la dist/
echo "Files in dist:"
find dist -type f | head -10
echo "================================"

cd ..

echo "Copying frontend build to backend..."
rm -rf backend/dist
cp -r frontend/dist backend/

echo "=== BACKEND DIST CONTENTS ==="
ls -la backend/dist/
echo "Files in backend/dist:"
find backend/dist -type f | head -10
echo "============================="

echo "Building Go backend..."
cd backend
go build -tags netgo -ldflags "-s -w" -o app .

echo "=== BUILD COMPLETE ==="