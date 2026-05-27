#!/bin/sh
# HiveHook CLI installer.
#   curl -fsSL https://hivehook.com/install.sh | sh
#
# Environment overrides:
#   HIVEHOOK_INSTALL   install directory (default: /usr/local/bin, fallback ~/.local/bin)
#   HIVEHOOK_VERSION   version to install (default: latest release)
set -eu

REPO="hivehook/cli"
BIN="hivehook"

# --- styling -----------------------------------------------------------------
if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
	BOLD=$(printf '\033[1m'); RESET=$(printf '\033[0m')
	GREEN=$(printf '\033[38;5;42m'); RED=$(printf '\033[38;5;203m'); GREY=$(printf '\033[38;5;245m')
	ROSE=$(printf '\033[38;2;244;63;94m')  # #f43f5e
else
	BOLD=''; RESET=''; GREEN=''; RED=''; GREY=''; ROSE=''
fi
CHECK="${GREEN}✓${RESET}"

banner() {
	printf '\n'
	printf '   %s⬢%s  %sHiveHook%s\n' "$ROSE" "$RESET" "$BOLD" "$RESET"
	printf '      %sWebhook infrastructure CLI%s\n\n' "$GREY" "$RESET"
}

step()  { printf '   %s%s%s %s\n' "$CHECK" "" "" "$1"; }
info()  { printf '   %s•%s %s\n' "$GREY" "$RESET" "$1"; }
fail()  { printf '   %s✗ %s%s\n' "$RED" "$1" "$RESET" >&2; exit 1; }

# --- preflight ---------------------------------------------------------------
need() { command -v "$1" >/dev/null 2>&1 || fail "required tool not found: $1"; }
need uname; need tar
DL=""
if command -v curl >/dev/null 2>&1; then DL=curl
elif command -v wget >/dev/null 2>&1; then DL=wget
else fail "need curl or wget"; fi

http_get() { # url -> stdout
	if [ "$DL" = curl ]; then curl -fsSL "$1"; else wget -qO- "$1"; fi
}

# --- detect platform ---------------------------------------------------------
detect_platform() {
	os=$(uname -s); arch=$(uname -m)
	case "$os" in
		Linux)  os=linux ;;
		Darwin) os=darwin ;;
		*) fail "unsupported OS: $os (use the Windows installer or 'go install')" ;;
	esac
	case "$arch" in
		x86_64|amd64) arch=amd64 ;;
		arm64|aarch64) arch=arm64 ;;
		*) fail "unsupported architecture: $arch" ;;
	esac
	PLATFORM="${os}_${arch}"
}

# --- resolve version ---------------------------------------------------------
resolve_version() {
	if [ -n "${HIVEHOOK_VERSION:-}" ]; then VERSION="$HIVEHOOK_VERSION"; return; fi
	tag=$(http_get "https://api.github.com/repos/${REPO}/releases/latest" \
		| grep '"tag_name":' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
	[ -n "$tag" ] || fail "could not resolve latest release"
	VERSION="$tag"
}

# --- progress bar ------------------------------------------------------------
# Streams the download to a file and renders a block bar from actual bytes,
# using the Content-Length advertised by the redirect target.
download() { # url dest
	url="$1"; dest="$2"
	total=$(curl -sIL "$url" 2>/dev/null \
		| awk 'BEGIN{IGNORECASE=1} /^content-length:/{v=$2} END{gsub(/\r/,"",v); print v}')

	if [ "$DL" = curl ]; then
		curl -fsSL "$url" -o "$dest" 2>/dev/null &
	else
		wget -qO "$dest" "$url" &
	fi
	pid=$!

	width=24
	while kill -0 "$pid" 2>/dev/null; do
		render_bar "$dest" "$total" "$width"
		sleep 0.08
	done
	wait "$pid" || fail "download failed"
	render_bar "$dest" "$total" "$width"
	printf '\n'
}

render_bar() { # dest total width
	now=0; [ -f "$1" ] && now=$(wc -c < "$1" 2>/dev/null | tr -d ' ')
	tot="$2"; w="$3"
	if [ -n "$tot" ] && [ "$tot" -gt 0 ] 2>/dev/null; then
		pct=$(( now * 100 / tot )); [ "$pct" -gt 100 ] && pct=100
	else
		pct=0
	fi
	fill=$(( pct * w / 100 ))
	filled=''; i=0
	while [ "$i" -lt "$fill" ]; do filled="${filled}⬢"; i=$((i+1)); done
	empty=''
	while [ "$i" -lt "$w" ]; do empty="${empty}⬡"; i=$((i+1)); done
	printf '\r   %sDownloading%s %s%s%s%s%s %s%3d%%%s' \
		"$BOLD" "$RESET" "$ROSE" "$filled" "$GREY" "$empty" "$RESET" "$BOLD" "$pct" "$RESET"
}

# --- install -----------------------------------------------------------------
choose_dir() {
	if [ -n "${HIVEHOOK_INSTALL:-}" ]; then DEST="$HIVEHOOK_INSTALL"; return; fi
	if [ -w /usr/local/bin ] 2>/dev/null; then DEST=/usr/local/bin; return; fi
	DEST="$HOME/.local/bin"
}

main() {
	banner
	detect_platform;                 step "Detected ${BOLD}${PLATFORM}${RESET}"
	resolve_version;                 step "Installing ${BOLD}${BIN} ${VERSION}${RESET}"

	v=${VERSION#v}
	asset="${BIN}_${v}_${PLATFORM}.tar.gz"
	base="https://github.com/${REPO}/releases/download/${VERSION}"
	tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT

	download "${base}/${asset}" "${tmp}/${asset}"

	http_get "${base}/checksums.txt" > "${tmp}/checksums.txt" 2>/dev/null || true
	if [ -s "${tmp}/checksums.txt" ] && command -v shasum >/dev/null 2>&1; then
		want=$(grep " ${asset}\$" "${tmp}/checksums.txt" | awk '{print $1}')
		got=$(shasum -a 256 "${tmp}/${asset}" | awk '{print $1}')
		[ -n "$want" ] && [ "$want" != "$got" ] && fail "checksum mismatch"
		step "Verified checksum"
	fi

	tar -xzf "${tmp}/${asset}" -C "$tmp"
	choose_dir
	mkdir -p "$DEST"
	install -m 0755 "${tmp}/${BIN}" "${DEST}/${BIN}" 2>/dev/null \
		|| { cp "${tmp}/${BIN}" "${DEST}/${BIN}" && chmod 0755 "${DEST}/${BIN}"; }
	step "Installed to ${BOLD}${DEST}/${BIN}${RESET}"

	printf '\n   %sHiveHook CLI is ready.%s\n' "$GREEN$BOLD" "$RESET"
	case ":$PATH:" in
		*":$DEST:"*) : ;;
		*) printf '   %sAdd %s to your PATH:%s\n     export PATH="%s:$PATH"\n' "$GREY" "$DEST" "$RESET" "$DEST" ;;
	esac
	printf '   Run %shivehook login%s to get started.\n\n' "$BOLD" "$RESET"
}

main "$@"
