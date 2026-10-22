package portal

// Service is one entry of the portal catalog. Everyone signed in with an
// ioprodz account sees the whole catalog; a service opens for the users who
// hold a role in its Zitadel project and shows as locked for the others.
type Service struct {
	ProjectId   string
	Name        string
	Url         string
	Description string
	Note        string
	Category    string
	// Icon is a monochrome SVG under /static/img/portal, drawn in white on
	// Color. Without one the tile shows the service's initial instead.
	Icon  string
	Color string
}

type Category struct {
	Key  string
	Name string
}

var categories = []Category{
	{Key: "devops", Name: "DevOps / Infra"},
	{Key: "data", Name: "Data & Analytics"},
	{Key: "productivity", Name: "Productivity"},
	{Key: "ai", Name: "AI"},
	{Key: "apps", Name: "Apps"},
}

var catalog = []Service{
	{
		ProjectId:   "395367342827110403",
		Name:        "Webmail",
		Url:         "https://webmail.ioprodz.com",
		Description: "Read and send your @ioprodz.com mail.",
		Category:    "productivity",
		Icon:        "/static/img/portal/webmail.svg",
		Color:       "#3D6AAE",
	},
	{
		ProjectId:   "395367342827110403",
		Name:        "Mail account",
		Url:         "https://mail.ioprodz.com/account",
		Description: "App passwords for mail apps, aliases, settings.",
		Category:    "productivity",
		Icon:        "/static/img/portal/mail-account.svg",
		Color:       "#2B4C7E",
	},
	{
		ProjectId:   "395378177318322179",
		Name:        "Nextcloud",
		Url:         "https://nextcloud.ioprodz.com",
		Description: "Files, calendar and contacts.",
		Category:    "productivity",
		Icon:        "/static/img/portal/nextcloud.svg",
		Color:       "#0082C9",
	},
	{
		ProjectId:   "395377653885960195",
		Name:        "Gitea",
		Url:         "https://git.ioprodz.com",
		Description: "Git repositories.",
		Category:    "devops",
		Icon:        "/static/img/portal/gitea.svg",
		Color:       "#609926",
	},
	{
		ProjectId:   "395379861817917443",
		Name:        "Coolify",
		Url:         "https://paas.ioprodz.com",
		Description: "Deploy and run apps.",
		Category:    "devops",
		Icon:        "/static/img/portal/coolify.svg",
		Color:       "#6B3FD4",
	},
	{
		ProjectId:   "395524418236841987",
		Name:        "Harbor",
		Url:         "https://harbor.ioprodz.com",
		Description: "Container image registry.",
		Category:    "devops",
		Icon:        "/static/img/portal/harbor.svg",
		Color:       "#4A9A28",
	},
	{
		ProjectId:   "395422360401870851",
		Name:        "Deploy",
		Url:         "https://deploy.ioprodz.com",
		Description: "Apps, previews, CI and dev machines.",
		Category:    "devops",
		Icon:        "/static/img/portal/deploy.svg",
		Color:       "#0F766E",
	},
	{
		ProjectId:   "395381905081499651",
		Name:        "pgAdmin",
		Url:         "https://pga.ioprodz.com",
		Description: "Postgres administration.",
		Category:    "data",
		Icon:        "/static/img/portal/postgresql.svg",
		Color:       "#336791",
	},
	{
		ProjectId:   "395380440615092227",
		Name:        "Grafana",
		Url:         "http://192.168.1.108:4060",
		Description: "Server and LLM monitoring.",
		Note:        "Home network only",
		Category:    "data",
		Icon:        "/static/img/portal/grafana.svg",
		Color:       "#E8590C",
	},
	{
		ProjectId:   "395381905769365507",
		Name:        "Langfuse",
		Url:         "https://lf.ioprodz.com",
		Description: "LLM traces and evaluations.",
		Category:    "ai",
		Icon:        "/static/img/portal/langfuse.svg",
		Color:       "#1F2937",
	},
	{
		ProjectId:   "396083245675970563",
		Name:        "Muzio",
		Url:         "https://muzio.ioprodz.com",
		Description: "Music streaming.",
		Category:    "apps",
		Icon:        "/static/img/portal/muzio.svg",
		Color:       "#DB2777",
	},
	{
		ProjectId:   "396278155368726531",
		Name:        "Jellyfin",
		Url:         "https://jellyfin.ioprodz.com",
		Description: "Films and series.",
		Category:    "apps",
		Icon:        "/static/img/portal/jellyfin.svg",
		Color:       "#AA5CC3",
	},
}
