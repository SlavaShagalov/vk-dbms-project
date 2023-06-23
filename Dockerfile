FROM golang:1.20.4-alpine3.17 as builder
COPY . /forum
WORKDIR /forum
RUN go mod download

RUN go build -o app ./cmd/backend/main.go

#FROM postgres:14-alpine3.18 as db
FROM postgres:14 as db
USER postgres
COPY ./scripts .
RUN pg_createcluster 14 main && \
    /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'password';" &&\
    createdb -O forum forum && \
    psql -f ./init.sql -d forum && \
    /etc/init.d/postgresql stop

#psql -f ./fill_data.sql -d forum && \

ENV POSTGRES_USER forum
ENV POSTGRES_DB forum
ENV POSTGRES_PASSWORD password
ENV POSTGRES_HOST localhost
ENV POSTGRES_PORT 5432

USER root
COPY --from=builder /forum/app app
ENV PGPASSWORD password
CMD service postgresql start && ./app
