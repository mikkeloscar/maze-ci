#!/bin/bash

workdir="maze_workdir"
maze_workspace="/var/lib/maze"
nspawn_container="$maze_workspace/nspawn_template"
src_dir="/home/vagrant/go"
container_user="maze"
mkdir -p $workdir
chown -R vagrant:vagrant $workdir
chown -R vagrant:vagrant $src_dir

# Update vm and install:
# git
# go
# vim
# arch-install-scripts
pacman -Syu vim git go arch-install-scripts --noconfirm

# install termite terminfo
curl -O https://raw.githubusercontent.com/thestinger/termite/master/termite.terminfo
tic -x termite.terminfo -o /usr/share/terminfo
rm termite.terminfo

# setup nspawn container template
mkdir -p $maze_workspace
btrfs subvolume create $nspawn_container

# setup fresh arch root system
pacstrap -c -d $nspawn_container base base-devel pkgbuild-introspection

# remove linux package
systemd-nspawn -D $nspawn_container pacman -Rns linux --noconfirm

# setup builduser
systemd-nspawn -D $nspawn_container useradd -m -s /bin/bash $container_user
/bin/bash -c "echo \"$container_user ALL=(ALL) NOPASSWD: ALL\" >> $nspawn_container/etc/sudoers.d/$container_user"


echo "" >> /home/vagrant/.bashrc
echo "export GOPATH=\$HOME/go" >> /home/vagrant/.bashrc
echo "cd go/src/github.com/mikkeloscar/maze-ci" >> /home/vagrant/.bashrc
