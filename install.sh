#!/bin/bash

if [ `command -v python` != "" ]; then
    sudo python -c "$(curl https://raw.githubusercontent.com/amblar/emailbomber/master/install.py -s)"
elif [ `command -v python3` != "" ]; then
    sudo python3 -c "$(curl https://raw.githubusercontent.com/amblar/emailbomber/master/install.py -s)"
else
    echo "Please install any version of python"
fi