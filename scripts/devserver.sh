#!/usr/bin/env bash
# devserver.sh — keep `shortlink serve` running on $SHORTLINK_ADDR (default
# 127.0.0.1:8787) for zurdo's [http:] criteria, rebuilding and restarting it
# whenever a .go file changes.
#
# It serves the newest of the checkout and every .carol/worktrees/<phase>, so
# while `carol run` drives zurdo in a worktree, the criteria hit that code.
# Each restart loads a fresh copy of testdata/seed.json (when it exists), so
# every check starts from the same links.
#
# Usage: scripts/devserver.sh [dir]   (pin one tree instead of following)
set -u

ROOT=$(cd "$(dirname "$0")/.." && pwd)
ADDR=${SHORTLINK_ADDR:-127.0.0.1:8787}
PIN=${1:-}
STATE=$(mktemp -d "${TMPDIR:-/tmp}/shortlink-dev.XXXXXX")
PID=""

cleanup() { [ -n "$PID" ] && kill "$PID" 2>/dev/null; rm -rf "$STATE"; exit 0; }
trap cleanup INT TERM

log() { printf '%s devserver: %s\n' "$(date +%H:%M:%S)" "$*"; }

# The tree to serve: the pinned dir, else the most recently changed tree.
pick_tree() {
	if [ -n "$PIN" ]; then echo "$PIN"; return; fi
	local best="$ROOT" best_t t d
	best_t=$(stamp "$ROOT" | cut -d' ' -f1)
	for d in "$ROOT"/.carol/worktrees/*/; do
		[ -f "$d/go.mod" ] || continue
		t=$(stamp "$d" | cut -d' ' -f1)
		[ "${t:-0}" -gt "${best_t:-0}" ] && { best=${d%/}; best_t=$t; }
	done
	echo "$best"
}

# "<newest mtime> <checksum of every .go/go.mod/seed mtime>" for a tree.
stamp() {
	find "$1" \( -path '*/.carol' -o -path '*/.git' -o -path '*/.zurdo' \) -prune -o \
		\( -name '*.go' -o -name go.mod -o -name seed.json \) -printf '%T@ %p\n' 2>/dev/null |
		sort -n | awk '{t=$1} {s=s $0} END {printf "%d %s\n", t, s}' | {
			read -r t rest; printf '%s %s\n' "${t:-0}" "$(printf '%s' "$rest" | cksum | cut -d' ' -f1)"
		}
}

start() {
	local tree=$1
	[ -n "$PID" ] && { kill "$PID" 2>/dev/null; wait "$PID" 2>/dev/null; PID=""; }
	if ! (cd "$tree" && go build -o "$STATE/shortlink" ./cmd/shortlink) 2>"$STATE/build.log"; then
		log "build failed in $tree (retrying on next change)"; sed 's/^/    /' "$STATE/build.log" | head -20
		return
	fi
	local args=(serve --addr "$ADDR" --data "$STATE/links.json")
	rm -f "$STATE/links.json"
	[ -f "$tree/testdata/seed.json" ] && cp "$tree/testdata/seed.json" "$STATE/links.json"
	"$STATE/shortlink" "${args[@]}" >"$STATE/serve.log" 2>&1 &
	PID=$!
	sleep 1
	if kill -0 "$PID" 2>/dev/null; then
		log "serving $tree on http://$ADDR (pid $PID)"
	else
		log "shortlink serve exited at once:"; sed 's/^/    /' "$STATE/serve.log" | head -10; PID=""
	fi
}

last=""
log "watching $ROOT and .carol/worktrees/* (Ctrl-C to stop)"
while :; do
	tree=$(pick_tree)
	cur="$tree $(stamp "$tree")"
	if [ "$cur" != "$last" ] || { [ -n "$PID" ] && ! kill -0 "$PID" 2>/dev/null; }; then
		last=$cur
		start "$tree"
	fi
	sleep 2
done
