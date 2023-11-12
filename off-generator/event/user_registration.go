package event

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bxcodec/faker/v3"
	"github.com/juanjoss/off-generator/model"
)

var registerEndpoint = "http://users:8080/api/register"

type UserRegistration struct{}

type UserRegistrationRequest struct {
	User    model.User  `json:"user"`
	Devices []model.SSD `json:"devices"`
}

func (ur *UserRegistration) Handle() {
	// limiting the number of devices generated by faker
	_ = faker.SetRandomMapAndSliceMinSize(1)
	_ = faker.SetRandomMapAndSliceMaxSize(3)

	request := UserRegistrationRequest{
		User:    model.User{},
		Devices: []model.SSD{},
	}

	// generating fake user
	err := faker.FakeData(&request.User)
	if err != nil {
		log.Printf("unable to generate fake user data: %v", err)
	}

	// generating fake device
	err = faker.FakeData(&request.Devices)
	if err != nil {
		log.Printf("unable to generate fake SSD data: %v", err)
	}

	// send registration to users HTTP service
	jsonBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("unable to marshal request: %v", err)
	}

	res, err := http.Post(
		registerEndpoint,
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Printf("unable to POST request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Println(registerEndpoint)
		log.Println("event", ur.Type(), "failed")
	}
}

func (ur *UserRegistration) Type() string {
	return "user-registration"
}