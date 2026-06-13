load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/dependency_failure_tree
}

@test "dependency_failure_tree: shows full 3-level chain on failure" {
    run env NO_COLOR=1 lets deploy
    assert_failure
    assert_output --partial "deploy"
    assert_output --partial "build"
    assert_output --partial "lint"
    assert_output --partial "<-- failed here"
    assert_output --partial "Exit status 1"
}

@test "dependency_failure_tree: single node when no depends" {
    run env NO_COLOR=1 lets lint
    assert_failure
    assert_output --partial "lint"
    assert_output --partial "<-- failed here"
    assert_output --partial "Exit status 1"
}
