load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/no_lets_file
    cleanup
}

NOT_EXISTED_LETS_FILE="lets-not-existed.yaml"
TEMP_FAKE_BIN_DIR=""
TEMP_OPENED_URL_FILE=""

teardown() {
    if [[ -n "${TEMP_FAKE_BIN_DIR}" ]]; then
        rm -rf "${TEMP_FAKE_BIN_DIR}"
    fi

    if [[ -n "${TEMP_OPENED_URL_FILE}" ]]; then
        rm -f "${TEMP_OPENED_URL_FILE}"
    fi

    TEMP_FAKE_BIN_DIR=""
    TEMP_OPENED_URL_FILE=""
    cleanup
}

@test "no_lets_file: should not create .lets dir" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets

    assert_failure
    [[ ! -d .lets ]]
}

@test "no_lets_file: when wrong config specified with LETS_CONFIG - show find config error" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets

    assert_failure
    assert_output --partial "lets: config error: file does not exist: ${NOT_EXISTED_LETS_FILE}"
}

@test "no_lets_file: show config read error (broken config)" {
    LETS_CONFIG=broken_lets.yaml run lets

    assert_failure

    assert_line --index 0 "lets: config error: failed to parse broken_lets.yaml: yaml: unmarshal errors:"
    assert_line --index 1 "  line 3: cannot unmarshal !!int \`1\` into config.Commands"
}

@test "no_lets_file: show help for 'lets help' even if no config file" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets help

    assert_success
    assert_line --index 0  "A CLI task runner"
}

@test "no_lets_file: show help for 'lets --help' even if no config file" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets --help
    assert_success
    assert_line --index 0  "A CLI task runner"
}

@test "no_lets_file: show help for 'lets -h' even if no config file" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets -h
    assert_success
    assert_line --index 0  "A CLI task runner"
}

@test "no_lets_file: lets self doc opens docs without config" {
    TEMP_FAKE_BIN_DIR="$(mktemp -d)"
    TEMP_OPENED_URL_FILE="$(mktemp)"
    rm -f "${TEMP_OPENED_URL_FILE}"

    cat > "${TEMP_FAKE_BIN_DIR}/xdg-open" <<'EOF'
#!/usr/bin/env bash
printf "%s" "$1" > "${LETS_TEST_OPENED_URL_FILE}"
EOF
    chmod +x "${TEMP_FAKE_BIN_DIR}/xdg-open"

    cat > "${TEMP_FAKE_BIN_DIR}/open" <<'EOF'
#!/usr/bin/env bash
printf "%s" "$1" > "${LETS_TEST_OPENED_URL_FILE}"
EOF
    chmod +x "${TEMP_FAKE_BIN_DIR}/open"

    PATH="${TEMP_FAKE_BIN_DIR}:${PATH}" \
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} \
    LETS_TEST_OPENED_URL_FILE="${TEMP_OPENED_URL_FILE}" \
    run lets self doc

    assert_success

    for _ in $(seq 1 20); do
        if [[ -f "${TEMP_OPENED_URL_FILE}" ]]; then
            break
        fi
        sleep 0.1
    done

    run cat "${TEMP_OPENED_URL_FILE}"
    assert_success
    assert_output "https://lets-cli.org/docs/config"
}
