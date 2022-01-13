load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_cmd
}

@test "command_cmd: should run as string" {
    run lets cmd-as-string
    assert_success
    assert_line --index 0 "Main"
}

@test "command_cmd: should run as multiline string" {
    run lets cmd-as-multiline-string
    assert_success
    assert_line --index 0 "Main 1 line"
    assert_line --index 1 "Main 2 line"
}

@test "command_cmd: should run as array" {
    run lets cmd-as-array Hello

    assert_success
    assert_line --index 0 "Hello"
}

@test "command_cmd: should run as map" {
    run lets cmd-as-map
    assert_success

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines
    assert_line --index 0 "1"
    assert_line --index 1 "2"
}

@test "command_cmd: cmd-as-map must exit with error if any of cmd exits with error" {
    run lets cmd-as-map-error

    assert_failure
    # as there is no guarantee in which order cmds runs
    # we can not guarantee that all commands will run and complete.
    # But error message must be in the output.
    assert_output --partial "failed to run command 'cmd-as-map-error': exit status 2"
}

@test "command_cmd: cmd-as-map must propagate env" {
    run lets cmd-as-map-env-propagated
    assert_success

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    assert_line --index 0 "1 hello"
    assert_line --index 1 "2 hello"
}

@test "command_cmd: cmd-as-map run with --only" {
    run lets --only two cmd-as-map

    assert_success
    assert_line --index 0 "2"
}

@test "command_cmd: cmd-as-map run with --exclude" {
    run lets --exclude one cmd-as-map

    assert_success
    assert_line --index 0 "2"
}

@test "command_cmd: cmd-as-map run with --only and command own flags" {
    run lets --only two cmd-as-map-with-options --hello

    assert_success
    assert_line --index 0 "2 --hello"
}