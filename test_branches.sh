#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status.

echo "Initializing repository..."
rm -rf test_repo
mkdir test_repo
cd test_repo
../mini-git init

echo "Creating file on master branch..."
echo "Content on master branch" > test.txt
../mini-git add test.txt
../mini-git commit "Add test.txt on master"

echo "Creating and switching to feature branch..."
../mini-git branch feature
../mini-git checkout feature

echo "Modifying file on feature branch..."
echo "Content on feature branch" > test.txt
../mini-git add test.txt
../mini-git commit "Modify test.txt on feature"

echo "Adding all files and subdirectories to feature branch..."
touch test1.txt
mkdir test_dir
mkdir test_dir/test_subdir
echo "Content on feature branch" > test_dir/test2.txt
echo "Content on feature branch" > test_dir/test3.txt
echo "Content on feature branch" > test_dir/test_subdir/test4.txt
../mini-git add test_dir

echo "Switching back to master branch..."
../mini-git checkout master

echo "Content on master branch:"
cat test.txt

if [ "$(cat test.txt)" != "Content on master branch" ]; then
    echo "Test failed: Unexpected content on master branch"
    exit 1
fi

echo "Switching to feature branch..."
../mini-git checkout feature

echo "Content on feature branch:"
cat test.txt

if [ "$(cat test.txt)" != "Content on feature branch" ]; then
    echo "Test failed: Unexpected content on feature branch"
    exit 1
fi

echo "Test passed: File contents changed correctly between branches"
