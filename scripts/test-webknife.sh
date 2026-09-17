#!/usr/bin/env bash
set -euo pipefail

# Webknife manual test harness.
# Each scenario starts a Webknife server and prints curl commands for manual testing.
# Press Ctrl+C to stop the server after inspecting.

WEBKNIFE="${WEBKNIFE:-./bin/webknife}"
HOST="${HOST:-127.0.0.1}"
PORT="${PORT:-8080}"
ADDR="${HOST}:${PORT}"

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

banner() {
    echo ""
    echo -e "${BOLD}${CYAN}============================================================${RESET}"
    echo -e "${BOLD}${CYAN} Webknife Test $1${RESET}"
    echo -e "${BOLD}${CYAN}============================================================${RESET}"
    echo ""
}

info() {
    echo -e "${GREEN}$1${RESET}"
}

cmd() {
    echo -e "  ${BOLD}$1${RESET}"
}

cleanup() {
    if [[ -n "${SERVER_PID:-}" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
        kill "$SERVER_PID" 2>/dev/null || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi
}

trap cleanup EXIT INT TERM

start_server() {
    "$WEBKNIFE" "$@" --listen "$ADDR" &
    SERVER_PID=$!
    sleep 0.5
}

show_help() {
    cat <<EOF
Usage:
  scripts/test-webknife.sh <test>

Manual test harness for Webknife.
Each scenario starts a server in the foreground. Use Ctrl+C to stop.

Tests:
  01  Static file server
  02  Static server with Basic Authentication
  03  Echo / request inspection
  04  Custom HTTP response
  05  Custom error response / HTTP 503
  06  HTTP redirect
  07  Request header manipulation
  08  Response header manipulation
  09  Reverse proxy
  10  Structured JSON logging

Environment variables:
  WEBKNIFE   Path to webknife binary (default: ./bin/webknife)
  HOST       Listen host (default: 127.0.0.1)
  PORT       Listen port (default: 8080)

Examples:
  scripts/test-webknife.sh 01
  PORT=9090 scripts/test-webknife.sh 03
  WEBKNIFE=/usr/local/bin/webknife scripts/test-webknife.sh 01
  scripts/test-webknife.sh help
EOF
}

test_01() {
    banner "01: Static file server"

    TEST_DIR="/tmp/webknife-test"
    mkdir -p "$TEST_DIR"
    echo "<html><body><h1>Hello from Webknife</h1></body></html>" > "$TEST_DIR/index.html"

    info "Start Webknife with:"
    cmd "  $WEBKNIFE serve --listen $ADDR --root $TEST_DIR"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl http://$ADDR/"
    cmd "  curl -I http://$ADDR/"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server serve --root "$TEST_DIR"
    wait "$SERVER_PID" 2>/dev/null || true
}

test_02() {
    banner "02: Static server with Basic Authentication"

    TEST_DIR="/tmp/webknife-test"
    mkdir -p "$TEST_DIR"
    echo "<html><body><h1>Protected content</h1></body></html>" > "$TEST_DIR/index.html"

    info "Start Webknife with:"
    cmd "  $WEBKNIFE serve --listen $ADDR --root $TEST_DIR --auth admin:secret"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl http://$ADDR/                  # 401 Unauthorized"
    cmd "  curl -u admin:secret http://$ADDR/  # 200 OK"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server serve --root "$TEST_DIR" --auth admin:secret
    wait "$SERVER_PID" 2>/dev/null || true
}

test_03() {
    banner "03: Echo / request inspection"

    info "Start Webknife with:"
    cmd "  $WEBKNIFE echo --listen $ADDR"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl http://$ADDR/"
    cmd "  curl -X POST http://$ADDR/test -d 'hello=world'"
    cmd "  curl http://$ADDR/foo?bar=baz"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server echo
    wait "$SERVER_PID" 2>/dev/null || true
}

test_04() {
    banner "04: Custom HTTP response"

    info "Start Webknife with:"
    cmd "  $WEBKNIFE respond --listen $ADDR --status 200 --body 'All systems go'"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl -v http://$ADDR/"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server respond --status 200 --body "All systems go"
    wait "$SERVER_PID" 2>/dev/null || true
}

test_05() {
    banner "05: Custom error response / HTTP 503"

    info "Start Webknife with:"
    cmd "  $WEBKNIFE respond --listen $ADDR --status 503 --body 'Service unavailable'"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl -v http://$ADDR/"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server respond --status 503 --body "Service unavailable"
    wait "$SERVER_PID" 2>/dev/null || true
}

test_06() {
    banner "06: HTTP redirect"

    info "Start Webknife with:"
    cmd "  $WEBKNIFE redirect --listen $ADDR --to https://example.com"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl -v http://$ADDR/"
    cmd "  curl -v http://$ADDR/any/path"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server redirect --to https://example.com
    wait "$SERVER_PID" 2>/dev/null || true
}

test_07() {
    banner "07: Request header manipulation"

    info "Start Webknife with:"
    cmd "  $WEBKNIFE echo --listen $ADDR --set-request-header 'X-Custom: injected'"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl http://$ADDR/"
    echo "  (look for X-Custom in the JSON output)"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server echo --set-request-header "X-Custom: injected"
    wait "$SERVER_PID" 2>/dev/null || true
}

test_08() {
    banner "08: Response header manipulation"

    info "Start Webknife with:"
    cmd "  $WEBKNIFE echo --listen $ADDR --set-response-header 'X-Debug: true'"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl -v http://$ADDR/"
    echo "  (look for X-Debug in response headers)"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server echo --set-response-header "X-Debug: true"
    wait "$SERVER_PID" 2>/dev/null || true
}

test_09() {
    banner "09: Reverse proxy"

    info "This test expects an upstream HTTP server on port 9000."
    info "For example:"
    cmd "  python3 -m http.server 9000 --directory /tmp"
    echo ""
    info "Start Webknife with:"
    cmd "  $WEBKNIFE proxy --listen $ADDR --upstream http://127.0.0.1:9000"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl http://$ADDR/"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server proxy --upstream http://127.0.0.1:9000
    wait "$SERVER_PID" 2>/dev/null || true
}

test_10() {
    banner "10: Structured JSON logging"

    info "Start Webknife with:"
    cmd "  $WEBKNIFE echo --listen $ADDR --log-format json"
    echo ""
    info "Test from another terminal with:"
    cmd "  curl http://$ADDR/"
    echo "  (observe JSON log output from Webknife)"
    echo ""
    info "Press Ctrl+C to stop Webknife."
    echo ""

    start_server echo --log-format json
    wait "$SERVER_PID" 2>/dev/null || true
}

# --- Main ---

if [[ $# -lt 1 ]]; then
    show_help
    exit 1
fi

case "$1" in
    01) test_01 ;;
    02) test_02 ;;
    03) test_03 ;;
    04) test_04 ;;
    05) test_05 ;;
    06) test_06 ;;
    07) test_07 ;;
    08) test_08 ;;
    09) test_09 ;;
    10) test_10 ;;
    help|-h|--help) show_help ;;
    *)
        echo "Unknown test: $1" >&2
        echo "Run 'scripts/test-webknife.sh help' for available tests." >&2
        exit 1
        ;;
esac
