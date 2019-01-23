package videoutil

import (
   "github.com/temp25/hdl/helper"
   "github.com/temp25/hdl/urlretriever"
   "fmt"
   "text/tabwriter"
   "io"
   "log"
   "bytes"
   "os"
   "os/exec"
)

func ListFormatsOrDownloadVideo(isOnlyDownload bool, videoUrl string, videoId string, format string, ffmpegLocation string, outputFileName string, metaDataFlag bool) {

    videoUrlPageContents := helper.GetPageContents(videoUrl, false)

    playbackUri, metadata := urlretriever.GetPlaybackUri(videoUrlPageContents, videoUrl, videoId)
	
    if metadata["drmProtected"] == "true" {
        fmt.Println("Error: The video is DRM Protected")
        os.Exit(-1)
    }
	
    playbackUriPageContents := helper.GetPageContents(playbackUri, true)

    masterPlaybackUrl := urlretriever.GetMasterPlaybackUrl(playbackUriPageContents)

    masterPlaybackPageContents := helper.GetPageContents(masterPlaybackUrl, true)
	
    videoFormats, keys := urlretriever.GetVideoFormats(masterPlaybackPageContents, masterPlaybackUrl)

    if !isOnlyDownload {
        //NewWriter(io.Writer, minWidth, tabWidth, padding, padchar, flags)
        tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0) //tabwriter.Debug

        fmt.Fprintln(tw, "format code\textension\tresolution\tbandwidth\tcodec & frame rate\t")

        for _, key := range keys {
            k := fmt.Sprintf("hls-%d", key)
            v := videoFormats[k].(map[string]interface{})
            if v["FRAME-RATE"] == nil {
                fmt.Fprintf( tw, "%s\tmp4\t%s\t%s\t%s\n", k, v["RESOLUTION"].(string), v["K-FORM"].(string), v["CODECS"].(string) )
            } else {
                fmt.Fprintf( tw, "%s\tmp4\t%s\t%s\t%s  %s fps\n", k, v["RESOLUTION"].(string), v["K-FORM"].(string),  v["CODECS"].(string), v["FRAME-RATE"].(string) )
            }
            
        }
        tw.Flush()
        os.Exit(0)
    } else {

        if videoFormatInfo, availablity := videoFormats[format].(map[string]interface{}); availablity {

            streamUrl := videoFormatInfo["STREAM-URL"].(string)
            
            if ffmpegLocation == "" {
                fmt.Println("ffmpeg location not specified in args. Downloading binaries....")
                ffmpegLocation = helper.DownloadFFBinaries()
            } else {
                //check if valid path to binary
                ffmpegFilePath, err := os.Stat(ffmpegLocation)
                if !ffmpegFilePath.Mode().IsDir() {
                    fmt.Println("Path provided to ffmpeg is not a filepath but directory")
                    os.Exit(-1)
                }
                if os.IsNotExist(err) {
                    fmt.Println("File doesn't exists in path provided to ffmpeg")
                    os.Exit(-1)
                }
            }

            if outputFileName == "" {
                outputFileName = "Output.mp4"
            }

            var stdoutBuf, stderrBuf bytes.Buffer

            if err := os.Chmod(ffmpegLocation, 0555); err != nil {
              log.Fatal(err)
            }

            metaArgs := []string{}
            metaArgs = append(metaArgs, "-i")
            metaArgs = append(metaArgs, streamUrl)
            
            if metaDataFlag {
                for metaDataName, metaDataValue := range metadata {
                     metaArgs = append(metaArgs, "-metadata")
                     meta_data := fmt.Sprintf("%s=\"%s\"", metaDataName, metaDataValue)
                     metaArgs = append(metaArgs, meta_data)
                }
            } else {
                fmt.Println("Skipping adding metedata for video file")
            }

            metaArgs = append(metaArgs, "-c")
            metaArgs = append(metaArgs, "copy")
            metaArgs = append(metaArgs, "-y")
            metaArgs = append(metaArgs, outputFileName)
            
            fmt.Println("Starting ffmpeg to download video...")
             
            cmd := exec.Command(ffmpegLocation, metaArgs...)

            stdoutIn, _ := cmd.StdoutPipe()
            stderrIn, _ := cmd.StderrPipe()
        
            var errStdout, errStderr error
            stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
            stderr := io.MultiWriter(os.Stderr, &stderrBuf)
            err := cmd.Start()
            if err != nil {
                log.Fatalf("cmd.Start() failed with '%s'\n", err)
            }
        
            go func() {
                _, errStdout = io.Copy(stdout, stdoutIn)
            }()
        
            go func() {
                _, errStderr = io.Copy(stderr, stderrIn)
            }()
        
            err = cmd.Wait()
            if err != nil {
                log.Fatalf("cmd.Run() failed with %s\n", err)
            }
            if errStdout != nil || errStderr != nil {
                log.Fatal("failed to capture stdout or stderr\n")
            }

        } else {
            fmt.Printf("The specified video format %s is not available. Specify existing format from the list", format)
            os.Exit(-2)
        }

    }

}
