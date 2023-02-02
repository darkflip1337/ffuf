package ffuf

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ConfigOptionsHistory struct {
	ConfigOptions
	Time time.Time `json:"time"`
}

func WriteHistoryEntry(conf *Config) (string, error) {
	options := ConfigOptionsHistory{
		ConfigOptions: conf.ToOptions(),
		Time:          time.Now(),
	}
	jsonoptions, err := json.Marshal(options)
	if err != nil {
		return "", err
	}
	hashstr := calculateHistoryHash(jsonoptions)
	err = createConfigDir(filepath.Join(HISTORYDIR, hashstr))
	if err != nil {
		return "", err
	}
	err = os.WriteFile(filepath.Join(HISTORYDIR, hashstr, "options"), jsonoptions, 0640)
	return hashstr, err
}

func calculateHistoryHash(options []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(options))
}

func SearchHash(hash string) ([]ConfigOptionsHistory, int, error) {
	coptions := make([]ConfigOptionsHistory, 0)
	if len(hash) < 6 {
		return coptions, 0, errors.New("bad FFUFHASH value")
	}
	historypart := hash[0:5]
	position, err := strconv.ParseInt(hash[5:], 16, 32)
	if err != nil {
		return coptions, 0, errors.New("bad positional value in FFUFHASH")
	}
	all_dirs, err := os.ReadDir(HISTORYDIR)
	if err != nil {
		return coptions, 0, err
	}
	matched_dirs := make([]string, 0)
	for _, filename := range all_dirs {
		if filename.IsDir() {
			if strings.HasPrefix(strings.ToLower(filename.Name()), strings.ToLower(historypart)) {
				matched_dirs = append(matched_dirs, filename.Name())
			}
		}
	}
	for _, dirname := range matched_dirs {
		copts, err := configFromHistory(filepath.Join(HISTORYDIR, dirname))
		if err != nil {
			continue
		}
		coptions = append(coptions, copts)

	}
	return coptions, int(position), err
}

func configFromHistory(dirname string) (ConfigOptionsHistory, error) {
	jsonOptions, err := os.ReadFile(filepath.Join(dirname, "options"))
	if err != nil {
		return ConfigOptionsHistory{}, err
	}
	tmpOptions := ConfigOptionsHistory{}
	err = json.Unmarshal(jsonOptions, &tmpOptions)
	return tmpOptions, err
	/*
		// These are dummy values for this use case
		ctx, cancel := context.WithCancel(context.Background())
		conf, err := ConfigFromOptions(&tmpOptions.ConfigOptions, ctx, cancel)
		job.Input, errs = input.NewInputProvider(conf)
		return conf, tmpOptions.Time, err
	*/
}
