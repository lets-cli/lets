cleanup() {
    rm -rf .lets
}

# Usage:

# my_array=(2,4,1)
# sort_array my_array
# printf "%s" "${my_array[@]}" # -- will print 1 2 4
sort_array() {
    local -n array_to_sort=$1
    IFS=$'\n' array_to_sort=($(sort <<<"${array_to_sort[*]}"))
    unset IFS
}