setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/default_env
    TEST_DIR=$(pwd)
}


@test "LETS_COMMAND_NAME: contains command name" {
    run lets print-command-name-from-env
    assert_success
    assert_line --index 0 "print-command-name-from-env"
}

@test "LETS_COMMAND_ARGS: contains all positional args" {
    run lets print-command-args-from-env --foo --bar=x y

    assert_success
    assert_line --index 0 "--foo --bar=x y"
}

@test "\$@: contains all positional args" {
    run lets print-shell-args --foo --bar=x y

    assert_success
    assert_line --index 0 "--foo --bar=x y"
}


@test "LETS_CONFIG: contains config filename" {
    run lets print-env LETS_CONFIG

    assert_success
    assert_line --index 0 "LETS_CONFIG=lets.yaml"
}

@test "LETS_CONFIG_DIR: contains config dir" {
    run lets print-env LETS_CONFIG_DIR

    assert_success
    assert_line --index 0 "LETS_CONFIG_DIR=${TEST_DIR}"
}

@test "LETS_CONFIG_DIR: specified, overrides config dir" {
    LETS_CONFIG_DIR=./a run lets print-env LETS_CONFIG_DIR

    assert_success
    assert_line --index 0 "LETS_CONFIG_DIR=${TEST_DIR}/a"
}

@test "LETS_COMMAND_WORK_DIR: contains work_dir path if specified for command (in dir with lets config)" {
    cd ./a
    run lets print-workdir

    assert_success
    assert_line --index 0 "LETS_COMMAND_WORK_DIR=${TEST_DIR}/a/b"
}


@test "LETS_COMMAND_WORK_DIR: fail if LETS_CONFIG_DIR specified and no work_dir exists in LETS_CONFIG_DIR path" {
    LETS_CONFIG_DIR=./a run lets print-workdir

    assert_failure
    assert_line --index 0 "failed to run command 'print-workdir': chdir ${TEST_DIR}/b: no such file or directory"
}