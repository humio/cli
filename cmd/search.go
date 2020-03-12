package cmd

import (
	"context"
	"fmt"
	"github.com/humio/cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

func newSearchCmd() *cobra.Command {
	var (
		start    string
		end      string
		live     bool
		complete bool
		fmt      string
	)

	cmd := &cobra.Command{
		Use:   "search <repo> <query>",
		Short: "Search",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repository := args[0]
			queryString := args[1]
			client := NewApiClient(cmd)

			if live && complete {
				cmd.Println("Cannot use both --live and --complete at the same time.")
				os.Exit(1)
			}

			ctx := contextCancelledOnInterrupt(context.Background())

			// run in lambda func to be able to defer and delete the query job
			err := func() error {
				id, err := client.QueryJobs().Create(repository, api.Query{
					QueryString: queryString,
					Start:       start,
					End:         end,
					Live:        live,
				})

				if err != nil {
					return err
				}

				defer func(id string) {
					// Humio will eventually delete the query when we stop polling and we can't do much about errors here.
					_ = client.QueryJobs().Delete(repository, id)
				}(id)

				var result api.QueryResult
				poller := queryJobPoller{
					queryJobs:  client.QueryJobs(),
					repository: repository,
					id:         id,
				}
				result, err = poller.WaitAndPollContext(ctx)

				if err != nil {
					return err
				}

				var printer interface {
					print(api.QueryResult)
				}

				if result.Metadata.IsAggregate {
					printer = newAggregatePrinter(cmd.OutOrStdout())
				} else {
					printer = newEventListPrinter(cmd.OutOrStdout(), fmt)
				}

				for {
					if !complete && len(result.Events) > 0 {
						printer.print(result)
					}

					if result.Done && !live {
						break
					}

					result, err = poller.WaitAndPollContext(ctx)
					if err != nil {
						return err
					}
				}

				if complete {
					printer.print(result)
				}

				return nil
			}()

			if err == context.Canceled {
				err = nil
			}

			exitOnError(cmd, err, "error running search")
		},
	}

	cmd.Flags().StringVar(&start, "start", "10m", "Query start time [default 10m]")
	cmd.Flags().StringVar(&end, "end", "", "Query end time")
	cmd.Flags().BoolVar(&live, "live", false, "Run a live search and keep outputting until interrupted.")
	cmd.Flags().BoolVar(&complete, "complete", false, "Wait for query to complete before printing result. Mostly useful for aggregates.")
	cmd.Flags().StringVarP(&fmt, "fmt", "f", "{@timestamp} {@rawstring}", "Format string if the result is an event list\n"+
		"Insert fields by wrapping field names in brackets, e.g. {@timestamp} [default: '{@timestamp} {@rawstring}']\n"+
		"Limited format modifiers are supported such as {@timestamp:40} which will right align and left pad @timestamp to 40 characters.\n"+
		"{@timestamp:-40} left aligns and right pads to 40 characters.")

	return cmd
}

func contextCancelledOnInterrupt(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigC
		cancel()
	}()

	return ctx
}

type queryJobPoller struct {
	queryJobs  *api.QueryJobs
	repository string
	id         string
	nextPoll   time.Time
}

func (q *queryJobPoller) WaitAndPollContext(ctx context.Context) (api.QueryResult, error) {
	select {
	case <-time.After(q.nextPoll.Sub(time.Now())):
	case <-ctx.Done():
		return api.QueryResult{}, ctx.Err()
	}

	result, err := q.queryJobs.PollContext(ctx, q.repository, q.id)
	if err != nil {
		return result, err
	}

	q.nextPoll = time.Now().Add(time.Duration(result.Metadata.PollAfter) * time.Millisecond)

	return result, err
}

var fieldPrinters = map[string]func(v interface{}) (string, bool){
	"@timestamp": func(v interface{}) (string, bool) {
		fv, ok := v.(float64)
		if !ok {
			return "", false
		}

		sec, msec := int64(fv)/1000, int64(fv)%1000

		t := time.Unix(sec, msec*1000000)

		return t.Format(time.RFC3339Nano), true
	},
}

type eventListPrinter struct {
	printedIds     map[string]bool
	printFields    []string
	w              io.Writer
	printEventFunc func(io.Writer, map[string]interface{})
	fmt            string
}

func newEventListPrinter(w io.Writer, fmt string) *eventListPrinter {
	e := &eventListPrinter{
		printedIds: map[string]bool{},
		w:          w,
	}

	re := regexp.MustCompile(`(\{[^\}]+\})`)
	e.fmt = re.ReplaceAllStringFunc(fmt, func(f string) string {
		field := f[1 : len(f)-1]
		arg := ""

		if strings.Contains(field, ":") {
			pieces := strings.SplitN(field, ":", 2)
			field, arg = pieces[0], pieces[1]
		}

		e.printFields = append(e.printFields, field)
		return "%" + arg + "s"
	})

	e.initPrintFunc()
	return e
}

func (p *eventListPrinter) initPrintFunc() {
	var printers []func(map[string]interface{}) string
	for _, f := range p.printFields {
		f := f
		if printer, hasPrinter := fieldPrinters[f]; hasPrinter {
			printers = append(printers, func(m map[string]interface{}) string {
				v := m[f]
				if str, ok := printer(v); ok {
					return str
				} else {
					return fmt.Sprint(v)
				}
			})
		} else {
			printers = append(printers, func(m map[string]interface{}) string {
				v := m[f]
				return fmt.Sprint(v)
			})
		}
	}

	p.printEventFunc = func(w io.Writer, m map[string]interface{}) {
		fmtArgs := make([]interface{}, len(printers))
		for i, printer := range printers {
			fmtArgs[i] = printer(m)
		}
		fmt.Fprintf(w, p.fmt+"\n", fmtArgs...)
	}
}

func (p *eventListPrinter) print(result api.QueryResult) {
	for _, e := range result.Events {
		id, hasID := e["@id"].(string)
		if hasID && !p.printedIds[id] {
			p.printEventFunc(p.w, e)
			p.printedIds[id] = true
		} else if !hasID {
			p.printEventFunc(p.w, e)
		}
	}
}

type aggregatePrinter struct {
	w       io.Writer
	columns []string
}

func newAggregatePrinter(w io.Writer) *aggregatePrinter {
	return &aggregatePrinter{
		w: w,
	}
}

func (p *aggregatePrinter) print(result api.QueryResult) {
	if p.columns == nil {
		var f []string
		for k := range result.Events[0] {
			f = append(f, k)
		}
		p.columns = f
	}

	if len(p.columns) == 0 {
		return
	}

	if len(p.columns) == 1 && len(result.Events) == 1 {
		// single column, single result, just print it
		fmt.Fprintln(p.w, result.Events[0][p.columns[0]])
		return
	}

	t := tablewriter.NewWriter(p.w)
	t.SetAutoFormatHeaders(false)
	t.SetBorder(false)
	t.SetHeader(p.columns)
	t.SetHeaderLine(false)

	for _, e := range result.Events {
		var v []string
		for _, i := range p.columns {
			v = append(v, fmt.Sprint(e[i]))
		}
		t.Append(v)
	}

	t.Render()
	fmt.Fprintln(p.w)
}
