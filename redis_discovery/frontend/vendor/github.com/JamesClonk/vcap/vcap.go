package vcap

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

type VCAP struct {
	Application struct {
		ID            string   `json:"application_id"`
		Name          string   `json:"application_name"`
		Version       string   `json:"application_version"`
		InstanceID    string   `json:"instance_id"`
		InstanceIndex int      `json:"instance_index"`
		Host          string   `json:"host"`
		Port          int      `json:"port"`
		Users         string   `json:"users"`
		URIs          []string `json:"application_uris"`
		Limits        struct {
			Memory int `json:"mem"`
			Disk   int `json:"disk"`
			Files  int `json:"fds"`
		} `json:"limits"`
		Started *Timestamp `json:"started_at_timestamp"`
		State   *Timestamp `json:"state_timestamp"`
	}
	Host            string
	Port            int
	Services        map[string][]Service
	InstanceAddress string
	InstanceIP      string
	InstancePort    int
}

type Service struct {
	Name        string                 `json:"name"`
	Label       string                 `json:"label"`
	Tags        []string               `json:"tags"`
	Plan        string                 `json:"plan"`
	Credentials map[string]interface{} `json:"credentials"`
}

type Timestamp time.Time

func (t *Timestamp) String() string {
	return time.Time(*t).String()
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(ts, 0))
	return nil
}

func New() (*VCAP, error) {
	vcap := &VCAP{}

	vcap.Host = os.Getenv("VCAP_APP_HOST")
	if port := os.Getenv("VCAP_APP_PORT"); port != "" {
		p, err := strconv.ParseInt(port, 10, 32)
		if err != nil {
			return nil, err
		}
		vcap.Port = int(p)
	}

	vcap.InstanceAddress = os.Getenv("CF_INSTANCE_ADDR")
	vcap.InstanceIP = os.Getenv("CF_INSTANCE_IP")
	if port := os.Getenv("CF_INSTANCE_PORT"); port != "" {
		p, err := strconv.ParseInt(port, 10, 32)
		if err != nil {
			return nil, err
		}
		vcap.InstancePort = int(p)
	}

	if app := os.Getenv("VCAP_APPLICATION"); app != "" {
		if err := json.Unmarshal([]byte(app), &(vcap.Application)); err != nil {
			return nil, err
		}
	}

	if serv := os.Getenv("VCAP_SERVICES"); serv != "" {
		if err := json.Unmarshal([]byte(serv), &(vcap.Services)); err != nil {
			return nil, err
		}
	}

	// set some defaults in case of local development / missing VCAP_APPLICATION
	if vcap.Application.ID == "" {
		vcap.Application.ID = "123-456-789"
	}
	if vcap.Application.Name == "" {
		vcap.Application.Name = "devapp"
	}
	if vcap.Application.InstanceID == "" {
		vcap.Application.InstanceID = "987-654-321"
	}
	if vcap.Application.InstanceIndex == 0 {
		vcap.Application.InstanceIndex = 1
	}
	if vcap.Application.Host == "" {
		vcap.Application.Host = "localhost"
	}
	if vcap.Application.Port == 0 {
		vcap.Application.Port = 4000
	}
	if vcap.Host == "" {
		vcap.Host = "localhost"
	}
	if vcap.Port == 0 {
		vcap.Port = 4000
	}

	return vcap, nil
}

func (v *VCAP) GetService(name string) *Service {
	for _, services := range v.Services {
		for _, service := range services {
			if service.Name == name {
				return &service
			}
		}
	}
	return nil
}
