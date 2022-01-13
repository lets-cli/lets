load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_work_dir
}

@test "command_work_dir: should run command in work_dir" {
    run lets print-file
    assert_success
    assert_line --index 0 "hi there"
}
