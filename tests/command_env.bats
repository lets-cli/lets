load test_helpers

setup() {
    cd ./tests/command_env
}

@test "command_env: should provide env to command" {
    run lets env
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "ONE=1" ]]
    [[ "${lines[1]}" = "TWO=two" ]]
}
