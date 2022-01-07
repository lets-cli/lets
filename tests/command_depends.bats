load test_helpers

setup() {
    cd ./tests/command_depends
}

@test "command_depends: should run all depends commands before main command" {
    run lets run-with-depends
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Hello Developer" ]]
    [[ "${lines[1]}" = "Bar" ]]
    [[ "${lines[2]}" = "Main" ]]
}
