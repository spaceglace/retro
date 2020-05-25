package retro

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"retro/config"
	"sort"

	"github.com/wailsapp/wails"
)

const (
	baseURL = "https://retroachievements.org/API"
)

type Game struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Icon         string                 `json:"icon"`
	Console      string                 `json:"console"`
	BoxArt       string                 `json:"boxArt"`
	Softcore     string                 `json:"softcore"`
	Hardcore     string                 `json:"hardcore"`
	Achievements map[string]Achievement `json:"achievements"`
	Order        []string               `json:"order"`
}

type Achievement struct {
	ID                 string `json:"ID"`
	Badge              string `json:"BadgeName"`
	Title              string `json:"Title"`
	Description        string `json:"Description"`
	Softcore           string `json:"NumAwarded"`
	Hardcore           string `json:"NumAwardedHardcore"`
	Points             string `json:"Points"`
	TrueRatio          string `json:"TrueRatio"`
	DisplayOrder       string `json:"DisplayOrder"`
	DateEarned         string `json:"DateEarned,omitempty"`
	DateEarnedHardcore string `json:"DateEarnedHardcore,omitempty"`
}

type Retro struct {
	Logger *wails.CustomLogger
}

func NewRetro() *Retro {
	return &Retro{}
}

func (r *Retro) WailsInit(runtime *wails.Runtime) error {
	r.Logger = runtime.Log.New("api")
	r.Logger.Info("Initialized API")
	return nil
}

func (r *Retro) GetLastGame(username, apikey string) (string, error) {
	url := fmt.Sprintf(
		"%s/API_GetUserRecentlyPlayedGames.php?z=%s&y=%s&u=%s&c=1",
		baseURL,
		username,
		apikey,
		username,
	)
	response, err := http.Get(url)
	if code := response.StatusCode; err != nil || code < 200 || code >= 300 {
		r.Logger.Errorf(
			"Error getting user '%s' recently played games (%d): %w",
			username,
			code,
			err,
		)
		return "", err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		r.Logger.Errorf(
			"Error reading response of user '%s' recently played games: %w",
			username,
			err,
		)
		return "", err
	}

	var container []struct {
		GameID string `json:"GameID"`
	}
	if err := json.Unmarshal(responseBody, &container); err != nil {
		r.Logger.Errorf(
			"Error parsing payload of user '%s' recently played games: %w",
			username,
			err,
		)
		return "", err
	}

	r.Logger.Debugf(
		"Collected user '%s' recently played games",
		username,
	)
	return container[0].GameID, nil
}

func (r *Retro) GetGameInformation(username, apikey, gameid string) (Game, error) {
	url := fmt.Sprintf(
		"%s/API_GetGameInfoAndUserProgress.php?z=%s&y=%s&u=%s&g=%s",
		baseURL,
		username,
		apikey,
		username,
		gameid,
	)
	response, err := http.Get(url)
	if code := response.StatusCode; err != nil || code < 200 || code >= 300 {
		r.Logger.Errorf(
			"Error getting game information for user '%s' (%d): %w",
			username,
			code,
			err,
		)
		return Game{}, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		r.Logger.Errorf(
			"Error reading game info response for user '%s': %w",
			username,
			err,
		)
		return Game{}, err
	}

	var container struct {
		Title        string                 `json:"Title"`
		Console      string                 `json:"ConsoleName"`
		Icon         string                 `json:"ImageIcon"`
		BoxArt       string                 `json:"ImageBoxArt"`
		Softcore     string                 `json:"NumDistinctPlayersCasual"`
		Hardcore     string                 `json:"NumDistinctPlayersHardcore"`
		Achievements map[string]Achievement `json:"Achievements,omitempty"`
	}
	if err := json.Unmarshal(responseBody, &container); err != nil {
		r.Logger.Errorf(
			"Error parsing payloads of game info for user '%s': %w",
			username,
			err,
		)
		return Game{}, err
	}
	r.Logger.Debugf(
		"Collected game info for user '%s', found %d achievements",
		username,
		len(container.Achievements),
	)
	result := Game{
		ID:           gameid,
		Title:        container.Title,
		Console:      container.Console,
		Icon:         container.Icon,
		BoxArt:       container.BoxArt,
		Softcore:     container.Softcore,
		Hardcore:     container.Hardcore,
		Achievements: container.Achievements,
	}
	// look for a saved game order for achievements
	if order, ok := config.GetSettings().Orders[result.ID]; ok {
		result.Order = order
	} else {
		// no saved order, sort in display order
		for _, achievement := range result.Achievements {
			result.Order = append(result.Order, achievement.ID)
		}
		sort.Strings(result.Order)
	}
	return result, nil
}
