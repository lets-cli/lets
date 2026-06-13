load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_group
}

@test "help: running 'lets help' should group commands by their group names" {
    run lets help
    assert_success
    assert_output --partial "A CLI task runner"
    assert_output --partial "A GROUP"
    assert_output --partial "B GROUP"
    assert_output --partial "COMMON"
    assert_output --partial "C command"
    assert_output --partial "A command"
    assert_output --partial "D command"
}

@test "help: running 'lets --help' should group commands by their group names" {
    run lets --help
    assert_success
    assert_output --partial "A CLI task runner"
    assert_output --partial "A GROUP"
    assert_output --partial "B GROUP"
    assert_output --partial "COMMON"
}

@test "help: running 'lets' should group commands by their group names" {
    run lets
    assert_success
    assert_output --partial "A CLI task runner"
    assert_output --partial "A GROUP"
    assert_output --partial "B GROUP"
    assert_output --partial "COMMON"
}
