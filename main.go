// main.go

package main

import (
	"encoding/json"
	"flag"
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
	configFile *string
)

func aclContains(str string) bool {
	for _, a := range conf.ACL {
		if a == str {
			return true
		}
	}
	return false
}

func acl() gin.HandlerFunc {
	return func(c *gin.Context) {
		if aclContains(c.ClientIP()) {
			c.Next()
		} else {
			c.String(http.StatusForbidden, "403 forbidden")
		}
	}
}

func main() {
	configFile = flag.String("conf", "/etc/webrunner.yaml", "")
	parseConfig()
	if conf.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.Default()
	router.Use(acl())

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

	router.Run(conf.Listen)
}
