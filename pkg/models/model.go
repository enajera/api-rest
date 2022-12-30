package models

type Request struct {
	SearchType string `json:"search_type"`
	Query      struct {
		Term string `json:"term"`
	} `json:"query"`
	From         int      `json:"from"`
	MaxResults   int      `json:"max_results"`
	SourceFields []string `json:"_source"`
}

type Response struct {
	// Took     int  `json:"took"`
	// TimedOut bool `json:"timed_out"`
	// Shards   struct {
	// 	Total      int `json:"total"`
	// 	Successful int `json:"successful"`
	// 	Skipped    int `json:"skipped"`
	// 	Failed     int `json:"failed"`
	// } `json:"_shards"`
	Hits struct {
		// Total struct {
		// 	Value int `json:"value"`
		// } `json:"total"`
		// MaxScore float64 `json:"max_score"`
		Hits []struct {
			// Index  string  `json:"_index"`
			// Type   string  `json:"_type"`
			// ID     string  `json:"_id"`
			// Score  float64 `json:"_score"`
			Source struct {
				Body string `json:"Body"`
				// ContentTransferEncoding string `json:"ContentTransferEncoding"`
				// ContentType             string `json:"ContentType"`
				Date string `json:"Date"`
				From string `json:"From"`
				// MessageID               string `json:"MessageID"`
				// MimeVersion             string `json:"MimeVersion"`
				Subject string `json:"Subject"`
				To      string `json:"To"`
				// XBcc                    string `json:"XBcc"`
				// XCc                     string `json:"XCc"`
				// XFileName               string `json:"XFileName"`
				XFolder string `json:"XFolder"`
				// XFrom                   string `json:"XFrom"`
				// XOrigin                 string `json:"XOrigin"`
				// XTo                     string `json:"XTo"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type Result struct {
	Body    string `json:"Body"`
	To      string `json:"To"`
	From    string `json:"From"`
	Subject string `json:"Subject"`
	Folder string `json:"Folder"`
}

func EmailFields(res Response) []Result {

	var results []Result

	for _, hit := range res.Hits.Hits {
		result := Result{
			Body:    hit.Source.Body,
			To:      hit.Source.To,
			From:    hit.Source.From,
			Subject: hit.Source.Subject,
			Folder: hit.Source.XFolder,
		}
		results = append(results, result)
	}

	return results
}
