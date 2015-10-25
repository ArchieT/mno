package wyniki

import (
	"bytes"
	"encoding/csv"
	"github.com/gocarina/gocsv"
	"html"
	"io"
	"os"
	"sort"
	"strconv"
)

type druzynaCSV struct {
	Nazwa             string `csv:"Nazwa"`
	Nazwisko1         string `csv:"Nazwisko1"`
	Nazwisko2         string `csv:"Nazwisko2"`
	Nazwisko3         string `csv:"Nazwisko3"`
	Nazwisko4         string `csv:"Nazwisko4"`
	Nazwisko5         string `csv:"Nazwisko5"`
	Nazwisko6         string `csv:"Nazwisko6"`
	BrakiPunktuZlyKod int    `csv:"BrakiPunktuZłyKod"`
	Mylne             int    `csv:"Mylne"`
	Spoznienie        int    `csv:"Spóźnienie"`
	Skreslenia        int    `csv:"Skreślenia"`
	BrakSpecjalnego   int    `csv:"brakspecjalnego"`
	ZmianaDecyzji     int    `csv:"zmianadecyzji"`
}

type Druzyna struct {
	Nazwa             string
	Nazwiska          []string
	BrakiPunktuZlyKod int
	Mylne             int
	Spoznienie        int
	Skreslenia        int
	BrakSpecjalnego   int
	ZmianaDecyzji     int
}

func (d Druzyna) ZaSkreślenia() int {
	return 10 * d.Skreslenia
}
func (d Druzyna) ZaBrakPunktuZłyKod() int {
	return 90 * d.BrakiPunktuZlyKod
}
func (d Druzyna) ZaBrakZadaniaSpecjalnego() int {
	return 10 * d.BrakSpecjalnego
}
func (d Druzyna) ZaStowarzyszony() int {
	return 25 * d.Mylne
}
func (d Druzyna) ZaZmianęDecyzji() int {
	return 10 * d.ZmianaDecyzji
}
func (d Druzyna) ZaSpóźnienie() int {
	if d.Spoznienie < 21 {
		return d.Spoznienie
	}
	return d.Spoznienie + ((d.Spoznienie - 20) * 9)
}
func (d Druzyna) PunktyKarne() int {
	return d.ZaSkreślenia() + d.ZaBrakPunktuZłyKod() + d.ZaBrakZadaniaSpecjalnego() + d.ZaStowarzyszony() + d.ZaZmianęDecyzji() + d.ZaSpóźnienie()
}

func (d druzynaCSV) Druzyna() Druzyna {
	nazwiska := make([]string, 0, 6)
	if len(d.Nazwisko1) > 0 {
		nazwiska = append(nazwiska, d.Nazwisko1)
	}
	if len(d.Nazwisko2) > 0 {
		nazwiska = append(nazwiska, d.Nazwisko2)
	}
	if len(d.Nazwisko3) > 0 {
		nazwiska = append(nazwiska, d.Nazwisko3)
	}
	if len(d.Nazwisko4) > 0 {
		nazwiska = append(nazwiska, d.Nazwisko4)
	}
	if len(d.Nazwisko5) > 0 {
		nazwiska = append(nazwiska, d.Nazwisko5)
	}
	if len(d.Nazwisko6) > 0 {
		nazwiska = append(nazwiska, d.Nazwisko6)
	}
	return Druzyna{d.Nazwa, nazwiska, d.BrakiPunktuZlyKod, d.Mylne, d.Spoznienie, d.Skreslenia, d.BrakSpecjalnego, d.ZmianaDecyzji}
}

type Wyniki []Druzyna

func (w Wyniki) Miejsca() map[int]int {
	miejsca := make(map[int]int)
	punkty := make([]int, 0, len(w))
	for _, j := range w {
		punkty = append(punkty, j.PunktyKarne())
	}
	sort.Ints(punkty)
	index := 1
	for _, j := range punkty {
		if _, ok := miejsca[j]; !ok {
			miejsca[j] = index
			index++
		}
	}
	return miejsca
}

