#!/bin/bash

cd "$(dirname "${BASH_SOURCE[0]}")"/../../migrations
set -e

if [ -z "$1" ]; then
  echo "USAGE: $0 <name>"
  exit 1
fi

awkcmd=$(
  cat <<-EOF
BEGIN { FS="[_/]" }
{ n=\$2 }
END {
    gsub(/[^A-Za-z0-9]/, "_", name);
    printf("%s_%s.up.sql\n",   n + 1, name);
    printf("%s_%s.down.sql\n", n + 1, name);
}
EOF
)

files=$(find -s . -type f -name '[0-9]*.sql' | sort -n -t '/' -k 2 | awk -v name="$1" "$awkcmd")

for f in $files; do
  cat >"$f" <<EOF
BEGIN;



COMMIT;
EOF

  echo "Created migrations/$f"
done
