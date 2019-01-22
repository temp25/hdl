package main

import (
    "flag"
    "fmt"
    "os"
    "strings"
    //"github.com/temp25/hotstar-dl/helper"
    "github.com/temp25/hotstar-dl/urlretriever"
    //"text/tabwriter"
    "github.com/temp25/hotstar-dl/videoutil"

    //"bytes"
    //"io"
    //"os/exec"
    //"log"
    //"reflect"
    //"strconv"
    //"time"
)

var helpFlagDesc = "Prints this help and exit"
var listFormatsFlagDesc = "List available video formats for given url"
var formatFlagDesc = "Video format to download video in specified resolution"
var ffmpegPathFlagDesc = "Location of the ffmpeg binary(absolute path)"
var metadataFlagDesc = "Add metadata to the video file"
var outputFileNameFlagDesc = "Output file name"

// note, that variables are pointers
var helpFlag = flag.Bool("help", false, helpFlagDesc)
var listFormatsFlag = flag.Bool("list", false, listFormatsFlagDesc)
var formatFlag = flag.String("format", "", formatFlagDesc)
var ffmpegPathFlag = flag.String("ffmpeg-location", "", ffmpegPathFlagDesc)
var metadataFlag = flag.Bool("add-metadata", false, metadataFlagDesc)
var outputFileNameFlag = flag.String("output", "", outputFileNameFlagDesc)

func init() {
    
    //shorthand notations
    flag.BoolVar(helpFlag, "h", false, helpFlagDesc)
    flag.BoolVar(listFormatsFlag, "l", false, listFormatsFlagDesc)
    flag.StringVar(formatFlag, "f", "", formatFlagDesc)
    flag.BoolVar(metadataFlag, "m", false, metadataFlagDesc)
    flag.StringVar(outputFileNameFlag, "o", "", outputFileNameFlagDesc)

    //custom flag usage
    flag.Usage = func() {
        fmt.Fprintf(os.Stdout, "Usage: %s [OPTIONS] URL\n\n", os.Args[0])
        fmt.Println("Options:")
        fmt.Fprintf(os.Stdout, "-h, --help\t\t%s\n", helpFlagDesc)
        fmt.Fprintf(os.Stdout, "-l, --list\t\t%s\n", listFormatsFlagDesc)
        fmt.Fprintf(os.Stdout, "-f, --format\t\t%s\n", formatFlagDesc)
        fmt.Fprintf(os.Stdout, "--ffmpeg-location\t%s\n", ffmpegPathFlagDesc)
        fmt.Fprintf(os.Stdout, "-m, --add-metadata\t%s\n", metadataFlagDesc)
        fmt.Fprintf(os.Stdout, "-o, --output\t\t%s\n", outputFileNameFlagDesc)
        os.Exit(0)
        //flag.PrintDefaults()
    }

}

/*
func padZeroRight(num int64) int64 {
    tmp := fmt.Sprintf("%-13d", num)
	tmp = strings.Replace(tmp, " ", "0", -1)
	paddedNum, err := strconv.ParseInt(tmp, 10, 64)
    if err != nil {
        panic(err)
    }
    return paddedNum
}

func CountDigits(i int64) (count int64) {
    for i != 0 {
        i /= 10
        count = count + 1
    }
    return count
}

func getDateStr(t int64) string {
    location, err := time.LoadLocation("Asia/Kolkata")
    if err != nil {
        panic(err)
    }
    return time.Unix(0, t * int64(time.Millisecond)).In(location).String()
}
*/

