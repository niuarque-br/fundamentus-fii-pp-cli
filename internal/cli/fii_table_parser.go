// Copyright 2026. Extrai dados estruturados da resposta do Crawl4AI.
//
// O Crawl4AI retorna JSON com a estrutura:
//
//	{
//	  "success": true,
//	  "results": [{
//	    "tables": [{
//	      "headers": ["Papel", "Segmento", ...],
//	      "rows": [["AAZQ11", "Outros", ...], ...]
//	    }]
//	  }]
//	}
//
// Este pacote extrai as linhas da tabela e converte para []map[string]any
// para ser usado pelo pipeline de output do CLI.

package cli

import (
	"encoding/json"
	"fmt"
)

// Crawl4AIResponse representa a estrutura completa da resposta do Crawl4AI.
type Crawl4AIResponse struct {
	Success bool               `json:"success"`
	Results []Crawl4AIResult   `json:"results"`
}

// Crawl4AIResult representa um resultado individual do Crawl4AI.
type Crawl4AIResult struct {
	Tables []CrawledTable `json:"tables"`
}

// CrawledTable representa uma tabela extraída pelo Crawl4AI.
type CrawledTable struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

// ExtractFiiTableRows extrai as linhas da tabela de FIIs da resposta JSON
// do Crawl4AI e retorna uma lista de mapas (header → valor) prontos para
// serem usados pelo pipeline de output (printAutoTable, printOutputWithFlags, etc.).
func ExtractFiiTableRows(data []byte) ([]map[string]any, error) {
	var resp Crawl4AIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing Crawl4AI response: %w", err)
	}
	if !resp.Success {
		return nil, fmt.Errorf("Crawl4AI returned success=false")
	}
	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("Crawl4AI returned no results")
	}

	tables := resp.Results[0].Tables
	if len(tables) == 0 {
		return nil, fmt.Errorf("no tables found in Crawl4AI response")
	}

	table := tables[0]
	headers := table.Headers
	rows := table.Rows

	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		item := make(map[string]any, len(headers))
		for i, header := range headers {
			if i < len(row) {
				item[header] = row[i]
			} else {
				item[header] = ""
			}
		}
		items = append(items, item)
	}

	return items, nil
}

// IsCrawl4AITableResponse verifica se o JSON parece uma resposta com tabelas do Crawl4AI.
func IsCrawl4AITableResponse(data []byte) bool {
	var probe struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return false
	}
	return probe.Success
}
