#!/bin/bash
INPUT=$1
NAME=${INPUT#*iso/}
NAME=${NAME%.*}
OUTPUT=mp4/${NAME%.*}.mp4
ffmpeg -i $INPUT $OUTPUT \
  && ./upload_video.py --file=${OUTPUT} --title=${NAME} --privacyStatus="unlisted" \
  && echo DONE
