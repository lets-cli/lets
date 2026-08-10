setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/mixins
}

@test "mixins: mixins works" {
    run lets hello-from-minix
    assert_success
    assert_line --index 0 "Hello"
}

@test "mixins: mixin path in a subdir is loaded" {
    run lets hello-from-subdir-mixin
    assert_success
    assert_line --index 0 "Hello from sub/outer.yaml"
}

@test "mixins: mixin paths resolve against the config dir, not the root dir" {
    # run from sub/, so the root dir is sub/ — the mixin paths in ../lets.yaml
    # ('sub/outer.yaml', 'lets.mix.yaml') only resolve if they are taken
    # relative to the config file rather than to where lets was invoked
    cd sub
    run lets -c ../lets.yaml hello-from-subdir-mixin
    assert_success
    assert_line --index 0 "Hello from sub/outer.yaml"

    run lets -c ../lets.yaml hello-from-minix
    assert_success
    assert_line --index 0 "Hello"
}
