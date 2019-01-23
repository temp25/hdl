package main

import (
	"flag"
	"fmt"
	"github.com/temp25/hdl/urlretriever"
	"github.com/temp25/hdl/videoutil"
	"log"
	"net/url"
	"os"
	"strings"
)

//flag descriptions
var helpFlagDesc = "Prints this help and exit"
var listFormatsFlagDesc = "List available video formats for given url"
var formatFlagDesc = "Video format to download video in specified resolution"
var ffmpegPathFlagDesc = "Location of the ffmpeg binary(absolute path)"
var metadataFlagDesc = "Add metadata to the video file"
var outputFileNameFlagDesc = "Output file name"

//flag declarations
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

func main() {
	flag.Parse()
	flagCount := len(flag.Args())
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
	} else if videoUrl := flag.Args()[0]; videoUrl != "" {
		parsedUrl, err := url.Parse(videoUrl)
		if err != nil {
			log.Fatal(err)
		}
		switch parsedUrl.Scheme {
		case "":
			fmt.Println("Replacing empty url scheme with https")
			parsedUrl.Scheme = "https"
		case "https":
			//do nothing
		case "http":
			fmt.Println("Replacing http url scheme with https")
			parsedUrl.Scheme = "https"
		default:
			fmt.Println("Invalid url scheme please enter valid one")
			os.Exit(-1)
		}

		videoUrl = fmt.Sprintf("%v", parsedUrl)

		fmt.Println("Parsed video url is", parsedUrl)

		isValidUrl, videoId := urlretriever.IsValidHotstarUrl(videoUrl)
		if isValidUrl {
			if *listFormatsFlag {
				videoutil.ListFormatsOrDownloadVideo(false, videoUrl, videoId, *formatFlag, *ffmpegPathFlag, *outputFileNameFlag, *metadataFlag)
			} else if *formatFlag != "" {
				if !strings.HasPrefix(*formatFlag, "hls-") {
					fmt.Println("Invalid format specified")
					os.Exit(-1)
				} else {
					videoutil.ListFormatsOrDownloadVideo(true, videoUrl, videoId, *formatFlag, *ffmpegPathFlag, *outputFileNameFlag, *metadataFlag)
				}
			} else {
				//Check for other flags if associated with url if any
			}
		} else {
			fmt.Println("Invalid hotstar url. Please enter a valid one")
			os.Exit(-1)
		}

	} else {
		fmt.Println("Invalid args specified")
		flag.Usage()
	}
}
