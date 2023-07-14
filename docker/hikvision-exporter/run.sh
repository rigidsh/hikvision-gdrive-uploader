#!/bin/sh

ffmpeg -i rtsp://$HIKVISION_LOGIN:$HIKVISION_PASSWORD@$HIKVISION_HOST:$HIKVISION_PORT/Streaming/Channels/$HIKVISION_CHANNEL$HIKVISION_STREAM -map 0:v -acodec copy -vcodec copy -f segment -strftime 1 -segment_time 60 -segment_format mp4 "/data/$HIKVISION_CHANNEL-%Y%m%d-%s.mp4"
