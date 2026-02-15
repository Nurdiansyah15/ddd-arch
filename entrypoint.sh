#!/bin/sh
set -e

# (Opsional) Generate ulang Swagger docs saat container start
/go/bin/swag init -g cmd/main.go --output docs --parseDependency --parseInternal=false || echo "Swagger init failed, maybe already generated."

# Jalankan Air (live reload)
exec air -c .air.toml
