setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_not_found
}

@test "command_not_found: exit code is 2 when command does not exist" {
    run lets no_such_command
    assert_failure 2
}

@test "command_not_found: exit code is 2 when self subcommand does not exist" {
    run lets self no_such_command
    assert_failure 2
}
