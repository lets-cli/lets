load test_helpers

setup() {
    cd ./tests/no_lets_file
    cleanup
}

NOT_EXISTED_LETS_FILE="lets-not-existed.yaml"

@test "no_lets_file: should not create .lets dir" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets
    printf "%s\n" "${lines[@]}"

    [[ $status != 0 ]]
    [[ ! -d .lets ]]
}

@test "no_lets_file: show find config error" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets
    printf "%s\n" "${lines[@]}"

    [[ $status != 0 ]]
    [[ "${lines[0]}" = "failed to find config file ${NOT_EXISTED_LETS_FILE}: can not find config" ]]
}

@test "no_lets_file: show config read error (broken config)" {
    LETS_CONFIG=broken_lets.yaml run lets
    printf "%s\n" "${lines[@]}"

    [[ $status != 0 ]]
    [[ "${lines[0]}" = "failed to load config file broken_lets.yaml: failed to parse config: field 'commands': must be a mapping" ]]
}

@test "no_lets_file: show help for 'lets help' even if no config file" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets help
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "A CLI command runner" ]]
    [[ "${lines[1]}" = "Usage:" ]]
    [[ "${lines[2]}" = "  lets [flags]" ]]
    [[ "${lines[3]}" = "Flags:" ]]
    [[ "${lines[4]}" = "  -E, --env stringToString    set env variable for running command KEY=VALUE (default [])" ]]
    [[ "${lines[5]}" = "      --exclude stringArray   run all but excluded command(s) described in cmd as map" ]]
    [[ "${lines[6]}" = "  -h, --help                  help for lets" ]]
    [[ "${lines[7]}" = "      --only stringArray      run only specified command(s) described in cmd as map" ]]
    [[ "${lines[8]}" = "      --upgrade               upgrade lets to latest version" ]]
    [[ "${lines[9]}" = "  -v, --version               version for lets" ]]
}

@test "no_lets_file: show help for 'lets -h' even if no config file" {
    LETS_CONFIG=${NOT_EXISTED_LETS_FILE} run lets -h
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "A CLI command runner" ]]
    [[ "${lines[1]}" = "Usage:" ]]
    [[ "${lines[2]}" = "  lets [flags]" ]]
    [[ "${lines[3]}" = "Flags:" ]]
    [[ "${lines[4]}" = "  -E, --env stringToString    set env variable for running command KEY=VALUE (default [])" ]]
    [[ "${lines[5]}" = "      --exclude stringArray   run all but excluded command(s) described in cmd as map" ]]
    [[ "${lines[6]}" = "  -h, --help                  help for lets" ]]
    [[ "${lines[7]}" = "      --only stringArray      run only specified command(s) described in cmd as map" ]]
    [[ "${lines[8]}" = "      --upgrade               upgrade lets to latest version" ]]
    [[ "${lines[9]}" = "  -v, --version               version for lets" ]]
}
