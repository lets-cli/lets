#!/usr/bin/env bash
set -e

LETS_HOME="${LETS_HOME:-$HOME/.lets}"
BIN_DIR="${LETS_HOME}/bin"
LETS_VERSION="${LETS_VERSION:-}"

usage() {
  this=$1
  cat <<EOF
$this: download lets binary for lets-cli/lets

Usage: $this [-d] [tag]
  -d turns on debug logging
   [tag] is a tag from
   https://github.com/lets-cli/lets/releases
   If tag is missing, then the latest will be used.
   LETS_VERSION environment variable overrides [tag].
   LETS_HOME sets installation home, Defaults to $HOME/.lets

EOF
  exit 2
}

parse_args() {
  while getopts "dh?x" arg; do
    case "$arg" in
      d) log_set_priority 10 ;;
      h | \?) usage "$0" ;;
      x) set -x ;;
    esac
  done
  shift $((OPTIND - 1))
  TAG=$1
  if [ -n "${LETS_VERSION}" ]; then
    TAG="${LETS_VERSION}"
  fi
}
# this function wraps all the destructive operations
# if a curl|bash cuts off the end of the script due to
# network, either nothing will happen or will syntax error
# out preventing half-done work
execute() {
  tmpdir=$(mktemp -d)
  log_debug "downloading files into ${tmpdir}"
  http_download "${tmpdir}/${TARBALL}" "${TARBALL_URL}"
  http_download "${tmpdir}/${CHECKSUM}" "${CHECKSUM_URL}"
  hash_sha256_verify "${tmpdir}/${TARBALL}" "${tmpdir}/${CHECKSUM}"
  srcdir="${tmpdir}"
  (cd "${tmpdir}" && untar "${TARBALL}")
  test ! -d "${BIN_DIR}" && install -d "${BIN_DIR}"
  for binexe in $BINARIES; do
    if [ "$OS" = "windows" ]; then
      binexe="${binexe}.exe"
    fi
    install "${srcdir}/${binexe}" "${BIN_DIR}/"
    log_info "installed ${BIN_DIR}/${binexe}"
  done
  rm -rf "${tmpdir}"
}
get_binaries() {
  case "$PLATFORM" in
    darwin/amd64) BINARIES="lets" ;;
    darwin/arm64) BINARIES="lets" ;;
    linux/386) BINARIES="lets" ;;
    linux/amd64) BINARIES="lets" ;;
    *)
      log_crit "platform $PLATFORM is not supported.  Make sure this script is up-to-date and file request at https://github.com/${PREFIX}/issues/new"
      exit 1
      ;;
  esac
}
tag_to_version() {
  if [ -z "${TAG}" ]; then
    log_info "checking GitHub for latest tag"
  else
    log_info "checking GitHub for tag '${TAG}'"
  fi
  REALTAG=$(github_release "$OWNER/$REPO" "${TAG}") && true
  if test -z "$REALTAG"; then
    log_crit "unable to find '${TAG}' - use 'latest' or see https://github.com/${PREFIX}/releases for details"
    exit 1
  fi
  # if version starts with 'v', remove it
  TAG="$REALTAG"
  VERSION=${TAG#v}
}
adjust_format() {
  # change format (tar.gz or zip) based on OS
  true
}
adjust_os() {
  # adjust archive name based on OS
  case ${OS} in
    386) OS=i386 ;;
    amd64) OS=x86_64 ;;
    darwin) OS=Darwin ;;
    linux) OS=Linux ;;
    windows) OS=Windows ;;
  esac
  true
}
adjust_arch() {
  # adjust archive name based on ARCH
  case ${ARCH} in
    386) ARCH=i386 ;;
    amd64) ARCH=x86_64 ;;
    darwin) ARCH=Darwin ;;
    linux) ARCH=Linux ;;
    windows) ARCH=Windows ;;
  esac
  true
}

