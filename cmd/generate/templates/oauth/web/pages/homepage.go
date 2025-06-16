package pages

import (
	"net/http"

	"oauth/views/pages"
)

func (p *Pages) Home(w http.ResponseWriter, r *http.Request) {
	err := home.Index().Render(r.Context(), w)
	if err != nil {
		p.logger.Err(err).Msg("Error occurred in rendering homepage")
		return
	}
}
