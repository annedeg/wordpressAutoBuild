#!/bin/sh
echo "1" &&
sudo apt-get install -y php7.2 php7.2-cli php7.2-common php7.2-mysql mysql-server &&
sudo ufw allow in "Apache Full" &&
sudo apt install mysql-server &&
sudo mysql_secure_installation &&
sudo mysqladmin -p create wordpress &&
sudo mysql -p -Bse "CREATE USER 'gouser'@'localhost' IDENTIFIED BY 'VERY-HARD-PASSWORD';GRANT ALL PRIVILEGES ON `wordpress`.* TO 'gouser'@'localhost';FLUSH PRIVILEGES;exit;" &&
echo "ben hier nu bro"