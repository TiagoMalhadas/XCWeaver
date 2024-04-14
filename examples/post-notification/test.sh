#!/bin/bash

for i in 1 2 3 4 5
do
   curl "localhost:12345/post_notification?post=post$1"
done