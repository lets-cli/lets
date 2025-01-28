load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/help
}

HELP_MESSAGE=$(cat <<EOF
A CLI task runner

Usage:
  lets [flags]
  lets [command]

Commands:
  bar         Print bar
  foo         Print foo

Internal commands:
  help        Help about any command
  self        Manage lets CLI itself

Flags:
      --all                   show all commands (including the ones with _)
  -c, --config string         config file (default is lets.yaml)
  -d, --debug count           show debug logs (or use LETS_DEBUG=1). If used multiple times, shows more verbose logs
  -E, --env stringToString    set env variable for running command KEY=VALUE (default [])
      --exclude stringArray   run all but excluded command(s) described in cmd as map
  -h, --help                  help for lets
      --init                  create a new lets.yaml in the current folder
      --no-depends            skip 'depends' for running command
      --only stringArray      run only specified command(s) described in cmd as map
      --upgrade               upgrade lets to latest version
  -v, --version               version for lets

Use "lets help [command]" for more information about a command.
EOF
)

HELP_MESSAGE_WITH_HIDDEN=$(cat <<EOF
A CLI task runner

Usage:
  lets [flags]
  lets [command]

Commands:
  _x          Hidden x
  bar         Print bar
  foo         Print foo

Internal commands:
  help        Help about any command
  self        Manage lets CLI itself

Flags:
      --all                   show all commands (including the ones with _)
  -c, --config string         config file (default is lets.yaml)
  -d, --debug count           show debug logs (or use LETS_DEBUG=1). If used multiple times, shows more verbose logs
  -E, --env stringToString    set env variable for running command KEY=VALUE (default [])
      --exclude stringArray   run all but excluded command(s) described in cmd as map
  -h, --help                  help for lets
      --init                  create a new lets.yaml in the current folder
      --no-depends            skip 'depends' for running command
      --only stringArray      run only specified command(s) described in cmd as map
      --upgrade               upgrade lets to latest version
  -v, --version               version for lets

Use "lets help [command]" for more information about a command.
EOF
)

@test "help: should create .lets dir" {
    run lets

    assert_success
    [[ -d .lets ]]
}

@test "help: run 'lets' as is" {
    run lets
    assert_success

    assert_output "$HELP_MESSAGE"
}

@test "help: run 'lets --help' (must be same as running lets as is)" {
    run lets --help
    assert_success

    assert_output "$HELP_MESSAGE"
}

@test "help: run 'lets help' (must be same as running lets as is)" {
    run lets help
    assert_success

    assert_output "$HELP_MESSAGE"
}

@test "help: show hidden commands" {
    run lets --all
    assert_success

    assert_output "$HELP_MESSAGE_WITH_HIDDEN"
}
