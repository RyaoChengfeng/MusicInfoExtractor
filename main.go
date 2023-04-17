package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wtolson/go-taglib"
)

type MusicFile struct {
	Title      string `json:"title,omitempty"`
	Artist     string `json:"artist,omitempty"`
	Album      Album  `json:"album,omitempty"`
	Genre      string `json:"genre,omitempty"`
	Year       int    `json:"year,omitempty"`
	Duration   string `json:"duration,omitempty"`
	Bitrate    string `json:"bitrate,omitempty"`
	Samplerate string `json:"samplerate,omitempty"`
	//FilePath   string `json:"file_path"`
}

type Album struct {
	Album string `json:"album,omitempty"`
	Track int    `json:"track,omitempty"`
}

func main() {
	args := os.Args[1:]
	var folderPaths []string
	var outputPath string

	for i, arg := range args {
		if arg == "-o" {
			if i+1 < len(args) {
				outputPath = args[i+1]
				break
			}
		} else {
			if arg == "." {
				wd, err := os.Getwd()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting current working directory: %v\n", err)
					os.Exit(1)
				}
				arg = wd
			}
			folderPaths = append(folderPaths, arg)
			//fmt.Printf("folderPaths: %v", folderPaths)
		}
	}

	//os.Exit(1)
	if len(folderPaths) == 0 {
		fmt.Println("Usage: ./main [folder_path1] [folder_path2] ... -o [output_path]")
		os.Exit(1)
	}
	if outputPath == "" {
		outputPath = "music.json"
		//fmt.Println("Usage: ./main [folder_path1] [folder_path2] ... -o [output_path]")
		//os.Exit(1)
	}

	// 处理文件夹目录
	output := make(map[string]interface{})
	sort.Strings(folderPaths)
	processSubFolders(output, folderPaths)

	// 将输出写入文件
	if outputPath != "" {
		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding output as JSON: %v\n", err)
			os.Exit(1)
		}
		err = ioutil.WriteFile(outputPath, data, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output to file %s: %v\n", outputPath, err)
			os.Exit(1)
		}
	}
	//else {
	//	data, err := json.MarshalIndent(output, "", "  ")
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "Error encoding output as JSON: %v\n", err)
	//		os.Exit(1)
	//	}
	//	fmt.Println(string(data))
	//}
}

// 处理子目录
func processSubFolders(parentOutput map[string]interface{}, subFolders []string) {
	if len(subFolders) == 0 {
		return
	}
	// 处理每个子目录
	for _, folderPath := range subFolders {
		// 获取目录信息
		dirInfo, err := os.Stat(folderPath)
		if err != nil {
			fmt.Printf("Error accessing folder %s: %v\n", folderPath, err)
			continue
		}

		// 跳过非目录文件和.开头的隐藏文件夹
		if !dirInfo.IsDir() || (strings.HasPrefix(dirInfo.Name(), ".") && dirInfo.Name() != ".") {
			continue
		}

		// 提取目录中的音乐文件信息
		musicFiles := map[string]MusicFile{}
		err = filepath.Walk(folderPath, func(filePath string, fileInfo os.FileInfo, err error) error {
			// 跳过非音乐文件和.开头的隐藏文件
			if fileInfo.IsDir() || strings.HasPrefix(fileInfo.Name(), ".") || !isMusicFile(fileInfo.Name()) {
				return nil
			}

			// 解析音乐文件
			tag, err := taglib.Read(filePath)
			if err != nil {
				fmt.Printf("Error reading tag from file %s: %v\n", filePath, err)
				return nil
			}
			defer tag.Close()

			// 创建音乐文件结构体
			musicFile := MusicFile{
				Title:      tag.Title(),
				Artist:     tag.Artist(),
				Album:      Album{Album: tag.Album(), Track: tag.Track()},
				Genre:      tag.Genre(),
				Year:       tag.Year(),
				Duration:   formatDuration(tag.Length().Seconds()),
				Bitrate:    fmt.Sprintf("%d kbps", tag.Bitrate()/1000),
				Samplerate: fmt.Sprintf("%d Hz", tag.Samplerate()),
				//FilePath:   filePath,
			}
			musicFiles[fileInfo.Name()] = musicFile

			return nil
		})
		if err != nil {
			fmt.Printf("Error accessing folder %s: %v\n", folderPath, err)
			continue
		}

		// 如果该目录中有音乐文件，则将
		dirName := filepath.Base(folderPath)
		parentOutput[dirName] = map[string]interface{}{
			"music": musicFiles,
		}

		// 递归处理子目录
		subFolders := []string{}
		fileInfos, err := ioutil.ReadDir(folderPath)
		if err != nil {
			fmt.Printf("Error reading subfolders of folder %s: %v\n", folderPath, err)
		} else {
			for _, fileInfo := range fileInfos {
				if fileInfo.IsDir() && !strings.HasPrefix(fileInfo.Name(), ".") {
					subFolders = append(subFolders, filepath.Join(folderPath, fileInfo.Name()))
				}
			}
		}
		processSubFolders(parentOutput[dirName].(map[string]interface{}), subFolders)
	}
}

// 判断文件是否是音乐文件
func isMusicFile(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".mp3" || ext == ".flac" || ext == ".ogg" || ext == ".wav" || ext == ".m4a" || ext == ".wma" || ext == ".aac" || ext == ".aiff"
}

func formatDuration(seconds float64) string {
	minutes := int(seconds) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%d:%02d", minutes, secs)
}