cat /dev/null <<EOF
------------------------------------------------------------------------
https://github.com/client9/shlib - portable posix shell functions
Public domain - http://unlicense.org
https://github.com/client9/shlib/blob/master/LICENSE.md
but credit (and pull requests) appreciated.
------------------------------------------------------------------------
EOF
is_command() {
  command -v "$1" >/dev/null
}
echoerr() {
  echo "$@" 1>&2
}
log_prefix() {
  echo "$0"
}
_logp=6
log_set_priority() {
  _logp="$1"
}
log_priority() {
  if test -z "$1"; then
    echo "$_logp"
    return
  fi
  [ "$1" -le "$_logp" ]
}
log_tag() {
  case $1 in
    0) echo "emerg" ;;
    1) echo "alert" ;;
    2) echo "crit" ;;
    3) echo "err" ;;
    4) echo "warning" ;;
    5) echo "notice" ;;
    6) echo "info" ;;
    7) echo "debug" ;;
    *) echo "$1" ;;
  esac
}
log_debug() {
  log_priority 7 || return 0
  echoerr "$(log_prefix)" "$(log_tag 7)" "$@"
}
log_info() {
  log_priority 6 || return 0
  echoerr "$(log_prefix)" "$(log_tag 6)" "$@"
}
log_err() {
  log_priority 3 || return 0
  echoerr "$(log_prefix)" "$(log_tag 3)" "$@"
}
log_crit() {
  log_priority 2 || return 0
  echoerr "$(log_prefix)" "$(log_tag 2)" "$@"
}
resolve_path() {
  source_path=$1

  while [ -L "$source_path" ]; do
    source_dir=$(cd -P "$(dirname "$source_path")" >/dev/null 2>&1 && pwd) || return 1
    source_path=$(readlink "$source_path")
    case "$source_path" in
      /*) ;;
      *) source_path="${source_dir}/${source_path}" ;;
    esac
  done

  source_dir=$(cd -P "$(dirname "$source_path")" >/dev/null 2>&1 && pwd) || return 1
  echo "${source_dir}/$(basename "$source_path")"
}
is_homebrew_lets() {
  resolved_path=$(resolve_path "$1") || return 1
  case "$resolved_path" in
    */Cellar/lets/*) return 0 ;;
    *) return 1 ;;
  esac
}
check_old_usr_local_install() {
  old_path="/usr/local/bin/lets"

  if [ ! -e "$old_path" ]; then
    return
  fi

  if is_homebrew_lets "$old_path"; then
    log_info "detected Homebrew-managed ${old_path}; leaving it untouched"
    return
  fi

  log_crit "found old system-wide lets installation at ${old_path}"
  log_crit "remove it before continuing:"
  log_crit "  sudo rm ${old_path}"
  log_crit "then run this installer again"
  exit 1
}
dir_in_path() {
  check_dir=$1

  if [ -d "$check_dir" ]; then
    check_dir=$(cd "$check_dir" >/dev/null 2>&1 && pwd) || return 1
  fi

  echo ":$PATH:" | grep -q ":$check_dir:"
}
try_symlink_in_path() {
  binary_name=$1
  preferred_dirs=(
    "$HOME/.local/bin"
    "$HOME/bin"
    "$HOME/.bin"
  )

  for dir in "${preferred_dirs[@]}"; do
    if ! dir_in_path "$dir"; then
      continue
    fi

    mkdir -p "$dir" 2>/dev/null || continue

    symlink_path="${dir}/${binary_name}"
    target_path="${BIN_DIR}/${binary_name}"

    if [ -L "$symlink_path" ]; then
      rm -f "$symlink_path"
    fi

    if ln -sf "$target_path" "$symlink_path" 2>/dev/null; then
      log_info "created symlink: ${symlink_path} -> ${target_path}"
      return 0
    fi
  done

  return 1
}
update_shell_profile() {
  binary_name="lets"

  if try_symlink_in_path "$binary_name"; then
    return
  fi

  local_bin_dir="$HOME/.local/bin"
  mkdir -p "$local_bin_dir" 2>/dev/null || true
  symlink_path="${local_bin_dir}/${binary_name}"
  target_path="${BIN_DIR}/${binary_name}"

  if [ -L "$symlink_path" ]; then
    rm -f "$symlink_path"
  fi

  if ln -sf "$target_path" "$symlink_path" 2>/dev/null; then
    log_info "created symlink: ${symlink_path} -> ${target_path}"
  else
    log_err "could not create symlink in ${local_bin_dir}"
    log_err "please add ${BIN_DIR} to your PATH manually:"
    echo "  export PATH=\"${BIN_DIR}:\$PATH\""
    return
  fi

  default_shell="bash"
  if [ "$(uname -s)" = "Darwin" ]; then
    default_shell="zsh"
  fi

  os_name=$(uname -s)
  shell_name=$(basename "${SHELL:-$default_shell}")
  shell_profile=""
  path_export=""

  case "$shell_name" in
    zsh)
      shell_profile="$HOME/.zshrc"
      path_export="export PATH=\"\$HOME/.local/bin:\$PATH\""
      ;;
    bash)
      if [ "$os_name" = "Darwin" ]; then
        if [ -f "$HOME/.bash_profile" ]; then
          shell_profile="$HOME/.bash_profile"
        elif [ -f "$HOME/.bashrc" ]; then
          shell_profile="$HOME/.bashrc"
        else
          shell_profile="$HOME/.bash_profile"
        fi
      else
        if [ -f "$HOME/.bashrc" ]; then
          shell_profile="$HOME/.bashrc"
        elif [ -f "$HOME/.bash_profile" ]; then
          shell_profile="$HOME/.bash_profile"
        else
          shell_profile="$HOME/.bashrc"
        fi
      fi
      path_export="export PATH=\"\$HOME/.local/bin:\$PATH\""
      ;;
    fish)
      shell_profile="$HOME/.config/fish/config.fish"
      path_export="fish_add_path \"\$HOME/.local/bin\""
      ;;
    *)
      log_err "unknown shell: ${shell_name}"
      log_err "please add ~/.local/bin to your PATH manually:"
      echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
      return
      ;;
  esac

  if [ -f "$shell_profile" ] && grep -v '^[[:space:]]*#' "$shell_profile" 2>/dev/null | grep -qE 'PATH=.*\.local/bin|fish_add_path .*\.local/bin'; then
    log_info "~/.local/bin already configured in ${shell_profile/#$HOME/\~}"
    echo ""
    log_info "to use lets immediately, run:"
    echo "  ${path_export}"
    return
  fi

  tilde_profile="${shell_profile/#$HOME/\~}"
  echo ""
  if [ -t 0 ]; then
    read -r -p "Add ~/.local/bin to your PATH in ${tilde_profile}? [y/n] " -n 1
    echo ""
    case "$REPLY" in
      [Yy]) ;;
      *)
        log_info "skipped modifying shell config"
        log_info "to use lets, add ~/.local/bin to your PATH manually:"
        echo "  ${path_export}"
        return
        ;;
    esac
  else
    log_info "adding ~/.local/bin to PATH in ${tilde_profile}"
  fi

  if [ ! -f "$shell_profile" ]; then
    mkdir -p "$(dirname "$shell_profile")"
    touch "$shell_profile"
  fi

  {
    echo ""
    echo "# lets"
    echo "$path_export"
  } >>"$shell_profile"

  log_info "added ~/.local/bin to PATH in ${tilde_profile}"
  echo ""
  log_info "to use lets immediately, run:"
  echo "  ${path_export}"
}
uname_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$os" in
    cygwin_nt*) os="windows" ;;
    mingw*) os="windows" ;;
    msys_nt*) os="windows" ;;
  esac
  echo "$os"
}
uname_arch() {
  arch=$(uname -m)
  case $arch in
    x86_64) arch="amd64" ;;
    x86) arch="386" ;;
    i686) arch="386" ;;
    i386) arch="386" ;;
    aarch64) arch="arm64" ;;
    armv5*) arch="armv5" ;;
    armv6*) arch="armv6" ;;
    armv7*) arch="armv7" ;;
  esac
  echo ${arch}
}
uname_os_check() {
  os=$(uname_os)
  case "$os" in
    darwin) return 0 ;;
    dragonfly) return 0 ;;
    freebsd) return 0 ;;
    linux) return 0 ;;
    android) return 0 ;;
    nacl) return 0 ;;
    netbsd) return 0 ;;
    openbsd) return 0 ;;
    plan9) return 0 ;;
    solaris) return 0 ;;
    windows) return 0 ;;
  esac
  log_crit "uname_os_check '$(uname -s)' got converted to '$os' which is not a GOOS value. Please file bug at https://github.com/client9/shlib"
  return 1
}
uname_arch_check() {
  arch=$(uname_arch)
  case "$arch" in
    386) return 0 ;;
    amd64) return 0 ;;
    arm64) return 0 ;;
    armv5) return 0 ;;
    armv6) return 0 ;;
    armv7) return 0 ;;
    ppc64) return 0 ;;
    ppc64le) return 0 ;;
    mips) return 0 ;;
    mipsle) return 0 ;;
    mips64) return 0 ;;
    mips64le) return 0 ;;
    s390x) return 0 ;;
    amd64p32) return 0 ;;
  esac
  log_crit "uname_arch_check '$(uname -m)' got converted to '$arch' which is not a GOARCH value.  Please file bug report at https://github.com/client9/shlib"
  return 1
}
untar() {
  tarball=$1
  case "${tarball}" in
    *.tar.gz | *.tgz) tar --no-same-owner -xzf "${tarball}" ;;
    *.tar) tar --no-same-owner -xf "${tarball}" ;;
    *.zip) unzip "${tarball}" ;;
    *)
      log_err "untar unknown archive format for ${tarball}"
      return 1
      ;;
  esac
}
http_download_curl() {
  local_file=$1
  source_url=$2
  header=$3
  if [ -z "$header" ]; then
    code=$(curl -w '%{http_code}' -sL -o "$local_file" "$source_url")
  else
    code=$(curl -w '%{http_code}' -sL -H "$header" -o "$local_file" "$source_url")
  fi
  if [ "$code" != "200" ]; then
    log_debug "http_download_curl received HTTP status $code"
    return 1
  fi
  return 0
}
http_download_wget() {
  local_file=$1
  source_url=$2
  header=$3
  if [ -z "$header" ]; then
    wget -q -O "$local_file" "$source_url"
  else
    wget -q --header "$header" -O "$local_file" "$source_url"
  fi
}
http_download() {
  log_debug "http_download $2"
  if is_command curl; then
    http_download_curl "$@"
    return
  elif is_command wget; then
    http_download_wget "$@"
    return
  fi
  log_crit "http_download unable to find wget or curl"
  return 1
}
http_copy() {
  tmp=$(mktemp)
  http_download "${tmp}" "$1" "$2" || return 1
  body=$(cat "$tmp")
  rm -f "${tmp}"
  echo "$body"
}
github_release() {
  owner_repo=$1
  version=$2
  test -z "$version" && version="latest"
  giturl="https://github.com/${owner_repo}/releases/${version}"
  json=$(http_copy "$giturl" "Accept:application/json")
  test -z "$json" && return 1
  version=$(echo "$json" | tr -s '\n' ' ' | sed 's/.*"tag_name":"//' | sed 's/".*//')
  test -z "$version" && return 1
  echo "$version"
}
hash_sha256() {
  TARGET=${1:-/dev/stdin}
  if is_command gsha256sum; then
    hash=$(gsha256sum "$TARGET") || return 1
    echo "$hash" | cut -d ' ' -f 1
  elif is_command sha256sum; then
    hash=$(sha256sum "$TARGET") || return 1
    echo "$hash" | cut -d ' ' -f 1
  elif is_command shasum; then
    hash=$(shasum -a 256 "$TARGET" 2>/dev/null) || return 1
    echo "$hash" | cut -d ' ' -f 1
  elif is_command openssl; then
    hash=$(openssl -dst openssl dgst -sha256 "$TARGET") || return 1
    echo "$hash" | cut -d ' ' -f a
  else
    log_crit "hash_sha256 unable to find command to compute sha-256 hash"
    return 1
  fi
}
hash_sha256_verify() {
  TARGET=$1
  checksums=$2
  if [ -z "$checksums" ]; then
    log_err "hash_sha256_verify checksum file not specified in arg2"
    return 1
  fi
  BASENAME=${TARGET##*/}
  want=$(grep "${BASENAME}" "${checksums}" 2>/dev/null | tr '\t' ' ' | cut -d ' ' -f 1)
  if [ -z "$want" ]; then
    log_err "hash_sha256_verify unable to find checksum for '${TARGET}' in '${checksums}'"
    return 1
  fi
  got=$(hash_sha256 "$TARGET")
  if [ "$want" != "$got" ]; then
    log_err "hash_sha256_verify checksum for '$TARGET' did not verify ${want} vs $got"
    return 1
  fi
}
cat /dev/null <<EOF
------------------------------------------------------------------------
End of functions from https://github.com/client9/shlib
------------------------------------------------------------------------
EOF

PROJECT_NAME="lets"
OWNER=lets-cli
REPO="lets"
BINARY=lets
FORMAT=tar.gz
OS=$(uname_os)
ARCH=$(uname_arch)
PREFIX="$OWNER/$REPO"

# use in logging routines
log_prefix() {
	echo "$PREFIX"
}
PLATFORM="${OS}/${ARCH}"
GITHUB_DOWNLOAD=https://github.com/${OWNER}/${REPO}/releases/download

uname_os_check "$OS"
uname_arch_check "$ARCH"

parse_args "$@"

check_old_usr_local_install

get_binaries

tag_to_version

adjust_format

adjust_os

adjust_arch

log_info "found version: ${VERSION} for ${TAG}/${OS}/${ARCH}"

NAME=${PROJECT_NAME}_${OS}_${ARCH}
TARBALL=${NAME}.${FORMAT}
TARBALL_URL=${GITHUB_DOWNLOAD}/${TAG}/${TARBALL}
CHECKSUM=${PROJECT_NAME}_checksums.txt
CHECKSUM_URL=${GITHUB_DOWNLOAD}/${TAG}/${CHECKSUM}


execute

update_shell_profile
