package openia

import (
	"context"
	"encoding/json"
	"fmt"

	"saas-medico/internal/config"
	models "saas-medico/internal/modules/nutricion/models"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type OpenIaService struct {
	Client *openai.Client
}
type APIResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewOpenIaService() *OpenIaService {
	cfg := config.AppConfig

	client := openai.NewClient(
		option.WithAPIKey(cfg.OpenIaKey),
	)
	return &OpenIaService{Client: &client}
}

func (ps *OpenIaService) AskModelIa(Alimentos []models.NutricionAlimento, prompt string, dieta []models.NutricionDietaPaciente) (*responses.Response, error) {
	alimentosJson, err := json.Marshal(Alimentos)
	dietasJson, err := json.Marshal(dieta)
	if err != nil {
		return nil, err
	}

	instructions := `
		Eres un nutriólogo experto en nutrición.

		Tu función NO es crear planes alimenticios completos ni menús estructurados.

		Tu función es:
		- conocer en detalle los alimentos proporcionados por el sistema,
		- responder preguntas nutricionales del usuario,
		- analizar calorías y macronutrientes,
		- sugerir qué podría comer el usuario basándote únicamente en los alimentos disponibles.

		Reglas importantes:
		- Debes saludar al usuario por el nombre que esta en las dietas.
		- SOLO puedes usar los alimentos proporcionados.
		- NO inventes alimentos.
		- Si un alimento existe pero no tiene valores nutricionales completos, indícalo claramente.
		- Puedes sugerir combinaciones simples de alimentos, pero NO menús completos.
		- Responde siempre en español.
		- Sé claro, directo y práctico.
		- Si faltan datos del usuario, indícalo.
		- No reemplazas a un médico.

		Si el usuario pregunta "¿qué puedo comer?", responde con sugerencias basadas únicamente en los alimentos disponibles.
		`
	finalPrompt := fmt.Sprintf(`
		Pregunta del usuario:
		%s

		Lista de alimentos disponibles:
		%s

		Dieta e informacion del Paciente:
		%s
		`, prompt, string(alimentosJson), string(dietasJson))
	resp, err := ps.Client.Responses.New(
		context.Background(),
		responses.ResponseNewParams{
			Model:        "gpt-5.4",
			Instructions: openai.String(instructions),
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(finalPrompt),
			},
		})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
