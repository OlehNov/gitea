FROM golang:1.22 AS build

WORKDIR /app

RUN apt-get update && \
    apt-get install -y curl && \
    curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . ./

RUN make backend

RUN make frontend

RUN GOOS=linux GOARCH=amd64 go build -o backend ./main.go

RUN chmod +x backend

COPY package*.json ./

RUN npm install

FROM alpine:latest AS production

WORKDIR /app

RUN apk add --no-cache libc6-compat

RUN adduser -D -g '' gitea_user

COPY --from=build /app/backend /app/

COPY --from=build /app/public /app/public
COPY --from=build /app/templates /app/templates

COPY --from=build /app/package*.json /app/
COPY --from=build /app/node_modules /app/node_modules

RUN chown -R gitea_user:gitea_user /app

USER gitea_user

CMD ["./backend"]