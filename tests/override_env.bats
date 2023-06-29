load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/override_env
}

@test "override_env: should use default global env for running command" {
    run lets say_hello_global_env
    assert_success
    assert_line --index 0 "Hello John"
}

@test "override_env: should use default command env for running command" {
    run lets say_hello_command_env
    assert_success
    assert_line --index 0 "Hello Rick"
}

@test "override_env: should override global env for running command with -E" {
    run lets -E NAME=Morty say_hello_global_env
    assert_success
    assert_line --index 0 "Hello Morty"
}

@test "override_env: should override command env for running command with -E" {
    run lets -E NAME=Morty say_hello_command_env
    assert_success
    assert_line --index 0 "Hello Morty"
}

@test "override_env: should override command env for running command with --env" {
    run lets --env NAME=Morty say_hello_command_env
    assert_success
    assert_line --index 0 "Hello Morty"
}

@test "override_env: should set env variable for command with -E even if there is no either global or command env var" {
    run lets -E FOO=BAR print-foo
    assert_success
    assert_line --index 0 "FOO=BAR"
}
