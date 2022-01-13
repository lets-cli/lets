setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/os_command_env
}


@test "LETS_COMMAND_NAME: contains command name" {
    run lets print-command-name-from-env
    assert_success
    assert_line --index 0 "print-command-name-from-env"
}

@test "LETS_COMMAND_ARGS: contains all positional args" {
    run lets print-command-args-from-env --foo --bar=x y

    assert_success
    assert_line --index 0 "--foo --bar=x y"
}

@test "\$@: contains all positional args" {
    run lets print-shell-args --foo --bar=x y

    assert_success
    assert_line --index 0 "--foo --bar=x y"
}
