#!/bin/sh
secret="`dirname "$0"`/.env"
echo $secret
source $secret
day=$1
if [[ -z $1 ]]; then
    echo "Missing day argument"
    exit 1
fi
#day=$(pwd | awk -F'/' '{print $NF}' | awk -F'.' '{print $1}' | sed 's/^0*//')

curl "https://adventofcode.com/2019/day/$day" -o /tmp/aoc --compressed -H "Cookie: session=$secret"
dir=$(cat /tmp/aoc | grep -Eo '<h2>--- Day[^<>]+</h2>' | cut -c13- | rev | cut -c 10- | rev | sed 's/:/./' | sed -e 's/^\([0-9]\)\./0\1./')
echo "Creating directory $dir"
mkdir "$dir"
cd "$dir"
curl "https://adventofcode.com/2019/day/$day/input" -o input -H "Cookie: session=$secret"
