// main.go

package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type (
	// TConf - struktura pro naparsovani konfiguracniho souboru
	TConf struct {
		Listen string
		ACL    []string
		Debug  bool
		Jobs   []struct {
			Name  string
			Tasks []struct {
				Name string
				File string
			}
		}
	}
)

var (
	router     *gin.Engine
	conf       TConf
	configFile = "config.yaml"
)

func main() {
	parseConfig()
	if conf.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.Default()

	for _, job := range conf.Jobs {
		router.GET("/"+job.Name+"/:domain", func(c *gin.Context) {
			var wg sync.WaitGroup
			responses := make(map[string]string)
			param := c.Param("domain")
			wg.Add(len(job.Tasks))
			for _, task := range job.Tasks {
				go func(name, file, param string) {
					response, _ := exec.Command(file, strings.Split(param, "+")...).Output()
					responses[name] = string(response)
					wg.Done()
				}(task.Name, task.File, param)
			}
			wg.Wait()
			jsonResponse, _ := json.Marshal(responses)
			c.String(http.StatusOK, string(jsonResponse))
		})
	}

	router.Run(":8080")
}
