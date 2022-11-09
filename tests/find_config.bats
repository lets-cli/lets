load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/find_config
    find . -type d -name ".lets" -delete
    cleanup
}

@test "find_config: should find lets.yaml in parent dir" {
    cd a/b
    run lets foo
    assert_success
    assert_line --index 0 "foo"
}

@test "find_config: .lets must be created in the same dir where lets.yaml placed" {
    cd a/b
    run lets foo
    assert_success

    [[ ! -d .lets ]]
    [[ -d ../../.lets ]]
}

@test "find_config: LETS_CONFIG changes which config file to read" {
    LETS_CONFIG=lets1.yaml run lets hi

    assert_success
    assert_line --index 0 "Hi from lets1.yaml"
}

@test "find_config: --config changes which config file to read" {
    # also check that root --config and subcommand --config works together
    run lets --config lets1.yaml hi --config=xxx

    assert_success
    assert_line --index 0 "Hi from lets1.yaml"
    assert_line --index 1 "Option --config=xxx"
}

@test "find_config: subcommand --config must not change which config file to read" {
    run lets hi --config=xxx

    assert_success
    assert_line --index 0 "Hi from lets.yaml"
    assert_line --index 1 "Option --config=xxx"
}