#!/bin/sh
# Log files on the PVC are 0660 group 65532 (the pod's fsGroup); run the workers with that primary group so they can read them.
sed -i 's/^user .*/user nginx 65532;/' /etc/nginx/nginx.conf
exec nginx -g "daemon off;"
