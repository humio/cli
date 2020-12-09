package shipper

import (
	"encoding/json"
	"fmt"
	"github.com/humio/cli/api"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type eventList struct {
	Type     string            `json:"type"`
	Fields   map[string]string `json:"fields"`
	Messages []string          `json:"messages"`
}

type ErrorBehaviour int

const (
	ErrorBehaviourDrop ErrorBehaviour = iota
	ErrorBehaviourPanic
)

type LineHandler interface {
	HandleLine(line string)
}

type LogShipper struct {
	APIClient           *api.Client
	URL                 string
	Fields              map[string]string
	ParserName          string
	MaxAttemptsPerBatch int
	ErrorBehaviour      ErrorBehaviour
	BatchSizeLines      int
	BatchSizeBytes      int
	BatchTimeout        time.Duration
	Logger              func(format string, v ...interface{})

	events          chan string
	finishedSending chan struct{}
}

func (s *LogShipper) HandleLine(line string) {
	s.events <- line
}

func (s *LogShipper) Finish() {
	close(s.events)
	<-s.finishedSending
}

func (s *LogShipper) Start() {
	s.events = make(chan string, s.BatchSizeLines)
	s.finishedSending = make(chan struct{})

	go func() {
		defer func() { close(s.finishedSending) }()
		var batch []string
		if s.BatchSizeLines != 0 {
			batch = make([]string, 0, s.BatchSizeLines)
		}

		for {
			bytes := 0

			e, more := <-s.events
			if !more {
				break
			}

			batch = append(batch, e)
			bytes += len(e)

			timeout := time.After(s.BatchTimeout)

			if s.BatchSizeBytes > 0 && bytes > s.BatchSizeBytes {
				goto send
			}

		loop:
			for {
				select {
				case e, more := <-s.events:
					if !more {
						break loop
					}
					batch = append(batch, e)
					bytes += len(e)
					if len(batch) >= s.BatchSizeLines || (s.BatchSizeBytes > 0 && bytes > s.BatchSizeBytes) {
						break loop
					}
				case <-timeout:
					break loop
				}
			}

		send:
			s.sendBatch(batch)

			batch = batch[:0]
			bytes = 0
		}
		if len(batch) > 0 {
			s.sendBatch(batch)
		}
	}()
}

func (s *LogShipper) sendBatch(messages []string) {
	ship := func() error {
		var eg errgroup.Group

		pr, pw := io.Pipe()

		jsonBody := []eventList{{
			Type:     s.ParserName,
			Fields:   s.Fields,
			Messages: messages,
		}}

		eg.Go(func() error {
			defer pw.Close()
			return json.NewEncoder(pw).Encode(jsonBody)
		})

		var resp *http.Response

		eg.Go(func() error {
			var err error
			resp, err = s.APIClient.HTTPRequest(http.MethodPost, s.URL, pr)
			return err
		})

		err := eg.Wait()

		if err != nil {
			return err
		}

		if resp.StatusCode > 400 {
			responseData, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				return fmt.Errorf("error reading http response body: %w", err)
			}

			return fmt.Errorf("bad response while sending events (status='%s'): %s", resp.Status, responseData)
		} else {
			// discard the response in order to re-use the connection
			_, _ = io.Copy(ioutil.Discard, resp.Body)
			_ = resp.Body.Close()
		}

		return nil
	}

	var err error
	for i := 0; i < s.MaxAttemptsPerBatch; i++ {
		if i > 0 {
			backOff := time.Duration(0.5*math.Pow(2, float64(i-1))*1000) * time.Millisecond
			if s.Logger != nil {
				s.Logger("Backoff for %v...", backOff)
			}
			time.Sleep(backOff)
		}
		err = ship()
		if err == nil {
			break
		}
		if s.Logger != nil {
			s.Logger("Error while sending logs to Humio. Retrying %d more times. Error message: %v", s.MaxAttemptsPerBatch-i-1, err)
		}
	}

	if err != nil {
		switch s.ErrorBehaviour {
		case ErrorBehaviourPanic:
			if s.Logger != nil {
				s.Logger("Error sending logs to Humio: %v", err)
			}
			panic(fmt.Sprintf("Error sending logs to Humio: %v", err))
		case ErrorBehaviourDrop:
			if s.Logger != nil {
				s.Logger("Error sending logs to Humio, dropping %d events: %v", len(messages), err)
			}
		}
	}
}
