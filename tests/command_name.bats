setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_name

}

@test "command name: can be y o yes" {
    run lets yes
    assert_success
    assert_line --index 0 "Hi from yes"
}
