package main

import (
	"time"
	"word-of-wisdom/internal/app"
	"word-of-wisdom/internal/config"
	"word-of-wisdom/internal/pow"
	"word-of-wisdom/internal/quotes"
	"word-of-wisdom/pkg/logger"
)

func main() {
	cfg := config.Config{
		Port:              ":9000",
		MaxConnections:    100,
		ConnectionTimeout: 2 * time.Second,
		ShutdownTimeout:   5 * time.Second,
	}

	s := app.NewServer(
		cfg,
		logger.GetLogger(),
		app.NewHandler(
			quotes.NewRandomQuoteProvider([]string{
				"We are not what we know but what we are willing to learn.",
				"Good people are good because they've come to wisdom through failure.",
				"Your word is a lamp for my feet, a light for my path.",
				"The first problem for all of us, men and women, is not to learn, but to unlearn.",
				"The only limit to our realization of tomorrow is our doubts of today.",
				"Do what you can, with what you have, where you are.",
				"The journey of a thousand miles begins with one step.",
				"Opportunities don't happen. You create them.",
			}),
			pow.NewSHA256PoW(4),
		),
	)

	s.Start()
}
