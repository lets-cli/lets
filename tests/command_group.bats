load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_group
}

@test "help: running 'lets help' should group commands by their group names" {
    run lets help
    assert_success
}

@test "help: running 'lets --help' should group commands by their group names" {
    run lets --help
    assert_success
}

@test "help: running 'lets' should group commands by their group names" {
    run lets
    assert_success
}
