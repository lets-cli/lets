load test_helpers

setup() {
    cd ./tests/global_env
}

@test "global_env: should provide env to command" {
    run lets global-env
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "ONE=1" ]]
    [[ "${lines[1]}" = "TWO=two" ]]
    [[ "${lines[2]}" = "THREE=3" ]]
}
