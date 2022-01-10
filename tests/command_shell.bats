load test_helpers

setup() {
    cd ./tests/command_shell
}

@test "command_shell: should run command using shell specified in command" {
    run lets show-shell

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "/bin/sh" ]]
}
