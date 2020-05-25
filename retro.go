package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

const (
	baseURL = "https://retroachievements.org/API"
)

type game struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Icon         string                 `json:"icon"`
	Console      string                 `json:"console"`
	BoxArt       string                 `json:"boxArt"`
	Softcore     string                 `json:"softcore"`
	Hardcore     string                 `json:"hardcore"`
	Achievements map[string]achievement `json:"achievements"`
}

type achievement struct {
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

func getLastGame(username, apikey string) (string, error) {
	url := fmt.Sprintf(
		"%s/API_GetUserRecentlyPlayedGames.php?z=%s&y=%s&u=%s&c=1",
		baseURL,
		username,
		apikey,
		username,
	)
	response, err := http.Get(url)
	if code := response.StatusCode; err != nil || code < 200 || code >= 300 {
		logger.Error("Error getting user recently played games",
			zap.Int("code", code),
			zap.String("username", username),
			zap.Error(err),
		)
		return "", err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response of user recently played games",
			zap.String("username", username),
			zap.Error(err),
		)
		return "", err
	}

	var container []struct {
		GameID string `json:"GameID"`
	}
	if err := json.Unmarshal(responseBody, &container); err != nil {
		logger.Error("Error parsing payload of user recently played games",
			zap.String("username", username),
			zap.String("payload", string(responseBody)),
			zap.Error(err),
		)
		return "", err
	}

	logger.Debug("Collected user recently played games",
		zap.String("username", username),
		zap.Any("result", container),
	)
	return container[0].GameID, nil
}

func getGameInformation(username, apikey, gameid string) (game, error) {
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
		logger.Error("Error getting game information",
			zap.Int("code", code),
			zap.String("username", username),
			zap.String("gameid", gameid),
			zap.Error(err),
		)
		return game{}, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response of game info",
			zap.String("username", username),
			zap.Error(err),
		)
		return game{}, err
	}

	var container struct {
		GameID       string                 `json:"GameID"`
		Title        string                 `json:"Title"`
		Console      string                 `json:"ConsoleName"`
		Icon         string                 `json:"ImageIcon"`
		BoxArt       string                 `json:"ImageBoxArt"`
		Softcore     string                 `json:"NumDistinctPlayersCasual"`
		Hardcore     string                 `json:"NumDistinctPlayersHardcore"`
		Achievements map[string]achievement `json:"Achievements,omitempty"`
	}
	if err := json.Unmarshal(responseBody, &container); err != nil {
		logger.Error("Error parsing payload of game info",
			zap.String("username", username),
			zap.String("payload", string(responseBody)),
			zap.Error(err),
		)
		return game{}, err
	}

	logger.Debug("Collected game info",
		zap.String("game", container.Title),
		zap.String("username", username),
		zap.Int("achievements", len(container.Achievements)),
	)
	result := game{
		ID:           container.GameID,
		Title:        container.Title,
		Console:      container.Console,
		Icon:         container.Icon,
		BoxArt:       container.BoxArt,
		Softcore:     container.Softcore,
		Hardcore:     container.Hardcore,
		Achievements: map[string]achievement{},
	}
	for _, v := range container.Achievements {
		result.Achievements[v.ID] = v
	}
	return result, nil
}

func getCurrentGame(username, apikey string) (game, error) {
	current, err := getLastGame(username, apikey)
	if err != nil {
		return game{}, err
	}

	result, err := getGameInformation(username, apikey, current)
	if err != nil {
		return game{}, err
	}

	return result, nil
}
