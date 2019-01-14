package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
	"github.com/gorilla/securecookie"
	"github.com/jung-kurt/gofpdf"
)

const errormessage = `<div class="alert alert-danger" role="alert">
  <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>
  <span class="sr-only">Error:</span>
  $MESSAGE$
</div>`

const cookieplace = "funkloch"

var navitems [][]template.HTML
var sc = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func createNavitem(name string, link string) []template.HTML {
	active := `<li class="active"><a href="/` + link + `">` + name + ` <span class="sr-only">(current)</span></a></li>`
	inactive := `<li><a href="/` + link + `">` + name + `</a></li>`
	activehtml := template.HTML(active)
	inactivehtml := template.HTML(inactive)
	return []template.HTML{activehtml, inactivehtml}
}

func init() {
	navitems = append(navitems, createNavitem("Overview", ""))
	navitems = append(navitems, createNavitem("Items", "item"))
	navitems = append(navitems, createNavitem("Boxes", "box"))
	navitems = append(navitems, createNavitem("Equipment", "equipment"))
	navitems = append(navitems, createNavitem("Events", "event"))
	navitems = append(navitems, createNavitem("Packinglist", "packinglist"))
	navitems = append(navitems, createNavitem("Stores", "store"))
	navitems = append(navitems, createNavitem("Faults", "fault"))
	navitems = append(navitems, createNavitem("Whishlists", "whishlist"))
}

const (
	OverviewActive int = 1 + iota
	ItemsActive
	BoxesActive
	EquipmentActive
	EventsActive
	PackinglistActive
	StoresActive
	FaultsActive
	WishlistsActive
)

func BuildSidebar(item int) template.HTML {
	var res template.HTML
	res = res + `<div class="col-sm-3 col-md-2 sidebar">`
	res = res + `<ul class="nav nav-sidebar">`
	for i, n := range navitems {
		var add template.HTML
		if i+1 == item {
			add = n[0]
		} else {
			add = n[1]
		}
		res = res + add + "\n"
	}
	res = res + `</ul>`
	res = res + `</div>`
	return res
}

func BuildMessage(tp string, message string) template.HTML {
	message = html.EscapeString(message)
	return template.HTML(strings.Replace(tp, "$MESSAGE$", message, -1))
}

func EncodeBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

type config struct {
	Port   int
	Master string
}

func (c *config) Load() error {
	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(c)
	return err
}

func getNewHTTPRequest(method string, endpoint string, in io.Reader) (*http.Request, error) {
	var req *http.Request
	var err error
	url := "http://" + conf.Master + "/api/v100/" + endpoint
	fmt.Println(url)
	req, err = http.NewRequest(method, url, in)
	if err != nil {
		return req, err
	}
	return req, nil
}

func sendauthorizedHTTPRequest(method string, endpoint string, token string, in io.Reader, v interface{}) error {
	req, err := getNewHTTPRequest(method, endpoint, in)
	if err != nil {
		return errors.New("Error creating request: " + err.Error())
	}

	req.Header.Set("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Error executing request: " + err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		fmt.Println(resp.StatusCode)
		decoder := json.NewDecoder(resp.Body)
		var er ErrorResponse
		err = decoder.Decode(&er)
		if err != nil {
			data, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(data))
			return errors.New("Error while decoding Error response. Only God can help you now:" + err.Error())
		}
		return errors.New("Got negativ status code: " + er.Errorcode + ":" + er.Errormessage)
	}

	if v != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(v)
		if err != nil {
			return errors.New("Unable to decode on given Interface: " + err.Error())
		}
	}
	return nil
}

func SetCookie(w http.ResponseWriter, key string, value string) error {

	cmap := make(map[string]string)
	cmap[key] = value

	encoded, err := sc.Encode(cookieplace, cmap)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  cookieplace,
		Value: encoded,
		Path:  "/",
	}
	cookie.Expires = time.Now().Add(10 * 365 * 24 * time.Hour)

	http.SetCookie(w, cookie)
	return nil
}

