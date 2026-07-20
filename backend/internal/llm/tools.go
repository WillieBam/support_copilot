package llm

import "encoding/json"

// Tool definition schema (OpenAI format)
type FunctionDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type Tool struct {
	Type     string             `json:"type"` // "function"
	Function FunctionDefinition `json:"function"`
}

// tool exposed to LLM
func GetValidateAlertTool() Tool {
	params := `{                                                                                                                                                          
            "type": "object",                                                                                                                                                     
            "properties": {                                                                                                                                                       
                "alert_id": {                                                                                                                                                         
                    "type": "string",                                                                                                                                                     
                    "description": "Unique identifier of the alert to validate"                                                                              
                }                                                                                                                                                                     
            },                                                                                                                                                                    
            "required": ["alert_id"]                                                                                                                                              
        }`

	return Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "validate_alert",
			Description: "Validates system alert by ID, retrieves telemetry metrics, and runs anomaly prediction model.",
			Parameters:  json.RawMessage(params),
		},
	}
}
