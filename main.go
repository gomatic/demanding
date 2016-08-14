package main

import (
	_ "expvar"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/elazarl/goproxy"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

//
const MAJOR = "0.1"

// DO NOT UPDATE. This is populated by the build. See the Makefile.
var VERSION = "0"

//
var settings struct {
	Debugger bool
	Verbose  bool

	Container bool
	Image     bool
}

//
func main() {
	app := cli.NewApp()
	app.Name = "demanding"
	app.Usage = "Demanding."
	app.Version = MAJOR + "." + VERSION
	app.EnableBashCompletion = true

	destination()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "container, c",
			Usage:       "Try restarting a container.",
			Destination: &settings.Container,
		},
		cli.BoolFlag{
			Name:        "image, i",
			Usage:       "Try launching an image.",
			Destination: &settings.Image,
		},
	}

	app.Action = func(ctx *cli.Context) error {
		from := ":1080"
		switch args := ctx.Args(); len(args) {
		default:
			fallthrough
		case 2:
			from = args[0]
		case 1:
		case 0:
			return cli.NewExitError(fmt.Sprintf("usage: %s OPTIONS image|container [[address]:port]", filepath.Base(os.Args[0])), 1)
		}
		if err := proxy(from); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}

	app.Run(os.Args)
}

//
func destination() ([]string, []string) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	if cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders); err != nil {
		panic(err)
	} else {
		return images(cli), containers(cli)
	}
}

//
func images(cli *client.Client) []string {
	options := types.ImageListOptions{All: true}
	if is, err := cli.ImageList(context.Background(), options); err != nil {
		panic(err)
	} else {
		var images sort.StringSlice
		for _, i := range is {
			if i.ParentID == "" || i.RepoTags[0] != "<none>:<none>" {
				images = append(images, i.ID[7:])
				fmt.Println(i.ID[7:14], i.RepoTags[0])
			}
		}
		return images
	}
}

//
func containers(cli *client.Client) []string {
	options := types.ContainerListOptions{All: true}
	if cs, err := cli.ContainerList(context.Background(), options); err != nil {
		panic(err)
	} else {
		var containers sort.StringSlice
		for _, c := range cs {
			containers = append(containers, c.Names[0])
			fmt.Println(c.Names[0])
		}
		return containers
	}
}

//
func proxy(from string) error {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return &http.Request{
				Method: "CONNECT",
				URL: &url.URL{
					Host: "",
				},
				Header: make(http.Header),
			}, nil
		})

	return http.ListenAndServe(from, proxy)
}
