FROM golang:1.21 as build
ENV CGO_ENABLED=0
WORKDIR /src
COPY . .
RUN make bin

FROM alpine:3.19
COPY --from=build /src/bin/ginkgo2allure /bin/ginkgo2allure
ENTRYPOINT ["/bin/ginkgo2allure"]