/*---------------------------------------------------------------------------------------------
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.TrimSpace(strings.Join([]string{
			consts.VERSION, consts.BuildInfo}, " ",
		)))
	},
}
