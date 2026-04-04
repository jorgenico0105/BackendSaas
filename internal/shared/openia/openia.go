package openia

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"saas-medico/internal/config"
	models "saas-medico/internal/modules/nutricion/models"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/redis/go-redis/v9"
)

type OpenIaService struct {
	Client      *openai.Client
	RedisClient *redis.Client
}

// 👇 Mensaje individual tipo chat
type ChatMessage struct {
	Role      string    `json:"role"`      // "user" | "assistant"
	Content   string    `json:"content"`   // mensaje
	CreatedAt time.Time `json:"createdAt"` // timestamp
}

func NewOpenIaService(redisClient *redis.Client) *OpenIaService {
	cfg := config.AppConfig

	client := openai.NewClient(
		option.WithAPIKey(cfg.OpenIaKey),
	)

	return &OpenIaService{
		Client:      &client,
		RedisClient: redisClient,
	}
}

func BuildConversationKey(patientID uint) string {
	return "conv:" + strconv.FormatUint(uint64(patientID), 10)
}

func AppendMessage(
	redisClient *redis.Client,
	ctx context.Context,
	key string,
	role string,
	content string,
) error {
	msg := ChatMessage{
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return redisClient.RPush(ctx, key, data).Err()
}

func GetConversation(
	redisClient *redis.Client,
	ctx context.Context,
	key string,
) ([]ChatMessage, error) {
	values, err := redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]ChatMessage, 0, len(values))

	for _, value := range values {
		var msg ChatMessage
		if err := json.Unmarshal([]byte(value), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func TrimConversation(
	redisClient *redis.Client,
	ctx context.Context,
	key string,
	max int64,
) error {
	return redisClient.LTrim(ctx, key, -max, -1).Err()
}

func ClearConversation(
	redisClient *redis.Client,
	ctx context.Context,
	key string,
) error {
	return redisClient.Del(ctx, key).Err()
}

func (ps *OpenIaService) AskModelIa(
	alimentos []models.NutricionAlimento,
	prompt string,
	dieta []models.NutricionDietaPaciente,
) ([]ChatMessage, error) {

	ctx := context.Background()

	if len(dieta) == 0 {
		return nil, fmt.Errorf("no hay dieta")
	}

	alimentosJSON, _ := json.Marshal(alimentos)
	dietasJSON, _ := json.Marshal(dieta)

	instructions := `
Eres un nutriólogo experto en nutrición.

- NO creas menús completos
- SOLO usa alimentos dados
- Responde en español
- Sé claro y práctico
`

	finalPrompt := fmt.Sprintf(`
Pregunta:
%s

Alimentos:
%s

Paciente:
%s
`, prompt, string(alimentosJSON), string(dietasJSON))
	key := BuildConversationKey(dieta[0].Paciente.ID)

	err := AppendMessage(ps.RedisClient, ctx, key, "user", prompt)
	if err != nil {
		return nil, err
	}

	resp, err := ps.Client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model:        "gpt-5.4",
			Instructions: openai.String(instructions),
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(finalPrompt),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	answer := resp.OutputText()

	err = AppendMessage(ps.RedisClient, ctx, key, "assistant", answer)
	if err != nil {
		return nil, err
	}
	_ = TrimConversation(ps.RedisClient, ctx, key, 50)

	_ = ps.RedisClient.Expire(ctx, key, 24*time.Hour).Err()

	return GetConversation(ps.RedisClient, ctx, key)
}
