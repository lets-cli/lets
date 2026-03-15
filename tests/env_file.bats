load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/env_file
}

@test "env_file: should load global env files and keep env precedence" {
    run lets print-global
    assert_success
    assert_line --index 0 "GLOBAL_FROM_FILE=from-global-file"
    assert_line --index 1 "GLOBAL_OVERRIDE=from-global-file"
    assert_line --index 2 "OS_FILE=1"
}

@test "env_file: should load command env files and keep env precedence" {
    run lets print-command
    assert_success
    assert_line --index 0 "COMMAND_FROM_FILE=from-command-file"
    assert_line --index 1 "COMMAND_OVERRIDE=from-command-file"
}

@test "env_file: should ignore optional missing file" {
    run lets print-optional
    assert_success
    assert_line --index 0 "OPTIONAL_OK=1"
}

@test "env_file: should fail on required missing global env file" {
    run lets -c lets.global-missing.yaml noop
    assert_failure
    assert_line --partial "env_file"
    assert_line --partial ".env.required.global.missing"
}

@test "env_file: should fail on required missing command env file" {
    run lets -c lets.command-missing.yaml fail
    assert_failure
    assert_line --partial "command 'fail'"
    assert_line --partial "env_file"
    assert_line --partial ".env.required.command.missing"
}

@test "env_file: should report invalid env file with filename" {
    run lets -c lets.global-invalid.yaml noop
    assert_failure
    assert_line --partial "failed to parse env_file"
    assert_line --partial ".env.invalid"
}
