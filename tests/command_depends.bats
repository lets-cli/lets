load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_depends
}

@test "command_depends: should run all depends commands before main command" {
    run lets run-with-depends
    assert_success
    assert_line --index 0 "Hello world with level INFO"
    assert_line --index 1 "Bar"
    assert_line --index 2 "Main"
}

@test "command_depends: should override args" {
    run lets override-args
    assert_success
    assert_line --index 0 "Hello Developer with level INFO"
    assert_line --index 1 "Bar"
    assert_line --index 2 "Override args"
}

@test "command_depends: should override env" {
    run lets override-env
    assert_success
    assert_line --index 0 "Hello world with level DEBUG"
    assert_line --index 1 "Override env"
}
