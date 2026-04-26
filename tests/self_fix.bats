load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/self_fix
    cp lets.old.yaml lets.yaml
    cp mixin.old.yaml mixin.yaml
}

teardown() {
    rm -f lets.yaml mixin.yaml
    cleanup
}

@test "self_fix: dry-run prints migrated config without changing configs" {
    run lets self fix --dry-run

    assert_success
    assert_line --partial "files:"
    assert_line --partial "persist: true"
    assert_line --partial "remote mixin not updated: https://example.com/lets.mixin.yaml"

    run grep -q "persist_checksum: true" lets.yaml
    assert_success
}

@test "self_fix: migrates root config and local mixins" {
    run lets self fix

    assert_success
    assert_line --partial "Migration 'checksum' applied successfully"
    assert_line --partial "remote mixin not updated: https://example.com/lets.mixin.yaml"

    run grep -q "persist_checksum" lets.yaml
    [[ "$status" = 1 ]]

    run grep -q "persist: true" lets.yaml
    assert_success

    run grep -q "files:" mixin.yaml
    assert_success
}
