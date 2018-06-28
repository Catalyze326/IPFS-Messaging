import os
#Install go for the backend
os.system("wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz")
os.system("tar -xzf go1.10.3.linux-amd64.tar.gz")
os.system("mv go /usr/local")
os.system("export GOROOT=/usr/local/go")

#download ipfs, and the two things we need for ipfs clustering
os.system("wget --no-check-certificate https://dist.ipfs.io/go-ipfs/v0.4.15/go-ipfs_v0.4.15_linux-amd64.tar.gz")
os.system("wget --no-check-certificate https://dist.ipfs.io/ipfs-cluster-ctl/v0.4.0/ipfs-cluster-ctl_v0.4.0_linux-amd64.tar.gz")
os.system("wget --no-check-certificate https://dist.ipfs.io/ipfs-cluster-service/v0.4.0/ipfs-cluster-service_v0.4.0_linux-amd64.tar.gz")

#Unzip all the previus things
os.system("tar -xzf go-ipfs_v0.4.15_linux-amd64.tar.gz")
os.system("tar -xzf ipfs-cluster-ctl_v0.4.0_linux-amd64.tar.gz")
os.system("tar -xzf ipfs-cluster-service_v0.4.0_linux-amd64.tar.gz")

#install and initalize ipfs
os.system("cd go-ipfs && sudo ./install.sh")
os.system("ipfs init")

#Copy the two things we need for ipfs clustering into /usr/bin so we can just use them as commands and initalize the clustering service
os.system("sudo cp ipfs-cluster-ctl/ipfs-cluster-ctl /usr/bin | sudo cp ipfs-cluster-service/ipfs-cluster-service /usr/bin")
os.system("ipfs-cluster-service init")

#Install openssl so that we can generate rsa and aes tokens
os.system("sudo apt-get -y install openssl")

#clean up the mess and delete all the files I just downloaded
os.system("rm -r ipfs*")
os.system("rm -r go-i*")
