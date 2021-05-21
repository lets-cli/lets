load test_helpers

setup() {
    cd ./tests/command_main_cmd_tpl
}

@test "command_main_cmd_tpl: should run with .yml alias" {
    # We can use yaml alias syntax to prevent repeatition of docsopt description
    # The placeholder string is \$\{LETS_COMMAND_NAME\} for now
    run lets main-cmd-tpl-with-adds posarg --config=some_path
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "posarg" ]]
    [[ "${lines[1]}" = "some_path" ]]
}

@test "command_main_cmd_tpl: should fail with yaml alias wo placeholder" {
    run lets -v main-cmd-tpl-with-adds-wo-command_main_cmd_tpl posarg --config=some_path
    printf "%s\n" "${lines[@]}"

    [[ $status != 0 ]]
    [[ "${lines[@]}" =~ "unknown flag: --config" ]]
}
