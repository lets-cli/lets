load test_helpers

clean_test_files() {
    cleanup
    rm -f foo_test.txt
}

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_persist_checksum
    clean_test_files
}

teardown() {
    clean_test_files
}

FIRST_CHECKSUM=833330f14e30e3ce1907f1e126e1ea4db1ec349f
CHANGED_CHECKSUM=95d4080082937fe50b8db90f0c21acc597c9d176

TEMP_FILE=foo_test.txt

@test "command_persist_checksum: should check if checksum has changed" {
    export CMD_NAME=persist-checksum

    run lets ${CMD_NAME}
    printf "first run: %s\n" "${lines[@]}"

    # first run, no stored checksum - lets should calculate checksum, store it to disk
    # 1. check checksum value
    # 2. check LETS_CHECKSUM_CHANGED has to be changed as there was now checksum at all and now we have new checksum
    # 3. check checksum persisted

    assert_success
    assert_line --index 0 "LETS_CHECKSUM=${FIRST_CHECKSUM}"
    assert_line --index 1 "LETS_CHECKSUM_CHANGED=true"

    # it creates .lets
    [[ -d .lets ]]
    # it creates checksums folder in .lets for storing commands checksums
    [[ -d .lets/checksums ]]
    # it creates "persist-checksum" folder - a folder with name of a command
    [[ -d .lets/checksums/${CMD_NAME} ]]
    # it creates "lets_default_checksum" file - a file with an actual checksum persisted after command has finished
    [[ -f .lets/checksums/${CMD_NAME}/lets_default_checksum ]]

    # second run, previous checksum persisted. lets must read it and check that its not changed
    run lets ${CMD_NAME}
    printf "second run: %s\n" "${lines[@]}"

    assert_success
    assert_line --index 0 "LETS_CHECKSUM=${FIRST_CHECKSUM}"
    assert_line --index 1 "LETS_CHECKSUM_CHANGED=false"

    # third run, there is stored checksum and we creating new file. checksum must be changed now

    # create file suiting glob pattern foo_*.txt
    touch ${TEMP_FILE} && printf "footemp" > ${TEMP_FILE}

    # 1. check checksum value has changed
    # 2. check LETS_CHECKSUM_CHANGED has changed to true
    # 2. check new checksum persisted
    run lets ${CMD_NAME}
    printf "third run: %s\n" "${lines[@]}"

    assert_success
    assert_line --index 0 "LETS_CHECKSUM=${CHANGED_CHECKSUM}"
    assert_line --index 1 "LETS_CHECKSUM_CHANGED=true"
}

@test "command_persist_checksum: should persist checksum for cmd-as-map" {
    export CMD_NAME=persist-checksum-for-cmd-as-map

    run lets ${CMD_NAME}

    # first run, no stored checksum
    # 1. check checksum value
    # 2. check LETS_CHECKSUM_CHANGED has to be changed as there was now checksum at all and now we have new checksum
    # 3. check checksum persisted

    assert_success

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    assert_line --index 0 "1 LETS_CHECKSUM=${FIRST_CHECKSUM}"
    assert_line --index 1 "2 LETS_CHECKSUM_CHANGED=true"

    # it creates .lets
    [[ -d .lets ]]
    # it creates checksums folder in .lets for storing commands checksums
    [[ -d .lets/checksums ]]
    # it creates "persist-checksum" folder - a folder with name of a command
    [[ -d .lets/checksums/${CMD_NAME} ]]
    # it creates "lets_default_checksum" file - a file with an actual checksum persisted after command has finished
    [[ -f .lets/checksums/${CMD_NAME}/lets_default_checksum ]]

    # second run, previous checksum persisted. lets must read it and check that its not changed
    run lets ${CMD_NAME}
    assert_success

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    assert_line --index 0 "1 LETS_CHECKSUM=${FIRST_CHECKSUM}"
    assert_line --index 1 "2 LETS_CHECKSUM_CHANGED=false"

    # third run, there is stored checksum and we creating new file. checksum must be changed now

    # create file suiting glob pattern foo_*.txt
    touch ${TEMP_FILE} && printf "footemp" > ${TEMP_FILE}

    # 1. check checksum value has changed
    # 2. check LETS_CHECKSUM_CHANGED has changed to true
    # 2. check new checksum persisted
    run lets ${CMD_NAME}

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    assert_success

    sort_array lines
    assert_line --index 0 "1 LETS_CHECKSUM=${CHANGED_CHECKSUM}"
    assert_line --index 1 "2 LETS_CHECKSUM_CHANGED=true"
}

@test "command_persist_checksum: should persist checksum only if exit code = 0" {
    run lets with-error-code-1

    [[ $status = 1 ]]
    assert_line --index 0 "LETS_CHECKSUM=${FIRST_CHECKSUM}"
    assert_line --index 1 "LETS_CHECKSUM_CHANGED=true"

    [[ -d .lets ]]
    [[ ! -d .lets/checksums ]]
    [[ ! -d .lets/checksums/with-error-code-1 ]]
}

@test "command_persist_checksum: should check if using persist_checksum without checksum will fail" {
    cd ./use_persist_without_checksum

    run lets use-persist-without-checksum

    [[ $status = 1 ]]

    assert_line --index 0 "lets: config error: failed to parse lets.yaml: 'persist_checksum' must be used with 'checksum'"
}
