load test_helpers

setup() {
    cd ./tests/command_cmd
}

@test "command_cmd: should run as string" {
    run lets cmd-as-string
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Main" ]]
}

@test "command_cmd: should run as multiline string" {
    run lets cmd-as-multiline-string
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Main 1 line" ]]
    [[ "${lines[1]}" = "Main 2 line" ]]
}

@test "command_cmd: should run as array" {
    run lets cmd-as-array Hello
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Hello" ]]
}

@test "command_cmd: should run as map" {
    run lets cmd-as-map
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    [[ "${lines[0]}" = "1" ]]
    [[ "${lines[1]}" = "2" ]]
}

@test "command_cmd: cmd-as-map must exit with error if any of cmd exits with error" {
    run lets cmd-as-map-error
    printf "%s\n" "${lines[@]}"

    [[ $status != 0 ]]
    # as there is no guarantee in which order cmds runs
    # we can not guarantee that all commands will run and complete.
    # But error message must be in the output.
    [[ "${lines[@]}" =~ "Error: failed to run command 'cmd-as-map-error': exit status 2" ]]
}

@test "command_cmd: cmd-as-map must propagate env" {
    run lets cmd-as-map-env-propagated
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    [[ "${lines[0]}" = "1 hello" ]]
    [[ "${lines[1]}" = "2 hello" ]]
}

@test "command_cmd: cmd-as-map run with --only" {
    run lets --only two cmd-as-map
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    [[ "${lines[0]}" = "2" ]]
}

@test "command_cmd: cmd-as-map run with --exclude" {
    run lets --exclude one cmd-as-map
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    [[ "${lines[0]}" = "2" ]]
}


@test "command_cmd: cmd-as-map run with --only and command own flags" {
    run lets --only two cmd-as-map-with-options --hello
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    [[ "${lines[0]}" = "2 --hello" ]]
}