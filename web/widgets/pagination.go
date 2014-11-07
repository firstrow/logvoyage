package widgets

import (
	"fmt"
	"html/template"
	_ "log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// Help to paginate elastic search
type pagination struct {
	totalRecords uint64
	perPage      int
	req          *http.Request
}

func (this *pagination) SetTotalRecords(total uint64) {
	this.totalRecords = total
}

func (this *pagination) GetTotalRecords() uint64 {
	return this.totalRecords
}

func (this *pagination) GetTotalPages() int {
	r := float64(this.GetTotalRecords()) / float64(this.GetPerPage())
	return int(math.Ceil(r))
}

func (this *pagination) SetPerPage(limit int) {
	this.perPage = limit
}

func (this *pagination) GetPerPage() int {
	return this.perPage
}

// Detect `from` number for elastic
func (this *pagination) DetectFrom() int {
	page, _ := strconv.Atoi(this.req.URL.Query().Get("p"))
	if page > 1 {
		return this.GetPerPage()*page - this.GetPerPage()
	}
	return 0
}

// Get current page number
func (this *pagination) GetPageNumber() int {
	number, err := strconv.Atoi(this.req.URL.Query().Get("p"))
	if err != nil {
		return 1
	} else {
		return number
	}
}

func (this *pagination) Render() template.HTML {
	prev := `<li><a href="#">&laquo;</a></li>`
	next := `<li><a href="#">&raquo;</a></li>`

	padding := 5
	currentPage := this.GetPageNumber()
	startPage := currentPage - padding
	stopPage := currentPage + padding

	if startPage < 1 {
		startPage = 1
	}

	if stopPage > this.GetTotalPages() {
		stopPage = this.GetTotalPages()
	}

	var result string
	for i := startPage; i <= stopPage; i++ {
		result += fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", this.buildUrl(i), i)
	}

	return template.HTML(prev + result + next)
}

func (this *pagination) buildUrl(p int) string {
	u := this.req.RequestURI
	v, _ := url.Parse(u)
	values := v.Query()
	values.Set("p", strconv.Itoa(p))
	v.RawQuery = values.Encode()

	return v.String()
}

// Create new pagination object.
// Pass http request to detect current page from request uri.
func NewPagination(req *http.Request) *pagination {
	return &pagination{req: req}
}
