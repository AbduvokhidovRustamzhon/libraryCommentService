package app

import (
	"commentsService/pkg/crud/models"
	"context"
	"encoding/json"
	"github.com/AbduvokhidovRustamzhon/mux2/pkg/mux"
	"github.com/AbduvokhidovRustamzhon/rest/pkg/rest"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

const multipartMaxBytes = 10 * 1024 * 1024

func (receiver *server) handleBooksList() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		list, err := receiver.booksSvc.Books(request.Context())
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		urlsJSON, err := json.Marshal(list)
		if err != nil {
			log.Print(err)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(urlsJSON)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (receiver *server) handleBooksSave() func(responseWriter http.ResponseWriter, request *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseMultipartForm(multipartMaxBytes)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		file, header, err := request.FormFile("media")
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Print(header.Filename)
		defer file.Close()

		contentType := header.Header.Get("Content-Type")

		fileName, err := receiver.filesSvc.Save(file, contentType)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		name := header.Filename
		log.Print(name)
		book := models.Book{Name: name, Pages: 1, FileName: fileName}
		receiver.booksSvc.Save(context.Background(), book)
		err = rest.WriteJSONBody(writer, &book)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

}

func (receiver *server) handleBooksRemove() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		book := models.Book{}
		err = json.Unmarshal(body, &book)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		log.Print(book.Id)

		err = receiver.booksSvc.RemoveById(context.Background(), int(book.Id))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		http.Redirect(writer, request, "/api/books", http.StatusFound)
		return
	}
}

func (receiver *server) handleFavicon() func(http.ResponseWriter, *http.Request) {
	file, err := ioutil.ReadFile(filepath.Join(receiver.assetsPath, "favicon.ico"))
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(file)
		if err != nil {
			log.Print(err)
		}
	}
}

func (receiver *server) show() {

}

func (s *server) saveFile(fileHeader *multipart.FileHeader) (name string, err error) {
	file, err := fileHeader.Open()
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()

	contentType := fileHeader.Header.Get("Content-Type")
	name, err = s.filesSvc.Save(file, contentType)
	if err != nil {
		return
	}

	return
}



func (receiver *server) handleCommentsList() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		_, err = ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		bookID, ok := mux.FromContext(request.Context(),  "BookId") //string(comment.BookId)
		if !ok {
			return
		}
		log.Print(bookID)
		list, err := receiver.commentsSvc.ShowCommentsById(bookID, request.Context())
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		urlsJSON, err := json.Marshal(list)
		if err != nil {
			log.Print(err)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(urlsJSON)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (receiver server) handleNewComments() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		comment := models.Comments{
			Id:              0,
			Comment:         "",
			BookId:          0,
			CommentatorName: "",
			Removed:         false,
		}
		err = json.Unmarshal(body, &comment)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		log.Print(comment)

		err = receiver.commentsSvc.SaveComments(comment,context.Background())
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		http.Redirect(writer, request, "/api/books", http.StatusFound)
		return
	}
	}

func (receiver *server) handleCommentsRemove() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		comment := models.Comments{
			Id:              0,
			Comment:         "",
			BookId:          0,
			CommentatorName: "",
			Removed:         false,
		}
		err = json.Unmarshal(body, &comment)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		log.Print(comment.Id)

		err = receiver.commentsSvc.RemoveCommentsById(context.Background(), int(comment.Id))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

			_, err = writer.Write([]byte("Comment Deleted!"))
			if err != nil {
			log.Print(err)
		}
	}
}