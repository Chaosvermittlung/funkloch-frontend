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

func init() {
	item1 := []template.HTML{`<li class="active"><a href="/">Overview <span class="sr-only">(current)</span></a></li>`, `<li><a href="/">Overview</a></li>`}
	item2 := []template.HTML{`<li class="active"><a href="/item">Items <span class="sr-only">(current)</span></a></li>`, `<li><a href="/item">Items</a></li>`}
	item3 := []template.HTML{`<li class="active"><a href="/equipment">Equipment <span class="sr-only">(current)</span></a></li>`, `<li><a href="/equipment">Equipment</a></li>`}
	item4 := []template.HTML{`<li class="active"><a href="/event">Event <span class="sr-only">(current)</span></a></li>`, `<li><a href="/event">Events</a></li>`}
	item5 := []template.HTML{`<li class="active"><a href="/store">Stores <span class="sr-only">(current)</span></a></li>`, `<li><a href="/store">Stores</a></li>`}
	item6 := []template.HTML{`<li class="active"><a href="/fault">Faults <span class="sr-only">(current)</span></a></li>`, `<li><a href="/fault">Faults</a></li>`}
	item7 := []template.HTML{`<li class="active"><a href="/wishlist">Wishlists <span class="sr-only">(current)</span></a></li>`, `<li><a href="/wishlist">Wishlists</a></li>`}
	navitems = append(navitems, item1)
	navitems = append(navitems, item2)
	navitems = append(navitems, item3)
	navitems = append(navitems, item4)
	navitems = append(navitems, item5)
	navitems = append(navitems, item6)
	navitems = append(navitems, item7)
}

const (
	OverviewActive int = 1 + iota
	ItemsActive
	EquipmentActive
	EventsActive
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
	bcode, err := ean.Encode(id)

	if err != nil {
		return err
	}

	bc, err := barcode.Scale(bcode, 620, 200)

	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, bc, nil)
	by := buf.Bytes()
	r := bytes.NewReader(by)
	pdf.RegisterImageReader("code", "JPEG", r)
	pdf.Image("code", 1, 0, 62, 20, false, "", 0, "")
	w := pdf.GetStringWidth(id)
	pos := (62 - w) / 2
	pdf.Text(pos, 20+4, id)

	fs := 1.0
	w = pdf.GetStringWidth(store)
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
	pos = (64 - w) / 2
	h := fs / 72 * 25.4
	y := (16 - h) / 2
	pdf.Text(pos, 40-y, store)
	err = pdf.Output(out)
	return err
}
