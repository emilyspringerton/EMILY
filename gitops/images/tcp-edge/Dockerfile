FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/nginx.conf
RUN printf '#!/bin/sh\nexec nginx -g "daemon off;"\n' > /run-nginx.sh && chmod +x /run-nginx.sh
