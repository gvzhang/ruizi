#! /bin/bash

work_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

for pid in $(pgrep -f "ruizi"); do
  kill "$pid"
done

rm -f $work_dir/data/*.bin

sh "$work_dir/data/init_link.sh"

