package router

import (
	"net/http"

	actor "github.com/DmitriyKomarovCoder/Film_library/internal/actor/delivery/http"
	movie "github.com/DmitriyKomarovCoder/Film_library/internal/movie/delivery/http"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/middleware"
)

func NewRouter(hMovie *movie.Handler, hActor *actor.Handler, logger *logger.Logger) *http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/actors", hActor.GetActor)
	r.HandleFunc("/actors/add", hActor.AddActor)
	r.HandleFunc("/actors/update", hActor.UpdateActor)
	r.HandleFunc("/actors/delete", hActor.DeleteActor)

	r.HandleFunc("/films/add", hMovie.AddMovie)
	r.HandleFunc("/films/update", hMovie.UpdateMovie)
	r.HandleFunc("/films", hMovie.GetMovie)
	r.HandleFunc("/films/delete", hMovie.DeleteMovie)
	r.HandleFunc("/films/search", hMovie.SearchMovie)

	handler := middleware.ValidateEndpoint(r, logger)
	handler = middleware.Auth(handler)
	handler = middleware.Logging(handler, logger)
	handler = middleware.PanicRecovery(handler, logger)

	return &handler
}
