docker run -p 3306:3306 --name some-mysql \
    -e MYSQL_RANDOM_ROOT_PASSWORD=yes \
    -e MYSQL_DATABASE=bad-data \
    -e MYSQL_USER=baduser \
    -e MYSQL_PASSWORD=terriblyinsecurepassword \
    mysql