setup() {
    cd ./tests/no_lets_file
}

# if we getting colored string we want to replace all color related symbols to get just string
strip_color() {
    echo $(echo "$1" | sed 's/\x1b\[[0-9;]*m//g')
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
    line0=$(strip_color "${lines[0]}")
    [[ "${line0}" = "[ERROR] failed to load config file lets.yaml: open /app/tests/no_lets_file/lets.yaml: no such file or directory" ]]
}

@test "no_lets_file: show help for 'lets help' even if no config file" {
    run lets help
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "A CLI command runner" ]]
    [[ "${lines[1]}" = "Usage:" ]]
    [[ "${lines[2]}" = "  lets [flags]" ]]
    [[ "${lines[3]}" = "Flags:" ]]
    [[ "${lines[4]}" = "  -h, --help   help for lets" ]]
}

# TODO why there is different outputs for lets help and lets --help
@test "no_lets_file: show help for 'lets -h' even if no config file" {
    run lets -h
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "A CLI command runner" ]]
    [[ "${lines[1]}" = "Usage:" ]]
    [[ "${lines[2]}" = "  lets [flags]" ]]
    [[ "${lines[3]}" = "Flags:" ]]
    [[ "${lines[4]}" = "  -h, --help      help for lets" ]]
    [[ "${lines[5]}" = "      --version   version for lets" ]]
}
