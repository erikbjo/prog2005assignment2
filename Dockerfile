FROM golang:1.22 AS builder

LABEL authors="erbj@stud.ntnu.no,simonhou@stud.ntnu.no"
LABEL stage=builder

WORKDIR /assignment-2

# Copy the entire source code into the container
COPY . .

# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o executable ./cmd/api/main.go

# Define exposed port
EXPOSE 8000 8001

# Create .env file
# TODO: Export local .env file to container, this is a temporary solution
RUN echo 'PORT="8000"' > .env

# Entrypoint command
ENTRYPOINT ["./executable"]