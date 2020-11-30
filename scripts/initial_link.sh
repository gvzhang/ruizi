#! /bin/bash

work_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

cmd_output="$work_dir/dist/ruizi-initial_link"
log_output="$work_dir/dist/ruizi-initial_link.error.log"
rm -f "$cmd_output"
rm -f "$log_output"
/usr/local/bin/go build -o "$cmd_output" "$work_dir/cmd/initial_link/main.go"
if [ $? -eq 0 ]; then
  "$cmd_output" 1>"$log_output" 2>&1 &
else
  echo "$cmd_output build fail"
fi
