FROM ubuntu:16.04
RUN apt-get update && apt-get install -y redis-server nginx


## Admin user
RUN useradd -c 'Admin' -m -d /home/admin -s /bin/bash admin
RUN mkdir -p /var/www/html

#Entrypoint files
ADD files/run.sh /root/run.sh
ADD files/persistence.sh /root/persistence.sh
RUN chmod +x /root/persistence.sh
RUN chmod +x /root/run.sh

#Add web configurations and default files
RUN chmod 0777 /var/www/html
ADD --chown='root:root' files/index.html /var/www/html/index.html
ADD files/default /etc/nginx/sites-available/default


EXPOSE 6379
EXPOSE 80
CMD ["/root/run.sh"]