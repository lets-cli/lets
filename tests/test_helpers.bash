cleanup() {
    rm -rf .lets
}

# if we getting colored string we want to replace all color related symbols to get just string
strip_color() {
    echo $(echo "$1" | sed 's/\x1b\[[0-9;]*m//g')
}