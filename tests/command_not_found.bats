setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_not_found
}

@test "command_not_found: exit code is 2 when command does not exist" {
    run lets no_such_command
    assert_failure 2
}

@test "command_not_found: exit code is 2 when self subcommand does not exist" {
    run lets self no_such_command
    assert_failure 2
}

@test "command_not_found: suggest root command for close typo" {
    run lets slef
    assert_failure 2
    assert_output --partial 'unknown command "slef" for "lets"'
    assert_output --partial 'Did you mean this?'
    assert_output --partial 'self'
}

@test "command_not_found: suggest self subcommand for close typo" {
    run lets self ls
    assert_failure 2
    assert_output --partial 'unknown command "ls" for "lets self"'
    assert_output --partial 'Did you mean this?'
    assert_output --partial 'lsp'
}

@test "command_not_found: no suggestions for completely unrelated command" {
    run lets zzzznotacommand
    assert_failure 2
    assert_output --partial 'unknown command "zzzznotacommand" for "lets"'
    refute_output --partial 'Did you mean this?'
}

@test "command_not_found: no suggestions for completely unrelated self subcommand" {
    run lets self zzzznotacommand
    assert_failure 2
    assert_output --partial 'unknown command "zzzznotacommand" for "lets self"'
    refute_output --partial 'Did you mean this?'
}
