@test "version: show lets version for -v" {
    run lets -v
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "lets version dev" ]]
}

@test "version: show lets version for --version" {
    run lets --version
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "lets version dev" ]]
}
