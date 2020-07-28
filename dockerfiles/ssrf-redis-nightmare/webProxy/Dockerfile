FROM ubuntu:16.04
RUN apt-get update && apt-get install -y nginx php7.0-fpm php-curl


RUN mkdir -p /var/www/html
ADD files/code/ /var/www/html/
ADD files/default /etc/nginx/sites-available/default
ADD files/run.sh /root/run.sh
RUN chmod +x /root/run.sh
EXPOSE 80
CMD ["/root/run.sh"]