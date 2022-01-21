FROM golang:alpine AS builder
RUN apk add --no-cache build-base
RUN mkdir -p /azureprocessor
WORKDIR /azureprocessor
COPY . .
RUN go mod download
WORKDIR /azureprocessor/scripts
RUN g++ cmp.c -o cmp
WORKDIR /azureprocessor/cmd/processor/azure
RUN CGO_ENABLED=0 GOOS=linux go build

FROM mcr.microsoft.com/azure-functions/dotnet:3.0-appservice
ENV AzureWebJobsScriptRoot=/home/site/wwwroot AzureFunctionsJobHost__Logging__Console__IsEnabled=true
RUN apt update
RUN apt-get -y install build-essential python2 python3 default-jdk bash
WORKDIR /home/site/wwwroot
COPY --from=builder /azureprocessor/cmd/processor/azure /home/site/wwwroot
COPY --from=builder /azureprocessor/scripts /home/site/wwwroot