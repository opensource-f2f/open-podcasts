# goplay

Listen podcast via CLI.

## Get started

Install via [hd](https://github.com/linuxsuren/http-downloader):

```shell
hd install goplay
```

Start to listen:

```shell
goplay 开源面对面
```

Or, you can also use it in a container (only for Linux):

```shell
docker run --device /dev/snd surenpi/goplay:latest@sha256:badde269814a4ed88737c6d5fc0bff15b0e0b9e3f43dd9fd0f3eeb281a7d05b1 sh
```