load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/help
}

HELP_MESSAGE=<<EOF
A CLI command runner

Usage:
  lets [flags]
  lets [command]

Available Commands:
  bar         Print bar
  foo         Print foo
  help        Help about any command

Flags:
  -E, --env stringToString    set env variable for running command KEY=VALUE (default [])
      --exclude stringArray   run all but excluded command(s) described in cmd as map
  -h, --help                  help for lets
      --init                  create a new lets.yaml in the current folder
      --no-depends            skip 'depends' for running command
      --only stringArray      run only specified command(s) described in cmd as map
      --upgrade               upgrade lets to latest version
  -v, --version               version for lets

Use "lets [command] --help" for more information about a command.
EOF

@test "help: should create .lets dir" {
    run lets

    assert_success
    [[ -d .lets ]]
}

@test "help: run 'lets' as is" {
    run lets
    assert_success

    assert_output $HELP_MESSAGE
}

@test "help: run 'lets help' (must be same as running lets as is)" {
    run lets --help
    assert_success

    assert_output $HELP_MESSAGE
}
