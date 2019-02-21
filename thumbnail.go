/*
 * siteshot is a website screenshot-capturing web application.
 * Copyright © 2016-2019 A Bunch Tell LLC.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// configuration constants
const (
	port = 3333
)

var (
	infoLog *log.Logger
	errLog  *log.Logger
)

var (
	wd string

	// Thumbnail parameters
	thumbWidth  int
	thumbHeight int
)

func main() {
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	errLog = log.New(os.Stderr, fmt.Sprintf("[%s] ", red("ERROR")), log.Ldate|log.Ltime|log.Lshortfile)
	infoLog = log.New(os.Stdout, fmt.Sprintf("[%s] ", blue("INFO")), log.Ldate|log.Ltime)

	// Get information about the environment
	var err error
	wd, err = os.Getwd()
	if err != nil {
		errLog.Printf("Couldn't get working dir: %v", err)
		return
	}

	// Set configuration
	dim := flag.String("dim", "", "Dimensions of final thumbnail")
	flag.Parse()
	if *dim != "" {
		c := strings.Split(*dim, "x")
		if len(c) != 2 {
			errLog.Printf("usage: siteshot --dim 320x240")
			return
		}
		thumbWidth, err = strconv.Atoi(c[0])
		if err != nil {
			errLog.Printf("usage: siteshot --dim 320x240")
			return
		}
		thumbHeight, err = strconv.Atoi(c[1])
		if err != nil {
			errLog.Printf("usage: siteshot --dim 320x240")
			return
		}
	} else {
		thumbWidth = 320
		thumbHeight = 240
	}

	// Start server
	infoLog.Printf("Listening on :%d", port)
	infoLog.Printf("Width: %d Height: %d", thumbWidth, thumbHeight)
	http.HandleFunc("/", makeThumbnail)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func makeThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	infoLog.Printf("Received request for %s", url)

	thumbFile := strings.Replace(url[strings.Index(url, "://")+3:]+".png", "/", ".", -1)

	// Fetch the thumbnail
	thumb := exec.Command("xvfb-run", "--server-args", "-screen 0 1366x768x24", "webkit2png", "-o", thumbFile, url)
	if err := thumb.Run(); err != nil {
		errLog.Printf("xvfb-run: %v", err)
		return
	}

	// Resize to width
	cmd := exec.Command("convert", filepath.Join(wd, thumbFile), "-define", "png:big-depth=16", "-define", "png:color-type=6", "-thumbnail", fmt.Sprintf("%d", thumbWidth), filepath.Join(wd, thumbFile))
	if err := cmd.Run(); err != nil {
		errLog.Printf("convert -thumbnail: %v", err)
		return
	}

	// Crop to height
	cmd = exec.Command("convert", filepath.Join(wd, thumbFile), "-define", "png:big-depth=16", "-define", "png:color-type=6", "-crop", fmt.Sprintf("%dx%d+0+0", thumbWidth, thumbHeight), filepath.Join(wd, thumbFile))
	if err := cmd.Run(); err != nil {
		errLog.Printf("convert -crop: %v", err)
		return
	}

	infoLog.Printf(color.GreenString("✓") + " Created " + thumbFile)

	fmt.Fprintf(w, thumbFile)
}
