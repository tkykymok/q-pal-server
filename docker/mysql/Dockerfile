FROM --platform=linux/arm64 mysql:latest
ENV TZ=Asia/Tokyo

COPY ./docker/mysql/conf.d/my.cnf /etc/mysql/conf.d/my.cnf

CMD ["mysqld"]
