package courier

import (
	"net/http"
)

func sendPatch(url string) error {
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
