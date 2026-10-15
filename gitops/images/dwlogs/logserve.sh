#!/bin/sh
# Log files on the PVC are 0660 group 65532 (the pod's fsGroup); run the workers with that primary group so they can read them.
# nginx's `user` directive resolves the group BY NAME (a bare "65532" fails getgrnam -> [emerg], crashloop), so create a named group first.
getent group 65532 >/dev/null 2>&1 || addgroup -g 65532 dwlogs
G="$(getent group 65532 | cut -d: -f1)"
sed -i "s/^user .*/user nginx $G;/" /etc/nginx/nginx.conf
exec nginx -g "daemon off;"
