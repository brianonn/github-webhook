package main

import (
	"fmt"
	"log"
    "os"
	"io"
	"bytes"
	"strings"
	"time"
	"os/exec"
	"io/ioutil"
	"net/http"
    "github.com/google/go-github/github"
)

const port = ":5000"
const watchedRef = "refs/heads/staging"
const script = "./push.sh"
var scriptArgs = []string{"push", "restart"}

var secret string = os.Getenv("GITHUB_WEBHOOK_SECRET")

// move all this semaphore stuff to a local library later
type empty struct{}
type semaphore chan empty

var jobSemaphore semaphore // global semaphore for job queue
const maxRequests = 1       // only handle one concurrent job at a time right now

func newSemaphore(count int) semaphore {
	return make(semaphore, count)
}

func (s semaphore) Aquire() {
    s <- empty{}
}

func (s semaphore) Release() {
    <- s
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {

	var (
		payload []byte
		err     error
	)

    // get a semaphore and schedule it's auto-release
    jobSemaphore.Aquire()
    defer jobSemaphore.Release()

    if (len(secret) > 0) {
        payload, err = github.ValidatePayload(r, []byte(secret))
        if err != nil {
            log.Printf("error validating secret: err=%s\n", err)
            return
        } else {
			log.Printf("secret header validates OK");
		}
    } else {
        payload, err = ioutil.ReadAll(r.Body)
        if err != nil {
            log.Printf("error reading request body: err=%s\n", err)
            return
        }
    }
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {

	case *github.PushEvent:
        log.Println("got a Push event")
		// this is a commit push, do something with it
        ref := e.GetRef()
        if (ref != "") { log.Printf("ref: %s\n", ref) }
        if (ref == watchedRef) {
            pusher := e.GetPusher().GetName()
            if (pusher != "" ) { log.Printf("pusher: %s\n", pusher) }

			deleted := e.GetDeleted()
			if (deleted) {
				// ignore push notifications to delete branches
				log.Println("ignoring deleted branch push")
				fmt.Fprintf(w, "ignoring deleted branch push")
				return
			}
			cmdName := script
			cmdArgs := scriptArgs
            cmdExec := exec.Command(cmdName, cmdArgs...)
			var outbuf, errbuf bytes.Buffer
			cmdExec.Stdout = &outbuf
			cmdExec.Stderr = &errbuf
            log.Printf("exec %s started", cmdName)
			err = cmdExec.Run()
			if (err != nil) {
                serr := fmt.Errorf("There was an error running the hook command: %s", err).Error()
				log.Println(serr)
				http.Error(w, serr, http.StatusInternalServerError)
			} else {
                log.Printf("exec %s finished", cmdName)
                soutput := strings.Replace(outbuf.String(), "\n", "\r\n", -1)
                serror := strings.Replace(errbuf.String(), "\n", "\r\n", -1)
                w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
                io.WriteString(w, soutput)
                io.WriteString(w, "===================\r\n")
                io.WriteString(w, serror)
            }
        }

	case *github.PullRequestEvent:
        log.Println("got a PullRequest  event")
		// this is a pull request, do something with it

	case *github.WatchEvent:
        log.Println("got a Watch event")
		// https://developer.github.com/v3/activity/events/types/#watchevent
		// someone starred our repository
		if e.Action != nil && *e.Action == "starred" {
			fmt.Printf("%s starred repository %s\n",
				*e.Sender.Login, *e.Repo.FullName)
		}

	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}

func main() {

    // initialize the semaphore
	jobSemaphore = newSemaphore(maxRequests)

	srv := &http.Server {
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 3 * 60 * time.Second, // 3 minutes for long running hooks
		Addr: port,
	}

	// routes
	http.HandleFunc("/webhook", handleWebhook)

    // start the server
    log.Println("server started on port", port)
	log.Fatal(srv.ListenAndServe())
}
