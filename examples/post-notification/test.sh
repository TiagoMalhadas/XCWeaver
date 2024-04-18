#!/bin/bash

duration=$1

# Function to print progress bar
print_progress() {
    # Calculate percentage completion
    elapsed=$1
    progress=$((elapsed * 100 / duration))

    printf "["
    for ((i=0; i<progress/2; i++)); do printf "#"; done
    for ((i=progress/2; i<50; i++)); do printf " "; done
    printf "] %d%%\r" $progress
}

# Start time
start=$(date +%s)

echo "Starting requests for $1 secounds..."
while true; do

   curl "localhost:12345/post_notification?post=my_first_post2"

   now=$(date +%s)
   elapsed=$((now - start))
   print_progress $elapsed

   # Check if the elapsed time exceeds the duration
   if [ $elapsed -ge $duration ]; then
      break
   fi
done

print_progress $duration
echo "Test completed."
