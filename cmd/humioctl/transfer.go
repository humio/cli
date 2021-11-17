package main

import (
	"fmt"
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func newTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "transfer",
		Short:  "Cluster transfers [Experimental]",
		Hidden: true,
	}

	cmd.AddCommand(newTransferCreateManagedExportGroupCmd())
	cmd.AddCommand(newTransferGetManagedExportGroupCmd())
	cmd.AddCommand(newTransferRemoveManagedExportGroupCmd())
	cmd.AddCommand(newTransferJobsCmd())

	return cmd
}

func newTransferCreateManagedExportGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-managed-export-group",
		Short: "Create a group that has an attached role with permission to export the organization through the cluster transfer API.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			groupID, err := client.Transfer().CreateManagedExportGroup()
			exitOnError(cmd, err, "Error creating managed export group")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created managed export group with group ID: %v\n", groupID)
		},
	}

	return cmd
}

func newTransferGetManagedExportGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-managed-export-group",
		Short: "Get the id of the managed export group.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			groupID, err := client.Transfer().GetManagedExportGroup()
			exitOnError(cmd, err, "Error retrieving managed export group")

			fmt.Fprintf(cmd.OutOrStdout(), "Group ID: %v\n", groupID)
		},
	}

	return cmd
}

func newTransferRemoveManagedExportGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-managed-export-group",
		Short: "Remove the managed export group.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			err := client.Transfer().RemoveManagedExportGroup()
			exitOnError(cmd, err, "Error removing managed export group")

			fmt.Fprintln(cmd.OutOrStdout(), "Successfully removed managed export group")
		},
	}

	return cmd
}

func newTransferJobsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Transfer jobs",
	}

	cmd.AddCommand(newTransferJobsListCmd())
	cmd.AddCommand(newTransferJobsAddCmd())
	cmd.AddCommand(newTransferJobsCancelCmd())
	cmd.AddCommand(newTransferJobsStatusCmd())

	return cmd
}

func newTransferJobsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List jobs",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			jobs, err := client.Transfer().ListTransferJobs()
			exitOnError(cmd, err, "Error listing transfer jobs")

			var rows [][]string
			for _, job := range jobs {
				at := ""
				state := "Active"
				if job.CompletedAt != nil {
					state = "Completed"
					at = job.CompletedAt.String()
				} else if job.CancelledAt != nil {
					state = "Cancelled"
					at = job.CancelledAt.String()
				}
				rows = append(rows, []string{job.ID, job.SourceClusterURL, strconv.Itoa(job.MaximumParallelDownloads), strconv.Itoa(len(job.Dataspaces)), state, at})
			}

			printOverviewTable(cmd, []string{"ID | Source | Parallel | No. Dataspaces | State | At"}, rows)
		},
	}

	return cmd
}

func detailTransferJob(cmd *cobra.Command, job interface{}) {
	var details [][]string

	switch j := job.(type) {
	case api.TransferJob:
		details = append(details, []string{"ID", j.ID})
		details = append(details, []string{"Source Cluster", j.SourceClusterURL})
		details = append(details, []string{"Maximum parallel downloads", strconv.Itoa(j.MaximumParallelDownloads)})
		details = append(details, []string{"Dataspaces", j.Dataspaces[0]})
		for _, ds := range j.Dataspaces[1:] {
			details = append(details, []string{"", ds})
		}
		if j.CompletedAt != nil {
			details = append(details, []string{"Completed At", j.CompletedAt.String()})
		}
		if j.CancelledAt != nil {
			details = append(details, []string{"Cancelled At", j.CancelledAt.String()})
		}
	case api.TransferJobStatus:
		details = append(details, []string{"Status", j.Status})
		details = append(details, []string{"Status Line", j.StatusLine})
		details = append(details, []string{"Running", fmt.Sprint(j.Running)})
		details = append(details, []string{"Error", j.Error})
		details = append(details, []string{"Progress", fmt.Sprintf("%d/%d", j.TransferredSegments, j.TotalSegments)})
	}

	printDetailsTable(cmd, details)
}

func newTransferJobsAddCmd() *cobra.Command {
	var (
		setTargetAsNewMaster     bool
		onlyTransferDataspaces   bool
		maximumParallelDownloads int
	)

	cmd := &cobra.Command{
		Use:   "add [flags] <source cluster url> <source cluster token> <destination organization id> <dataspaceID[,dataspaceID]*>",
		Short: "Add jobs",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			jobs, err := client.Transfer().AddTransferJob(args[0], args[1], args[2], strings.Split(args[3], ","), maximumParallelDownloads, setTargetAsNewMaster, onlyTransferDataspaces)
			exitOnError(cmd, err, "Error creating transfer job")

			fmt.Fprintf(cmd.OutOrStdout(), "Added transfer job with ID: %s\n", jobs.ID)
		},
	}

	cmd.Flags().BoolVar(&setTargetAsNewMaster, "set-as-new-master", false, "Configures the source cluster to proxy ingest to the new cluster when segments have been transferred.")
	cmd.Flags().BoolVar(&onlyTransferDataspaces, "only-transfer-dataspaces", false, "Only transfer dataspaces and segments, and skip metadata such as users and dashboards.")
	cmd.Flags().IntVar(&maximumParallelDownloads, "max-parallel-downloads", 0, "Limit maximum parallel transfers of segment files.")

	return cmd
}

func newTransferJobsCancelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel <transfer job id>",
		Short: "Cancel an ongoing transfer job",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			job, err := client.Transfer().CancelTransferJob(args[0])
			exitOnError(cmd, err, "Error cancelling transfer job")

			detailTransferJob(cmd, job)
		},
	}

	return cmd
}

func newTransferJobsStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <transfer job id>",
		Short: "Get status of an ongoing transfer job",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			job, err := client.Transfer().GetTransferJobStatus(args[0])
			exitOnError(cmd, err, "Error getting status of transfer job")

			detailTransferJob(cmd, job)
		},
	}

	return cmd
}
