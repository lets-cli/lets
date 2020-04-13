load test_helpers

setup() {
    cd ./tests/command_options
}

@test "command_options: no options is passed" {
    run lets test-options
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Flags command" ]]
    [[ "${lines[1]}" = "LETSOPT_KV_OPT=" ]]
    [[ "${lines[2]}" = "LETSOPT_BOOL_OPT=" ]]
    [[ "${lines[3]}" = "LETSOPT_ARGS=" ]]
    [[ "${lines[4]}" = "LETSOPT_ATTR=" ]]
    [[ "${lines[5]}" = "LETSCLI_KV_OPT=" ]]
    [[ "${lines[6]}" = "LETSCLI_BOOL_OPT=" ]]
    [[ "${lines[7]}" = "LETSCLI_ARGS=" ]]
    [[ "${lines[8]}" = "LETSCLI_ATTR=" ]]
}

@test "command_options: should parse --kv-opt value with equal sign" {
    run lets test-options --kv-opt=hello
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Flags command" ]]
    [[ "${lines[1]}" = "LETSOPT_KV_OPT=hello" ]]
    [[ "${lines[2]}" = "LETSOPT_BOOL_OPT=" ]]
    [[ "${lines[3]}" = "LETSOPT_ARGS=" ]]
    [[ "${lines[4]}" = "LETSOPT_ATTR=" ]]
    [[ "${lines[5]}" = "LETSCLI_KV_OPT=--kv-opt hello" ]]
    [[ "${lines[6]}" = "LETSCLI_BOOL_OPT=" ]]
    [[ "${lines[7]}" = "LETSCLI_ARGS=" ]]
    [[ "${lines[8]}" = "LETSCLI_ATTR=" ]]
}

@test "command_options: should parse --kv-opt value with space" {
    run lets test-options --kv-opt hello
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Flags command" ]]
    [[ "${lines[1]}" = "LETSOPT_KV_OPT=hello" ]]
    [[ "${lines[2]}" = "LETSOPT_BOOL_OPT=" ]]
    [[ "${lines[3]}" = "LETSOPT_ARGS=" ]]
    [[ "${lines[4]}" = "LETSOPT_ATTR=" ]]
    [[ "${lines[5]}" = "LETSCLI_KV_OPT=--kv-opt hello" ]]
    [[ "${lines[6]}" = "LETSCLI_BOOL_OPT=" ]]
    [[ "${lines[7]}" = "LETSCLI_ARGS=" ]]
    [[ "${lines[8]}" = "LETSCLI_ATTR=" ]]
}

@test "command_options: should parse --kv-opt and --bool-opt" {
    run lets test-options --kv-opt hello --bool-opt
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Flags command" ]]
    [[ "${lines[1]}" = "LETSOPT_KV_OPT=hello" ]]
    [[ "${lines[2]}" = "LETSOPT_BOOL_OPT=true" ]]
    [[ "${lines[3]}" = "LETSOPT_ARGS=" ]]
    [[ "${lines[4]}" = "LETSOPT_ATTR=" ]]
    [[ "${lines[5]}" = "LETSCLI_KV_OPT=--kv-opt hello" ]]
    [[ "${lines[6]}" = "LETSCLI_BOOL_OPT=--bool-opt" ]]
    [[ "${lines[7]}" = "LETSCLI_ARGS=" ]]
    [[ "${lines[8]}" = "LETSCLI_ATTR=" ]]
}

@test "command_options: should parse --kv-opt, --bool-opt and positional args" {
    run lets test-options --kv-opt hello --bool-opt myarg1 myarg2
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Flags command" ]]
    [[ "${lines[1]}" = "LETSOPT_KV_OPT=hello" ]]
    [[ "${lines[2]}" = "LETSOPT_BOOL_OPT=true" ]]
    [[ "${lines[3]}" = "LETSOPT_ARGS=myarg1 myarg2" ]]
    [[ "${lines[4]}" = "LETSOPT_ATTR=" ]]
    [[ "${lines[5]}" = "LETSCLI_KV_OPT=--kv-opt hello" ]]
    [[ "${lines[6]}" = "LETSCLI_BOOL_OPT=--bool-opt" ]]
    [[ "${lines[7]}" = "LETSCLI_ARGS=myarg1 myarg2" ]]
    [[ "${lines[8]}" = "LETSCLI_ATTR=" ]]
}

