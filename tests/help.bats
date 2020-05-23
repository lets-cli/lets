load test_helpers

setup() {
    cd ./tests/help
}

@test "help: should create .lets dir" {
    run lets
    printf "%s\n" "${lines[@]}"

    [[ $status == 0 ]]
    [[ -d .lets ]]
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
    [[ "${lines[9]}" =  "  -E, --env stringToString    set env variable for running command KEY=VALUE (default [])" ]]
    [[ "${lines[10]}" = "      --exclude stringArray   run all but excluded command(s) described in cmd as map" ]]
    [[ "${lines[11]}" = "  -h, --help                  help for lets" ]]
    [[ "${lines[12]}" = "      --only stringArray      run only specified command(s) described in cmd as map" ]]
    [[ "${lines[13]}" = "  -v, --version               version for lets" ]]
    [[ "${lines[14]}" = 'Use "lets [command] --help" for more information about a command.' ]]
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
    [[ "${lines[9]}" =  "  -E, --env stringToString    set env variable for running command KEY=VALUE (default [])" ]]
    [[ "${lines[10]}" = "      --exclude stringArray   run all but excluded command(s) described in cmd as map" ]]
    [[ "${lines[11]}" = "  -h, --help                  help for lets" ]]
    [[ "${lines[12]}" = "      --only stringArray      run only specified command(s) described in cmd as map" ]]
    [[ "${lines[13]}" = "  -v, --version               version for lets" ]]
    [[ "${lines[14]}" = 'Use "lets [command] --help" for more information about a command.' ]]
}