func main() {
    flag.Parse()

    flagCount := len(flag.Args())
    //fmt.Println("flag_args :", flag.Args())
    //fmt.Println("flagLen :", flagLen)

    /*for index, flag_args := range flag.Args() {
        fmt.Println("flag_args :", flag_args)
    }*/

    /*fmt.Println("helpFlag :",*helpFlag)
    fmt.Println("listFormatsFlag :",*listFormatsFlag)
    fmt.Println("formatFlag :",*formatFlag)
    fmt.Println("ffmpegPathFlag :",*ffmpegPathFlag)
    fmt.Println("metadataFlag :",*metadataFlag)*/

    
    if *helpFlag {
        flag.Usage()
    } else if flagCount == 0 {
        fmt.Println("Must provide atleast one url at end")
        flag.Usage()
        os.Exit(-1)
    } else if flagCount > 1 {
        fmt.Println("Url must be provided at end before options")
        flag.Usage()
        os.Exit(-1)
    } else if !urlretriever.IsValidHotstarUrl(flag.Args()[0]) {
        fmt.Println("Invalid hotstar url")
        os.Exit(-1)
    } else if *listFormatsFlag {
        videoutil.ListFormatsOrDownloadVideo(false, flag.Args()[0], *formatFlag, *ffmpegPathFlag, *outputFileNameFlag, *metadataFlag)
    } else if *formatFlag != "" {

        if !strings.HasPrefix(*formatFlag, "hls-") {
            fmt.Println("Invalid format specified")
            os.Exit(-1)
        } else {
            videoutil.ListFormatsOrDownloadVideo(true, flag.Args()[0], *formatFlag, *ffmpegPathFlag, *outputFileNameFlag, *metadataFlag)
        }

    } else {
       fmt.Println("Invalid args specified")
       flag.Usage()
    }
    
    
}

/*
func listFormatsOrDownloadVideo(isOnlyDownload bool) {
    
    videoUrl := flag.Args()[0]

    videoUrlPageContents := helper.GetPageContents(videoUrl, false)

    playbackUri, metadata := urlretriever.GetPlaybackUri(videoUrlPageContents, videoUrl)

    if metadata["drmProtected"] == "true" {
        fmt.Println("Error: The video is DRM Protected")
        os.Exit(-1)
    }

    fmt.Println("\n\n")

    playbackUriPageContents := helper.GetPageContents(playbackUri, true)

    masterPlaybackUrl := urlretriever.GetMasterPlaybackUrl(playbackUriPageContents)

    masterPlaybackPageContents := helper.GetPageContents(masterPlaybackUrl, true)

    videoFormats := urlretriever.GetVideoFormats(masterPlaybackPageContents, masterPlaybackUrl)

    if !isOnlyDownload {
        //NewWriter(io.Writer, minWidth, tabWidth, padding, padchar, flags)
        tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0) //tabwriter.Debug

        fmt.Fprintln(tw, "format code\textension\tresolution\tbandwidth\tcodec & frame rate\t")

        for k, v := range videoFormats {
            v := v.(map[string]interface{})
            if v["FRAME-RATE"] == nil {
                fmt.Fprintf( tw, "%s\tmp4\t%s\t%s\t%s\n", k, v["RESOLUTION"].(string), v["K-FORM"].(string), v["CODECS"].(string) )
            } else {
                fmt.Fprintf( tw, "%s\tmp4\t%s\t%s\t%s  %s fps\n", k, v["RESOLUTION"].(string), v["K-FORM"].(string),  v["CODECS"].(string), v["FRAME-RATE"].(string) )
            }
            
        }
        tw.Flush()
        os.Exit(0)
    } else {
        format := *formatFlag

        if videoFormatInfo, availablity := videoFormats[format].(map[string]interface{}); availablity {

            streamUrl := videoFormatInfo["STREAM-URL"].(string)
            //fmt.Println("Urrl to be streamed is %s", streamUrl)
            ffmpegLocation := *ffmpegPathFlag
            if ffmpegLocation == "" {
                fmt.Println("ffmpeg location not specified in args. Downloading binaries....")
                ffmpegLocation = helper.DownloadFFBinaries()
            }

            outputFileName := *outputFileNameFlag
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
            
            if *metadataFlag {
                for k2, v2 := range metadata {
                     metaArgs = append(metaArgs, "-metadata")
                     meta_data := fmt.Sprintf("%s=\"%s\"", k2, v2)
                     metaArgs = append(metaArgs, meta_data)
                }
            }

            metaArgs = append(metaArgs, "-c")
            metaArgs = append(metaArgs, "copy")
            metaArgs = append(metaArgs, "-y")
            metaArgs = append(metaArgs, outputFileName)
            
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
*/


