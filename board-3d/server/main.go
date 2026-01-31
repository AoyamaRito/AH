package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "time"
)

type charsResponse struct {
    Board struct {
        Cols  int    `json:"cols"`
        Rows  int    `json:"rows"`
        Image string `json:"image"`
    } `json:"board"`
    Files []string `json:"files"`
}

func main() {
    addr := flag.String("addr", ":8000", "listen address (overridden by $PORT if set)")
    root := flag.String("root", ".", "root directory to serve")
    flag.Parse()

    // Respect PORT from environment (Railway/Heroku style)
    if p := os.Getenv("PORT"); p != "" {
        if !strings.HasPrefix(p, ":") { p = ":" + p }
        *addr = p
    }

    // Normalize root
    absRoot, err := filepath.Abs(*root)
    if err != nil {
        log.Fatalf("resolve root: %v", err)
    }

    mux := http.NewServeMux()

    // Dynamic chars manifest at the path frontend expects
    mux.HandleFunc("/AH/board-3d/data/chars.json", func(w http.ResponseWriter, r *http.Request) {
        writeCharsJSON(w, absRoot)
    })

    // Also provide an API alias
    mux.HandleFunc("/AH/board-3d/api/chars", func(w http.ResponseWriter, r *http.Request) {
        writeCharsJSON(w, absRoot)
    })

    // Simple health endpoint
    mux.HandleFunc("/AH/board-3d/api/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    // Redirect helper to ensure trailing slash for directory
    mux.HandleFunc("/AH/board-3d", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/AH/board-3d/", http.StatusTemporaryRedirect)
    })

    // Static files from root
    fs := http.FileServer(http.Dir(absRoot))
    mux.Handle("/", logRequests(securityHeaders(fs)))

    srv := &http.Server{
        Addr:              *addr,
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
    }

    log.Printf("Serving %s on http://localhost%s/", absRoot, *addr)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}

func writeCharsJSON(w http.ResponseWriter, root string) {
    // Scan for character images with fallback roots
    // 1) <root>/AH/img/Chara  (when root is project root containing AH/)
    // 2) <root>/img/Chara     (when root is AH repo root)
    charDir := filepath.Join(root, "AH", "img", "Chara")
    if _, err := os.Stat(charDir); err != nil {
        alt := filepath.Join(root, "img", "Chara")
        if _, err2 := os.Stat(alt); err2 == nil {
            charDir = alt
        }
    }
    var files []string
    exts := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".webp": true}

    _ = filepath.Walk(charDir, func(path string, info os.FileInfo, err error) error {
        if err != nil || info == nil || info.IsDir() {
            return nil
        }
        if !exts[strings.ToLower(filepath.Ext(info.Name()))] {
            return nil
        }
        // Build relative path from /AH/board-3d to the image file
        base := filepath.Join(root, "AH", "board-3d")
        if _, errb := os.Stat(base); errb != nil {
            base = filepath.Join(root, "board-3d")
        }
        rel, err := filepath.Rel(base, path)
        if err != nil {
            return nil
        }
        rel = filepath.ToSlash(rel)
        if !strings.HasPrefix(rel, "..") {
            rel = "../" + rel
        }
        files = append(files, rel)
        return nil
    })

    // Board defaults
    resp := charsResponse{}
    resp.Board.Cols = 10
    resp.Board.Rows = 10
    resp.Board.Image = "../img/Board_AH2.jpg"
    resp.Files = files

    // JSON
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    // Basic CORS for convenience (dev only)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    _ = enc.Encode(resp)
}

// Middleware: log requests
func logRequests(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// Middleware: minimal security headers
func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "SAMEORIGIN")
        // Allow SharedArrayBuffer on Chrome if needed (site-isolation hint)
        if runtime.GOOS == "darwin" || runtime.GOOS == "linux" || runtime.GOOS == "windows" {
            w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
            w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
        }
        next.ServeHTTP(w, r)
    })
}
