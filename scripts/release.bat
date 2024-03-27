@echo off
SET /p v=Version:
release ./cmd/release --ldflags="-X main.Version=%v%" --soft-float