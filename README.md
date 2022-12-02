# Compress PPTX

Well, it can also compress `docx`, `odp`, `odt` and other openxml format documents.

It use `ffmpeg` to transcode all media files in document.

- For images, it usually reduce 50% of the file size and without lossing to much image detail.

- For videos, it perform transcode and **audio normalization**.

This program has not been rigorously tested at the moment, so be sure to check the compressed file.

## Demo

<https://yongyuancv.cn/compress-pptx/>

## Self-host

First of all, clone and cd this repo

```shell
git clone https://github.com/heimoshuiyu/compress-pptx.git
cd compress-pptx
```

Then basically you have 2 ways to deploy:

### The docker way

Run

```shell
docker build -t [image_name] .
docker run -it --rm -p 8888:8888 [image_name]
```

Then goto <http://localhost:8888> and you will see the app running.

### The manually way

First make sure your `ffmpeg` is installed and set the `$PATH` envirable correctly.

Then

```shell
go build -v -o compress-pptx
./compress-pptx
```

Then goto [http://localhost:8888](http://localhost:8888) and you will see the app running.
