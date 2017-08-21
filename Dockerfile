FROM python:3.5-jessie

# install stuff
RUN apt-get update
RUN apt-get install -y  nginx
RUN apt-get install -y  supervisor

#make nginx run in the foreground
RUN echo "daemon off;" >> /etc/nginx/nginx.conf

#copy nginx stuff
COPY nginx.conf                /etc/nginx/sites-enabled/skeletor.conf
COPY supervisord-web.conf      /etc/supervisor/conf.d/supervisord.conf

EXPOSE 5000

#kick it all off
CMD ["/usr/bin/supervisord"]