package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"hackathon/db"
	"hackathon/model"
)

type GeminiRequest struct {
	Contents []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}


type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func FactCheckHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDStr := vars["id"]
	log.Println("🔍 postIDStr:", postIDStr)

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("❌ Invalid post ID:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post model.Post
	if err := db.GetDB().Preload("User").First(&post, postID).Error; err != nil {
		log.Println("❌ Post not found in DB:", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	log.Printf("✅ DB post loaded: %+v\n", post)

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("❌ GEMINI_API_KEY not set in env")
		http.Error(w, "GEMINI_API_KEY not set", http.StatusInternalServerError)
		return
	}

	url := "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-pro-002:generateContent?key=" + apiKey



	prompt := fmt.Sprintf(`
	You are an AI fact-checking assistant.
	Your job is to analyze the factual accuracy of a given social media post.
	Use a concise, witty, and informative tone — similar to the style of X’s Grok AI.

	For the post below:
	「%s」

	Please:
	1. Indicate whether the main claim is likely true, false, or uncertain.
	2. Provide a brief explanation with evidence or reasoning.
	3. Add a touch of personality (e.g. light sarcasm or humor) while staying factual.
	4. Write your answer in Japanese.

	Format:
	- ✅ 正しい場合 → 「○ 本当です。理由：...」
	- ❌ 間違っている場合 → 「× 間違いです。理由：...」
	- ❓ 判断が難しい場合 → 「？ 判断が分かれます。理由：...」
	`, post.Content)

	log.Println("📤 Prompt sending to Gemini:\n", prompt)

	payload := GeminiRequest{
	Contents: []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{
		{
			Role: "user",
			Parts: []struct {
				Text string `json:"text"`
			}{
				{Text: prompt},
			},
		},
	},
}


	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("❌ Failed to marshal payload:", err)
		http.Error(w, "Failed to encode request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("❌ Failed to call Gemini API:", err)
		http.Error(w, "Failed to call Gemini API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("❌ Failed to read Gemini response:", err)
		http.Error(w, "Failed to read Gemini response", http.StatusInternalServerError)
		return
	}

	log.Println("📩 Gemini raw response:", string(body))

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		log.Println("❌ Failed to unmarshal Gemini response:", err)
		http.Error(w, "Failed to parse Gemini response", http.StatusInternalServerError)
		return
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		log.Println("❌ Gemini response missing candidates/parts")
		http.Error(w, "No analysis result", http.StatusInternalServerError)
		return
	}

	result := geminiResp.Candidates[0].Content.Parts[0].Text
	log.Println("✅ Gemini result:", result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"factcheck": result})
}
