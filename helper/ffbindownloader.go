package helper

import (
    "fmt"
    "os"
    "io"
    "runtime"
    "bytes"
    "strings"
    "os/exec"
	"log"
	"github.com/mholt/archiver"
	"path/filepath"
)

var staticBuildMap = map[string]string {
	"android,arm": "https://github.com/WritingMinds/ffmpeg-android/releases/download/v0.3.4/prebuilt-binaries.zip",
	"android": "https://github.com/WritingMinds/ffmpeg-android/releases/download/v0.3.4/prebuilt-binaries.zip",
	"windows,amd64": "https://ffmpeg.zeranoe.com/builds/win64/static/ffmpeg-20190116-282a471-win64-static.zip",
	"windows,386": "https://ffmpeg.zeranoe.com/builds/win32/static/ffmpeg-20190116-282a471-win32-static.zip",
	"windows": "https://ffmpeg.zeranoe.com/builds/win32/static/ffmpeg-20190116-282a471-win32-static.zip",
	"darwin,amd64": "https://ffmpeg.zeranoe.com/builds/macos64/static/ffmpeg-20190116-282a471-macos64-static.zip",
	"darwin,arm64": "https://ffmpeg.zeranoe.com/builds/macos64/static/ffmpeg-20190116-282a471-macos64-static.zip",
	"darwin": "https://ffmpeg.zeranoe.com/builds/macos64/static/ffmpeg-20190116-282a471-macos64-static.zip",
	"linux,amd64": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz",
	"linux,386": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-i686-static.tar.xz",
	"linux,arm64": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz",
	"linux": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-i686-static.tar.xz",
	"dragonfly": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-i686-static.tar.xz",
	"freebsd": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-i686-static.tar.xz",
	"netbsd": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-i686-static.tar.xz",
	"openbsd": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-i686-static.tar.xz",
	"plan9": "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-i686-static.tar.xz",
	
}

func fileCopy(src, dst string) (int64, error) {
        sourceFileStat, err := os.Stat(src)
        if err != nil {
                return 0, err
        }

        if !sourceFileStat.Mode().IsRegular() {
                return 0, fmt.Errorf("%s is not a regular file", src)
        }

        source, err := os.Open(src)
        if err != nil {
                return 0, err
        }
        defer source.Close()

        destination, err := os.Create(dst)
        if err != nil {
                return 0, err
        }
        defer destination.Close()
        nBytes, err := io.Copy(destination, source)
        return nBytes, err
}

func getStaticBuildUrl() string {
    
    osArch := fmt.Sprintf("%s,%s", runtime.GOOS, runtime.GOARCH)
    if val, ok := staticBuildMap[osArch]; ok {
        return val
    } else if val, ok := staticBuildMap[runtime.GOOS]; ok {
        return val
    }
    
    return ""
}

func DownloadFFBinaries() string {
   
	var stdoutBuf, stderrBuf bytes.Buffer
	
	var ffmpegBinLocation string
    
    staticBuildUrl := getStaticBuildUrl()
    staticBuildName := "binaries"
    currentDirectory, _ := os.Getwd()

    homeDir, err := HomeDir()
	if err!=nil {
        panic(err)
	}

    if _, err := os.Stat(filepath.Join(homeDir, "ffmpeg")); !os.IsNotExist(err) {
    	fmt.Println("ffmpeg binary present in user home directory from previous run. skipping download...")
    	ffmpegBinLocation = filepath.Join(homeDir, "ffmpeg")
    } else if _, err := os.Stat(filepath.Join(homeDir, "ffmpeg.exe")); !os.IsNotExist(err) {
    	fmt.Println("ffmpeg.exe executable present in user home directory from previous run. skipping download...")
    	ffmpegBinLocation = filepath.Join(homeDir, "ffmpeg.exe")
    } else {
    	if strings.HasSuffix(staticBuildUrl, "zip") {
	        staticBuildName = fmt.Sprintf("%s.zip", staticBuildName)
	    } else if strings.HasSuffix(staticBuildUrl, "tar") {
	        staticBuildName = fmt.Sprintf("%s.tar", staticBuildName)
	    } else if strings.HasSuffix(staticBuildUrl, "tar.xz") {
	        staticBuildName = fmt.Sprintf("%s.tar.xz", staticBuildName)
	    } else if strings.HasSuffix(staticBuildUrl, "tar.gz") {
	        staticBuildName = fmt.Sprintf("%s.tar.gz", staticBuildName)
	    } else if strings.HasSuffix(staticBuildUrl, "tar.bz2") {
	        staticBuildName = fmt.Sprintf("%s.tar.bz2", staticBuildName)
	    } else if strings.HasSuffix(staticBuildUrl, "tar.lz4") {
	        staticBuildName = fmt.Sprintf("%s.tar.lz4", staticBuildName)
	    } else if strings.HasSuffix(staticBuildUrl, "tar.sz") {
	        staticBuildName = fmt.Sprintf("%s.tar.sz", staticBuildName)
	    } else if strings.HasSuffix(staticBuildUrl, "rar") {
	        staticBuildName = fmt.Sprintf("%s.rar", staticBuildName)
	    }

	    
		cmd := exec.Command("wget", staticBuildUrl, "-O", staticBuildName)
		
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
		
		//remove existing directory
	    os.RemoveAll("extracted")

            fmt.Println("Extracting binaries to extracted folder");
	    
	    err = archiver.Unarchive(staticBuildName, "extracted")
	    
	    if err != nil {
	        panic(err)
		}

            fmt.Println("Extracted binaries to extracted folder");
		
		filepath.Walk(filepath.Join(currentDirectory, "extracted"), func(path string, f os.FileInfo, _ error) error {
			if !f.IsDir() {
				if f.Name() == "ffmpeg" || f.Name() == "ffmpeg.exe" {
				    
	            	bytesCopied, err := fileCopy(path, filepath.Join(homeDir, f.Name()))
	            	if err != nil {
	            	    panic(err)
	            	}
	            	fmt.Printf("\nffmpeg binary of %d bytes copied to user home %s\n\n", bytesCopied, homeDir)
	            	ffmpegBinLocation = filepath.Join(homeDir, f.Name())
	            	return io.EOF
				   
				}
			}
			return nil
		})
		
		//remove existing directory
		os.RemoveAll("extracted")
		//remove downloaded file
		os.Remove(staticBuildName)
    }
	
	return ffmpegBinLocation

}

