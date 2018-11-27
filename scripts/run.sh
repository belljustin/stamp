#!/bin/bash
set -e
parent_path=$(dirname ${BASH_SOURCE[0]})/..
echo $parent_path

docker-compose -f build/docker-compose.yml up -d db
rambler -c sql/rambler.json apply -all

# var for session name (to avoid repeated occurences)
sn=stamp

# Start the session and window 0 in src folder
#   This will also be the default cwd for new windows created
#   via a binding unless overridden with default-path.
cd $parent_path
tmux new-session -s "$sn" -n stamp -d


# TODO: pass in db creds
tmux new-window -c "$parent_path" \
    -n postgres \
    "psql -h 127.0.0.1 db user"

tmux new-window -c "$parent_path/truffle" \
    -n geth \
    "$parent_path/scripts/startGeth.sh"

tmux new-window -c "$parent_path" \
    -n stamp \
    "make go; ./build/stamp -config config.json; bash"

tmux attach-session -t "$sn"
