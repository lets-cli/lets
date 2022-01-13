load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/find_config/child/another_child
    rm -rf ../../.lets
    cleanup
}

@test "find_lets_file: should find lets.yaml in parent dir" {
    run lets foo
    assert_success
    assert_line --index 0 "foo"
}

@test "find_lets_file: .lets must be created in the same dir where lets.yaml placed" {
    run lets foo
    assert_success

    [[ ! -d .lets ]]
    [[ -d ../../.lets ]]
}
