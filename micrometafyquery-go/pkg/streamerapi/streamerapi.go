package streamerapi

import (
	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/model"
)

// StreamerAPI is an interface reprensenting the possible actions on a music streaming service web API
type StreamerAPI interface {
	SearchTrack(queryString string) ([]model.Track, error)
}
