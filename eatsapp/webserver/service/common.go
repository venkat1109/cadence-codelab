package service

import (
	"errors"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"net/http"
)

type (
	// Item models a food item on the menu.
	Item struct {
		ID          string
		Name        string
		Description string
		Image       string
		Price       float32
	}

	// Menu models a restaurant menu.
	Menu struct {
		Items []*Item
	}
)

const (
	// TemplatesGlob stores the value used to glob for templates.
	TemplatesGlob = "eatsapp/webserver/assets/tmpl/*"
)

// Templates stores the pre-processed templates.
var Templates *template.Template

func LoadTemplates() {
	Templates = template.Must(template.ParseGlob(TemplatesGlob))
}

// ViewHandler renders a http response using a template based on "page" param or request path.
func ViewHandler(w http.ResponseWriter, r *http.Request, data interface{}) error {
	page := r.URL.Query().Get("page")
	if len(page) == 0 {
		page = r.URL.Path[1:]
	}

	err := Templates.ExecuteTemplate(w, page, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

// NewMenu returns a new Menu object whose
// contents are loaded from the specified
// file path
func NewMenu(file string) (*Menu, error) {
	return loadMenu(file)
}

// GetItemByID returns the item matching the ID value passed as a param.
func (m *Menu) GetItemByID(id string) (*Item, error) {
	for _, v := range m.Items {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, errors.New("Invalid menu item: " + id)
}

// load populates the fields in the receiver from the file passed as parameter.
func loadMenu(file string) (*Menu, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var menu Menu
	err = yaml.Unmarshal(data, &menu)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}
