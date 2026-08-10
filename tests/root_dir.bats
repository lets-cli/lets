load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/root_dir
    find . -type d -name ".lets" -exec rm -rf {} +
}

# The root dir is the dir lets was invoked from, never the dir the config lives in.
# Each test below pins one way of pointing lets at a config.

@test "root_dir: config in cwd, no -c" {
    run lets pwd
    assert_success
    assert_line --index 0 "<root_dir>"
}

@test "root_dir: -c naming a config in cwd" {
    run lets -c lets.yaml pwd
    assert_success
    assert_line --index 0 "<root_dir>"
}

@test "root_dir: -c pointing into a child dir does not move the root" {
    run lets -c sub/lets.yaml pwd
    assert_success
    assert_line --index 0 "<root_dir>"
}

@test "root_dir: config found recursively up the tree does not move the root" {
    cd deep/nested
    run lets pwd
    assert_success
    assert_line --index 0 "<root_dir>/deep/nested"
}

@test "root_dir: -c pointing at a parent config does not move the root" {
    cd sub
    run lets -c ../lets.yaml pwd
    assert_success
    assert_line --index 0 "<root_dir>/sub"
}

@test "root_dir: LETS_CONFIG_DIR steers discovery but not the root" {
    LETS_CONFIG_DIR=sub run lets pwd
    assert_success
    assert_line --index 0 "<root_dir>"
}

@test "root_dir: work_dir resolves against the root" {
    run lets pwd-with-work-dir
    assert_success
    assert_line --index 0 "<root_dir>/wd"
}

@test "root_dir: work_dir resolves against the root, not the config dir" {
    # invoked from deep/, so work_dir 'wd' means deep/wd — which does not exist
    cd deep
    run lets -c ../lets.yaml pwd-with-work-dir
    assert_failure
    assert_output --partial "deep/wd"
}
