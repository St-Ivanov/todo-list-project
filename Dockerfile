FROM ubuntu:latest

WORKDIR /todo_app

COPY todo_app ./

RUN mkdir ./web

COPY web ./web

RUN chmod +x todo_app

EXPOSE 7540

CMD ["./todo_app"]