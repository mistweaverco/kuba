export type SearchEntry = {
	title: string;
	href: string;
	/** Extra terms users might type */
	keywords: string[];
	/** Short hint shown in results */
	excerpt?: string;
};

export const SEARCH_INDEX: SearchEntry[] = [
	{
		title: "Home",
		href: "/",
		keywords: ["home"],
		excerpt: "The home page.",
	},
	{
		title: "Installation",
		href: "/installation",
		keywords: [
			"arch",
			"pkgbuild",
			"kuba-bin",
			"pacman",
			"makepkg",
			"aur",
			"paru",
			"yay",
			"install",
			"download",
			"linux",
			"macos",
			"windows",
			"powershell",
			"curl",
			"wget",
			"zsh",
			"bash",
		],
		excerpt: "Install kuba CLI on Arch (kuba-bin), other Linux distros, macOS, and Windows.",
	},
	{
		title: "Usage",
		href: "/usage",
		keywords: ["run", "contain", "test", "show", "update", "debug", "docker", "ci", "cd"],
	},
	{
		title: "Configuration",
		href: "/configuration",
		keywords: [
			"kuba.yaml",
			"schema",
			"env",
			"secret-key",
			"secret-path",
			"value",
			"interpolation",
			"convert",
			"init",
		],
		excerpt:
			"How to configure kuba with kuba.yaml, including schema, env interpolation, and secret management.",
	},
	{
		title: "Getting Started",
		href: "/configuration#getting-started",
		keywords: [
			"convert",
			"ksvc",
			"dotenv",
			".env",
			"cloud run",
			"knative",
			"remote import",
			"provider",
			"gcp",
			"aws",
			"azure",
			"app runner",
			"container apps",
		],
		excerpt:
			"Getting started guide for converting existing .env files or remote secrets into a kuba.yaml configuration.",
	},
	{
		title: "Providers",
		href: "/providers",
		keywords: [
			"providers",
			"gcp",
			"aws",
			"azure",
			"openbao",
			"bitwarden",
			"local",
			"auth",
			"permissions",
		],
		excerpt:
			"Configure providers to fetch secrets from external sources like GCP Secret Manager, AWS Secrets Manager, Azure Key Vault, OpenBAO, Bitwarden, or local files.",
	},
	{
		title: "Examples",
		href: "/examples",
		keywords: ["examples"],
		excerpt: "Find example kuba.yaml configurations",
	},
	{
		title: "Cross provider examples",
		href: "/examples#cross-provider-configuration",
		keywords: ["examples", "cross provider"],
		excerpt:
			"Example kuba.yaml showing how to use multiple providers together to fetch secrets from GCP, AWS, and Azure in the same configuration.",
	},
	{
		title: "Github Actions examples",
		href: "/examples#github-actions",
		keywords: ["examples", "github"],
		excerpt: "Example github-actions workflow using kuba to inject secrets into a CI job.",
	},
	{
		title: "GitLab CI examples",
		href: "/examples#gitlab-ci",
		keywords: ["examples", "gitlab"],
		excerpt: "Example .gitlab-ci.yml using kuba to inject secrets into a CI job.",
	},
	{
		title: "Node.js examples",
		href: "/examples#nodejs-express-application",
		keywords: ["examples", "nodejs", "express", "node", "typescript"],
		excerpt: "Example kuba.yaml for running a Node.js Express application with secrets injected.",
	},
	{
		title: "Python examples",
		href: "/examples#python-flask-application",
		keywords: ["examples", "python"],
		excerpt: "Example kuba.yaml for running a Python Flask application with secrets injected.",
	},
	{
		title: "Docker examples",
		href: "/examples#docker-and-container-examples",
		keywords: ["examples", "docker"],
		excerpt: "Example kuba.yaml for running a Docker container with secrets injected.",
	},
	{
		title: "Docker Compose examples",
		href: "/examples#docker-compose-integration",
		keywords: ["examples", "docker", "compose"],
		excerpt: "Example kuba.yaml for running a Docker container with secrets injected.",
	},
	{
		title: "Interactive TUI",
		href: "/usage#tui",
		keywords: ["tui", "interactive", "secrets", "edit", "add", "terminal ui"],
		excerpt: "Use kuba tui to view, add, and edit secrets interactively.",
	},
	{
		title: "Changelog (CLI)",
		href: "/usage#changelog",
		keywords: ["changelog", "release notes", "latest", "version"],
		excerpt: "Use kuba changelog to view the baked-in changelog in your terminal.",
	},
	{
		title: "Create Template",
		href: "/usage#create-template",
		keywords: ["create", "template", "editor", "VISUAL", "EDITOR", "kuba create template"],
		excerpt: "Use kuba create template to create or edit your user template in $EDITOR.",
	},
];
