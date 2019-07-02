FROM scratch
EXPOSE 80
WORKDIR /app
COPY . /app
ENTRYPOINT ["/app/main"]
