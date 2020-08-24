package graph

import (
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/graph/model"
)

func draftEntitiesToSSEntities(draft model.EventDraft) supersense.Entities {
	entities := supersense.Entities{
		Tags:  []string{},
		Media: []supersense.MediaEntity{},
		Urls:  []supersense.URLEntity{},
	}

	if draft.Entities != nil {
		for _, tag := range draft.Entities.Tags {
			entities.Tags = append(entities.Tags, tag)
		}
		for _, media := range draft.Entities.Media {
			if media == nil {continue}
			entities.Media = append(entities.Media, supersense.MediaEntity{
				URL:  media.URL,
				Type: media.Type,
			})
		}

		for _, url := range draft.Entities.Urls {
			if url == nil {continue}
			entities.Urls = append(entities.Urls, supersense.URLEntity{
				URL:        url.URL,
				DisplayURL: url.DisplayURL,
			})
		}
	}

	return entities
}