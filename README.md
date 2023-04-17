# MusicInfoExtractor

提取音乐文件的标签信息，并将它们输出为JSON文件

## Usage
```bash
go build main.go
```
```bash
./main [folder_path1] [folder_path2] ... -o [output_path]
```
Or just `./main .`

输出的格式为

```json
{
    "folder_path1":{
        "music":{
            "V.A. - Last Swan.mp3": {
                "title": "Last Swan",
                "artist": "V.A.",
                "album": {
                  "album": "Swan Song OST",
                  "track": 1
                },
                "year": "2005",
                "duration": "5:59",
                "bitrate": "128 kbps",
                "samplerate": "44100 Hz",
              }
        },
        "folder":{
            "subfolder_path1":{
                "music":[],
                "folder":{
                    ...
                }
            },
            "subfolder_path2":{
                ...
            }
        }
    }
}
```

