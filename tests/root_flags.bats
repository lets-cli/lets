setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/root_flags
}

@test "root_flags: --config works (no subcommand)" {
    run lets --config lets1.yaml
    assert_success
    assert_line --index 0 "A CLI task runner"
}

@test "root_flags: --debug works (no subcommand)" {
    run lets --debug
    assert_success
    assert_line --index 0 "lets: found lets.yaml config file in $(pwd) directory"
}

@test "root_flags: --config works (subcommand)" {
    run lets --config lets1.yaml bar
    assert_success
    assert_line --index 0 "DEBUG="
    assert_line --index 1 "CONFIG="
    assert_line --index 2 "ONLY="
}

@test "root_flags: no root flags (subcommand with flags works)" {
    run lets foo --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "DEBUG=true"
    assert_line --index 1 "CONFIG=xxx"
    assert_line --index 2 "ONLY=yyy"
}

@test "root_flags: --config works (subcommand with flags)" {
    run lets --config lets1.yaml bar --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "DEBUG=true"
    assert_line --index 1 "CONFIG=xxx"
    assert_line --index 2 "ONLY=yyy"
}

@test "root_flags: --debug works (subcommand)" {
    run lets --debug foo
    assert_success
    assert_line --index 0 "lets: found lets.yaml config file in $(pwd) directory"
    assert_line --index 8 "DEBUG="
    assert_line --index 9 "CONFIG="
    assert_line --index 10 "ONLY="
}

@test "root_flags: --debug works (subcommand with flags)" {
    run lets --debug foo --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "lets: found lets.yaml config file in $(pwd) directory"
    assert_line --index 8 "DEBUG=true"
    assert_line --index 9 "CONFIG=xxx"
    assert_line --index 10 "ONLY=yyy"
}

@test "root_flags: --debug and --config works (subcommand with flags)" {
    run lets --debug --config lets1.yaml bar --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "lets: found lets1.yaml config file in $(pwd) directory"
    assert_line --index 8 "DEBUG=true"
    assert_line --index 9 "CONFIG=xxx"
    assert_line --index 10 "ONLY=yyy"
}

@test "root_flags: --debug and --config= works (subcommand with flags)" {
    run lets --debug --config=lets1.yaml bar --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "lets: found lets1.yaml config file in $(pwd) directory"
    assert_line --index 8 "DEBUG=true"
    assert_line --index 9 "CONFIG=xxx"
    assert_line --index 10 "ONLY=yyy"
}

@test "root_flags: --debug and -c works (subcommand with flags)" {
    run lets --debug -c lets1.yaml bar --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "lets: found lets1.yaml config file in $(pwd) directory"
    assert_line --index 8 "DEBUG=true"
    assert_line --index 9 "CONFIG=xxx"
    assert_line --index 10 "ONLY=yyy"
}

@test "root_flags: --config and --debug works (subcommand with flags)" {
    run lets --config lets1.yaml --debug bar --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "lets: found lets1.yaml config file in $(pwd) directory"
    assert_line --index 8 "DEBUG=true"
    assert_line --index 9 "CONFIG=xxx"
    assert_line --index 10 "ONLY=yyy"
}

@test "root_flags: --config= and --debug works (subcommand with flags)" {
    run lets --config=lets1.yaml --debug bar --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "lets: found lets1.yaml config file in $(pwd) directory"
    assert_line --index 8 "DEBUG=true"
    assert_line --index 9 "CONFIG=xxx"
    assert_line --index 10 "ONLY=yyy"
}

@test "root_flags: -c and --debug works (subcommand with flags)" {
    run lets -c lets1.yaml --debug bar --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "lets: found lets1.yaml config file in $(pwd) directory"
    assert_line --index 8 "DEBUG=true"
    assert_line --index 9 "CONFIG=xxx"
    assert_line --index 10 "ONLY=yyy"
}

@test "root_flags: -c and -d works (subcommand with flags)" {
    run lets -c lets1.yaml -d bar --debug --config=xxx --only=yyy
    assert_success
    assert_line --index 0 "lets: found lets1.yaml config file in $(pwd) directory"
    assert_line --index 8 "DEBUG=true"
    assert_line --index 9 "CONFIG=xxx"
    assert_line --index 10 "ONLY=yyy"
}