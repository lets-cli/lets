setup() {
    cd /app/tests/help
}

@test "help: run 'lets' as is" {
    run lets
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "A CLI command runner" ]]
    [[ "${lines[1]}" = "Usage:" ]]
    [[ "${lines[2]}" = "  lets [flags]" ]]
    [[ "${lines[3]}" = "  lets [command]" ]]
    [[ "${lines[4]}" = "Available Commands:" ]]
    [[ "${lines[5]}" = "  bar         Print bar" ]]
    [[ "${lines[6]}" = "  foo         Print foo" ]]
    [[ "${lines[7]}" = "  help        Help about any command" ]]
    [[ "${lines[8]}" = "Flags:" ]]
    [[ "${lines[9]}" = "  -h, --help      help for lets" ]]
    [[ "${lines[10]}" = "      --version   version for lets" ]]
    [[ "${lines[11]}" = 'Use "lets [command] --help" for more information about a command.' ]]
}

@test "help: run 'lets help'" {
    run lets
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "A CLI command runner" ]]
    [[ "${lines[1]}" = "Usage:" ]]
    [[ "${lines[2]}" = "  lets [flags]" ]]
    [[ "${lines[3]}" = "  lets [command]" ]]
    [[ "${lines[4]}" = "Available Commands:" ]]
    [[ "${lines[5]}" = "  bar         Print bar" ]]
    [[ "${lines[6]}" = "  foo         Print foo" ]]
    [[ "${lines[7]}" = "  help        Help about any command" ]]
    [[ "${lines[8]}" = "Flags:" ]]
    [[ "${lines[9]}" = "  -h, --help      help for lets" ]]
    [[ "${lines[10]}" = "      --version   version for lets" ]]
    [[ "${lines[11]}" = 'Use "lets [command] --help" for more information about a command.' ]]
}
