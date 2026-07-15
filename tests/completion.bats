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
    assert_output --partial "Generates completion scripts for bash, zsh"
    [[ ! -d .lets ]]
}

@test "completion: should return completion if lets.yaml exists" {
    run lets completion
    assert_success
    assert_output --partial "Generates completion scripts for bash, zsh"
    [[ -d .lets ]]
}

@test "completion: should return list of commands" {
    run lets completion --commands
    assert_success
    assert_line --index 0 "bar"
    assert_line --index 1 "foo"
}

@test "completion: should return verbose list of commands" {
    run lets completion --commands --verbose
    assert_success
    assert_line --index 0 "bar:Print bar"
    assert_line --index 1 "foo:Print foo"
}

@test "completion: should return list of options for command" {
    run lets completion --options bar
    assert_success
    assert_line --index 0 "--debug"
    assert_line --index 1 "--env"
}

@test "completion: should return verbose list of options for command" {
    run lets completion --options bar --verbose
    assert_success
    assert_line --index 0 "--debug[Run with debug]"
    assert_line --index 1 "--env[Set env]"
}

# Regression: querying the terminal background color from a background process
# group suspends the process with SIGTTOU, so `source <(lets completion -s zsh)`
# in an interactive (job-control) shell hung forever. `script` allocates a PTY,
# `set -m` enables job control like an interactive shell.
@test "completion: source <(lets completion -s zsh) must not hang under job control" {
    run timeout 10 script -qec "zsh -c 'set -m; source <(lets completion -s zsh); echo SOURCED_OK'" /dev/null
    assert_success
    assert_output --partial "SOURCED_OK"
}
