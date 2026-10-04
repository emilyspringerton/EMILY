FROM nginx:1.27-alpine
RUN apk add --no-cache curl
RUN rm -f /etc/nginx/conf.d/default.conf
COPY logserve.conf /etc/nginx/conf.d/logserve.conf
COPY logserve.sh logsync.sh /
RUN chmod +x /logserve.sh /logsync.sh
