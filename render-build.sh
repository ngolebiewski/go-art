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

echo "Moving frontend build to backend for embedding..."
rm -rf backend/dist
mv frontend/dist backend/

echo "=== BACKEND DIST CONTENTS (FOR EMBEDDING) ==="
ls -la backend/dist/
echo "Files in backend/dist:"
find backend/dist -type f | head -10
echo "=============================================="

echo "Building Go backend with embedded files..."
cd backend

# Set production environment for hybrid detection
export NODE_ENV=production
export RENDER=true

echo "Environment variables set:"
echo "  NODE_ENV=$NODE_ENV"
echo "  RENDER=$RENDER"

go build -tags netgo -ldflags "-s -w" -o app .

echo "=== VERIFYING GO BINARY ==="
ls -la app
file app
echo "=========================="

echo "=== BUILD COMPLETE ==="
echo "✅ Frontend built and moved to backend/dist/"
echo "✅ Go binary built with embedded static files"
echo "✅ Single binary contains everything!"