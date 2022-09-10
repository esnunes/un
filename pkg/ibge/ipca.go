package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DefaultBaseURL = "https://servicodados.ibge.gov.br/api/v3"

	YearMonthLayout = "200601"
)

type Serie struct {
	Date  time.Time
	Value float64
}

type Client struct {
	BaseURL string
	Context context.Context
	Log     *logrus.Logger
	HTTP    *http.Client
}

func (c *Client) IPCA(dateFrom, dateTo time.Time) ([]Serie, error) {
	if dateTo.IsZero() {
		c.Log.Infof("ibge: retrieving IPCA index for %s", dateFrom)
	} else {
		c.Log.Infof("ibge: retrieving IPCA index from %s to %s", dateFrom, dateTo)
	}

	path := fmt.Sprintf(
		"/agregados/1737/periodos/%s/variaveis/63?localidades=N1[all]",
		yearMonthSequence(dateFrom, dateTo),
	)
	c.Log.Infof("ibge: making http request to: %v", path)

	body, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	var data []struct {
		Resultados []struct {
			Series []struct {
				Serie map[string]string
			}
		}
	}
	if err := json.Unmarshal(body, &data); err != nil {
		c.Log.Errorf("ibge: failed to unmarshal response: %v", err)
		return nil, err
	}

	series := []Serie{}
	if len(data) == 0 || len(data[0].Resultados) == 0 || len(data[0].Resultados[0].Series) == 0 {
		return series, nil
	}

	for k, v := range data[0].Resultados[0].Series[0].Serie {
		date, err := time.Parse(YearMonthLayout, k)
		if err != nil {
			c.Log.Errorf("ibge: failed to parse serie date: %v", err)
			return nil, err
		}
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			c.Log.Errorf("ibge: failed to parse serie value: %v", err)
			return nil, err
		}
		series = append(series, Serie{
			Date:  date,
			Value: value,
		})
	}
	sort.Slice(series, func(i, j int) bool {
		return series[i].Date.Before(series[j].Date)
	})

	return series, nil
}

func (c *Client) Get(path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, c.baseURL()+path, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client().Do(req)
	if err != nil {
		return nil, err
	}
	if res.Body == nil {
		return nil, nil
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (c *Client) baseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return DefaultBaseURL
}

func (c *Client) context() context.Context {
	if c.Context != nil {
		return c.Context
	}
	return context.Background()
}

func (c *Client) client() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

func yearMonthSequence(dateFrom, dateTo time.Time) string {
	if dateTo.IsZero() {
		dateTo = dateFrom
	}
	b := strings.Builder{}
	for d := dateFrom; d.Before(dateTo); d = d.AddDate(0, 1, 0) {
		b.WriteString(d.Format(YearMonthLayout))
		b.WriteString("|")
	}
	b.WriteString(dateTo.Format(YearMonthLayout))
	return b.String()
}
