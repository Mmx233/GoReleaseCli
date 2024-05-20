@echo off
SET /p v=Version:
go run ./cmd/release ./cmd/release --ldflags="-X main.Version=%v%" --extra-arches -c tar.gz