package util

import (
	"io/ioutil"
	"net/http"
	"ruizi/pkg/logger"

	"github.com/avast/retry-go"
)

func RetryGet(url string) ([]byte, error) {
	var body []byte
	err := retry.Do(
		func() error {
			resp, err := http.Get(url)
			if err == nil {
				defer func() {
					if err := resp.Body.Close(); err != nil {
						panic(err)
					}
				}()
				body, err = ioutil.ReadAll(resp.Body)
			}

			return err
		},
		retry.Attempts(3),
		retry.OnRetry(func(n uint, err error) {
			logger.Sugar.Info("%s retry #%d: %s", url, n, err)
		}))
	if err != nil {
		return nil, err
	}
	return body, nil
}
