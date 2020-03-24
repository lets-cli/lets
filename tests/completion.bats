load test_helpers

setup() {
    cd ./tests/completion
    cleanup
}

@test "completion: should return completion if no lets.yaml" {
    cd ./no_lets_file
    cleanup

    run lets completion
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ "${lines[0]}" == "Generates completion scripts for bash, zsh" ]]
    [[ ! -d .lets ]]
}

@test "completion: should return completion if lets.yaml exists" {
    run lets completion
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ "${lines[0]}" == "Generates completion scripts for bash, zsh" ]]
    [[ -d .lets ]]
}

@test "completion: should return list of commands" {
    run lets completion --list
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ "${lines[0]}" == "bar" ]]
    [[ "${lines[1]}" == "foo" ]]
}

@test "completion: should return verbose list of commands" {
    run lets completion --list --verbose
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ "${lines[0]}" == "bar:Print bar" ]]
    [[ "${lines[1]}" == "foo:Print foo" ]]
}
