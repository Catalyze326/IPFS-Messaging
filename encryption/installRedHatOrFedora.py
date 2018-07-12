import os
#Install go for the backend
os.system("wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz")
os.system("tar -xzf go1.10.3.linux-amd64.tar.gz")
os.system("mv go /usr/local")
os.system("export GOROOT=/usr/local/go")

#Install the ipfs api
os.system("go get -u github.com/ipfs/go-ipfs-api")

#download ipfs, and the two things we need for ipfs clustering
os.system("wget --no-check-certificate https://dist.ipfs.io/go-ipfs/v0.4.15/go-ipfs_v0.4.15_linux-amd64.tar.gz")

#Unzip all the previus things
os.system("tar -xzf go-ipfs_v0.4.15_linux-amd64.tar.gz")

#install and initalize ipfs
os.system("cd go-ipfs && sudo ./install.sh")
os.system("ipfs init")

#clean up the mess and delete all the files I just downloaded
os.system("rm -r ipfs*")
os.system("rm -r go-i*")
