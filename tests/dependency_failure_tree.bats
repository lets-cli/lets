load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/dependency_failure_tree
}

@test "dependency_failure_tree: shows full 3-level chain on failure" {
    run env NO_COLOR=1 lets deploy
    assert_failure
    assert_line --index 0 "lets: command failed:"
    assert_line --index 1 "  deploy"
    assert_line --index 2 "  └─ build"
    assert_line --index 3 "    └─ lint  <-- failed here"
    assert_line --index 4 "lets: exit status 1"
}

@test "dependency_failure_tree: single node when no depends" {
    run env NO_COLOR=1 lets lint
    assert_failure
    assert_line --index 0 "lets: command failed:"
    assert_line --index 1 "  lint  <-- failed here"
    assert_line --index 2 "lets: exit status 1"
}
