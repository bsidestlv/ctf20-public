FROM ubuntu:16.04
RUN apt-get update && apt-get install -y redis-server nginx php-fpm


## Admin user
RUN useradd -c 'Admin' -m -d /home/admin -s /bin/bash admin
RUN mkdir -p /var/www/html
ADD --chown='admin:admin' files/flag.txt /home/admin/flag.txt

#Entrypoint files
ADD files/run.sh /root/run.sh
ADD files/persistence.sh /root/persistence.sh
ADD files/index.php /root/index.php
RUN chmod +x /root/persistence.sh
RUN chmod +x /root/run.sh

#Add web configurations and default files
RUN chmod 0777 /var/www/html
ADD --chown='root:root' files/index.php /var/www/html/index.php
ADD files/default /etc/nginx/sites-available/default


EXPOSE 6379
EXPOSE 80
CMD ["/root/run.sh"]