@test "command_options: should parse only positional args" {
    run lets test-options myarg1 myarg2
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Flags command" ]]
    [[ "${lines[1]}" = "LETSOPT_KV_OPT=" ]]
    [[ "${lines[2]}" = "LETSOPT_BOOL_OPT=" ]]
    [[ "${lines[3]}" = "LETSOPT_ARGS=myarg1 myarg2" ]]
    [[ "${lines[4]}" = "LETSOPT_ATTR=" ]]
    [[ "${lines[5]}" = "LETSCLI_KV_OPT=" ]]
    [[ "${lines[6]}" = "LETSCLI_BOOL_OPT=" ]]
    [[ "${lines[7]}" = "LETSCLI_ARGS=myarg1 myarg2" ]]
    [[ "${lines[8]}" = "LETSCLI_ATTR=" ]]
}

@test "command_options: should parse repeated kv flags --attr" {
    run lets test-options --attr=myarg1 --attr=myarg2
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "Flags command" ]]
    [[ "${lines[1]}" = "LETSOPT_KV_OPT=" ]]
    [[ "${lines[2]}" = "LETSOPT_BOOL_OPT=" ]]
    [[ "${lines[3]}" = "LETSOPT_ARGS=" ]]
    [[ "${lines[4]}" = "LETSOPT_ATTR=myarg1 myarg2" ]]
    [[ "${lines[5]}" = "LETSCLI_KV_OPT=" ]]
    [[ "${lines[6]}" = "LETSCLI_BOOL_OPT=" ]]
    [[ "${lines[7]}" = "LETSCLI_ARGS=" ]]
    [[ "${lines[8]}" = "LETSCLI_ATTR=--attr myarg1 myarg2" ]]
}

@test "command_options: option without required argument" {
    run lets test-options --kv-opt
    printf "%s\n" "${lines[@]}"

    [[ $status != 0 ]]
    [[ "${lines[0]}" = "Error: failed to parse docopt options for cmd test-options: --kv-opt requires argument" ]]
    [[ "${lines[1]}" = "Usage:" ]]
    [[ "${lines[2]}" = "  lets test-options [--kv-opt=<kv-opt>] [--bool-opt] [--attr=<attr>...] [<args>...]" ]]
    [[ "${lines[3]}" = "Options:" ]]
    [[ "${lines[4]}" = "  <args>...                Positional args in the end" ]]
    [[ "${lines[5]}" = "  --bool-opt, -b           Boolean opt" ]]
    [[ "${lines[6]}" = "  --kv-opt=<kv-opt>, -K    Key value opt" ]]
    [[ "${lines[7]}" = "  --attr=<attr>...         Repeated kv args" ]]
}

@test "command_options: wrong usage" {
    run lets options-wrong-usage
    printf "%s\n" "${lines[@]}"

    [[ $status != 0 ]]
    [[ "${lines[0]}" = "Error: failed to parse docopt options for cmd options-wrong-usage: no such option" ]]
    [[ "${lines[1]}" = "Usage: lets options-wrong-usage-xxx" ]]
}

@test "command_options: should not break json argument" {
    run lets test-proxy-options \
        start \
        path.to.pythonModule.py \
        --kwargs='{"sobaka": true}' \
        '--json={"x": 25, "y": [1, 2, 3]}'
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    linesLen="${#lines[@]}"
    [[ $linesLen = 4 ]]
    [[ "${lines[0]}" = 'start' ]]
    [[ "${lines[1]}" = 'path.to.pythonModule.py' ]]
    [[ "${lines[2]}" = '--kwargs={"sobaka": true}' ]]
    [[ "${lines[3]}" = '--json={"x": 25, "y": [1, 2, 3]}' ]]
}

@test "command_options: should not break string argument with whitespace" {
    run lets test-proxy-options generate somethingUseful -m "my message contains whitespace!!"
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    linesLen="${#lines[@]}"
    [[ $linesLen = 4 ]]
    [[ "${lines[0]}" = "generate" ]]
    [[ "${lines[1]}" = "somethingUseful" ]]
    [[ "${lines[2]}" = "-m" ]]
    [[ "${lines[3]}" = "my message contains whitespace!!" ]]
}
