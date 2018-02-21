package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nlepage/hashcode-cli/config"
)

const (
	tokenEnvVar = "HASHCODE_TOKEN"
	submitURL   = "https://hashcode-judge.appspot.com/api/judge/v1/submissions"
	createURL   = "https://hashcode-judge.appspot.com/api/judge/v1/upload/createUrl"
)

func upload(args []string) error {
	token, err := config.Token()
	if err != nil {
		return err
	}

	datasets, err := config.Datasets()
	if err != nil {
		return err
	}

	sourceFile, err := archiveSource()
	if err != nil {
		return err
	}

	sourceKey, err := uploadFile(token, sourceFile)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		for name, id := range datasets {
			if err = submitDataset(token, sourceKey, name, id); err != nil {
				return err
			}
		}
	} else {
		for _, name := range args {
			id, ok := datasets[name]
			if !ok {
				return fmt.Errorf("Unknown dataset %s", name)
			}

			if err = submitDataset(token, sourceKey, name, id); err != nil {
				return err
			}
		}
	}

	return nil
}

func submitDataset(token string, sourceKey string, name string, id string) error {
	fmt.Printf("Submitting dataset %s...\n", name)

	fName := filepath.Join(config.DatasetsDir(), name+".txt")

	if _, err := os.Stat(fName); os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", fName)
		return nil
	}

	key, err := uploadFile(token, fName)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", submitURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Set("dataSet", id)
	q.Set("sourcesBlobKey", sourceKey)
	q.Set("submissionBlobKey", key)
	req.URL.RawQuery = q.Encode()

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, res.Body)
		return fmt.Errorf("Error while submitting dataset %s", name)
	}

	fmt.Printf("Submitted dataset %s...\n", name)

	return nil
}

type createPayload struct {
	Value string `json:"value"`
}

type uploadPayload struct {
	File []string `json:"file"`
}

func uploadFile(token string, name string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", createURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, res.Body)
		return "", fmt.Errorf("Error while uploading dataset %s", name)
	}

	createPayload := createPayload{}
	err = json.NewDecoder(res.Body).Decode(&createPayload)
	if err != nil {
		return "", err
	}

	uploadURL := createPayload.Value

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return "", err
	}
	w.Close()

	req, err = http.NewRequest("POST", uploadURL, &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err = client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, res.Body)
		return "", fmt.Errorf("Error while uploading dataset %s", name)
	}

	uploadPayload := uploadPayload{}
	err = json.NewDecoder(res.Body).Decode(&uploadPayload)
	if err != nil {
		return "", err
	}

	return uploadPayload.File[0], nil
}
