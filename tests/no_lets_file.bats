load test_helpers

setup() {
    cd ./tests/no_lets_file
    cleanup
}

@test "no_lets_file: should not create .lets dir" {
    run lets
    printf "%s\n" "${lines[@]}"

    [[ $status != 0 ]]
    [[ ! -d .lets ]]
}

@test "no_lets_file: show config read error" {
    run lets
    printf "%s\n" "${lines[@]}"

    [[ $status != 0 ]]
    [[ "${lines[0]}" = "[ERROR] failed to load config file lets.yaml: open $(pwd)/lets.yaml: no such file or directory" ]]
}

@test "no_lets_file: show help for 'lets help' even if no config file" {
    run lets help
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "A CLI command runner" ]]
    [[ "${lines[1]}" = "Usage:" ]]
    [[ "${lines[2]}" = "  lets [flags]" ]]
    [[ "${lines[3]}" = "Flags:" ]]
    [[ "${lines[4]}" = "  -h, --help      help for lets" ]]
    [[ "${lines[5]}" = "  -v, --version   version for lets" ]]
}

@test "no_lets_file: show help for 'lets -h' even if no config file" {
    run lets -h
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "A CLI command runner" ]]
    [[ "${lines[1]}" = "Usage:" ]]
    [[ "${lines[2]}" = "  lets [flags]" ]]
    [[ "${lines[3]}" = "Flags:" ]]
    [[ "${lines[4]}" = "  -h, --help      help for lets" ]]
    [[ "${lines[5]}" = "  -v, --version   version for lets" ]]
}
