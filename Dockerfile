FROM golang:1.22 as builder

LABEL authors="erbj@stud.ntnu.no,simonhou@stud.ntnu.no"
LABEL stage=builder

WORKDIR /assignment-2

# Copy the entire source code into the container
COPY . .

# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o executable ./cmd/api/main.go

# Define exposed port
EXPOSE 8080

# Create .env file
# TODO: Export local .env file to container, this is a temporary solution
RUN echo "PORT=8080" > .env

# Entrypoint command
ENTRYPOINT ["./executable"]