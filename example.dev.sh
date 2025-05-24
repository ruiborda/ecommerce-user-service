#!/bin/bash

set -e

# Secret para firmar y verificar tokens JWT de autenticación
export JWT_SECRET="your_jwt_secret_here"

# Credenciales de Firebase/GCP en formato base64
export GCP_CREDENTIAL_JSON_BASE64="your_credential_json_base64_here"

# Puerto en el que se ejecutará el servidor (por defecto 8080)
export PORT="${PORT:-8080}"

# Modo de ejecución de Gin (debug/release)
export GIN_MODE=debug

# Ejecutar la aplicación
go run main.go