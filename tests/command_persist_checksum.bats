load test_helpers

clean_test_files() {
    cleanup
    rm -f foo_test.txt
}

setup() {
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

    # first run, no stored checksum
    # 1. check checksum value
    # 2. check LETS_CHECKSUM_CHANGED has not changed
    # 3. check checksum persisted

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "LETS_CHECKSUM=${FIRST_CHECKSUM}" ]]
    [[ "${lines[1]}" = "LETS_CHECKSUM_CHANGED=false" ]]

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

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "LETS_CHECKSUM=${FIRST_CHECKSUM}" ]]
    [[ "${lines[1]}" = "LETS_CHECKSUM_CHANGED=false" ]]

    # third run, there is stored checksum and we creating new file. checksum must be changed now

    # create file suiting glob pattern foo_*.txt
    touch ${TEMP_FILE} && printf "footemp" > ${TEMP_FILE}

    # 1. check checksum value has changed
    # 2. check LETS_CHECKSUM_CHANGED has changed to true
    # 2. check new checksum persisted
    run lets ${CMD_NAME}
    printf "third run: %s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "LETS_CHECKSUM=${CHANGED_CHECKSUM}" ]]
    [[ "${lines[1]}" = "LETS_CHECKSUM_CHANGED=true" ]]
}

@test "command_persist_checksum: should persist checksum for cmd-as-map" {
    export CMD_NAME=persist-checksum-for-cmd-as-map

    run lets ${CMD_NAME}
    printf "%s\n" "${lines[@]}"

    # first run, no stored checksum
    # 1. check checksum value
    # 2. check LETS_CHECKSUM_CHANGED has not changed
    # 3. check checksum persisted

    [[ $status = 0 ]]

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    [[ "${lines[0]}" = "1 LETS_CHECKSUM=${FIRST_CHECKSUM}" ]]
    [[ "${lines[1]}" = "2 LETS_CHECKSUM_CHANGED=false" ]]

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
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    sort_array lines

    [[ "${lines[0]}" = "1 LETS_CHECKSUM=${FIRST_CHECKSUM}" ]]
    [[ "${lines[1]}" = "2 LETS_CHECKSUM_CHANGED=false" ]]

    # third run, there is stored checksum and we creating new file. checksum must be changed now

    # create file suiting glob pattern foo_*.txt
    touch ${TEMP_FILE} && printf "footemp" > ${TEMP_FILE}

    # 1. check checksum value has changed
    # 2. check LETS_CHECKSUM_CHANGED has changed to true
    # 2. check new checksum persisted
    run lets ${CMD_NAME}
    printf "%s\n" "${lines[@]}"

    # there is no guarantee in which order cmds will finish, so we sort output on our own
    [[ $status = 0 ]]

    sort_array lines
    [[ "${lines[0]}" = "1 LETS_CHECKSUM=${CHANGED_CHECKSUM}" ]]
    [[ "${lines[1]}" = "2 LETS_CHECKSUM_CHANGED=true" ]]
}

@test "command_persist_checksum: should persist checksum only if exit code = 0" {
    run lets with-error-code-1
    printf "%s\n" "${lines[@]}"

    [[ $status = 1 ]]
    [[ "${lines[0]}" = "LETS_CHECKSUM=${FIRST_CHECKSUM}" ]]
    [[ "${lines[1]}" = "LETS_CHECKSUM_CHANGED=false" ]]

    [[ -d .lets ]]
    [[ ! -d .lets/checksums ]]
    [[ ! -d .lets/checksums/with-error-code-1 ]]
}

@test "command_persist_checksum: should check if using persist_checksum without checksum will fail" {
    cd ./use_persist_without_checksum

    run lets use-persist-without-checksum
    printf "%s\n" "${lines[@]}"

    [[ $status = 1 ]]

    [[ "${lines[0]}" = "[ERROR] failed to load config file lets.yaml: failed to parse 'use-persist-without-checksum' command: field persist_checksum: you must declare 'checksum' for command to use 'persist_checksum'" ]]
}
