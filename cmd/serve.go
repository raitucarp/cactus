package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func init() {
	serveCmd.Flags().Int("port", 9090, "server port")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:       "serve",
	Short:     "Serve directory as a data source",
	Long:      `Cactus serve will read some path as data source.`,
	Example:   "cactus serve ./some/data/path",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{"./"},
	Run: func(cmd *cobra.Command, args []string) {

		dataSource := args[0]

		if _, err := os.Stat(dataSource); os.IsNotExist(err) {
			cmd.PrintErrln("Error: " + dataSource + " does not exist.")
			cmd.Println("----------------------")
			cmd.Usage()
			os.Exit(0)
		}

		portFlag := cmd.Flag("port").Value
		port := fmt.Sprintf(":%v", portFlag)

		os.Setenv("GIN_MODE", "release")
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
		}

		r := gin.New()
		r.Use(gin.Recovery())
		r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		}))

		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		fmt.Printf("Cactus now serve %s directory, on port %s\n", dataSource, port)
		r.Run(port)
	},
}
