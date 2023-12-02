FROM golang:1.20

WORKDIR /go_recipe

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

WORKDIR /go_recipe/cmd/recipe_app
RUN CGO_ENABLED=0 GOOS=linux go build -o /greenlight

EXPOSE 4000
CMD ["/greenlight"]
