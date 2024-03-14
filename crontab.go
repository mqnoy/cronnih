package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Options struct {
	Enable      bool          `yaml:"enable"`
	Schedule    string        `yaml:"schedule"`
	Periodic    time.Duration `yaml:"periodic"`
	Endpoint    string        `yaml:"endpoint"`
	Credentials Credentials   `yaml:"credentials"`
}

type CronJob struct {
	Jobs map[string]Options `yaml:"jobs"`
}

type Crontab struct {
	cc *cron.Cron
}

func (ct *Crontab) Setup() {
	yamlFile, err := os.ReadFile("jobs.yaml")
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	var cronJob CronJob
	err = yaml.Unmarshal(yamlFile, &cronJob)
	if err != nil {
		log.Fatalf("error parsing YAML: %v", err)
	}

	for jobName, jobOptions := range cronJob.Jobs {
		// skip disable job
		if !jobOptions.Enable {
			continue
		}

		jobName := jobName
		jobOptions := jobOptions

		go func(jobName string, job Options) {
			if job.Periodic.Seconds() != 0 {
				for {
					time.Sleep(job.Periodic)
					log.Infof("Executing job %s for endpoint %s\n", jobName, jobOptions.Endpoint)
					err := makeRequest(job.Endpoint, job.Credentials)
					if err != nil {
						log.Warnf("Error makeRequest %s: %v", jobName, err)
					}
				}
			}

			if _, err := cron.ParseStandard(jobOptions.Schedule); err == nil {
				_, err := ct.cc.AddFunc(job.Schedule, func() {
					fmt.Printf("Executing job %s for endpoint %s\n", jobName, job.Endpoint)
					err := makeRequest(job.Endpoint, job.Credentials)
					if err != nil {
						log.Warnf("Error makeRequest %s: %v", jobName, err)
					}
				})
				if err != nil {
					log.Fatalf("error scheduling job %s: %v", jobName, err)
				}
			}

		}(jobName, jobOptions)
	}

}

func makeRequest(url string, creds Credentials) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	if creds.Username != "" && creds.Password != "" {
		req.SetBasicAuth(creds.Username, creds.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK status code received: %d", resp.StatusCode)
	}

	return nil
}
