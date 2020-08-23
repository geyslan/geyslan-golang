
FROM golang:latest
LABEL maintainer="Geyslan G. Bem <geyslan@gmail.com>"

RUN apt-get -y update
RUN apt-get -y install lsb-release

# Postgresql install
RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get -y install postgresql
RUN update-rc.d postgresql enable

WORKDIR /oowlish
COPY . .

USER postgres
RUN service postgresql start && \
    psql -f /oowlish/db.sql

#RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/11/main/pg_hba.conf
#RUN echo "listen_addresses='*'" >> /etc/postgresql/11/main/postgresql.conf
#CMD ["/usr/lib/postgresql/11/bin/postgres", "-D", "/var/lib/postgresql/11/main", "-c", "config_file=/etc/postgresql/11/main/postgresql.conf"]

USER root
RUN go mod download
RUN go build

ENV OOWPORT 4000
ENV OOWFILE test.access.log

ENV DBHOST localhost
ENV DBPORT 5432
ENV DBUSER oowlish
ENV DBPASS oowlish
ENV DBDATABASE oowlish

CMD service postgresql start && ./oowlish $OOWFILE
