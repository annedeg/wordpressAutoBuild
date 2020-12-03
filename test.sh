#!/bin/sh
sudo apt-get install -y apache2 php7.2 php7.2-cli php7.2-common php7.2-mysql mysql-server &&
sudo ufw allow in "Apache Full" &&
sudo apt install mysql-server &&
sudo mysqladmin -p create wordpress || true &&
sudo mysql -p -Bse "drop user 'gouser'@localhost; flush privileges; CREATE USER 'gouser'@'localhost' IDENTIFIED BY 'VERY-HARD-PASSWORD'; GRANT ALL PRIVILEGES ON wordpress.* TO 'gouser'@'localhost';FLUSH PRIVILEGES;"
