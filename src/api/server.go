package api

import (
	"encoding/json"
	"net/http"

	"vineelsai.com/rce/src/docker"
	"vineelsai.com/rce/src/utils"
)

type Payload struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

func Serve(PORT string) {
	// handle the /run endpoint
	http.HandleFunc("/run", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			var payload Payload

			err := json.NewDecoder(req.Body).Decode(&payload)
			if err != nil {
				panic(err)
			}

			filePath, runId := utils.CreateFile(payload.Code, payload.Language, "runs")

			// convert the output to json
			output, err := json.Marshal(map[string]string{
				"output": string(docker.Run(filePath, payload.Language, runId)),
				"runId":  runId,
			})
			if err != nil {
				panic(err)
			}

			// set the content type to json, enable CORS and write the output
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("Access-Control-Allow-Origin", "*")
			res.Write(output)
		case http.MethodGet:
			res.Write([]byte("Hello World"))
		}
	})

	// start the server
	if err := http.ListenAndServe("0.0.0.0:"+PORT, nil); err != nil {
		panic(err)
	}
}
