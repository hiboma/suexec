%w[
  git
].each do |p|
  package p do
    action :install
  end
end

file "/etc/profile.d/go.sh" do
  mode 0755
  content <<'EOS'
export PATH=/usr/local/go/bin:/home/vagrant/go/bin:$PATH
export GOPATH=/home/vagrant/go
EOS
  action :create
end

directory "/home/vagrant/go/src/github.com/hiboma" do
  owner "vagrant"
  group "vagrant"
  action :create
  recursive true
end

[
 "/home/vagrant/go/src/github.com/",
 "/home/vagrant/go/src/",
 "/home/vagrant/go/pkg",
 "/home/vagrant/go/"
].each do |path|
  directory path do
    owner "vagrant"
    group "vagrant"
  end
end

execute "install golang" do
  cwd "/usr/local/src"
  command <<EOS
wget https://go.googlecode.com/files/go1.1.2.linux-amd64.tar.gz &&
tar zxvf go1.1.2.linux-amd64.tar.gz                             &&
mv go /usr/local/go
EOS
  creates "/usr/local/go/bin/go"
end
