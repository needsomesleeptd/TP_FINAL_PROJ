# Official Nginx image (Yes, in this article I always use the `latest`. Kill me!)
FROM nginx:latest

RUN rm  /etc/nginx/conf.d/*.conf
# Copy Nginx configuration file to the container
COPY nginx.conf /etc/nginx/conf.d

# Expose port 80
EXPOSE 80