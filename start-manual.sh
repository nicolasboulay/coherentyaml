#!/usr/bin/env bash
set -euo pipefail

exec docker compose run --rm --build --service-ports codex bash
