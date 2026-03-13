load test_helpers

TEST_VERSION=1.2.3
TEST_BUILD_DATE=2024-01-15T10:30:00Z

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
}

@test "version: show lets version for -v" {
    run lets -v
    assert_success
    assert_line --index 0 "lets version 0.0.0-dev"
}

@test "version: show lets version for --version" {
    run lets --version
    assert_success
    assert_line --index 0 "lets version 0.0.0-dev"
}

@test "version: show build date when set via ldflags" {
    go build -ldflags="-X main.Version=${TEST_VERSION} -X main.BuildDate=${TEST_BUILD_DATE}" -o /tmp/lets-version-test cmd/lets/main.go
    run /tmp/lets-version-test --version
    assert_success
    assert_line --index 0 "lets version ${TEST_VERSION} (${TEST_BUILD_DATE})"
    rm -f /tmp/lets-version-test
}
