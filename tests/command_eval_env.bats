load test_helpers

setup() {
    cd ./tests/command_eval_env
}

@test "command_eval_env: should compute and provide env to command" {
    run lets eval-env
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "ONE=1" ]]
    [[ "${lines[1]}" = "TWO=two" ]]
    [[ "${lines[2]}" = "COMPUTED=Computed env" ]]
}
