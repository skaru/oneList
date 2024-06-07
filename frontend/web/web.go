package web

import (
	"bytes"
	"crypto/tls"
	"embed"
	"log"
	"net/http"
	"one-list/item"
	"one-list/storage"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

//go:embed server.crt
var certFile []byte

//go:embed server.key
var keyFile []byte

//go:embed templates/*
var htmlFiles embed.FS

const TEMPLATE_DIR = "templates/"

type Web struct {
	storage    storage.Storage
	authCookie http.Cookie
	md         goldmark.Markdown
	//templates  map[string]*template.Template
}

func (web *Web) Init(storage storage.Storage, username string, password string) {
	web.storage = storage

	web.md = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)

	web.authCookie = http.Cookie{
		Name:     username,
		Value:    password,
		MaxAge:   2147483647,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.HandleFunc("/", web.viewAll)
	http.HandleFunc("/view", web.view)
	http.HandleFunc("/edit", web.edit)
	http.HandleFunc("/save", web.save)
	http.HandleFunc("/create", web.create)
	http.HandleFunc("/delete", web.delete)
	http.HandleFunc("/hierinloggengraag", web.login)

	log.Println("Starting webserver")
	cert, _ := tls.X509KeyPair(certFile, keyFile)

	// Construct a tls.config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// Other options
	}

	// Build a server:
	server := http.Server{
		// Other options
		Addr:      ":8080",
		TLSConfig: tlsConfig,
	}

	log.Fatal(server.ListenAndServeTLS("", ""))
}

func (web Web) Close() {
	log.Println("Shutting down webserver")
}

func (web *Web) login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	hash := r.Form.Get("string")

	time.Sleep(6 * time.Second)
	if hash == web.authCookie.Value {
		http.SetCookie(w, &web.authCookie)
		http.Redirect(w, r, "/", http.StatusFound)
	} else if hash != "" {
		//w.WriteHeader(http.StatusForbidden)
		return
	}

	tmpl := template.Must(template.ParseFS(htmlFiles, TEMPLATE_DIR+"login.html", TEMPLATE_DIR+"header.html", TEMPLATE_DIR+"footer.html")) //remove into one call after testing
	tmpl.Execute(w, nil)
}

func (web *Web) viewAll(w http.ResponseWriter, r *http.Request) {
	if web.validate(w, r) {
		var data struct {
			Items       []item.Item
			Description string
		}

		items := web.storage.FetchAllItems()

		if len(items) > 0 {
			item.UpdateAndSortItems(items)

			data = struct {
				Items       []item.Item
				Description string
			}{
				Items:       items,
				Description: web.parseDescription(items[0]),
			}
		}
		tmpl := template.Must(template.ParseFS(htmlFiles, TEMPLATE_DIR+"viewAll.html", TEMPLATE_DIR+"header.html", TEMPLATE_DIR+"footer.html")) //remove into one call after testing
		tmpl.Execute(w, data)
	}
}

func (web *Web) view(w http.ResponseWriter, r *http.Request) {
	if web.validate(w, r) {
		ID, _ := strconv.Atoi(r.URL.Query().Get("ID"))
		if ID <= 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		selectedItem := web.storage.FetchItem(ID)

		var buffer bytes.Buffer
		web.md.Convert([]byte(selectedItem.Description), &buffer)

		data := struct {
			Item        item.Item
			Description string
		}{
			Item:        selectedItem,
			Description: web.parseDescription(selectedItem),
		}

		tmpl := template.Must(template.ParseFS(htmlFiles, TEMPLATE_DIR+"view.html", TEMPLATE_DIR+"header.html", TEMPLATE_DIR+"footer.html")) //remove into one call after testing
		tmpl.Execute(w, data)
	}
}

func (web *Web) edit(w http.ResponseWriter, r *http.Request) {
	if web.validate(w, r) {
		ID, _ := strconv.Atoi(r.URL.Query().Get("ID"))
		selectedItem := web.storage.FetchItem(ID)
		data := struct {
			Item item.Item
			Now  time.Time
		}{
			Item: selectedItem,
			Now:  time.Now().AddDate(0, 0, -1),
		}

		tmpl := template.Must(template.ParseFS(htmlFiles, TEMPLATE_DIR+"edit.html", TEMPLATE_DIR+"header.html", TEMPLATE_DIR+"footer.html")) //remove into one call after testing
		tmpl.Execute(w, data)
	}
}

func (web *Web) save(w http.ResponseWriter, r *http.Request) {
	if web.validate(w, r) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		ID, _ := strconv.Atoi(r.Form.Get("id"))
		selectedItem := web.storage.FetchItem(ID)

		due, _ := time.Parse("2006-01-02", r.Form.Get("due"))

		if dateEqual(due, time.Now().AddDate(0, 0, -1)) {
			due = selectedItem.Due
		}

		selectedItem.Name = r.Form.Get("name")
		selectedItem.Description = r.Form.Get("description")
		selectedItem.Due = due
		selectedItem.Reminder_interval, _ = strconv.Atoi(r.Form.Get("reminder_interval"))
		status, _ := strconv.Atoi(r.Form.Get("Display_status"))
		selectedItem.Display_status = item.Progress(status)
		selectedItem.Last_update = time.Now()

		web.storage.UpdateItem(selectedItem)

		referer := r.Header.Get("Referer")
		if strings.Contains(referer, "edit") {
			referer = "/view?ID=" + strconv.Itoa(selectedItem.ID)
		} else if referer == "" {
			referer = "/"
		}
		http.Redirect(w, r, referer, http.StatusFound)
	}
}

func (web *Web) create(w http.ResponseWriter, r *http.Request) {
	if web.validate(w, r) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		name := r.Form.Get("name")
		if name != "" {
			index := len(web.storage.FetchAllItems())
			if index <= 0 {
				index = 1
			}

			item := item.NewItem(index, r.Form.Get("name"))
			log.Println(item)
			web.storage.AddItem(item)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (web *Web) delete(w http.ResponseWriter, r *http.Request) {
	if web.validate(w, r) {
		ID, _ := strconv.Atoi(r.URL.Query().Get("ID"))
		web.storage.DeleteItem(ID)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (web *Web) validate(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie(web.authCookie.Name)
	if err != nil || cookie.Value != web.authCookie.Value {
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	return true
}

func (web *Web) parseDescription(item item.Item) string {
	var buffer bytes.Buffer
	web.md.Convert([]byte(item.Description), &buffer)

	return buffer.String()
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
