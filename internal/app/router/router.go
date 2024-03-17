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
	r.HandleFunc("/actors", hActor.GetActors)
	r.HandleFunc("/actors/add", hActor.AddActor)
	r.HandleFunc("/actors/update", hActor.UpdateActor)
	r.HandleFunc("/actors/delete", hActor.DeleteActor)

	r.HandleFunc("/movies/add", hMovie.AddMovie)
	r.HandleFunc("/movies/update", hMovie.UpdateMovie)
	r.HandleFunc("/movies", hMovie.GetMovie)
	r.HandleFunc("/movies/delete", hMovie.DeleteMovie)
	r.HandleFunc("/movies/search", hMovie.SearchMovie)

	handler := middleware.ValidateEndpoint(r)
	handler = middleware.Auth(handler)
	handler = middleware.Logging(handler, logger)
	handler = middleware.PanicRecovery(handler, logger)

	return &handler
}
