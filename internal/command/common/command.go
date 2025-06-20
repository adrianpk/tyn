// Package common proporciona estructuras e interfaces comunes para los comandos de la aplicación
package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adrianpk/tyn/internal/bkg"
	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Command interface {
	GetCobraCommand() *cobra.Command
	ExecuteDirect(ctx context.Context, args []string, flags map[string]interface{}) error
	ExecuteViaIPC(args []string, flags map[string]interface{}) error
	Name() string
}

type BaseCommand struct {
	Svc         *svc.Svc
	CobraCmd    *cobra.Command
	CommandName string
}

func (b *BaseCommand) GetCobraCommand() *cobra.Command {
	return b.CobraCmd
}

func (b *BaseCommand) Name() string {
	return b.CommandName
}

func (b *BaseCommand) Execute(ctx context.Context, args []string, flags map[string]interface{},
	directFn func(context.Context, []string, map[string]interface{}) error,
	ipcFn func([]string, map[string]interface{}) error) error {

	if b.Svc != nil {
		log.Printf("Executing command '%s' directly (service available)", b.CommandName)
		return directFn(ctx, args, flags)
	}

	log.Printf("Executing command '%s' via IPC (using daemon)", b.CommandName)
	return ipcFn(args, flags)
}

func SendToIPC(commandName string, params interface{}) error {
	resp, err := bkg.SendCommand(commandName, params)
	if err != nil {
		return fmt.Errorf("error communicating with daemon: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("daemon returned error: %s", resp.Error)
	}

	return nil
}

func ExtractFlagsFromCommand(cmd *cobra.Command) map[string]interface{} {
	flags := make(map[string]interface{})
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flags[flag.Name] = flag.Value.String()
	})
	return flags
}

func UnmarshalResponse(resp *bkg.Response, result interface{}) error {
	if !resp.Success {
		return fmt.Errorf("daemon returned error: %s", resp.Error)
	}

	if result == nil {
		return nil
	}

	err := json.Unmarshal(resp.Data, result)
	if err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	return nil
}