func GetCookie(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(cookieplace)

	if err != nil {
		return "", err
	}
	cmap := make(map[string]string)
	err = sc.Decode(cookieplace, cookie.Value, &cmap)

	if err != nil {
		return "", err
	}
	return cmap[key], nil

}

func RemoveCookie(w http.ResponseWriter, r *http.Request) {
	expire := time.Now().AddDate(0, 0, 1)

	cookieMonster := &http.Cookie{
		Name:    cookieplace,
		Expires: expire,
		Value:   "",
	}

	// http://golang.org/pkg/net/http/#SetCookie
	// add Set-Cookie header
	http.SetCookie(w, cookieMonster)
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func createlabel(id string, store string, out io.Writer) error {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: 62, Ht: 42},
		OrientationStr: "P",
	})
	pdf.SetMargins(1, 1, 1)
	pdf.SetFont("Helvetica", "", 14)
	pdf.AddPage()
	pdf.SetXY(1, 0)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	store = tr(store)
	fs := 1.0
	w := pdf.GetStringWidth(store)
	for (w < 60) && (fs < 45) {
		fs = fs + 1
		pdf.SetFontSize(fs)
		w = pdf.GetStringWidth(store)
	}
	fs = fs - 1
	if fs < 1 {
		fs = 1
	}
	pdf.SetFontSize(fs)
	pos := (64 - w) / 2
	h := fs / 72 * 25.4
	y := (16 - h) + 3
	pdf.Text(pos, y, store)

	bcode, err := ean.Encode(id)

	if err != nil {
		return err
	}

	bc, err := barcode.Scale(bcode, 420, 200)

	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, bc, nil)
	by := buf.Bytes()
	r := bytes.NewReader(by)
	pdf.RegisterImageReader("code", "JPEG", r)
	imageh := 42.0
	pdf.Image("code", 11, y+3, imageh, 20, false, "", 0, "")

	pdf.SetFontSize(14)
	w = pdf.GetStringWidth(id)
	pos = (64 - w) / 2
	//fh := 14 / 72 * 25.4
	pdf.Text(pos, y+imageh-15, id)

	err = pdf.Output(out)
	return err
}

func formatIndex(index int) string {
	Result := ""
	if index < 10 {
		Result = "0"
	}
	Result = Result + strconv.Itoa(index)
	return Result
}

func createContentlabel(items []itemResponse, out io.Writer) error {
	count := float64(len(items))
	height := count*4 + 12.0
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: 62, Ht: height},
		OrientationStr: "P",
	})
	pdf.SetMargins(1, 1, 1)
	pdf.SetFont("Helvetica", "", 20)
	pdf.AddPage()
	pdf.SetXY(1, 0)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	content := tr("Content:")
	y := 6.0
	pdf.Text(1.0, y, content)
	pdf.SetFont("Helvetica", "", 10)
	y = 12.0
	for index, i := range items {
		line := formatIndex(index+1) + ": " + strconv.Itoa(i.Item.Code) + " - " + i.Equipment.Name
		pdf.Text(1.0, y, line)
		y = y + 4
	}
	err := pdf.Output(out)
	return err
}

const cellHeight = 7.0

func createPackinglistHeader(pdf *gofpdf.Fpdf, header string, subheader string, event string) {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Helvetica", "B", 30)
	pdf.Text(5.0, 15, tr(header))
	pdf.SetFont("Helvetica", "", 24)
	pdf.Text(5.0, 25, tr(subheader))
	pdf.SetFont("Helvetica", "", 16)
	pdf.Text(5.0, 32, tr(event))
	oldwidth := pdf.GetLineWidth()
	pdf.SetLineWidth(1)
	w, _ := pdf.GetPageSize()
	pdf.Line(5.0, 35, w-5, 35)
	pdf.SetLineWidth(oldwidth)
}

