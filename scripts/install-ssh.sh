#!/usr/bin/env bash

set -e

echo "Installing SSH for your system (Aptible requires an older version)"

mkdir ./tmp
pushd tmp
if [[ "$OSTYPE" == "darwin"* ]]; then
  wget https://aptible-ssh-binaries.s3.us-east-2.amazonaws.com/ssh.zip
  unzip ssh.zip
  pushd openssh-7.3p1
fi

for file in "ssh" "ssh-keygen"; do
  mv ./$file ../../public/
done

pushd -0
dirs -c
rm -rf ./tmp
