import os
import time
os.system("sudo apt-get install golang-go")

os.system("wget --no-check-certificate https://dist.ipfs.io/go-ipfs/v0.4.15/go-ipfs_v0.4.15_linux-amd64.tar.gz")

os.system("tar -xzf go-ipfs_v0.4.15_linux-amd64.tar.gz")

os.system("cd go-ipfs/")

#os.system("ls")

os.system("cd go-ipfs && sudo ./install.sh")

os.system("sudo apt-get -y install npm")

os.system("npm i nohup")

os.system("ipfs init")

os.system("nohup ipfs daemon")

os.system("go run Encrypt.go createRSA.go signRSA.go")
