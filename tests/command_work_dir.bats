load test_helpers

setup() {
    cd ./tests/command_work_dir
}

@test "command_work_dir: should run command in work_dir" {
    run lets print-file

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "hi there" ]]
}
