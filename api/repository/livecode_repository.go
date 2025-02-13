package repository

import (
	"context"
	"encoding/json"
	"errors"
	"livecode_tribalworldwide/api/entities"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var ApiUrl = "https://randomuser.me/api/?results=1000&inc=gender,name,email,location,login" //Se le puede modificar el numero para pedir mas usuarios por petición pero con cuidado para que no tire rate limit

const (
	baseTimeout = 3 * time.Second //Se puede modificar para cambiar el tiempo base de espera por peticion
	workerNums  = 3               // Se pueden aumentar los workers pero hay que tener cuidado por el mismo tema de que la api que se consume tiene un rate limit y saca error 429 si son muchos requests o pidiendo muchos datos
	maxRetries  = 5
)

var httpClient = &http.Client{
	Timeout: baseTimeout,
}

type LivecodeRepository interface {
	GetUsers() ([]entities.User, error)
}

type LivecodeRepo struct {
	logger logrus.FieldLogger
}

func NewLivecodeRepository(logger logrus.FieldLogger) *LivecodeRepo {
	return &LivecodeRepo{logger: logger}
}

type apiResponse struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Email    string `json:"email"`
		Location struct {
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"location"`
		Login struct {
			UUID string `json:"uuid"`
		} `json:"login"`
	} `json:"results"`
}

func (r *LivecodeRepo) fetchUsers() ([]entities.User, error) {
	var data apiResponse

	for attempt := 1; attempt <= maxRetries; attempt++ {
		timeout := baseTimeout + time.Duration(attempt)*time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", ApiUrl, nil)
		if err != nil {
			return nil, err
		}
		r.logger.Infof("Fetching users, attempt %d...", attempt)

		start := time.Now()
		resp, err := httpClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			r.logger.Warnf("Error fetching users: %v", err)
		} else {
			defer resp.Body.Close()
			r.logger.Infof("Response received in %v", duration)
			if resp.StatusCode == http.StatusOK {
				if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
					users := make([]entities.User, 0, len(data.Results))
					for _, item := range data.Results {
						users = append(users, entities.User{
							UUID:      item.Login.UUID,
							FirstName: item.Name.First,
							LastName:  item.Name.Last,
							Email:     item.Email,
							City:      item.Location.City,
							Country:   item.Location.Country,
							Gender:    item.Gender,
						})
					}
					return users, nil
				}
			}
			if resp.StatusCode == http.StatusTooManyRequests {
				r.logger.Warnf("Rate limited (attempt %d), retrying...", attempt)
				if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
					if seconds, err := strconv.Atoi(retryAfter); err == nil {
						time.Sleep(time.Duration(seconds)*time.Second + time.Duration(rand.Intn(1000))*time.Millisecond)
						continue
					}
				}
				sleepTime := time.Duration(2<<attempt+rand.Intn(500)) * time.Millisecond
				time.Sleep(sleepTime)
				continue
			}
			return nil, errors.New("invalid response status")
		}
	}
	return nil, errors.New("failed to fetch users after retries")
}

func (r *LivecodeRepo) GetUsers() ([]entities.User, error) {
	var wg sync.WaitGroup
	userChan := make(chan []entities.User, workerNums)

	for i := 0; i < workerNums; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			users, err := r.fetchUsers()
			if err != nil {
				r.logger.Warn("Worker failed to fetch users")
				return
			}
			userChan <- users
		}()
	}

	go func() {
		wg.Wait()
		close(userChan)
	}()

	allUsers := make(map[string]entities.User)
	for users := range userChan {
		for _, user := range users {
			allUsers[user.UUID] = user
		}
	}

	result := make([]entities.User, 0, len(allUsers))
	for _, user := range allUsers {
		result = append(result, user)
	}

	return result, nil
}
