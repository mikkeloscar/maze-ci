# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "arch_btrfs"

  config.vm.synced_folder "./", "/home/vagrant/go/src/github.com/mikkeloscar/maze-ci",
      owner: "vagrant",
      group: "vagrant"

  config.vm.provision "shell", path: ".vagrant_bootstrap.sh"
end
