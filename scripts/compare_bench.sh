#!/usr/bin/env bash
set -euo pipefail

THRESHOLD=${THRESHOLD:-1.10}

# Extract all benchmark names dynamically
benchmarks=$(grep '^Benchmark' bench.txt | awk '{print $1}')

for name in $benchmarks; do
  echo "🔍 Checking $name ..."

  baseline=$(influx query '
    from(bucket:"'"$INFLUX_BUCKET"'")
    |> range(start: -30d)
    |> filter(fn:(r)=>r.branch=="main" and r.name=="'"$name"'")
    |> last()
    |> findRecord(fn:(key)=> true, idx:0)
  ' --raw | tail -n1 | awk '{print $2}' || echo 0)

  current=$(grep "$name" bench.txt | awk '{print $3}' | tr -d 'ns/op' | head -n1)

  if [[ "$baseline" == "0" || -z "$baseline" ]]; then
    echo "ℹ️ No baseline found for $name, skipping comparison."
    continue
  fi

  if (( $(echo "$current > $baseline * $THRESHOLD" | bc -l) )); then
    echo "❌ Regression detected in $name: ${current}ns vs baseline ${baseline}ns"
    exit 1
  else
    echo "✅ $name OK (${current}ns vs ${baseline}ns)"
  fi
done
