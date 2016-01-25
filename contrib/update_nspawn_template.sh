#!/bin/bash

nspawn_container=${1:-"/var/lib/maze/nspawn_template"}

# update container
systemd-nspawn -D $nspawn_container pacman -Syu --noconfirm

# clear any .pacnew files
find $nspawn_container -name "*.pacnew" -exec rename .pacnew '' '{}' \;
