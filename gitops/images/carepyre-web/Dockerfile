FROM nginx:1.27-alpine
RUN rm -f /etc/nginx/conf.d/default.conf
COPY default.conf /etc/nginx/conf.d/default.conf
COPY site/ /usr/share/nginx/html/
RUN printf '#!/bin/sh\nexec nginx -g "daemon off;"\n' > /run-nginx.sh && chmod +x /run-nginx.sh
