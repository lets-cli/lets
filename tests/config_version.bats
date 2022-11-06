load test_helpers

TEST_VERSION=0.0.2

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    # NOTICE to test this functionality properly we building lets with specified version ${TEST_VERSION}
    go build -ldflags="-X main.version=${TEST_VERSION}" -o ./tests/config_version/lets *.go
    cd ./tests/config_version
}

teardown() {
    rm -f ./lets
}

@test "config_version: if config version lower than lets version - its ok" {
    LETS_CONFIG=lets-with-version-0.0.1.yaml run ./lets

    assert_success
    assert_line --index 0 "A CLI task runner"
}

@test "config_version: if config version greater than lets version - fail - require upgrade" {
    LETS_CONFIG=lets-with-version-0.0.3.yaml run ./lets

    assert_failure
    assert_line --index 0 "config version '0.0.3' is not compatible with 'lets' version '0.0.2'. Please upgrade 'lets' to '0.0.3' using 'lets --upgrade' command or following documentation at https://lets-cli.org/docs/installation'"
}

@test "config_version: no version specified" {
    LETS_CONFIG=lets-without-version.yaml run ./lets
    assert_success
    assert_line --index 0 "A CLI task runner"
}
