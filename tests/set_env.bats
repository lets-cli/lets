load test_helpers

setup() {
    cd ./tests/set_env
}

@test "set_env: should use default global env for running command" {
    run lets say_hello_global
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ "${lines[0]}" = "Hello John" ]]
}

@test "set_env: should use default command env for running command" {
    run lets say_hello
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ "${lines[0]}" = "Hello Rick" ]]
}

@test "set_env: should override env for running command with -E" {
    run lets -E NAME=Morty say_hello
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ "${lines[0]}" = "Hello Morty" ]]
}

@test "set_env: should override env for running command with --env" {
    run lets --env NAME=Morty say_hello
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ "${lines[0]}" = "Hello Morty" ]]
}

@test "say_command: should set env var $LETS_COMMAND_NAME" {
    run lets say_command
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ "${lines[0]}" = "say_command" ]]
}
