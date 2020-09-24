load test_helpers

setup() {
    cd ./tests/global_before
}

@test "global_before: should insert before script for each cmd" {
    run lets hello
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Hello" ]]
}

@test "global_before: should merge before scripts from mixins" {
    run lets world
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "World" ]]
}
