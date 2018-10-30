package service

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/vitelabs/go-vite/common/helper"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/log15"
	"io/ioutil"
	"net/http"
)

var (
	powClientLog = log15.New("module", "pow_request")
	requestUrl   string
)

func InitUrl(url string) {
	requestUrl = url
}

func GetData(addr types.Address, preHash types.Hash) types.Hash {
	return types.DataListHash(addr.Bytes(), preHash.Bytes())
}

type workGenerate struct {
	Action    string `json:"action"`
	DataHash  string `json:"hash"`
	Threshold string `json:"threshold"`
}

type workGenerateResult struct {
	Work []byte `json:"work"`
}

func GenerateWork(dataHash string, threshold string) ([]byte, error) {
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
	return workResult.Work, nil
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

func VaildateWork(dataHash string, threshold string, work []byte) (bool, error) {
	wg := &workValidate{
		Action:    "work_validate",
		DataHash:  dataHash,
		Threshold: threshold,
		Work:      hex.EncodeToString(work),
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
		return nil, err
	}
	powClientLog.Info("Response Status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(helper.BytesToString(body))
	return body, nil
}
