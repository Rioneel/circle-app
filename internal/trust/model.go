package trust

import "time"

type TrustEdge struct {
	FromId string
	ToId string
	Weight float64
	CreatedAt time.Time
	// UpdatedAt time.Time
}

type TrustedUser struct {
	ID string 
	Name string 
	Weight float64
}

