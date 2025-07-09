// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

type Message struct {
	Content string `json:"content"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ToolMetadata struct {
	Sources []string `json:"sources"`
}

type AIResponse struct {
	Choices      []Choice     `json:"choices"`
	ToolMetadata ToolMetadata `json:"tool_metadata"`
}

type RespOutput struct {
	Explanation string   `json:"explanation"`
	Sources     []string `json:"sources"`
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/{pipeline}/explain pipelines ExplainPipeline
//
// Explain a pipeline
//
// ---
// produces:
// - application/yaml
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repository
//   required: true
//   type: string
// - in: path
//   name: pipeline
//   description: Commit SHA for pipeline to retrieve
//   required: true
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format
//   type: string
//   default: yaml
//   enum:
//   - json
//   - yaml
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved and explained the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/PipelineBuild"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// ExplainPipeline represents the API handler to capture and
// explain a pipeline configuration.
func ExplainPipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	p := pipeline.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	var prompt string

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), p.GetCommit())

	l.Debugf("explaining templates for pipeline %s", entry)

	// ensure we use the expected pipeline type when compiling
	r.SetPipelineType(p.GetType())

	// create the compiler object
	compiler := compiler.FromContext(c).Duplicate().WithCommit(p.GetCommit()).WithMetadata(m).WithRepo(r).WithUser(u)

	ruleData := prepareRuleData(c)

	if ruleData != nil {
		rdBytes, err := yaml.Marshal(ruleData)
		if err != nil {
			retErr := fmt.Errorf("unable to explain pipeline %s: %w", entry, err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		prompt = fmt.Sprintf("RULEDATA: %s", string(rdBytes))
	}

	// explain the templates in the pipeline
	pipeline, _, err := compiler.CompileLite(ctx, p.GetData(), ruleData, false)
	if err != nil {
		retErr := fmt.Errorf("unable to explain pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	pBytes, err := yaml.Marshal(pipeline)
	if err != nil {
		retErr := fmt.Errorf("unable to explain pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if len(prompt) > 0 {
		prompt = fmt.Sprintf("%s\n%s", string(pBytes), prompt)
	} else {
		prompt = string(pBytes)
	}

	aiPayload := map[string]interface{}{
		"model":             "gpt-4o-mini-paygo",
		"temperature":       1,
		"max_new_tokens":    4095,
		"top_p":             1,
		"frequency_penalty": 0.5,
		"presence_penalty":  0,
		"timeout":           120,
		"tools": []map[string]interface{}{
			{
				"type": "rag_search",
				"rag_search": map[string]interface{}{
					"index": "bullsai",
					"top_k": 5,
				},
			},
		},
		"stream": false,
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `You are a question and answer bot designed to receive pipelines as a prompt followed 
				OPTIONALLY by the ruledata (or context) for the pipeline to run, such as 'push to main' or 'tag v1.0'. 
				Your job is to tell the user what the pipeline does in simple and precise terms. Be as detailed as possible. 
				If you don't know the answer, you can say 'I don't know'. Format your response with bullets and lists and paragraphs, 
				but do NOT use other markdown styling such as bold (surrounding words with **) or italics (surrounding words with _).
				If you find ways to improve the pipeline, mention those improvements in a "suggestion" section at the bottom of the explanation.`,
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	payloadBytes, err := json.Marshal(aiPayload)
	if err != nil {
		retErr := fmt.Errorf("unable to explain pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	aiClient := http.DefaultClient
	aiClient.Timeout = 30 * time.Second

	// set the API endpoint path we send the request to
	url := os.Getenv("OPENAI_API_URL")

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		retErr := fmt.Errorf("unable to form a request to %s: %w", u, err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	req.Header.Set("Authorization", "Bearer <INSERT OPEN AI TOKEN HERE>")
	req.Header.Set("Content-Type", "application/json")

	resp, err := aiClient.Do(req)
	if err != nil {
		retErr := fmt.Errorf("unable to explain pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			retErr := fmt.Errorf("unable to explain pipeline %s: %w", entry, err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		retErr := fmt.Errorf("unable to explain pipeline %s: %w", entry, body)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		retErr := fmt.Errorf("unable to explain pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	aiResp := new(AIResponse)

	err = json.Unmarshal(responseBody, aiResp)
	if err != nil {
		retErr := fmt.Errorf("unable to explain pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	respOutput := new(RespOutput)
	respOutput.Explanation = fmt.Sprintf("%s", aiResp.Choices[0].Message.Content)
	respOutput.Sources = aiResp.ToolMetadata.Sources

	c.JSON(http.StatusOK, respOutput)
}
