// Cliente HTTP simples para leitura de JSON de dependencias do edge.
// O gateway usa esta camada para compor respostas sem acoplar handlers a detalhes de transporte.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONReader interface {
	GetJSON(ctx context.Context, requestURL string, target any) error
}

func (checker *HTTPHealthChecker) GetJSON(ctx context.Context, requestURL string, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	response, err := checker.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected_status:%d", response.StatusCode)
	}

	return json.NewDecoder(response.Body).Decode(target)
}
