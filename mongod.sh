#!/bin/bash
sudo systemctl stop mongod
sudo apt remove -f mongod* -y 
sudo apt purge --allow-change-held-packages mongod* -y 
sudo apt autoremove -y 
sudo rm /etc/apt/sources.list.d/mongodb*.list 
sudo rm /usr/share/keyrings/mongodb-archive-keyring.gpg
sudo rm -r /var/log/mongodb 
sudo rm -r /var/lib/mongodb 
curl -fsSL https://pgp.mongodb.com/server-4.4.asc | gpg --dearmor | sudo tee /usr/share/keyrings/mongodb-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/mongodb-archive-keyring.gpg] https://repo.mongodb.org/apt/debian buster/mongodb-org/4.4 main" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.4.list
sudo apt-get update -y
sudo apt install -y mongodb-org=4.4.18 mongodb-org-server=4.4.18 mongodb-org-shell=4.4.18 mongodb-org-mongos=4.4.18 mongodb-org-tools=4.4.18 
echo "mongodb-org hold" | sudo dpkg --set-selections 
echo "mongodb-org-server hold" | sudo dpkg --set-selections  
echo "mongodb-org-shell hold" | sudo dpkg --set-selections 
echo "mongodb-org-mongos hold" | sudo dpkg --set-selections
echo "mongodb-org-tools hold" | sudo dpkg --set-selections
sudo mkdir -p /var/log/mongodb 
sudo mkdir -p /var/lib/mongodb 
sudo chown -R mongodb:mongodb /var/log/mongodb
sudo chown -R mongodb:mongodb /var/lib/mongodb
sudo chown mongodb:mongodb /tmp/mongodb-27017.sock 
sudo systemctl enable mongod 
sudo systemctl start  mongod 
