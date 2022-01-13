load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/help
}

@test "help: should create .lets dir" {
    run lets

    assert_success
    [[ -d .lets ]]
}

@test "help: run 'lets' as is" {
    run lets
    assert_success

    assert_line --index 0 "A CLI command runner"
    assert_line --index 1 "Usage:"
    assert_line --index 2 "  lets [flags]"
    assert_line --index 3 "  lets [command]"
    assert_line --index 4 "Available Commands:"
    assert_line --index 5 "  bar         Print bar"
    assert_line --index 6 "  foo         Print foo"
    assert_line --index 7 "  help        Help about any command"
    assert_line --index 8 "Flags:"
    assert_line --index 9 "  -E, --env stringToString    set env variable for running command KEY=VALUE (default [])"
    assert_line --index 10 "      --exclude stringArray   run all but excluded command(s) described in cmd as map"
    assert_line --index 11 "  -h, --help                  help for lets"
    assert_line --index 12 "      --init                  creates a new lets.yaml in the current folder"
    assert_line --index 13 "      --only stringArray      run only specified command(s) described in cmd as map"
    assert_line --index 14 "      --upgrade               upgrade lets to latest version"
    assert_line --index 15 "  -v, --version               version for lets"
}

@test "help: run 'lets help' (must be same as running lets as is)" {
    run lets --help
    assert_success

    assert_line --index 0 "A CLI command runner"
    assert_line --index 1 "Usage:"
    assert_line --index 2 "  lets [flags]"
    assert_line --index 3 "  lets [command]"
    assert_line --index 4 "Available Commands:"
    assert_line --index 5 "  bar         Print bar"
    assert_line --index 6 "  foo         Print foo"
    assert_line --index 7 "  help        Help about any command"
    assert_line --index 8 "Flags:"
    assert_line --index 9 "  -E, --env stringToString    set env variable for running command KEY=VALUE (default [])"
    assert_line --index 10 "      --exclude stringArray   run all but excluded command(s) described in cmd as map"
    assert_line --index 11 "  -h, --help                  help for lets"
    assert_line --index 12 "      --init                  creates a new lets.yaml in the current folder"
    assert_line --index 13 "      --only stringArray      run only specified command(s) described in cmd as map"
    assert_line --index 14 "      --upgrade               upgrade lets to latest version"
    assert_line --index 15 "  -v, --version               version for lets"
}
