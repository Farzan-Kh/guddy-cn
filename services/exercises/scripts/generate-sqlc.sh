#!/bin/sh
docker run --rm -v "$PWD":/src -w /src sqlc/sqlc:1.29.0 generate -f services/exercises/sqlc.yaml