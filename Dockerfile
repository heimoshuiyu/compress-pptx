FROM golang AS build
WORKDIR /src
COPY . ./
ENV GOPROXY="goproxy.cn"

RUN go build -v -o /compress-pptx

FROM jrottenberg/ffmpeg
WORKDIR /
COPY --from=build /src/compress-pptx /compress-pptx
EXPOSE 8888
ENTRYPOINT ["/compress-pptx"]
