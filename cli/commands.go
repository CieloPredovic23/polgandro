package cli

import "github.com/urfave/cli"

const (
	// JSONParamsKey ...
	JSONParamsKey = "json-params"
	// JSONParamsBase64Key ...
	JSONParamsBase64Key = "json-params-base64"
	// WorkflowKey ...
	WorkflowKey = "workflow"
	// PatternKey ...
	PatternKey = "pattern"
)

var (
	commands = []cli.Command{
		{
			Name:    "setup",
			Aliases: []string{"s"},
			Usage:   "Setup the current host. Install every required tool to run Workflows.",
			Action:  setup,
			Flags: []cli.Flag{
				flMinimalSetup,
				flFullModeSteup,
			},
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Generates a Workflow/app config file in the current directory, which then can be run immediately.",
			Action:  initConfig,
		},
		{
			Name:   "version",
			Usage:  "Prints the version",
			Action: printVersionCmd,
			Flags: []cli.Flag{
				flOutputFormat,
				cli.BoolFlag{Name: "full", Usage: "Prints the build number as well."},
			},
		},
		{
			Name:   "validate",
			Usage:  "Validates a specified bitrise config.",
			Action: validate,
			Flags: []cli.Flag{
				flPath,
				flConfig,
				flConfigBase64,
				flInventory,
				flInventoryBase64,
				flFormat,
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Runs a specified Workflow.",
			Action:  run,
			Flags: []cli.Flag{
				// deprecated
				flPath,
				// should deprecate
				cli.StringFlag{Name: ConfigBase64Key, Usage: "base64 decoded config data."},
				cli.StringFlag{Name: InventoryBase64Key, Usage: "base64 decoded inventory data."},
				// cli params
				cli.StringFlag{Name: WorkflowKey, Usage: "workflow id to run"},
				cli.StringFlag{Name: ConfigKey + ", " + configShortKey, Usage: "Path where the workflow config file is located."},
				cli.StringFlag{Name: InventoryKey + ", " + inventoryShortKey, Usage: "Path of the inventory file."},
				// cli params used in CI mode
				cli.StringFlag{Name: JSONParamsKey, Usage: "specify command flags with json string-string hash"},
				cli.StringFlag{Name: JSONParamsBase64Key, Usage: "specify command flags with base64 encoded json string-string hash"},
			},
		},
		{
			Name:   "trigger-check",
			Usage:  "Prints out which workflow will triggered by specified pattern.",
			Action: triggerCheck,
			Flags: []cli.Flag{
				flPath,
				flConfig,
				flConfigBase64,
				flInventory,
				flInventoryBase64,
				flFormat,
			},
		},
		{
			Name:    "trigger",
			Aliases: []string{"t"},
			Usage:   "Triggers a specified Workflow.",
			Action:  trigger,
			Flags: []cli.Flag{
				// deprecated
				flPath,
				// should deprecate
				cli.StringFlag{Name: ConfigBase64Key, Usage: "base64 decoded config data."},
				cli.StringFlag{Name: InventoryBase64Key, Usage: "base64 decoded inventory data."},
				// cli params
				cli.StringFlag{Name: PatternKey, Usage: "trigger pattern"},
				cli.StringFlag{Name: ConfigKey + ", " + configShortKey, Usage: "Path where the workflow config file is located."},
				cli.StringFlag{Name: InventoryKey + ", " + inventoryShortKey, Usage: "Path of the inventory file."},
				// cli params used in CI mode
				cli.StringFlag{Name: JSONParamsKey, Usage: "specify command flags with json string-string hash"},
				cli.StringFlag{Name: JSONParamsBase64Key, Usage: "specify command flags with base64 encoded json string-string hash"},
			},
		},
		{
			Name:   "export",
			Usage:  "Export the bitrise configuration.",
			Action: export,
			Flags: []cli.Flag{
				flPath,
				flConfig,
				flConfigBase64,
				flFormat,
				flOutputPath,
				flPretty,
			},
		},
		{
			Name:   "normalize",
			Usage:  "Normalize the bitrise configuration.",
			Action: normalize,
			Flags: []cli.Flag{
				flPath,
				flConfig,
				flConfigBase64,
			},
		},
		{
			Name:   "step-list",
			Usage:  "List of available steps.",
			Action: stepList,
			Flags: []cli.Flag{
				flCollection,
				flFormat,
			},
		},
		{
			Name:    "step-info",
			Aliases: []string{"i"},
			Usage:   "Provides information (step ID, last version, given version) about specified step.",
			Action:  stepInfo,
			Flags: []cli.Flag{
				flCollection,
				flVersion,
				flFormat,
				flShort,
				flStepYML,
			},
		},
		{
			Name:   "workflows",
			Usage:  "List of available workflows in config.",
			Action: workflowList,
			Flags: []cli.Flag{
				flPath,
				flConfig,
				flConfigBase64,
				flFormat,
				cli.BoolFlag{
					Name:  MinimalModeKey,
					Usage: "Print only workflow summary.",
				},
			},
		},
		{
			Name:   "share",
			Usage:  "Publish your step.",
			Action: share,
			Subcommands: []cli.Command{
				{
					Name:   "start",
					Usage:  "Preparations for publishing.",
					Action: start,
					Flags: []cli.Flag{
						flCollection,
					},
				},
				{
					Name:   "create",
					Usage:  "Create your change - add it to your own copy of the collection.",
					Action: create,
					Flags: []cli.Flag{
						flTag,
						flGit,
						flStepID,
					},
				},
				{
					Name:   "audit",
					Usage:  "Validates the step collection.",
					Action: shareAudit,
				},
				{
					Name:   "finish",
					Usage:  "Finish up.",
					Action: finish,
				},
			},
		},
		{
			Name:  "plugin",
			Usage: "Plugin handling.",
			Subcommands: []cli.Command{
				{
					Name:   "install",
					Usage:  "Intsall bitrise plugin.",
					Action: pluginInstall,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "source",
							Usage: "Plugin source url.",
						},
						cli.StringFlag{
							Name:  "bin-source",
							Usage: "Plugin binary url.",
						},
						cli.StringFlag{
							Name:  "version",
							Usage: "Plugin version tag.",
						},
					},
				},
				{
					Name:   "delete",
					Usage:  "Delete bitrise plugin.",
					Action: pluginDelete,
				},
				{
					Name:   "list",
					Usage:  "List installed bitrise plugins.",
					Action: pluginList,
				},
			},
		},
	}
)
