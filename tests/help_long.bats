load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/help_long
}

@test "help: run 'lets' as is" {
    run lets
    assert_success
    assert_output --partial "A CLI task runner"
    assert_output --partial "super_long_command_longer_than_usual"
    assert_output --partial "bar"
    assert_output --partial "foo"
}
