#!/usr/bin/env bash

VERSION='v0.0.2'

packageName='todo_x86-64-linux'

URL="https://github.com/IridiumNan/project-todo/releases/download/$VERSION/$packageName"

echo "using wget command to download"

mkdir -p ~/.local/bin/

wget -O ~/.local/bin/todo "$URL"

echo "todo cli will be saved on ~/.local/bin/"

echo 'you should run export PATH=$PATH:$HOME/.local/bin/ to make it visiable'
echo "exec: chmod +x ~/.local/bin/todo"

chmod +x ~/.local/bin/todo

echo "Run   todo --help  or visit https://github.com/IridiumNan/project-todo for more details thanks"
