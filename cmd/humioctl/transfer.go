package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func newTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Cluster transfers",
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
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)

			groupID, err := client.Transfer().CreateManagedExportGroup()

			if err != nil {
				return nil, fmt.Errorf("error creating managed group: %w", err)
			}

			return groupID, nil
		}),
	}

	return cmd
}

func newTransferGetManagedExportGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-managed-export-group",
		Short: "Get the id of the managed export group.",
		Args:  cobra.ExactArgs(0),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)

			groupID, err := client.Transfer().GetManagedExportGroup()
			if err != nil {
				return nil, fmt.Errorf("error retrieving managed group: %w", err)
			}

			return groupID, nil
		}),
	}

	return cmd
}

func newTransferRemoveManagedExportGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-managed-export-group",
		Short: "Remove the managed export group.",
		Args:  cobra.ExactArgs(0),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)

			err := client.Transfer().RemoveManagedExportGroup()
			if err != nil {
				return nil, fmt.Errorf("error removing managed group: %w", err)
			}

			return fmt.Sprintf("Removed managed export group."), nil
		}),
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
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)

			jobs, err := client.Transfer().ListTransferJobs()
			if err != nil {
				return nil, fmt.Errorf("error listing transfer jobs: %w", err)
			}

			return jobs, nil
		}),
	}

	return cmd
}

func newTransferJobsAddCmd() *cobra.Command {
	var (
		setTargetAsNewMaster     bool
		maximumParallelDownloads int
	)

	cmd := &cobra.Command{
		Use:   "add <source cluster url> <source cluster token> <destination organization id> <dataspaceID[,dataspaceID]*>",
		Short: "Add jobs",
		Args:  cobra.ExactArgs(4),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)

			job, err := client.Transfer().AddTransferJob(args[0], args[1], args[2], strings.Split(args[3], ","), maximumParallelDownloads, setTargetAsNewMaster)
			if err != nil {
				return nil, fmt.Errorf("error adding transfer jobs: %w", err)
			}

			return job, nil
		}),
	}

	cmd.Flags().BoolVar(&setTargetAsNewMaster, "set-as-new-master", false, "Configures the source cluster to proxy ingest to the new cluster when segments have been transferred.")
	cmd.Flags().IntVar(&maximumParallelDownloads, "max-parallel-downloads", 0, "Limit maximum parallel transfers of segment files.")

	return cmd
}

func newTransferJobsCancelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel <transfer job id>",
		Short: "Cancel an ongoing transfer job",
		Args:  cobra.ExactArgs(1),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)

			job, err := client.Transfer().CancelTransferJob(args[0])
			if err != nil {
				return nil, fmt.Errorf("error cancelling transfer jobs: %w", err)
			}

			return job, nil
		}),
	}

	return cmd
}

func newTransferJobsStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <transfer job id>",
		Short: "Get status of an ongoing transfer job",
		Args:  cobra.ExactArgs(1),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)

			job, err := client.Transfer().GetTransferJobStatus(args[0])
			if err != nil {
				return nil, fmt.Errorf("error getting transfer jobs status: %w", err)
			}

			return job, nil
		}),
	}

	return cmd
}
