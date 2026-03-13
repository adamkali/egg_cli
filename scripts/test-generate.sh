#!/usr/bin/env bash
# test-generate.sh
#
# End-to-end test for egg_cli generate for both JWT and Auth0 templates.
# Spins up a postgres container, generates both project types, runs
# go build ./... to confirm everything compiles, then cleans up.
#
# Usage:
#   chmod +x scripts/test-generate.sh
#   ./scripts/test-generate.sh

set -e

EGG_CLI="$(cd "$(dirname "$0")/.." && pwd)/egg_cli"
TESTROOT="/tmp/egg-gen-test"
PG_CONTAINER="egg-gen-test-pg"
PG_USER="postgres"
PG_PASS="eggtest123"
PG_PORT="5499"      # non-standard port to avoid collisions

# ── Colours ──────────────────────────────────────────────────────────────────
GREEN='\033[0;32m'; RED='\033[0;31m'; YELLOW='\033[1;33m'; NC='\033[0m'
info()  { echo -e "${YELLOW}[test]${NC} $*"; }
ok()    { echo -e "${GREEN}[pass]${NC} $*"; }
fail()  { echo -e "${RED}[fail]${NC} $*"; exit 1; }

# ── Cleanup on exit ───────────────────────────────────────────────────────────
cleanup() {
  info "Stopping postgres container..."
  docker stop "$PG_CONTAINER" 2>/dev/null || true
  docker rm   "$PG_CONTAINER" 2>/dev/null || true
}
trap cleanup EXIT

# ── Ensure binary exists ──────────────────────────────────────────────────────
if [ ! -x "$EGG_CLI" ]; then
  info "Building egg_cli..."
  make -C "$(dirname "$EGG_CLI")" build
fi

# ── Start postgres ────────────────────────────────────────────────────────────
info "Starting postgres on port $PG_PORT..."
docker run -d \
  --name "$PG_CONTAINER" \
  -e POSTGRES_USER="$PG_USER" \
  -e POSTGRES_PASSWORD="$PG_PASS" \
  -e POSTGRES_DB="postgres" \
  -p "${PG_PORT}:5432" \
  postgres:17-alpine \
  > /dev/null

info "Waiting for postgres to be ready..."
for i in $(seq 1 30); do
  docker exec "$PG_CONTAINER" pg_isready -U "$PG_USER" -q && break
  [ "$i" -eq 30 ] && fail "Postgres did not become ready in time"
  sleep 1
done
ok "Postgres is ready"

# ── Helper: create a database ─────────────────────────────────────────────────
create_db() {
  docker exec "$PG_CONTAINER" \
    psql -U "$PG_USER" -c "CREATE DATABASE \"$1\";" > /dev/null
}

# ── Helper: run a single project test ─────────────────────────────────────────
# $1 = test name (jwt|auth0)
# $2 = extra flag for generate example (empty or "--auth0")
# $3 = database name
run_test() {
  local name="$1" extra_flag="$2" dbname="$3"
  local dir="$TESTROOT/$name"

  info "── $name test ────────────────────────────────────────"
  rm -rf "$dir" && mkdir -p "$dir"
  create_db "$dbname"

  # Generate the example egg.yaml
  (cd "$dir" && "$EGG_CLI" generate example $extra_flag) || fail "$name: generate example failed"

  # Patch the DB URL to point at our test postgres
  local db_url="postgres://${PG_USER}:${PG_PASS}@localhost:${PG_PORT}/${dbname}?sslmode=disable"
  sed -i "s|url: postgres://.*|url: ${db_url}|" "$dir/egg.yaml"

  # Run the full generate pipeline
  info "Running egg_cli generate for $name..."
  (cd "$dir" && "$EGG_CLI" generate --config egg.yaml) \
    || fail "$name: egg_cli generate failed"

  local proj="$dir/example"

  # Verify no unreplaced placeholders in Go/config files
  if grep -r "__EGG_" "$proj" --include="*.go" --include="*.yml" \
       --include="*.yaml" --include="*.mod" --include="*.env*" 2>/dev/null | grep -v ".git"; then
    fail "$name: found unreplaced __EGG_*__ placeholders"
  fi
  ok "$name: no unreplaced placeholders"

  # Verify go build
  info "Running go build ./... in $proj..."
  (cd "$proj" && go build ./...) || fail "$name: go build ./... failed"
  ok "$name: go build ./... succeeded"

  ok "── $name PASSED ────────────────────────────────────────"
}

# ── Run tests ─────────────────────────────────────────────────────────────────
rm -rf "$TESTROOT" && mkdir -p "$TESTROOT"

run_test "jwt"   ""       "example_jwt"
run_test "auth0" "--auth0" "example_auth0"

echo ""
ok "All tests passed!"
