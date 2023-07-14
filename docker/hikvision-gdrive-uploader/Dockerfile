FROM ubuntu:23.04
RUN apt-get update && apt-get install -y ca-certificates
ADD hikvision-gdrive-uploader /hikvision-gdrive-uploader
WORKDIR /
ENTRYPOINT ["./hikvision-gdrive-uploader"]
