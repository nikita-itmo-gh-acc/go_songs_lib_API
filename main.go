package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"bytes"

	_ "songsapi/docs"

	"github.com/lib/pq"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	"io"

	"songsapi/logger"
	"songsapi/middleware"
	"songsapi/query"
	"songsapi/storage"

	"time"

	"os"
)

// @title Songs Library API
// @version 1.0
// @description This is a project implementing a service, providing songs information.
// @host localhost:8080
// @BasePath /api/v1

// @contact.email nickita-ananiev@yandex.ru


type SongResponse struct {
	Songs 		[]storage.Song
	Page 		int		
	Limit 		int
}

type SongTextResponse struct {
	Couplets	[]string
	Page 		int
	Limit 		int
}

type SongSearchHandler struct {
	SongsTable 	storage.Storage[storage.Song]
}

type SongOperationsHandler struct {
	SongsTable 	storage.Storage[storage.Song]
}

type TextPaginationHandler struct {
	Song		*storage.Song
}

type SongDeleteHandler struct {
	Song		*storage.Song
	SongsTable 	storage.Storage[storage.Song]
}

type SongUpdateHandler struct {
	SongId		int
	SongsTable 	storage.Storage[storage.Song]
}

type SongAddHandler struct {
	SongsTable 	storage.Storage[storage.Song]
	GroupsTable storage.Storage[storage.Group]
	DebugApiURL	string
}

type SongAddRequest struct {
	Song 	string	`json:"song"`
	Group 	string	`json:"group"`
}


// @Summary Returns a songs search result 
// @Description This endpoint parses url query params and do SQL select request based on them.
// @Tags songs search
// @Produce json
// @Router /songs [get]
// @Param limit query int false "Maximum number of songs to return"
// @Param page query int false "Page"
// @Success 200 {object} SongResponse
// @Failure 400 
// @Failure 404 
// @Failure 500 
func (h *SongSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	songQuery := new(query.SongQuery)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(songQuery, r.URL.Query()); err != nil {
		logger.Err.Println("bad request query")
		http.Error(w, "Can't parse query params", http.StatusInternalServerError)
		return 
	}

	var buf bytes.Buffer

	err := songQuery.Validate()

	if err != nil {
		if errList, ok := err.(govalidator.Errors); ok {
			for _, field := range errList {
				fmt.Fprintf(&buf, "Validation error in field %v\n", field)
			} 
		}
		logger.Err.Println("query params didn't pass validation")
		http.Error(w, buf.String(), http.StatusBadRequest)
		return
	}

	foundSongs, err := h.SongsTable.Find(songQuery)
	if err != nil {
		HandleDBSearchFail(w, err)
		return
	}

	response := SongResponse{
		Songs: make([]storage.Song, len(foundSongs)),
		Page: songQuery.Page,
		Limit: songQuery.Limit,
	}

	for i, song := range foundSongs {
		response.Songs[i] = *song
	}
	
	RenderJSON(w, response)
}


// @Summary Gets song by Id  
// @Tags songs operations
// @Router /songs/{id} [get]
// @Param id path int true "Song ID"
// @Success 200 {object} storage.Song
// @Failure 400
// @Failure 404
// @Failure 500 
func (h *SongOperationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		logger.Err.Println("no id provided")
		http.Error(w, "id url variable is required", http.StatusBadRequest)
		return
	}
	songId, _ := strconv.Atoi(id)
	foundSong, err := h.SongsTable.Get(songId)
	if err != nil {
		HandleDBSearchFail(w, err)
		return
	}

	switch r.Method {

	case http.MethodGet:
		if strings.HasSuffix(r.URL.Path, "/text") {
			textHandler := &TextPaginationHandler{ Song: foundSong }
			textHandler.ServeHTTP(w, r)
			return
		}
		RenderJSON(w, *foundSong)

	case http.MethodDelete:
		deleteHandler := &SongDeleteHandler{ Song: foundSong, SongsTable: h.SongsTable }
		deleteHandler.ServeHTTP(w, r)
		return

	case http.MethodPut:
		updateHandler := &SongUpdateHandler{ SongId: songId, SongsTable: h.SongsTable }
		updateHandler.ServeHTTP(w, r)
		return
	}
}

