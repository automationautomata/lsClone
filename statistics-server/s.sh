#!/bin/bash

sudo mkdir $(pwd)/html/stat.info
sudo chown www-data:www-data -R $(pwd)/html/stat.info
sudo chmod 775 -R $(pwd)/html/stat.info

sudo touch /etc/apache2/sites-available/stat.info.conf
sudo echo "
 <VirtualHost *:80>
    ServerAdmin webmaster@stat.info
    ServerName stat.info
    DocumentRoot  $(pwd)/html/stat.info

    ErrorLog ${APACHE_LOG_DIR}/error.log
    CustomLog ${APACHE_LOG_DIR}/access.log combined
</VirtualHost> 
" > /etc/apache2/sites-available/stat.info.conf

sudo a2ensite stat.info.conf
sudo systemctl restart apache2 