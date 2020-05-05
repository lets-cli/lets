load test_helpers

setup() {
    cd ./tests/find_config/child/another_child
    rm -rf ../../.lets
    cleanup
}

@test "find_lets_file: should find lets.yaml in parent dir" {
    run lets foo
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]
    [[ "${lines[0]}" = "foo" ]]
}

@test "find_lets_file: .lets must be created in the same dir where lets.yaml placed" {
    run lets foo
    printf "%s\n" "${lines[@]}"

    [[ $status = 0 ]]

    [[ ! -d .lets ]]
    [[ -d ../../.lets ]]
}
