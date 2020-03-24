load test_helpers

setup() {
    cd ./tests/global_eval_env
}

@test "global_eval_env: should compute env from eval_env and provide env to command" {
    run lets global-eval_env
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "ONE=1" ]]
    [[ "${lines[1]}" = "TWO=two" ]]
    [[ "${lines[2]}" = "THREE=3" ]]
    [[ "${lines[3]}" = "FROM_EVAL_ENV=computed in eval env" ]]
}
