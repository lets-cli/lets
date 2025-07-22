load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_depends
}

@test "command_depends: should run all depends commands before main command" {
    run lets run-with-depends
    assert_success
    assert_line --index 0 "Hello World with level INFO"
    assert_line --index 1 "Bar"
    assert_line --index 2 "Main"
}

@test "command_depends: should override args" {
    run lets override-args
    assert_success
    assert_line --index 0 "Hello Developer with level INFO"
    assert_line --index 1 "Bar"
    assert_line --index 2 "Override args"
}

@test "command_depends: should override env" {
    run lets override-env
    assert_success
    assert_line --index 0 "Hello World with level DEBUG"
    assert_line --index 1 "Override env"
}

@test "command_depends: ref works in depends" {
    # checks that original command does not overrides ref to original
    # command. The order in depends is essential to test behavior.
    run lets with-ref-in-depends
    assert_success
    assert_line --index 0 "Hello World with level INFO"
    # World -> Developer by ref.args
    # INFO -> DEBUG by depends[1].env.INFO
    assert_line --index 1 "Hello Developer with level DEBUG"
    # World -> Bar (because dep args has more priority over ref args)
    assert_line --index 2 "Hello Bar with level INFO"
    assert_line --index 3 "I have ref in depends"
}


@test "command_depends: disallow parallel cmd in depends" {
    LETS_CONFIG=lets-parallel-in-depends.yaml run lets parallel-in-depends
    assert_failure
    assert_line --index 0 "lets: config error: command 'parallel-in-depends' depends on command 'parallel', but parallel cmd is not allowed in depends yet"
}

@test "command_depends: should show dependency tree on failure" {
    run lets run-with-failing-dep
    assert_failure 1
    assert_output --partial "'fail-command' failed: exit status 1"
    assert_output --partial "'run-with-failing-dep' ->"
    assert_output --partial "'fail-command' ⚠️"
}

@test "command_depends: should show dependency tree with multiple levels" {
    run lets level2-dep
    assert_failure 1
    assert_output --partial "'fail-command' failed: exit status 1"
    assert_output --partial "'level2-dep' ->"
    assert_output --partial "'run-with-failing-dep'"
    assert_output --partial "'fail-command' ⚠️"
}

@test "command_depends: should run successful deps before showing failure tree" {
    run lets multiple-deps-one-fail
    assert_failure 1
    assert_output --partial "Hello World with level INFO"
    assert_output --partial "Bar"
    assert_output --partial "'fail-command' failed: exit status 1"
    assert_output --partial "'multiple-deps-one-fail' ->"
    assert_output --partial "'fail-command' ⚠️"
}