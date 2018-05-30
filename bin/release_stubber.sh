#!/bin/bash -ex

# This requires that all original releases be downloaded in the $PWD/releases directory.
# This script will create stubbed versions of releases with the release.MF file in
# test/releases directory.

FILES=$PWD/releases/*

mkdir -p $PWD/test/releases
for f in $FILES
do
  tar --extract --file=$f release.MF
  tar -czvf $PWD/test/releases/$(basename $f) release.MF
  rm $PWD/release.MF
  echo $f
done
