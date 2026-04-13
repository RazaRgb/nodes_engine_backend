package engine

import (
	"context"
	"fmt"
	"google.golang.org/genai"
)

func geminiService(systemPrompt string, userPrompt string, timeout float64) (string, error) {

	timeout = min(timeout, 1800)
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		fmt.Printf("error creating client %+v\n", err)
		return "Error generating text", err
	}

	config := &genai.GenerateContentConfig{
		//Temperature: genai.Ptr[float32](0.1),
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{
					Text: systemPrompt,
				},
			},
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(userPrompt),
		config,
	)
	if err != nil {
		fmt.Printf("error connecting client %+v\n", err)
		return "Error generating text", err
	}

	fmt.Println(result.Text())
	return result.Text(), nil
}
