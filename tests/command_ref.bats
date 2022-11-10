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

@test "command ref: ref points to non-existing command" {
    run lets -c lets.no-command.yaml hi
    assert_failure
    assert_line --index 0 "lets: config error: failed to parse lets.no-command.yaml: ref 'hi' points to command 'hello' which is not exist"
}