// @Tags songs operations
// @Summary Deletes song by Id
// @Router /songs/{id} [delete]
// @Param id path int true "Song ID"
// @Success 204 
// @Failure 400
// @Failure 404
// @Failure 500 
func (h *SongDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.SongsTable.Delete(h.Song)
	if err != nil {
		logger.Err.Println("delete failed - ", err)
		http.Error(w, fmt.Sprintf("Can't delete song with id = %d, Error: %v", h.Song.Id, err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Song succesfully deleted"))
}

// @Tags songs operations
// @Summary Updates song by Id
// @Router /songs/{id} [put]
// @Param id path int true "Song ID"
// @Param request body storage.Song true "Song update request"
// @Success 202
// @Failure 400
// @Failure 404
// @Failure 500 
func (h *SongUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var updatedSong storage.Song
	PasreJSON(r.Body, &updatedSong)

	err := h.SongsTable.Update(&updatedSong)
	if err != nil {
		logger.Err.Println("update failed - ", err)
		http.Error(w, fmt.Sprintf("Can't update song with id = %d, Error: %v", h.SongId, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Song succesfully updated"))
}

// @Tags text pagination
// @Summary Returns song text fragment
// @Description Does search for the song in database, then splits its text to couplets
// @Router /songs/{id}/text [get]
// @Param id path int true "Song ID"
// @Param limit query int false "Maximum number of couplets to return"
// @Param page query int false "Page"
// @Success 200 {object} SongTextResponse
// @Failure 400
// @Failure 404
// @Failure 500
func (h *TextPaginationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page, pageErr := ToInt(r.URL.Query().Get("page"))
	limit, limitErr := ToInt(r.URL.Query().Get("limit"))
	if pageErr != nil || limitErr != nil {
		logger.Err.Println("bad request query")
		http.Error(w, "Can't parse query", http.StatusInternalServerError)
		return
	}
	couplets := strings.Split(h.Song.Text, "\n\n")

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = len(couplets)
	}

	offset := limit * (page - 1)
	end := min(offset + limit, len(couplets))
	if offset >= len(couplets) {
		logger.Err.Println("lyrics not found")
		http.Error(w, "Couplets not found", http.StatusNotFound)
		return
	}
	RenderJSON(w, &SongTextResponse{
		Couplets: couplets[offset:end],
		Page: page,
		Limit: limit,
	})
}

// @Tags songs operations
// @Summary Adds new song
// @Router /songs/add [post]
// @Param request body SongAddRequest true "Song creation request"
// @Success 201 
// @Failure 404
// @Failure 500 
func (h *SongAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var newSong storage.Song
	PasreJSON(r.Body, &newSong)
	defer r.Body.Close()
	
	params := url.Values{}
	params.Add("song", newSong.Name)
	params.Add("group", newSong.Group)

	fullURL := fmt.Sprintf("%s?%s", h.DebugApiURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		logger.Err.Println("request to info api failed - ", err)
		http.Error(w, "Error during info api request", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	
	PasreJSON(resp.Body, &newSong)

	parsedDate, err := time.Parse("02.01.2006", newSong.ReleaseDate)
	if err != nil {
		logger.Err.Println("can't parse release date - ", err)
		http.Error(w, "Can't parse data from info API", http.StatusInternalServerError)
		return
	}

	newSong.ReleaseDate = parsedDate.Format("2006-01-02")

	err = h.GroupsTable.Create(&storage.Group{ Name: newSong.Group })
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if !(ok && pgErr.Code == "23505") {
			logger.Err.Println("group creation failed - ", err)
			http.Error(w, "Can't add group into database", http.StatusInternalServerError)
			return
		}
	}

	groups, err := h.GroupsTable.Find(&query.GroupQuery{ Name: newSong.Group })
	if err != nil {
		logger.Err.Println("group search failed - ", err)
		http.Error(w, "Can't add new song due to database error", http.StatusInternalServerError)
		return
	}

	foundGroup := groups[0]
	newSong.GroupId = foundGroup.Id
	err = h.SongsTable.Create(&newSong)
	if err != nil {
		logger.Err.Println("song creation failed - ", err)
		http.Error(w, "Can't add new song into database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Song succesfully added"))
}

func HandleDBSearchFail(w http.ResponseWriter, e error) {
	msg := fmt.Sprintf("Search failed Error: %v", e)
	if e == sql.ErrNoRows {
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	http.Error(w, msg, http.StatusInternalServerError)
}

func ToInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func RenderJSON(w http.ResponseWriter, object interface{}) {
	js, err := json.Marshal(object)
	if err != nil {
		http.Error(w, "Can't render JSON from object:", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func PasreJSON(r io.ReadCloser, object interface{}) {
	body, _ := io.ReadAll(r)
	if err := json.Unmarshal(body, object); err != nil {
		fmt.Println("Can't unmarshal json:", err)
		return
	}
}

func init() {
	logger.DoConsoleLog()
	// logger.LogToFile("app.log")
}

func main() {
	if err := godotenv.Load(); err != nil {
    	logger.Err.Fatalln("can't find .env file")
    }

	dbConn := storage.GetDBConnection()
	defer dbConn.Close()

	migrator, err := storage.CreateMigrator(dbConn)
	if err != nil {
		logger.Err.Printf("can't create migrator - %v\n", err) 
	}

	logger.Debug.Println("start migration process...")
	migrator.MakeMigrations()

	songs := &storage.SongStorage{DB: dbConn}
	groups := &storage.GroupStorage{DB: dbConn}

	query.SetQueryValidators()

	router := mux.NewRouter()

	apiSongs := router.PathPrefix("/api/v1/songs").Subrouter()
	apiSongs.Handle("", &SongSearchHandler{ SongsTable: songs }).Methods("GET")
	apiSongs.Handle("/add", &SongAddHandler{ 
		SongsTable: songs, GroupsTable: groups, DebugApiURL: os.Getenv("Debug_API_URL") }).Methods("POST")

	apiSongOps := apiSongs.PathPrefix("/{id:[0-9]+}").Subrouter()
	opsHandler := &SongOperationsHandler{ SongsTable: songs }
	apiSongOps.Handle("", opsHandler).Methods("GET", "DELETE", "PUT")
	apiSongOps.Handle("/text", opsHandler).Methods("GET")
	
	port := os.Getenv("SERV_PORT")
	logger.Debug.Printf("start listening on %s port...\n", port)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	songApiRouter := middleware.AccessLogMiddleware(middleware.CORSMiddware(router))
	http.ListenAndServe(":" + port, songApiRouter)
}
