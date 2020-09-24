load test_helpers

setup() {
    cd ./tests/command_after
}

@test "command_after: should run after script if cmd string" {
    run lets cmd-with-after
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Main" ]]
    [[ "${lines[1]}" = "After" ]]
}

@test "command_after: should run after script if cmd as map" {
    run lets cmd-as-map-with-after
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Main" ]]
    [[ "${lines[1]}" = "After" ]]
}
