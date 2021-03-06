package displays

import (
	"github.com/ugent-library/biblio-backend/internal/app/helpers"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
)

func dissertationDetails(l *locale.Locale, p *models.Publication) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: l.T("builder.type"),
				Value: l.TS("publication_types", p.Type),
			},
			&display.Text{
				Label:         l.T("builder.doi"),
				Value:         p.DOI,
				ValueTemplate: "format/doi",
			},
			&display.Text{
				Label: l.T("builder.classification"),
				Value: l.TS("publication_classifications", p.Classification),
			},
		).
		AddSection(
			&display.Text{
				Label:    l.T("builder.title"),
				Value:    p.Title,
				Required: true,
			},
			&display.List{
				Label:  l.T("builder.alternative_title"),
				Values: p.AlternativeTitle,
			},
			&display.Text{
				Label: l.T("builder.publication_abbreviation"),
				Value: p.PublicationAbbreviation,
			},
		).
		AddSection(
			&display.List{
				Label:  l.T("builder.language"),
				Values: localize.LanguageNames(l, p.Language)},
			&display.Text{
				Label: l.T("builder.publication_status"),
				Value: l.TS("publication_publishing_statuses", p.PublicationStatus),
			},
			&display.Text{
				Label: l.T("builder.extern"),
				Value: helpers.FormatBool(p.Extern, "✓", "-"),
			},
			&display.Text{
				Label:    l.T("builder.year"),
				Value:    p.Year,
				Required: true,
			},
			&display.Text{
				Label: l.T("builder.place_of_publication"),
				Value: p.PlaceOfPublication,
			},
			&display.Text{
				Label: l.T("builder.publisher"),
				Value: p.Publisher,
			},
		).
		AddSection(
			&display.Text{
				Label: l.T("builder.volume"),
				Value: p.Volume,
			},
			&display.Text{
				Label: l.T("builder.page_count"),
				Value: p.PageCount,
			},
			&display.Text{
				Label: l.T("builder.series_title"),
				Value: p.SeriesTitle,
			},
		).
		AddSection(
			&display.Text{
				Label:    l.T("builder.defense_date"),
				Value:    p.DefenseDate,
				Required: true,
			},
			&display.Text{
				Label:    l.T("builder.defense_time"),
				Value:    p.DefenseTime,
				Required: true,
			},
			&display.Text{
				Label:    l.T("builder.defense_place"),
				Value:    p.DefensePlace,
				Required: true,
			},
		).
		AddSection(
			&display.Text{
				Label: l.T("builder.has_confidential_data"),
				Value: l.TS("confirmations", p.HasConfidentialData),
			},
			&display.Text{
				Label: l.T("builder.has_patent_application"),
				Value: l.TS("confirmations", p.HasPatentApplication),
			},
			&display.Text{
				Label: l.T("builder.has_publications_planned"),
				Value: l.TS("confirmations", p.HasPublicationsPlanned),
			},
			&display.Text{
				Label: l.T("builder.has_published_material"),
				Value: l.TS("confirmations", p.HasPublishedMaterial),
			},
		).
		AddSection(
			&display.Text{
				Label:   l.T("builder.wos_type"),
				Value:   l.TS("tooltip.publication", p.WOSType),
				Tooltip: l.T("tooltip.publication.wos_type"),
			},
			&display.Text{
				Label: l.T("builder.wos_id"),
				Value: p.WOSID,
			},
			&display.List{
				Label:  l.T("builder.issn"),
				Values: p.ISSN,
			},
			&display.List{
				Label:  l.T("builder.eissn"),
				Values: p.EISSN,
			},
			&display.List{
				Label:  l.T("builder.isbn"),
				Values: p.ISBN,
			},
			&display.List{
				Label:  l.T("builder.eisbn"),
				Values: p.EISBN,
			},
		)
}
