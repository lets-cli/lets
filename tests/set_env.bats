load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/set_env
}

@test "set_env: should use default global env for running command" {
    run lets say_hello_global
    assert_success
    assert_line --index 0 "Hello John"
}

@test "set_env: should use default command env for running command" {
    run lets say_hello
    assert_success
    assert_line --index 0 "Hello Rick"
}

@test "set_env: should override env for running command with -E" {
    run lets -E NAME=Morty say_hello
    assert_success
    assert_line --index 0 "Hello Morty"
}

@test "set_env: should override env for running command with --env" {
    run lets --env NAME=Morty say_hello
    assert_success
    assert_line --index 0 "Hello Morty"
}

@test "say_command: should set env var LETS_COMMAND_NAME" {
    run lets say_command
    assert_success
    assert_line --index 0 "say_command"
}
