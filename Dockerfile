FROM golang AS build
WORKDIR /src
COPY . ./
ENV GOPROXY="goproxy.cn"

RUN go mod download
ENV CGO_ENABLED=0
RUN go build -v -o /compress-pptx

FROM jrottenberg/ffmpeg
WORKDIR /
COPY --from=build /compress-pptx /compress-pptx
EXPOSE 8888
ENTRYPOINT ["/compress-pptx"]
