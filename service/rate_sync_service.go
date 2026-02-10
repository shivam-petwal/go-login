package service

import (
	"encoding/json"
	"fmt"
	"go-login/repository"
	"log"
	"net/http"
	"time"
)


type RateSyncService struct {
	exchangeRateRepo repository.ExchangeRateRepository
	client           *http.Client
}


func NewRateSyncService(repo repository.ExchangeRateRepository) *RateSyncService {
	return &RateSyncService{
		exchangeRateRepo: repo,
		client:           &http.Client{Timeout: 10 * time.Second},
	}
}

//   { "base":"USD", "date":"...", "rates":{"EUR":0.92} }
type frankfurterResponse struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

func (s *RateSyncService) SyncAll() {
	rates, err := s.exchangeRateRepo.GetAllActive()
	if err != nil {
		log.Printf("[RateSync] failed to fetch active pairs: %v", err)
		return
	}
	if len(rates) == 0 {
		log.Println("[RateSync] no active pairs found, skipping")
		return
	}

	log.Printf("[RateSync] syncing %d pair(s)...", len(rates))

	var updated, failed int
	for _, r := range rates {
		from, to := r.FromCurrency.Code, r.ToCurrency.Code

		latest, err := s.fetchRate(from, to)
		if err != nil {
			log.Printf("[RateSync] %s→%s: %v", from, to, err)
			failed++
			continue
		}

		old := r.Rate
		r.Rate = latest
		r.UpdatedAt = time.Now()

		if err := s.exchangeRateRepo.Update(&r); err != nil {
			log.Printf("[RateSync] %s→%s: db update failed: %v", from, to, err)
			failed++
			continue
		}

		log.Printf("[RateSync] %s→%s: %.6f → %.6f", from, to, old, latest)
		updated++
	}

	log.Printf("[RateSync] done: %d updated, %d failed", updated, failed)
}


func (s *RateSyncService) fetchRate(from, to string) (float64, error) {
	url := fmt.Sprintf("https://api.frankfurter.app/latest?from=%s&to=%s", from, to)

	resp, err := s.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("api call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	var data frankfurterResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("decode failed: %w", err)
	}

	rate, ok := data.Rates[to]
	if !ok || rate <= 0 {
		return 0, fmt.Errorf("invalid rate for %s in response", to)
	}

	return rate, nil
}
