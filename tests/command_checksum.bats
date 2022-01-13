load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_checksum
}

ALL_CHECKSUM="be48892c650a32df361202a3662f31e5eac2b83c"
FOO_CHECKSUM="833330f14e30e3ce1907f1e126e1ea4db1ec349f"
BAR_CHECKSUM="7917368d518c031517855672acf2ef82b9cb6836"

CHECKSUM_FROM_FOO_AND_BAR_CHECKSUMS="b778d48759ad4e6e9a755bd595d23eeaa2f7ff65"

@test "command_checksum: should calculate checksum as list of files" {
    run lets as-list-of-files
    assert_success
    assert_line --index 0 ${ALL_CHECKSUM}
}

@test "command_checksum: should calculate checksum as list of globs" {
    run lets as-list-of-globs
    assert_success
    assert_line --index 0 ${ALL_CHECKSUM}
}

@test "command_checksum: should calculate checksum as map of list of files" {
    run lets as-map-of-list-of-files
    assert_success
    assert_line --index 0 "LETS_CHECKSUM_FOO=${FOO_CHECKSUM}"
    assert_line --index 1 "LETS_CHECKSUM_BAR=${BAR_CHECKSUM}"
    assert_line --index 2 "LETS_CHECKSUM=${CHECKSUM_FROM_FOO_AND_BAR_CHECKSUMS}"
}

@test "command_checksum: should calculate checksum as map of list of globs" {
    run lets as-map-of-list-of-globs
    assert_success
    assert_line --index 0 "LETS_CHECKSUM_FOO=${FOO_CHECKSUM}"
    assert_line --index 1 "LETS_CHECKSUM_BAR=${BAR_CHECKSUM}"
    assert_line --index 2 "LETS_CHECKSUM=${CHECKSUM_FROM_FOO_AND_BAR_CHECKSUMS}"
}

@test "command_checksum: checksum from named key in map must be same as from list if files are the same" {
    run lets as-map-all-in-one
    assert_success
    assert_line --index 0 "LETS_CHECKSUM_ALL=${ALL_CHECKSUM}"
    assert_line --index 1 "LETS_CHECKSUM=794b73672fd1259d6fc742cb86713e769d723920"
}


@test "command_checksum: should calculate checksum from sub-dir" {
    cd ./subdir
    run lets as-list-of-files
    assert_success
    assert_line --index 0 ${ALL_CHECKSUM}
}
