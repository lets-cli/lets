load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/no_lets_file
    cleanup
}

NOT_EXISTED_LETS_FILE="lets-not-existed.yaml"

@test "no_lets_file: should not create .lets dir" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets


    assert_failure
    [[ ! -d .lets ]]
}

@test "no_lets_file: when wrong config specified with LETS_CONFIG - show find config error" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets

    assert_failure
    assert_output --partial "failed to find config file ${NOT_EXISTED_LETS_FILE} in"
}

@test "no_lets_file: show config read error (broken config)" {
    LETS_CONFIG=broken_lets.yaml run lets

    assert_failure

    assert_line --index 0 "lets: failed to parse broken_lets.yaml: yaml: unmarshal errors:"
    assert_line --index 1 "  line 3: cannot unmarshal !!int \`1\` into config.Commands"
}

@test "no_lets_file: show help for 'lets help' even if no config file" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets help

    assert_success
    assert_line --index 0  "A CLI command runner"
}

@test "no_lets_file: show help for 'lets -h' even if no config file" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets -h
    assert_success
    assert_line --index 0  "A CLI command runner"
}
