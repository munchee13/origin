package openshift

import (
	"fmt"
	"os"

	kubeversion "github.com/GoogleCloudPlatform/kubernetes/pkg/version"
	"github.com/spf13/cobra"

	"github.com/openshift/origin/pkg/cmd/cli"
	"github.com/openshift/origin/pkg/cmd/experimental/config"
	"github.com/openshift/origin/pkg/cmd/experimental/generate"
	"github.com/openshift/origin/pkg/cmd/experimental/login"
	"github.com/openshift/origin/pkg/cmd/experimental/policy"
	"github.com/openshift/origin/pkg/cmd/experimental/project"
	"github.com/openshift/origin/pkg/cmd/experimental/tokens"
	"github.com/openshift/origin/pkg/cmd/flagtypes"
	"github.com/openshift/origin/pkg/cmd/infra/builder"
	"github.com/openshift/origin/pkg/cmd/infra/deployer"
	"github.com/openshift/origin/pkg/cmd/infra/router"
	"github.com/openshift/origin/pkg/cmd/server"
	"github.com/openshift/origin/pkg/cmd/templates"
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"
	"github.com/openshift/origin/pkg/version"
)

const longDescription = `
OpenShift for Admins

OpenShift helps you build, deploy, and manage your applications. To start an all-in-one server, run:

    $ openshift start &

OpenShift is built around Docker and the Kubernetes cluster container manager.  You must have
Docker installed on this machine to start your server.

Note: This is an alpha release of OpenShift and will change significantly.  See
    https://github.com/openshift/origin for the latest information on OpenShift.
`

// CommandFor returns the appropriate command for this base name,
// or the global OpenShift command
func CommandFor(basename string) *cobra.Command {
	var cmd *cobra.Command

	switch basename {
	case "openshift-router":
		cmd = router.NewCommandTemplateRouter(basename)
	case "openshift-deploy":
		cmd = deployer.NewCommandDeployer(basename)
	case "openshift-sti-build":
		cmd = builder.NewCommandSTIBuilder(basename)
	case "openshift-docker-build":
		cmd = builder.NewCommandDockerBuilder(basename)
	case "osc":
		cmd = cli.NewCommandCLI(basename, basename)
	default:
		cmd = NewCommandOpenShift()
	}

	flagtypes.GLog(cmd.PersistentFlags())

	return cmd
}

// NewCommandOpenShift creates the standard OpenShift command
func NewCommandOpenShift() *cobra.Command {
	root := &cobra.Command{
		Use:   "openshift",
		Short: "OpenShift helps you build, deploy, and manage your cloud applications",
		Long:  longDescription,
		Run: func(c *cobra.Command, args []string) {
			c.SetOutput(os.Stdout)
			c.Help()
		},
	}

	root.SetUsageTemplate(templates.MainUsageTemplate())
	root.SetHelpTemplate(templates.MainHelpTemplate())

	root.AddCommand(server.NewCommandStartServer("start"))
	root.AddCommand(cli.NewCommandCLI("cli", "openshift cli"))
	root.AddCommand(cli.NewCmdKubectl("kube"))
	root.AddCommand(newExperimentalCommand("openshift", "ex"))
	root.AddCommand(newVersionCommand("version"))

	// infra commands are those that are bundled with the binary but not displayed to end users
	// directly
	infra := &cobra.Command{
		Use: "infra", // Because this command exposes no description, it will not be shown in help
	}

	infra.AddCommand(
		router.NewCommandTemplateRouter("router"),
		deployer.NewCommandDeployer("deploy"),
		builder.NewCommandSTIBuilder("sti-build"),
		builder.NewCommandDockerBuilder("docker-build"),
	)
	root.AddCommand(infra)

	return root
}

func newExperimentalCommand(parentName, name string) *cobra.Command {
	experimental := &cobra.Command{
		Use:   name,
		Short: "Experimental commands under active development",
		Long:  "The commands grouped here are under development and may change without notice.",
		Run: func(c *cobra.Command, args []string) {
			c.SetOutput(os.Stdout)
			c.Help()
		},
	}

	f := clientcmd.New(experimental.PersistentFlags())

	subName := fmt.Sprintf("%s %s", parentName, name)
	experimental.AddCommand(project.NewCmdNewProject(f, subName, "new-project"))
	experimental.AddCommand(config.NewCmdConfig(subName, "config"))
	experimental.AddCommand(tokens.NewCmdTokens(f, subName, "tokens"))
	experimental.AddCommand(policy.NewCommandPolicy(f, subName, "policy"))
	experimental.AddCommand(generate.NewCmdGenerate(f, subName, "generate"))
	experimental.AddCommand(login.NewCmdLogin(f, subName, "login"))
	//experimental.AddCommand(exrouter.NewCmdRouter(f, subName, "router", os.Stdout))
	return experimental
}

// newVersionCommand creates a command for displaying the version of this binary
func newVersionCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "Display version",
		Run: func(c *cobra.Command, args []string) {
			fmt.Printf("openshift %v\n", version.Get())
			fmt.Printf("kubernetes %v\n", kubeversion.Get())
		},
	}
}
