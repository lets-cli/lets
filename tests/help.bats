load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/help
}

@test "help: should create .lets dir" {
    run lets

    assert_success
    [[ -d .lets ]]
}

@test "help: run 'lets' as is" {
    run lets
    assert_success
    assert_output --partial "A CLI task runner"
    assert_output --partial "bar"
    assert_output --partial "foo"
    assert_output --partial "--config"
}

@test "help: run 'lets --help' (must be same as running lets as is)" {
    run lets --help
    assert_success
    assert_output --partial "A CLI task runner"
    assert_output --partial "bar"
    assert_output --partial "foo"
}

@test "help: run 'lets help' (must be same as running lets as is)" {
    run lets help
    assert_success
    assert_output --partial "A CLI task runner"
    assert_output --partial "bar"
    assert_output --partial "foo"
}

@test "help: show hidden commands" {
    run lets --all
    assert_success
    assert_output --partial "_x"
    assert_output --partial "bar"
    assert_output --partial "foo"
}
