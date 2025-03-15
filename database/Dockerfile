FROM postgres:latest

COPY *.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_DB=main
ENV POSTGRES_USER=tester
ENV POSTGRES_PASSWORD=testing

EXPOSE 5432