func createCheckBox(pdf *gofpdf.Fpdf, label string, fill bool) {
	pdf.Rect(pdf.GetX()+1, pdf.GetY()+(cellHeight-5.0), 5.0, 5.0, "D")
	pdf.CellFormat(5.0+1, cellHeight, "", "", 0, "RB", fill, 0, "")
	w := pdf.GetStringWidth(label)
	pdf.CellFormat(w+3, cellHeight, label, "", 0, "LB", fill, 0, "")
}

func addCounterOffset(pdf *gofpdf.Fpdf, fill bool) {
	pdf.CellFormat(4.0, cellHeight, "", "", 0, "RB", fill, 0, "")
	pdf.CellFormat(10.0, cellHeight, "", "", 0, "RB", fill, 0, "")
}

func addCommentLine(pdf *gofpdf.Fpdf, fill bool) {
	addCounterOffset(pdf, fill)
	pdf.CellFormat(pdf.GetStringWidth("Comment:"), cellHeight, "Comment:", "", 0, "LB", fill, 0, "")
	pdf.Line(pdf.GetX(), pdf.GetY()+cellHeight, pdf.GetX()+100, pdf.GetY()+cellHeight)
}

func addCheckboxes(pdf *gofpdf.Fpdf, check []string, fill bool) {
	addCounterOffset(pdf, fill)
	for _, c := range check {
		createCheckBox(pdf, c, fill)
	}
}

func createListEntry(pdf *gofpdf.Fpdf, counter, code, description string, check []string, fill bool) {
	pdf.CellFormat(4.0, cellHeight, counter, "", 0, "RB", fill, 0, "")
	pdf.CellFormat(10.0, cellHeight, "", "", 0, "RB", fill, 0, "")
	pdf.CellFormat(50.0, cellHeight, code, "", 0, "LB", fill, 0, "")
	pdf.CellFormat(0, cellHeight, description, "", 1, "LB", fill, 0, "")
	addCheckboxes(pdf, check, fill)
	pdf.Ln(-1)
	addCommentLine(pdf, fill)
	pdf.Ln(-1)
	pdf.Ln(-1)
}

func createPackinglistCheckStrings() []string {
	var result []string
	result = append(result, "Loaded from Store")
	result = append(result, "Unloaded at Event")
	result = append(result, "Loaded at Event")
	result = append(result, "Unloaded at Store")
	return result
}

func createBoxCheckStrings() []string {
	var result []string
	result = append(result, "Unloaded from Box")
	result = append(result, "Put back into Box")
	return result
}

func addnewPage(pdf *gofpdf.Fpdf, header string, subheader string, event string) {
	pdf.AddPage()
	createPackinglistHeader(pdf, header, subheader, event)
	pdf.SetY(40)
	pdf.SetFont("Helvetica", "", 13)
}

func createPacklistPDF(p Packinglist, out io.Writer) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	addnewPage(pdf, "Packinglist: "+p.Name, "Boxeslist", p.Event.Name)
	c := createPackinglistCheckStrings()
	_, pageheight := pdf.GetPageSize()
	for i, b := range p.Boxes {
		if pdf.GetY()+5*cellHeight > pageheight {
			addnewPage(pdf, "Packinglist: "+p.Name, "Boxeslist", p.Event.Name)
		}
		createListEntry(pdf, strconv.Itoa(i+1), strconv.Itoa(b.Code), b.Description, c, false)
	}
	c = createBoxCheckStrings()
	for _, b := range p.Boxes {
		addnewPage(pdf, "Packinglist: "+p.Name, "Items of Box: "+strconv.Itoa(b.Code), p.Event.Name)
		for ii, i := range b.Items {
			if pdf.GetY()+5*cellHeight > pageheight {
				addnewPage(pdf, "Packinglist: "+p.Name, "Boxeslist", p.Event.Name)
			}
			createListEntry(pdf, strconv.Itoa(ii+1), strconv.Itoa(i.Code), i.Equipment.Name, c, false)
		}
	}
	err := pdf.Output(out)
	return err
}
