load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/completion
    cleanup
}

@test "completion: should return completion if no lets.yaml" {
    cd ./no_lets_file
    cleanup

    LETS_CONFIG_DIR="no_lets_file" run lets completion
    assert_success
    assert_line --index 0 "Generates completion scripts for bash, zsh"
    [[ ! -d .lets ]]
}

@test "completion: should return completion if lets.yaml exists" {
    run lets completion
    assert_success
    assert_line --index 0 "Generates completion scripts for bash, zsh"
    [[ -d .lets ]]
}

@test "completion: should return list of commands" {
    run lets completion --list
    assert_success
    assert_line --index 0 "bar"
    assert_line --index 1 "foo"
}

@test "completion: should return verbose list of commands" {
    run lets completion --list --verbose
    assert_success
    assert_line --index 0 "bar:Print bar"
    assert_line --index 1 "foo:Print foo"
}
