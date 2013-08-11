package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gedex/imgdownloader/provider"
)

var (
	tag            = flag.String("tag", "", "Image tag")
	n              = flag.Uint("n", 10, "Number of images to download")
	providerString = flag.String("from", "flickr", "Image provider (flickr or instagram)")
	out            = flag.String("out", "./", "Path to downloaded images")

	// Active config is stored here
	activeConf config
)

type config map[string]string

func main() {
	flag.Parse()

	err := checkParams()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v.\n", err)

		printUsage(1)
	}

	activeConf, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v.\n", err)
		os.Exit(1)
	}

	p, err := provider.Get(*providerString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v.\n", err)
		os.Exit(1)
	}
	p.Configure(activeConf)

	listToDownload, err := p.Request(*tag, *n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v.\n", err)
		os.Exit(1)
	}

	downloads := make(chan string, len(listToDownload))
	for filename, link := range listToDownload {
		go download(link, filename, downloads)
	}

	for d := range downloads {
		fmt.Println(d)
	}
	close(downloads)
}

func checkParams() error {
	if *tag == "" {
		return fmt.Errorf("Tag is not supplied")
	}
	return nil
}

func printUsage(code int) {
	format := fmt.Sprintf("Usage of %s:\n", os.Args[0])
	if code > 0 {
		fmt.Fprintf(os.Stderr, format)
	} else {
		fmt.Fprintf(os.Stdout, format)
	}

	flag.PrintDefaults()
	os.Exit(code)
}

func getConfig() (config, error) {
	file, err := getFileConfig("./imgdownloader.json")
	if err != nil {
		file, err = getFileConfig("~/imgdownloader.json")
		return nil, err
	}

	c, err2 := readFileConfig(file)
	if err2 != nil {
		return nil, err2
	}
	return c, nil
}

func getFileConfig(path string) (file *os.File, err error) {
	file, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	return
}

func readFileConfig(file *os.File) (config, error) {
	defer file.Close()

	var data []byte

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		data = buf[:n]
	}

	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	vv := v.(map[string]interface{})
	if _, ok := vv[*providerString]; !ok {
		return nil, fmt.Errorf("%v provider is not found in imgdownaloder.json", *providerString)
	}
	vvv := vv[*providerString].(map[string]interface{})

	c := config{}
	for key, val := range vvv {
		vstr, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("invalid %v provider config", *providerString)
		}
		c[key] = vstr
	}
	return c, nil
}

func download(link, filename string, downloads chan<- string) {
	resp, err := http.Get(link)
	defer resp.Body.Close()
	if err != nil {
		downloads <- fmt.Sprintf("Error GET %v: %v", link, err)
		return
	}

	outFile := path.Join(*out, filename)
	out, err := os.Create(outFile)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		downloads <- fmt.Sprintf("Error copying %v to %v", filename, outFile)
		return
	}
	downloads <- fmt.Sprintf("Successfully downloading %v to %v", link, outFile)
}