type Zawody struct {
	*Wyniki
	Data     string
	Miejsce  string
	Czas     int
	CzasPlus int
}

func (z Zawody) Present(out io.Writer) {
	var buffer bytes.Buffer
	buffer.WriteString(`<html><head><meta charset="utf-8"><style> table, th, td { border: 1px solid black; padding: 5px; align: center; } </style><title>MnO — wyniki</title></head><body><p align="center"><font size="6"><b>Wyniki szkolnych Marszów na Orientację</b><br></font><font size="5"><b>`)
	buffer.WriteString(z.Data)
	buffer.WriteString(`<br>`)
	buffer.WriteString(z.Miejsce)
	buffer.WriteString(`</b></font></p><br><p align="center"><table><tr><th>Msc</th><th>PK</th><th>Nazwa</th><th>Imiona i nazwiska</th><th>BP</th><th>BP E</th><th>PS</th><th>Sþoźnienie</th><th>Poprawki</th><th>Skreślenia</th></tr>`)
	miejsca := z.Wyniki.Miejsca()
	td := func() {
		buffer.WriteString("</td><td>")
	}
	w := func(s string) {
		buffer.WriteString(s)
	}
	wi := func(naszint int) {
		w(strconv.Itoa(naszint))
	}
	wnaw := func(nawias int) {
		w(" (")
		wi(nawias)
		w(")")
	}
	wn := func(zwycz int, nawias int) {
		wi(zwycz)
		wnaw(nawias)
	}
	wnb := func(zwycz bool, nawias int) {
		if zwycz {
			w("TAK")
		} else {
			w("NIE")
		}
		wnaw(nawias)
	}
	for _, j := range *z.Wyniki {
		buffer.WriteString("<tr><td>")
		buffer.WriteString(strconv.Itoa(miejsca[j.PunktyKarne()]))
		buffer.WriteString("</td><td>")
		buffer.WriteString(strconv.Itoa(j.PunktyKarne()))
		buffer.WriteString("</td><td>")
		buffer.WriteString(html.EscapeString(j.Nazwa))
		buffer.WriteString("</td><td>")
		for ki, k := range j.Nazwiska {
			buffer.WriteString(html.EscapeString(k))
			if ki != len(j.Nazwiska)-1 {
				buffer.WriteString("<br>")
			}
		}
		buffer.WriteString("</td><td>")
		wn(j.BrakiPunktuZlyKod, j.ZaBrakPunktuZłyKod())
		td()
		wnb(j.BrakSpecjalnego != 0, j.ZaBrakZadaniaSpecjalnego())
		td()
		wn(j.Mylne, j.ZaStowarzyszony())
		td()
		wn(j.Spoznienie, j.ZaSpóźnienie())
		td()
		wn(j.ZmianaDecyzji, j.ZaZmianęDecyzji())
		td()
		wn(j.Skreslenia, j.ZaSkreślenia())
		w("</td></tr>")
	}
	w(`</table><br><hr><font size="1"></p><p align="right"><i>Wygenerowano za pomocą <a href="https://github.com/ArchieT/mno">github.com/ArchieT/mno</a></i></p></font></body></html>`)
	buffer.WriteTo(out)
}

type wynikiCSV []druzynaCSV

func (w wynikiCSV) Wyniki() Wyniki {
	lista := make(Wyniki, 0, len(w))
	for i := range w {
		lista = append(lista, w[i].Druzyna())
	}
	return lista
}

func Daj(in *os.File) Wyniki {
	defer in.Close()
	gocsv.SetCSVReader(func(in io.Reader) *csv.Reader {
		read := gocsv.DefaultCSVReader(in)
		read.TrimLeadingSpace = true
		return read
	})
	lista := make(wynikiCSV, 0, 50)
	if err := gocsv.UnmarshalFile(in, &lista); err != nil {
		panic(err)
	}
	return lista.Wyniki()
}
