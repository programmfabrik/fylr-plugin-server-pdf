package main

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/security"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/rpcc"
)

type pdfCreatorBody struct {
	Document   string              `json:"document"`
	FileName   string              `json:"file_name"`
	Properties page.PrintToPDFArgs `json:"properties"`
}

func createPdf(ectx context.Context, urlToRender string, port int, printToPDFArgs page.PrintToPDFArgs) ([]byte, error) {
	ctx, cancel := context.WithCancel(ectx)
	defer func() {
		// Ensure to all executions on context close
		cancel()
	}()

	devtoolUrl := url.URL{
		Scheme: "http",
		Host:   "localhost:" + strconv.Itoa(port),
	}

	// Use the DevTools API to manage targets
	pt, err := devtool.New(devtoolUrl.String()).Version(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Tried to connect to headless chrome. Could not create new devtool connection")
	}

	// Open a new RPC connection to the Chrome Debugging Protocol target
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return nil, errors.Wrap(err, "Tried to connect to headless chrome. Could not create DialContext")
	}
	defer conn.Close()

	// Create new browser context
	baseBrowser := cdp.NewClient(conn)
	err = baseBrowser.Security.SetIgnoreCertificateErrors(ctx, &security.SetIgnoreCertificateErrorsArgs{Ignore: true})
	if err != nil {
		return nil, errors.Wrap(err, "Could not set ignore certificate error to true")
	}
	newContextTarget, err := baseBrowser.Target.CreateBrowserContext(ctx, &target.CreateBrowserContextArgs{})
	if err != nil {
		panic(err)
	}
	defer baseBrowser.Target.DisposeBrowserContext(context.Background(), &target.DisposeBrowserContextArgs{BrowserContextID: newContextTarget.BrowserContextID})

	// Create a new blank target with the new browser context
	newTargetArgs := target.NewCreateTargetArgs("about:blank").SetBrowserContextID(newContextTarget.BrowserContextID)
	newTarget, err := baseBrowser.Target.CreateTarget(ctx, newTargetArgs)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open new blank target")
	}
	defer baseBrowser.Target.CloseTarget(context.Background(), &target.CloseTargetArgs{TargetID: newTarget.TargetID})

	// Connect the client to the new targetDialContext
	devtoolUrl.Scheme = "ws"
	devtoolUrl.Path = fmt.Sprintf("/devtools/page/%s", newTarget.TargetID)
	newContextConn, err := rpcc.DialContext(ctx, devtoolUrl.String(), rpcc.WithWriteBufferSize(104857586), rpcc.WithCompression())
	if err != nil {
		return nil, errors.Wrap(err, "Could not create dial context to new target")
	}
	defer newContextConn.Close()
	c := cdp.NewClient(newContextConn)

	// Close the target when done
	// (In development, skip this step to leave tabs open!)
	closeTargetArgs := target.NewCloseTargetArgs(newTarget.TargetID)
	defer func() { _, _ = baseBrowser.Target.CloseTarget(ctx, closeTargetArgs) }()

	// Enable the runtime
	err = c.Runtime.Enable(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Could not enable runtime")
	}

	// Enable the network
	err = c.Network.Enable(ctx, network.NewEnableArgs())
	if err != nil {
		return nil, errors.Wrap(err, "Could not enable network")
	}

	// Enable events on the Page domain
	err = c.Page.Enable(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Could not enable events")
	}

	// CSP bypass
	err = c.Page.SetBypassCSP(ctx, &page.SetBypassCSPArgs{Enabled: true})
	if err != nil {
		return nil, errors.Wrap(err, "Could not set bypass CSP to enabled")
	}

	// Listen for lifecycle events, which include loading external resources (CSS, JS, etc)
	if err := c.Page.SetLifecycleEventsEnabled(ctx, page.NewSetLifecycleEventsEnabledArgs(true)); err != nil {
		return nil, errors.Wrap(err, "Could not enable lifecycle events")
	}
	lifecycleEvent, err := c.Page.LifecycleEvent(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create lifecycle events listener")
	}
	defer lifecycleEvent.Close()

	// // Create a client to listen for the load event to be fired
	// loadEventFiredClient, err := c.Page.LoadEventFired(ctx)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Could not create load fire event listener")
	// }
	// defer loadEventFiredClient.Close()

	// Tell the page to navigate to the file

	navArgs := page.NewNavigateArgs(urlToRender)
	_, err = c.Page.Navigate(ctx, navArgs)
	if err != nil {
		return nil, errors.Wrap(err, "Could not navigate to page")
	}

	// Wait for lifecycle events to finish
	for {
		ev, err := lifecycleEvent.Recv()
		if err != nil {
			return nil, errors.Wrap(err, "Could not wait for finish lifecycle events")
		}
		if ev.Name == "networkIdle" {
			break
		}
	}

	// // Wait for the page to finish loading
	// _, err = loadEventFiredClient.Recv()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Could not wait for finish loading event")
	// }

	printDa, err := c.Page.PrintToPDF(ctx, &printToPDFArgs)
	if err != nil {
		return nil, errors.Wrap(err, "Could not print to pdf")
	}
	return printDa.Data, nil
}
