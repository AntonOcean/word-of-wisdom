package app

//go:generate mockery --name=powChallenge --filename pow_challenge.go --exported --with-expecter=True
//go:generate mockery --name=quoteProvider --filename quote_provider.go --exported --with-expecter=True
//go:generate mockery --name=Conn --filename conn.go --exported --with-expecter=True

import "net"

type (
	Handler interface {
		HandleConnection(conn Conn) error
	}

	Conn interface {
		net.Conn
	}

	powChallenge interface {
		GenerateChallenge() string
		ValidateChallenge(challenge, response string) bool
	}

	quoteProvider interface {
		GetQuote() string
	}
)
