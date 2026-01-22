#!/bin/sh

set -e

echo "Starting chat application..."

echo "Waiting for database..."
sleep 3

echo "Running database migrations..."
make migrate

echo "Building application..."
make build

echo "Starting server..."
exec make run