setup() {
    cd ./tests/command_cmd
}

@test "command_cmd: should run as string" {
    run lets cmd-as-string
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Main" ]]
}

@test "command_cmd: should run as multiline string" {
    run lets cmd-as-multiline-string
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Main 1 line" ]]
    [[ "${lines[1]}" = "Main 2 line" ]]
}

@test "command_cmd: should run as array" {
    run lets cmd-as-array Hello
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Hello" ]]
}
