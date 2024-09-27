package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Logger and Recovery middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load environment variables
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

	log.Println("Server running: http://localhost:8080")

	r.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/auth/:provider", SignInWithProvider)
	r.GET("/auth/:provider/callback", CallbackHandler)
	r.GET("/success", Success)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to run: http://localhost:8080")
	}
}

// Home function to handle the home route
func Home(c *gin.Context) {
	// Render the index template
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
	})
}

func SignInWithProvider(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func CallbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	_, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/success")
}

func Success(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(`
		<div style="
			background-color: #fff;
			padding: 40px;
			border-radius: 8px;
			box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
			text-align: center;
		">
			<h1 style="
				color: #333;
				margin-bottom: 20px;
			">You have Successfully signed in!</h1>
		</div>
	`)))
}
