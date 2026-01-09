load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/help_long
}

HELP_MESSAGE=$(cat <<EOF
A CLI task runner

Usage:
  lets [flags]
  lets [command]

Commands:
  bar                                   Print bar
  foo                                   Print foo
  super_long_command_longer_than_usual  Super long command

Internal commands:
  help                                  Help about any command
  self                                  Manage lets CLI itself

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


@test "help: run 'lets' as is" {
    run lets
    assert_success

    assert_output "$HELP_MESSAGE"
}
