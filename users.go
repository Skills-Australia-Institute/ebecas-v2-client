package ebecasv2client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// User represents an eBECAS user.
type User struct {
	ID        int64  `json:"id"`
	Initials  string `json:"initials"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Me gets the details of the API user making the request.
func (c *Client) Me(ctx context.Context) (User, int, error) {
	var user User

	requestURL := fmt.Sprintf("%s/users/me", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return user, http.StatusInternalServerError, fmt.Errorf("create get me request: %w", err)
	}

	data, statusCode, err := c.do(req)
	if err != nil {
		return user, statusCode, fmt.Errorf("get me, status %d: %w", statusCode, err)
	}

	if statusCode != http.StatusOK {
		return user, statusCode, fmt.Errorf(
			"get me returned status %d: %s",
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	if err := json.Unmarshal(data, &user); err != nil {
		return user, statusCode, fmt.Errorf("decode get me response: %w", err)
	}

	return user, statusCode, nil
}
