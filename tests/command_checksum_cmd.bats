load test_helpers

reset_test_files() {
    cleanup
    printf "first-checksum" > project/checksum.txt
}

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_checksum_cmd
    reset_test_files
}

teardown() {
    reset_test_files
}

@test "command_checksum_cmd: should use command shell and work_dir and persist checksum" {
    run lets checksum-cmd

    assert_success
    assert_line --index 0 "LETS_CHECKSUM=first-checksum"
    assert_line --index 1 "LETS_CHECKSUM_CHANGED=true"
    [[ -f .lets/checksums/checksum-cmd/lets_default_checksum ]]

    run lets checksum-cmd

    assert_success
    assert_line --index 0 "LETS_CHECKSUM=first-checksum"
    assert_line --index 1 "LETS_CHECKSUM_CHANGED=false"

    printf "second-checksum" > project/checksum.txt

    run lets checksum-cmd

    assert_success
    assert_line --index 0 "LETS_CHECKSUM=second-checksum"
    assert_line --index 1 "LETS_CHECKSUM_CHANGED=true"
}
