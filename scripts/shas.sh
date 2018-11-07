shasum -a 256 "$@" | sed -ne 's/\([a-zA-Z0-9]*\) .*/\1/p'
