FROM nginx

RUN rm /etc/nginx/nginx.conf
RUN rm /etc/nginx/conf.d/default.conf
ADD confs/default.conf /etc/nginx/conf.d/default.conf
ADD confs/nginx.conf /etc/nginx/nginx.conf
ADD site/index.html /usr/share/nginx/html/index.html
