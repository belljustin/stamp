#!/bin/bash
set -e
SRCPATH=$(
  cd $(dirname "$0")/..
  pwd
)

echo $SRCPATH

docker-compose -f build/docker-compose.yml up -d db
rambler -c sql/rambler.json apply -all

# var for session name (to avoid repeated occurences)
sn=stamp

# Start the session and window 0 in src folder
#   This will also be the default cwd for new windows created
#   via a binding unless overridden with default-path.
cd $SRCPATH
tmux new-session -d -s "$sn" \
    -n postgres \
    "psql -h 127.0.0.1 db user"

tmux new-window -c "$SRCPATH/truffle" \
    -n geth \
    "$SRCPATH/truffle/scripts/startGeth.sh; bash"

tmux new-window -c "$SRCPATH" \
    -n stamp \
    "make go; $SRCPATH/build/stamp -config config.json; bash"

tmux set-option -t $sn \
    remain-on-exit on

tmux attach-session -t "$sn"
