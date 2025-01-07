package params

var MimeTypes = map[string]string{
	// Images
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".ico":  "image/x-icon",

	// Documents
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",

	// Text
	".txt":  "text/plain",
	".csv":  "text/csv",
	".html": "text/html",
	".htm":  "text/html",
	".css":  "text/css",
	".js":   "text/javascript",
	".json": "application/json",
	".xml":  "application/xml",
	".yaml": "text/yaml",
	".yml":  "text/yaml",

	// Archives
	".zip": "application/zip",
	".gz":  "application/gzip",
	".tar": "application/x-tar",
	".7z":  "application/x-7z-compressed",
	".rar": "application/x-rar-compressed",

	// Audio
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".ogg":  "audio/ogg",
	".m4a":  "audio/mp4",
	".flac": "audio/flac",

	// Video
	".mp4": "video/mp4",
	".avi": "video/x-msvideo",
	".mov": "video/quicktime",
	".wmv": "video/x-ms-wmv",
	".mkv": "video/x-matroska",

	// Programming
	".go":   "text/x-go",
	".py":   "text/x-python",
	".java": "text/x-java",
	".rb":   "text/x-ruby",
	".php":  "text/x-php",
	".c":    "text/x-c",
	".cpp":  "text/x-c++",
	".rs":   "text/x-rust",

	// Binary formats
	".bin":   "application/octet-stream",
	".exe":   "application/octet-stream",
	".dll":   "application/octet-stream",
	".so":    "application/octet-stream",
	".dylib": "application/octet-stream",
}
