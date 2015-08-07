package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
)

var VERSION string

var shouldInit = kingpin.Flag("init", "Initialize a Bladefile").Bool()
var contentsTemplate = kingpin.Flag("template", "Contents json template to generate for.").Short('t').Default("Contents.json").String()
var interpolation = kingpin.Flag("interpolation", "Interpolation: l2, l3 (Lanczos 2 and 3), n (nearest neighbor), bc (bicubic), bl (bilinear), mn (Mitchell-Netravali)").Short('i').Default("l3").String()
var source = kingpin.Flag("source", "Source image to use. For optimal results supply a highest size PNG.").Short('s').String()
var out = kingpin.Flag("out", "Out folder. Use current folder if none given.").Short('o').String()
var includeContents = kingpin.Flag("catalog", "Include generation of a new image catalog Contents.json").Short('c').Bool()
var dryRun = kingpin.Flag("dryrun", "Do a dry run, don't write out anything").Short('d').Bool()
var mount = kingpin.Flag("mount", "Mount on an existing image catalog (take its Contents.json and output images to it)").Short('m').String()
var verbose = kingpin.Flag("verbose", "Verbose output").Bool()
var version = kingpin.Flag("version", "Current version").Short('v').Bool()

var BLADEFILE = `#
# Uncomment below to specify your own resources.
#

#blades:
#  - source: iTunesArtwork@2x.png
#    mount: out/iphone
#  - source: iTunesArtwork@2x.png
#    template: templates/watch.json
#    out: out/newwatch
#    contents: true
#  - source: iTunesArtwork@2x.png
#    template: templates/iphone-appicon.json
#    out: out/newiphone
#    contents: true
`

func init() {
	log.SetLevel(log.FatalLevel)
}
func main() {
	kingpin.Parse()

	if *version {
		println(VERSION)
		os.Exit(0)
	}

	if *shouldInit {
		ioutil.WriteFile("Bladefile", []byte(BLADEFILE), 0644)
		fmt.Println("Wrote Bladefile.")
		os.Exit(0)
	}

	if *verbose {
		log.SetLevel(log.InfoLevel)
	}
	bladefile := Bladefile{}
	if bladefile.Exists() {
		log.Infof("Found a local Bladefile.")
		err := bladefile.Load()
		if err != nil {
			log.Fatalf("Cannot load bladefile (%s)", err)
		}

		log.Infof("Bladefile contains %d blade defs.", len(bladefile.Blades))
		for _, blade := range bladefile.Blades {
			blade.Run()
		}
	} else {
		blade := Blade{
			Template:        *contentsTemplate,
			Interpolation:   *interpolation,
			Source:          *source,
			Out:             *out,
			IncludeContents: *includeContents,
			DryRun:          *dryRun,
			Mount:           *mount,
		}
		blade.Run()

	}
}
