#!/bin/bash
sudo systemctl stop mongod > /dev/null 2>&1 
sudo apt remove -f mongod* -y > /dev/null 2>&1 
sudo apt purge --allow-change-held-packages mongod* -y > /dev/null 2>&1 
sudo apt autoremove -y > /dev/null 2>&1 
sudo rm /etc/apt/sources.list.d/mongodb*.list > /dev/null 2>&1
sudo rm /usr/share/keyrings/mongodb-archive-keyring.gpg > /dev/null 2>&1 
curl -fsSL https://pgp.mongodb.com/server-4.4.asc | gpg --dearmor | sudo tee /usr/share/keyrings/mongodb-archive-keyring.gpg > /dev/null 2>&1
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/mongodb-archive-keyring.gpg] https://repo.mongodb.org/apt/debian bionic/mongodb-org/4.4 main" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.4.list > /dev/null 2>&1
sudo apt-get update -y > /dev/null 2>&1
sudo apt install -y mongodb-org=4.4.18 mongodb-org-server=4.4.18 mongodb-org-shell=4.4.18 mongodb-org-mongos=4.4.18 mongodb-org-tools=4.4.18 > /dev/null 2>&1 && sleep 2
echo "mongodb-org hold" | sudo dpkg --set-selections > /dev/null 2>&1
echo "mongodb-org-server hold" | sudo dpkg --set-selections > /dev/null 2>&1 
echo "mongodb-org-shell hold" | sudo dpkg --set-selections > /dev/null 2>&1 
echo "mongodb-org-mongos hold" | sudo dpkg --set-selections > /dev/null 2>&1 
echo "mongodb-org-tools hold" | sudo dpkg --set-selections > /dev/null 2>&1 
sudo chown -R mongodb:mongodb /var/lib/mongodb > /dev/null 2>&1
sudo chown mongodb:mongodb /tmp/mongodb-27017.sock > /dev/null 2>&1
sudo systemctl enable mongod > /dev/null 2>&1
sudo systemctl start  mongod > /dev/null 2>&1
