# Step 1: Build the UI
FROM node:18 as build-ui
WORKDIR /app

COPY . .
RUN npm install
RUN npm run build:about 

# Step 2: Build the executable
FROM golang:1.20 AS build
WORKDIR /app

COPY . .
COPY --from=build-ui /app/internal/apps/about/assets/css ./internal/apps/about/assets/css

ENV CGO_ENABLED=0

RUN go build -o letsfixoss .

# Step 3: Create a minimal final image
FROM debian:bullseye-slim

COPY --from=build /app/letsfixoss /usr/local/bin/letsfixoss
ENV LISTEN_ADDR=0.0.0.0:8080

ENTRYPOINT ["letsfixoss"]
