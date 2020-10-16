	Use:    "start",
	Short:  "Starting node",
	PreRun: loadConfigWKey,
	RunE: func(cmd *cobra.Command, args []string) error {
		return loadStartRun()
	},
}
