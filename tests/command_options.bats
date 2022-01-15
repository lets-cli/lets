load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_options
}

@test "command_options: no options is passed" {
    run lets test-options
    assert_success
    assert_line --index 0 "Flags command"
    assert_line --index 1 "LETSOPT_KV_OPT="
    assert_line --index 2 "LETSOPT_BOOL_OPT="
    assert_line --index 3 "LETSOPT_ARGS="
    assert_line --index 4 "LETSOPT_ATTR="
    assert_line --index 5 "LETSCLI_KV_OPT="
    assert_line --index 6 "LETSCLI_BOOL_OPT="
    assert_line --index 7 "LETSCLI_ARGS="
    assert_line --index 8 "LETSCLI_ATTR="
}

@test "command_options: should parse --kv-opt value with equal sign" {
    run lets test-options --kv-opt=hello
    assert_success
    assert_line --index 0 "Flags command"
    assert_line --index 1 "LETSOPT_KV_OPT=hello"
    assert_line --index 2 "LETSOPT_BOOL_OPT="
    assert_line --index 3 "LETSOPT_ARGS="
    assert_line --index 4 "LETSOPT_ATTR="
    assert_line --index 5 "LETSCLI_KV_OPT=--kv-opt hello"
    assert_line --index 6 "LETSCLI_BOOL_OPT="
    assert_line --index 7 "LETSCLI_ARGS="
    assert_line --index 8 "LETSCLI_ATTR="
}

@test "command_options: should parse --kv-opt value with space" {
    run lets test-options --kv-opt hello
    assert_success
    assert_line --index 0 "Flags command"
    assert_line --index 1 "LETSOPT_KV_OPT=hello"
    assert_line --index 2 "LETSOPT_BOOL_OPT="
    assert_line --index 3 "LETSOPT_ARGS="
    assert_line --index 4 "LETSOPT_ATTR="
    assert_line --index 5 "LETSCLI_KV_OPT=--kv-opt hello"
    assert_line --index 6 "LETSCLI_BOOL_OPT="
    assert_line --index 7 "LETSCLI_ARGS="
    assert_line --index 8 "LETSCLI_ATTR="
}

@test "command_options: should parse --kv-opt and --bool-opt" {
    run lets test-options --kv-opt hello --bool-opt
    assert_success
    assert_line --index 0 "Flags command"
    assert_line --index 1 "LETSOPT_KV_OPT=hello"
    assert_line --index 2 "LETSOPT_BOOL_OPT=true"
    assert_line --index 3 "LETSOPT_ARGS="
    assert_line --index 4 "LETSOPT_ATTR="
    assert_line --index 5 "LETSCLI_KV_OPT=--kv-opt hello"
    assert_line --index 6 "LETSCLI_BOOL_OPT=--bool-opt"
    assert_line --index 7 "LETSCLI_ARGS="
    assert_line --index 8 "LETSCLI_ATTR="
}

@test "command_options: should parse --kv-opt, --bool-opt and positional args" {
    run lets test-options --kv-opt hello --bool-opt myarg1 myarg2
    assert_success
    assert_line --index 0 "Flags command"
    assert_line --index 1 "LETSOPT_KV_OPT=hello"
    assert_line --index 2 "LETSOPT_BOOL_OPT=true"
    assert_line --index 3 "LETSOPT_ARGS=myarg1 myarg2"
    assert_line --index 4 "LETSOPT_ATTR="
    assert_line --index 5 "LETSCLI_KV_OPT=--kv-opt hello"
    assert_line --index 6 "LETSCLI_BOOL_OPT=--bool-opt"
    assert_line --index 7 "LETSCLI_ARGS=myarg1 myarg2"
    assert_line --index 8 "LETSCLI_ATTR="
}

@test "command_options: should parse only positional args" {
    run lets test-options myarg1 myarg2
    assert_success
    assert_line --index 0 "Flags command"
    assert_line --index 1 "LETSOPT_KV_OPT="
    assert_line --index 2 "LETSOPT_BOOL_OPT="
    assert_line --index 3 "LETSOPT_ARGS=myarg1 myarg2"
    assert_line --index 4 "LETSOPT_ATTR="
    assert_line --index 5 "LETSCLI_KV_OPT="
    assert_line --index 6 "LETSCLI_BOOL_OPT="
    assert_line --index 7 "LETSCLI_ARGS=myarg1 myarg2"
    assert_line --index 8 "LETSCLI_ATTR="
}

@test "command_options: should parse repeated kv flags --attr" {
    run lets test-options --attr=myarg1 --attr=myarg2
    assert_success
    assert_line --index 0 "Flags command"
    assert_line --index 1 "LETSOPT_KV_OPT="
    assert_line --index 2 "LETSOPT_BOOL_OPT="
    assert_line --index 3 "LETSOPT_ARGS="
    assert_line --index 4 "LETSOPT_ATTR=myarg1 myarg2"
    assert_line --index 5 "LETSCLI_KV_OPT="
    assert_line --index 6 "LETSCLI_BOOL_OPT="
    assert_line --index 7 "LETSCLI_ARGS="
    assert_line --index 8 "LETSCLI_ATTR=--attr myarg1 myarg2"
}

@test "command_options: option without required argument" {
    run lets test-options --kv-opt

    assert_failure
    assert_line --index 0 "failed to parse docopt options for cmd test-options: --kv-opt requires argument"
    assert_line --index 1 "Usage:"
    assert_line --index 2 "  lets test-options [--kv-opt=<kv-opt>] [--bool-opt] [--attr=<attr>...] [<args>...]"
    assert_line --index 3 "Options:"
    assert_line --index 4 "  <args>...                Positional args in the end"
    assert_line --index 5 "  --bool-opt, -b           Boolean opt"
    assert_line --index 6 "  --kv-opt=<kv-opt>, -K    Key value opt"
    assert_line --index 7 "  --attr=<attr>...         Repeated kv args"
}

@test "command_options: wrong usage" {
    run lets options-wrong-usage

    assert_failure
    assert_line --index 0 "failed to parse docopt options for cmd options-wrong-usage: no such option"
    assert_line --index 1 "Usage: lets options-wrong-usage-xxx"
}

@test "command_options: should not break json argument" {
    run lets test-proxy-options \
        start \
        path.to.pythonModule.py \
        --kwargs='{"sobaka": true}' \
        '--json={"x": 25, "y": [1, 2, 3]}'
    assert_success
    linesLen="${#lines[@]}"
    [[ $linesLen = 4 ]]
    assert_line --index 0 'start'
    assert_line --index 1 'path.to.pythonModule.py'
    assert_line --index 2 '--kwargs={"sobaka": true}'
    assert_line --index 3 '--json={"x": 25, "y": [1, 2, 3]}'
}

@test "command_options: should not break string argument with whitespace" {
    run lets test-proxy-options generate somethingUseful -m "my message contains whitespace!!"
    assert_success
    linesLen="${#lines[@]}"
    [[ $linesLen = 4 ]]
    assert_line --index 0 "generate"
    assert_line --index 1 "somethingUseful"
    assert_line --index 2 "-m"
    assert_line --index 3 "my message contains whitespace!!"
}

@test "command_options: param with same name as command name itself" {
    run lets say Bro
    assert_success
    assert_line --index 0 "Hi Bro"
}
