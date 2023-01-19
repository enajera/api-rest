package models

type SearchRequest struct {
	Index string `json:"index"`
	SearchType string `json:"search_type"`
	Query      struct {
		Term string `json:"term"`
	} `json:"query"`
	From         int      `json:"from"`
	MaxResults   int      `json:"max_results"`
	
}

type Request struct {
	SearchType string `json:"search_type"`
	Query      struct {
		Term string `json:"term"`
	} `json:"query"`
	From         int      `json:"from"`
	MaxResults   int      `json:"max_results"`
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
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		// MaxScore float64 `json:"max_score"`
		Hit []struct {
			// Index  string  `json:"_index"`
			// Type   string  `json:"_type"`
			// ID     string  `json:"_id"`
			// Score  float64 `json:"_score"`
			Source struct {
				Body string `json:"Body"`
				// ContentTransferEncoding string `json:"ContentTransferEncoding"`
				// ContentType             string `json:"ContentType"`
				Date      string `json:"Date"`
				From      string `json:"From"`
				MessageID string `json:"MessageID"`
				// MimeVersion             string `json:"MimeVersion"`
				Subject string `json:"Subject"`
				To      string `json:"To"`
				// XBcc                    string `json:"XBcc"`
				// XCc                     string `json:"XCc"`
				// XFileName               string `json:"XFileName"`
				// XFolder string `json:"XFolder"`
				// XFrom                   string `json:"XFrom"`
				// XOrigin                 string `json:"XOrigin"`
				// XTo                     string `json:"XTo"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type Result struct {
	Body      string `json:"Body"`
	To        string `json:"To"`
	From      string `json:"From"`
	Subject   string `json:"Subject"`
	// Folder    string `json:"Folder"`
	MessageID string `json:"MessageID"`
	Date      string `json:"Date"`
	
}

type Results struct {
    Total int        `json:"total"`
    Data  []Result `json:"data"`
}

func EmailFields(res Response) []Result {

	var results []Result

	for _, hit := range res.Hits.Hit {
		result := Result{
			Body:      hit.Source.Body,
			To:        hit.Source.To,
			From:      hit.Source.From,
			Subject:   hit.Source.Subject,
			// Folder:    hit.Source.XFolder,
			MessageID: hit.Source.MessageID,
			Date: hit.Source.Date,
		}
		results = append(results, result)
	}

	return results
}

type Index struct {
    List []struct {
        Name string `json:"name"`
    }`json:"list"`
}

type Login struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginResponse struct {
	Success bool `json:"success"`
}
