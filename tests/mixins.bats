load test_helpers

setup() {
    cd ./tests/mixins
}

@test "mixins: mixins works" {
    run lets hello-from-minix
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Hello" ]]
}
