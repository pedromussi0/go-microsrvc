FROM alpine:latest

RUN mkdir /app

COPY --from=build brokerApp /app

CMD ["/app/brokerApp"]