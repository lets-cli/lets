load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_work_dir
    find . -type d -name ".lets" -exec rm -rf {} +
}

PROJECT_CHECKSUM="8d856bf15117bf927d7a12caf3fb427f1cdd600f"
ROOT_CHECKSUM="10bb3cf83b3ba687e54b9d15ab8d15e282ccd6f0"

@test "command_work_dir: should run command in work_dir" {
    run lets print-file
    assert_success
    assert_line --index 0 "hi there"
}

@test "command_work_dir: checksum, env_file and env.sh all follow work_dir" {
    # ./input.txt and ./.env exist in both dirs with different contents, so each
    # line below would differ if that directive resolved against the root instead
    run lets everything-follows-work-dir
    assert_success
    assert_line --index 0 "cmd=project"
    assert_line --index 1 "file=project-content"
    assert_line --index 2 "env_file=project"
    assert_line --index 3 "env_sh=project"
    assert_line --index 4 "checksum=${PROJECT_CHECKSUM}"
}

@test "command_work_dir: without work_dir everything resolves against the root" {
    run lets everything-follows-root
    assert_success
    assert_line --index 0 "cmd=command_work_dir"
    assert_line --index 1 "file=other-content"
    assert_line --index 2 "env_file=root"
    assert_line --index 3 "env_sh=command_work_dir"
    assert_line --index 4 "checksum=${ROOT_CHECKSUM}"
}

@test "command_work_dir: checksum differs between work_dir and root" {
    # guards the two checksums above against both collapsing to the same value
    [[ "${PROJECT_CHECKSUM}" != "${ROOT_CHECKSUM}" ]]
}
