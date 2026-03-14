load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/env_dependency
}

@test "env_dependency: global env sh can use previously resolved global env" {
    run lets global-env-dependency
    assert_success
    assert_line --index 0 "GLOBAL_COMPOSE=docker-compose"
}

@test "env_dependency: command env sh can use previously resolved command env" {
    run lets command-env-dependency
    assert_success
    assert_line --index 0 "COMMAND_COMPOSE=podman-compose"
}

@test "env_dependency: command env sh can use global env" {
    run lets command-env-uses-global
    assert_success
    assert_line --index 0 "COMMAND_COMPOSE=docker-compose"
}

@test "env_dependency: forward references stay unresolved with sequential evaluation" {
    run env -u LETS_TEST_FORWARD_VAR lets command-forward-reference
    assert_success
    assert_line --index 0 "COMMAND_COMPOSE="
    assert_line --index 1 "LETS_TEST_FORWARD_VAR=from-command-env"
}
