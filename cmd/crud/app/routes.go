package app

import (
	_"github.com/AbduvokhidovRustamzhon/mux2/pkg/mux"
	"net/http"
)
func (receiver *server) InitRoutes() {
	mux := receiver.router

	mux.GET("/", receiver.handleBooksList())

	mux.GET("/api/books", receiver.handleBooksList())
	mux.POST("/api/books/save", receiver.handleBooksSave())
	mux.POST("/api/books/remove", receiver.handleBooksRemove())

	mux.GET("/api/books/comment/{BookId}", receiver.handleCommentsList())

	mux.POST("/api/books/comment/save", receiver.handleNewComments())

	mux.POST("/api/books/comment/remove", receiver.handleCommentsRemove())

	mux.GET("/favicon.ico", receiver.handleFavicon())
	mux.GET("/media", http.StripPrefix("/media", http.FileServer(http.Dir(receiver.mediaPath))).ServeHTTP)
}
