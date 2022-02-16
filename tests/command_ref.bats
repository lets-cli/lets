setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_ref

}

@test "command ref: run existing command with args from ref" {
    run lets hello-world
    assert_success
    assert_line --index 0 "Hello World"
}

@test "command ref: run existing command with args as list from ref" {
    run lets hello-list
    assert_success
    assert_line --index 0 "Hello Fellow friend"
}
