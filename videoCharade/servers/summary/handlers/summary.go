package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	if r.Method == "GET" {

		query := r.FormValue("url")
		if query == "" {
			http.Error(w, "Query parameter cannot be an empty string.", http.StatusBadRequest)
			return
		}
		body, bErr := fetchHTML(query)
		if bErr == nil {
			summary, sErr := extractSummary(query, body)
			header.Set("Content-Type", "application/json")
			if sErr == nil {
				data, _ := json.Marshal(summary)
				w.Write(data)
			} else {
				w.Write([]byte("[]"))
			}
		} else {
			http.Error(w, bErr.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Only GET Method allowed.", http.StatusBadRequest)
		return
	}
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	//GET the URL
	resp, err := http.Get(pageURL)

	//if there was an error, report it and exit
	if err != nil {
		return nil, err
	}

	//check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, errors.New("it is not html")
	}
	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	defer htmlStream.Close()
	tokenizer := html.NewTokenizer(htmlStream)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if "head" == token.Data {
				var tagTitle string
				summary := PageSummary{}
				imgCount := -1
				for {
					tokenType = tokenizer.Next()
					if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
						token := tokenizer.Token()
						if "meta" == token.Data || "link" == token.Data {
							attr := token.Attr
							attrDic := make(map[string]string)
							for i := range attr {
								key := string(attr[i].Key)
								val := string(attr[i].Val)
								attrDic[string(key)] = string(val)
							}
							attrKey := ""
							nameVal, nameExists := attrDic["name"]
							propVal, propExists := attrDic["property"]
							relVal, relExists := attrDic["rel"]
							if nameExists || propExists {
								if nameExists {
									attrKey = nameVal
								} else {
									attrKey = propVal
								}
								switch attrKey {
								case "og:type":
									summary.Type = attrDic["content"]
								case "twitter:card":
									if noOg(summary.Type) {
										summary.Type = attrDic["content"]
									}
								case "og:url":
									summary.URL = attrDic["content"]
								case "url":
									if noOg(summary.URL) {
										summary.URL = attrDic["content"]
									}
								case "og:title":
									summary.Title = attrDic["content"]
								case "twitter:title":
									if noOg(summary.Title) {
										summary.Title = attrDic["content"]
									}
								case "og:site_name":
									summary.SiteName = attrDic["content"]
								case "og:description":
									summary.Description = attrDic["content"]
								case "twitter:description":
									if noOg(summary.Description) {
										summary.Description = attrDic["content"]
									}
								case "description":
									if noOgOrTwit(summary.Description) {
										summary.Description = attrDic["content"]
									}
								case "author":
									summary.Author = attrDic["content"]
								case "keywords":
									keyString := attrDic["content"]
									keyArray := strings.Split(keyString, ",")
									for i := range keyArray {
										keyArray[i] = strings.Trim(keyArray[i], " ")
									}
									summary.Keywords = keyArray
								case "og:image":
									href := attrDic["content"]
									href = toAbsolute(href, pageURL)
									if !hasImgAlready(summary.Images, href) {
										imgCount++
										summary.Images = append(summary.Images, &PreviewImage{})
										summary.Images[imgCount].URL = href
									}
								case "twitter:image":
									href := attrDic["content"]
									href = toAbsolute(href, pageURL)
									if !hasImgAlready(summary.Images, href) {
										imgCount++
										summary.Images = append(summary.Images, &PreviewImage{})
										summary.Images[imgCount].URL = href
									}

								case "og:image:secure_url":
									summary.Images[imgCount].SecureURL = attrDic["content"]
								case "og:image:type":
									summary.Images[imgCount].Type = attrDic["content"]
								case "og:image:width":
									wString := attrDic["content"]
									width, err := strconv.Atoi(wString)
									if err == nil {
										summary.Images[imgCount].Width = width
									}
								case "og:image:height":
									height, err := strconv.Atoi(attrDic["content"])
									if err == nil {
										summary.Images[imgCount].Height = height
									}
								case "og:image:alt":
									summary.Images[imgCount].Alt = attrDic["content"]
								}
							} else if relExists {
								if relVal == "icon" {
									icon := PreviewImage{}
									href, hrefExists := attrDic["href"]
									sizes, sizeExists := attrDic["sizes"]
									imgType, typeExists := attrDic["type"]
									if hrefExists {
										icon.URL = toAbsolute(href, pageURL)
									}
									if sizeExists {
										if sizes != "any" {
											deminsions := strings.Split(sizes, "x")
											height, hErr := strconv.Atoi(deminsions[0])
											width, wErr := strconv.Atoi(deminsions[1])
											if hErr == nil && wErr == nil {
												icon.Height = height
												icon.Width = width
											}
										}
									}
									if typeExists {
										icon.Type = imgType
									}
									summary.Icon = &icon
								}
							}
						} else if "title" == token.Data {
							tokenType = tokenizer.Next()
							tagTitle = string(tokenizer.Token().Data)
						}
					}
					if tokenType == html.EndTagToken {
						token := tokenizer.Token()
						if token.Data == "head" {
							if len(summary.Title) < 1 {
								summary.Title = tagTitle
							}
							return &summary, nil
						}
					}
				}
			}
		}
	}
	return &PageSummary{}, nil
}

//toAbsolute conver the given path into an absolute path if the given path does not contain
//the given pageURL
func toAbsolute(path string, pageURL string) string {
	pathSlice := strings.SplitAfter(pageURL, ".com")
	pageURL = pathSlice[0]
	if !strings.Contains(path, "http") {
		path = pageURL + path
	}
	return path
}

//noOg checks if a given string is an OG meta tag. Returns a boolean indicating if it is
func noOg(s string) bool {
	return !strings.Contains(s, "og")
}

//noOgOrTwit checks if a given string is OG or Twitter meta tag. Returns a boolean indicating if it is
func noOgOrTwit(s string) bool {
	return !strings.Contains(s, "og") && !strings.Contains(s, "twitter")
}

//hasImgAlready iterates through the image slice to check if the slice already contains
//the given href. Returns a boolean indicating if it does
func hasImgAlready(imgs []*PreviewImage, href string) bool {
	for i := range imgs {
		if imgs[i].URL == href {
			return true
		}
	}
	return false
}
