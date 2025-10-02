package static

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

// SetupStaticFiles configures static file serving (hybrid approach)
func SetupStaticFiles(r *mux.Router, staticFiles embed.FS) {
	isDev := os.Getenv("NODE_ENV") != "production" && os.Getenv("RENDER") == ""

	if isDev {
		log.Println("🔧 Static files: DEVELOPMENT mode")
		setupDevServer(r)
	} else {
		log.Println("📦 Static files: PRODUCTION mode (embedded)")
		setupProdServer(r, staticFiles)
	}
}

// setupProdServer serves embedded static files for production
func setupProdServer(r *mux.Router, staticFiles embed.FS) {
	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		log.Fatal("❌ Failed to create sub filesystem:", err)
	}

	// Serve static files with SPA fallback
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		serveSPA(w, req, distFS, true)
	})

	log.Println("✅ Static files embedded and ready")
}

// setupDevServer serves static files from filesystem for development
func setupDevServer(r *mux.Router) {
	buildDir := "./dist"

	// Check if dist directory exists
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		log.Println("⚠️  No dist/ folder found")
		log.Println("   Run 'cd frontend && npm run build' to create React build")
		log.Println("   Or use React dev server: 'cd frontend && npm run dev'")

		// Serve a helpful dev page
		r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/" {
				serveDevPage(w)
				return
			}
			http.NotFound(w, req)
		})
		return
	}

	// Serve from filesystem
	distFS := os.DirFS(buildDir)
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		serveSPA(w, req, distFS, false)
	})

	log.Printf("✅ Static files serving from %s", buildDir)
}

// serveSPA serves the Single Page Application with proper fallback
func serveSPA(w http.ResponseWriter, r *http.Request, distFS fs.FS, isEmbedded bool) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}

	// Try to serve the requested file
	if file, err := distFS.Open(path); err == nil {
		file.Close()

		// Set appropriate cache headers
		setCacheHeaders(w, path)

		// Serve the file
		if isEmbedded {
			http.ServeFileFS(w, r, distFS, path)
		} else {
			http.ServeFile(w, r, filepath.Join("dist", path))
		}
		return
	}

	// File not found - serve index.html for SPA routing
	w.Header().Set("Cache-Control", "no-cache")

	if isEmbedded {
		http.ServeFileFS(w, r, distFS, "index.html")
	} else {
		http.ServeFile(w, r, "dist/index.html")
	}
}

// setCacheHeaders sets appropriate cache headers for different file types
func setCacheHeaders(w http.ResponseWriter, path string) {
	if strings.HasPrefix(path, "assets/") || strings.HasPrefix(path, "static/") {
		// Cache static assets for 1 year
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else if strings.HasSuffix(path, ".html") {
		// Don't cache HTML files
		w.Header().Set("Cache-Control", "no-cache")
	} else if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") {
		// Cache JS/CSS for 1 day (in case they're not in assets folder)
		w.Header().Set("Cache-Control", "public, max-age=86400")
	} else {
		// Default cache for other files
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}
}

// serveDevPage serves a helpful development page
func serveDevPage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Go Art API - Development</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px; 
            margin: 50px auto; 
            padding: 20px;
            line-height: 1.6;
            color: #333;
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 10px;
            text-align: center;
            margin-bottom: 30px;
        }
        .api-link {
            display: inline-block;
            background: #28a745;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
            margin: 5px;
        }
        .api-link:hover { background: #218838; }
        .section {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .endpoint {
            background: white;
            padding: 10px;
            border-left: 4px solid #007bff;
            margin: 10px 0;
            font-family: monospace;
        }
        pre {
            background: #2d3748;
            color: #e2e8f0;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎨 Go Art API</h1>
        <p>Development Mode - API Server Running</p>
    </div>
    
    <div class="section">
        <h2>🚀 Quick Test</h2>
        <a href="/api/hello" class="api-link">Test API (/api/hello)</a>
        <a href="/api/health" class="api-link">Health Check</a>
        <a href="/api/users" class="api-link">Users API</a>
    </div>
    
    <div class="section">
        <h2>📋 Available Endpoints</h2>
        <div class="endpoint">GET /api/users - List all users</div>
        <div class="endpoint">POST /api/users - Create user</div>
        <div class="endpoint">POST /api/auth/register - Register new user</div>
        <div class="endpoint">POST /api/auth/login - Login user</div>
        <div class="endpoint">GET /api/artists - List artists (placeholder)</div>
        <div class="endpoint">GET /api/artworks - List artworks (placeholder)</div>
    </div>
    
    <div class="section">
        <h2>🔧 Frontend Development</h2>
        <p><strong>Option 1:</strong> Use React dev server (recommended)</p>
        <pre>cd frontend && npm run dev</pre>
        <p>React will run on <code>http://localhost:5173</code> with hot reload</p>
        
        <p><strong>Option 2:</strong> Build React and serve from Go</p>
        <pre>cd frontend && npm run build</pre>
        <p>Then refresh this page to see the React app</p>
    </div>
    
    <div class="section">
        <h2>📚 API Documentation</h2>
        <p>Test the API with curl:</p>
        <pre>curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"fname":"John","lname":"Doe","email":"john@example.com","password":"password123"}'</pre>
    </div>
</body>
</html>
    `

	fmt.Fprint(w, html)
}
