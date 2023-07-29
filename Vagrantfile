Vagrant.configure(2) do |config|

  config.vm.box = "oss4u/opnsense"
  #config.vm.box_version = "23.1"

  config.winrm.timeout = 120
  config.winrm.retry_limit = 100
  config.ssh.username = "root"
  config.ssh.password = "opnsense"


  config.vm.provider 'virtualbox' do |vb|
    vb.memory = 1024
    vb.cpus = 1
    vb.gui = false # want gui for testing
    vb.customize ['modifyvm', :id, '--nic1', 'nat'] # don't touch this interface!

    # Setup firewall port assignments
    vb.customize ['modifyvm', :id, '--nic2', 'intnet']
    vb.customize ['modifyvm', :id, '--intnet2', 'intnet']
    vb.customize ['modifyvm', :id, '--nic3', 'intnet']
    vb.customize ['modifyvm', :id, '--intnet3', 'intnet']
    vb.customize ['modifyvm', :id, '--nic4', 'intnet']
    vb.customize ['modifyvm', :id, '--intnet4', 'intnet']
  end

  config.vm.network :forwarded_port, guest: 22, host: 10022, id: "ssh-orig", auto_correct: true
  config.vm.network :forwarded_port, guest: 80, host: 10080, id: "http", auto_correct: true
  #config.vm.network :forwarded_port, guest: 443, host: 10443, id: "https", auto_correct: true
  config.vm.network :forwarded_port, guest: 22, host: 2222, id: "ssh", auto_correct: true

  #config.vm.provision "file", source: "config.xml", destination: "/conf/config.xml" # copy default config to firewall
  #config.vm.provision "shell", inline: "opnsense-shell reload" # apply configuration

end