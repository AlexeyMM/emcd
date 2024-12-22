package jobs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/config"
)

type APIJobsClient struct {
	basePath string
	client   *resty.Client
}

func NewJobsClient(cfg config.APIJobsConfig) *APIJobsClient {
	baseURL := cfg.Host + cfg.Path

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 20
	transport.MaxConnsPerHost = 30

	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: transport,
	}

	cc := APIJobsClient{
		basePath: baseURL,
		client: resty.NewWithClient(client).
			SetHeader("x-access-key", cfg.Token).
			SetBaseURL(baseURL),
	}

	cc.defaultJobs(cfg.Token)

	return &cc
}

type DeleteJobRequest struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
}

func (c *APIJobsClient) DeleteJobSet(namespace, key string) error {
	path := "/jobs"
	request := DeleteJobRequest{
		Namespace: namespace,
		Key:       key,
	}

	resp, err := c.client.SetDebug(true).R().
		SetBody(request).
		Delete(path)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(string(resp.Body()))
	}

	return nil
}

type GetJobResponse struct {
	Current  string    `json:"current"`
	NextTime time.Time `json:"next_time"`
}

func (c *APIJobsClient) GetActiveJob(namespace, key string) (*GetJobResponse, error) {
	path := "/jobs/%s"
	jobKey := namespace + key

	var response GetJobResponse

	resp, err := c.client.SetDebug(true).R().
		SetResult(&response).
		Get(fmt.Sprintf(path, jobKey))
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New(string(resp.Body()))
	}

	return &response, nil
}

type CreateJobRequest struct {
	Namespace   string         `json:"namespace"`
	Key         string         `json:"key"`
	Start       string         `json:"start"`
	Token       string         `json:"token"`
	FirstCall   int64          `json:"first_call,omitempty"`
	ShouldExist bool           `json:"should_exist,omitempty"`
	UpdateExist bool           `json:"update_exist,omitempty"`
	Jobs        map[string]Job `json:"jobs"`
}

type Job struct {
	Name     string `json:"name"`
	Interval int    `json:"interval"`
	MaxTries int    `json:"tries"    default:"0"`
	Failure  string `json:"failure"`
	Success  string `json:"success"`
	Callback string `json:"callback"`
	Payload  string `json:"payload"`
}

func (c *APIJobsClient) createJobSet(jobs *CreateJobRequest) error {
	path := "/jobs"

	resp, err := c.client.SetDebug(true).R().
		SetBody(jobs).
		Post(path)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(string(resp.Body()))
	}

	return nil
}

func (c *APIJobsClient) defaultJobs(token string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(context.Background(), "recovered in defaultJobs: %+v", r)
		}
	}()
	checkBalance := Job{
		Name:     "check-balance",
		Interval: 60 * 60,
		MaxTries: 0,
		Callback: "/v1/fiat/2/balance/notify",
		Payload:  "",
	}

	jobsMap := map[string]Job{
		"check-balance": checkBalance,
	}

	jobRequest := CreateJobRequest{
		Namespace:   "fiat",
		Key:         "checkBalance",
		Start:       "check-balance",
		Jobs:        jobsMap,
		Token:       token,
		FirstCall:   30,
		ShouldExist: true,
	}

	if err := c.createJobSet(&jobRequest); err != nil {
		log.Error(context.Background(), "job name: %s, err: %s", jobRequest.Key, err.Error())
	}
}
