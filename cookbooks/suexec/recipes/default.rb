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

# Why patch/rebuild ?
# https://code.google.com/p/go/issues/detail?id=2617&q=group&colspec=ID%20Status%20Stars%20Priority%20Owner%20Reporter%20Summary
execute "patch & rebuild golang" do
  cwd "/usr/local/go/src"
  command <<EOS
curl https://codereview.appspot.com/download/issue13454043_3001_4001.diff | patch -p2 &&
curl https://codereview.appspot.com/download/issue13454043_3001_4002.diff | patch -p2 &&
curl https://codereview.appspot.com/download/issue13454043_3001_4003.diff | patch -p2 &&
curl https://codereview.appspot.com/download/issue13454043_3001_4004.diff | patch -p2 &&
curl https://codereview.appspot.com/download/issue13454043_3001_4005.diff | patch -p2 &&
./make.bash
EOS
  command
end
