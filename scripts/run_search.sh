#! /bin/bash

work_dir=$(cd "$(dirname "$0")"; cd ..; pwd)
for pid in $(pgrep -f "ruizi-search"); do
  kill "$pid"
done

cp -r "$work_dir/internal/search/api/templates" "$work_dir/dist"

run_cmd=("search")
for cmd in "${run_cmd[@]}"; do
  cmd_output="$work_dir/dist/ruizi-$cmd"
  log_output="$work_dir/dist/ruizi-$cmd.error.log"
  rm -f "$cmd_output"
  rm -f "$log_output"
  /usr/local/bin/go build -o "$cmd_output" "$work_dir/cmd/$cmd/main.go"
  if [ $? -eq 0 ]; then
    "$cmd_output" 1>"$log_output" 2>&1 &
  else
    echo "$cmd_output build fail"
  fi
done
