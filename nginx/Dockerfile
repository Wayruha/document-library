# Docker image for nginx

FROM alpine:latest

MAINTAINER Saif Alam <aust.saif@gmail.com>

# update packages 
RUN apk update

# install nginx 
RUN apk --no-cache add nginx 
RUN mkdir -pv /run/nginx

# Copy .conf file 
COPY nginx.conf /etc/nginx/

# Used for testing purpose 
#COPY ./www/ /var/www/nginx

#Expose port
EXPOSE 80

# Define default command.
CMD ["nginx"]
