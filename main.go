package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jetsetilly/supercharge/supercharge"
)

// the context type defines the command line parameters for the program and is
// also a valid io.Writer, suitable for verbose logging
type context struct {
	verbose   bool
	overwrite bool
}

func (ctx context) Write(p []byte) (n int, err error) {
	if ctx.verbose {
		os.Stdout.Write(p)
		return len(p), nil
	}
	return 0, nil
}

func main() {
	var ctx context

	// parse command line arguments
	flag.BoolVar(&ctx.verbose, "v", false, "verbose messages")
	flag.BoolVar(&ctx.overwrite, "o", false, "overwrite existing wav files")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [ROM files]\n\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Println("\nconverted WAV files will be saved in the same directory as the ROM file")
	}
	flag.Parse()

	// display usage if no rom files have been specified
	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	// process all files specified on the command line
	for _, f := range flag.Args() {
		f = filepath.Clean(f)
		err := process(ctx, f)
		if err != nil {
			ctx.Write([]byte(fmt.Sprintf("%s\n", err.Error())))
		}
	}
}

func process(ctx context, romFile string) error {
	// create filename for wav file
	wavFile, _ := strings.CutSuffix(romFile, filepath.Ext(romFile))
	wavFile = fmt.Sprintf("%s.wav", wavFile)

	// check whether wav file already exists
	if !ctx.overwrite {
		_, err := os.Stat(wavFile)
		if err == nil || !os.IsNotExist(err) {
			return fmt.Errorf("%s already exists", filepath.Base(wavFile))
		}
	}

	// open rom file and read the data in its entirety
	r, err := os.Open(romFile)
	if err != nil {
		return fmt.Errorf("%s: %w", filepath.Base(romFile), err)
	}
	defer r.Close()

	rom, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("%s: %w", filepath.Base(romFile), err)
	}

	// validate with the supercharge package that this rom data is okay
	err = supercharge.Validate(rom)
	if err != nil {
		return fmt.Errorf("%s skipped", filepath.Base(romFile))
	}

	// create wav file
	w, err := os.Create(wavFile)
	if err != nil {
		return fmt.Errorf("%s: %w", filepath.Base(romFile), err)
	}
	defer w.Close()

	// convert rom data to wav file
	var results bytes.Buffer
	err = supercharge.Convert(rom, w, &results)
	if err != nil {
		return fmt.Errorf("%s: %w", filepath.Base(romFile), err)
	}

	// display results
	ctx.Write([]byte(fmt.Sprintf("%s converted\n", filepath.Base(romFile))))
	ctx.Write(results.Bytes())

	return nil
}
