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
	padding      int // How many links display before and after current
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
	currentPage := this.GetPageNumber()
	startPage := currentPage - this.padding

	if startPage < 1 {
		startPage = 1
	}
	stopPage := startPage + (this.padding * 2) - 1

	if stopPage > this.GetTotalPages() {
		stopPage = this.GetTotalPages()
	}

	prevPage := currentPage - 1
	nextPage := currentPage + 1
	var next string
	var prev string
	if prevPage < 1 {
		prevPage = 1
		prev = fmt.Sprintf("<li class=\"disabled\"><span>&laquo;<span></li>")
	} else {
		prev = fmt.Sprintf("<li><a href=\"%v\">&laquo;</a></li>", this.buildUrl(prevPage))
	}
	if nextPage > this.GetTotalPages() {
		nextPage = this.GetTotalPages()
		next = "<li class=\"disabled\"><span>&raquo;</span></li>"
	} else {
		next = fmt.Sprintf("<li><a href=\"%v\">&raquo;</a></li>", this.buildUrl(nextPage))
	}

	var active string
	var result string
	for i := startPage; i <= stopPage; i++ {
		if i == currentPage {
			active = "active"
		} else {
			active = ""
		}
		result += fmt.Sprintf("<li class=\"%v\"><a href=\"%v\">%v</a></li>", active, this.buildUrl(i), i)
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
	return &pagination{req: req, padding: 5}
}
