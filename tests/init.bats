setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/init
    rm -f lets.yaml
}

@test "--init: init config if not exist" {
    [[ ! -f lets.yaml ]]
    run lets --init
    assert_success
    [[ -f lets.yaml ]]
    assert_line --index 0 "lets.yaml created in the current directory"
}

@test "--init: do not init config if already exist" {
    cd ./exists
    [[ -f lets.yaml ]]
    run lets --init
    assert_failure
    assert_output --partial "lets.yaml already exists in"
}
