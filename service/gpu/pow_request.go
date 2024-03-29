package gpu

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	requestUrl string
)

func InitUrl(ip string) {
	requestUrl = "http://" + ip
}

type workGenerate struct {
	Action    string `json:"action"`
	DataHash  string `json:"hash"`
	Threshold string `json:"threshold"`
}

type workGenerateResult struct {
	Work string `json:"work"`
}

func GenerateWork(dataHash string, threshold string) (*string, error) {
	wg := &workGenerate{
		Action:    "work_generate",
		Threshold: threshold,
		DataHash:  dataHash,
	}
	bytesData, err := json.Marshal(wg)
	if err != nil {
		return nil, err
	}

	body, err := httpRequest(bytesData)
	if err != nil {
		return nil, err
	}

	workResult := &workGenerateResult{}
	if err := json.Unmarshal(body, workResult); err != nil {
		return nil, err
	}
	return &workResult.Work, nil
}

type workValidate struct {
	Action    string `json:"action"`
	DataHash  string `json:"hash"`
	Threshold string `json:"threshold"`
	Work      string `json:"work"`
}

type workValidateResult struct {
	Valid string `json:"valid"`
}

func VaildateWork(dataHash string, threshold string, work string) (bool, error) {
	wg := &workValidate{
		Action:    "work_validate",
		DataHash:  dataHash,
		Threshold: threshold,
		Work:      work,
	}
	bytesData, err := json.Marshal(wg)
	if err != nil {
		return false, err
	}

	body, err := httpRequest(bytesData)
	if err != nil {
		return false, err
	}

	validateResult := &workValidateResult{}
	if err := json.Unmarshal(body, validateResult); err != nil {
		return false, err
	}
	if validateResult.Valid == "1" {
		return true, nil
	} else {
		return false, nil
	}
}

type workCancel struct {
	Action   string `json:"work_cancel"`
	DataHash string `json:"hash"`
}

func CancelWork(dataHash string) error {
	wg := &workCancel{
		Action:   "work_cancel",
		DataHash: dataHash,
	}
	bytesData, err := json.Marshal(wg)
	if err != nil {
		return err
	}

	if _, err := httpRequest(bytesData); err != nil {
		return err
	}
	return nil
}

func httpRequest(bytesData []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(bytesData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Info("request failed", "error", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}
