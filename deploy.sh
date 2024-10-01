#!/bin/bash

set -e

git pull
if [ $? -ne 0 ]; then
    echo "Git pull failed"
    exit 1
fi

go build
if [ $? -ne 0 ]; then
    echo "Go build failed"
    exit 1
fi

systemctl restart goorder
if [ $? -ne 0 ]; then
    echo "Restarting goorder failed"
    exit 1
fi

systemctl restart nginx
if [ $? -ne 0 ]; then
    echo "Restarting nginx failed"
    exit 1
fi

echo "Deployed successfully!"