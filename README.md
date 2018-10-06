# Tool to add creation time to video & image filenames

## Prerequisite
* ffprobe  (brew install ffmpeg)
* exiftool (brew install exiftool)

## Install
go install .

## Example
```
$ rename-media -f /Volumes/Video/SomeVideos/
Renaming: /Volumes/Video/SomeVideos/00000.MTS to /Volumes/Video/SomeVideos/2014-03-05_10.09_00000.MTS
Renaming: /Volumes/Video/SomeVideos/00001.MTS to /Volumes/Video/SomeVideos/2014-03-05_10.12_00001.MTS
Renaming: /Volumes/Video/SomeVideos/00002.MTS to /Volumes/Video/SomeVideos/2014-03-05_21.08_00002.MTS
Renaming: /Volumes/Video/SomeVideos/C0001.MP4 to /Volumes/Video/SomeVideos/2018-03-14_14.11_C0001.MP4
```
