package src

var extensions = map[string]string{
	".go":    "Go",
	".py":    "Python",
	".js":    "JavaScript",
	".ts":    "TypeScript",
	".c":     "C",
	".cpp":   "C++",
	".h":     "C Header",
	".hpp":   "C++ Header",
	".java":  "Java",
	".cs":    "C#",
	".php":   "PHP",
	".rb":    "Ruby",
	".rs":    "Rust",
	".swift": "Swift",
	".kt":    "Kotlin",
	".html":  "HTML",
	".css":   "CSS",
	".sh":    "Shell",
	".sql":   "SQL",
}

var ignoreDirs = []string{
	".git",
	"node_modules",
	"vendor",
}
