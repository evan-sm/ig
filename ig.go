package ig

import (
	//"fmt"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/k0kubun/pp"
	"github.com/tidwall/gjson"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36"
	postUrl   = "https://instagram.com/%v/?__a=1"
	xIgAppId  = "936619743392459"
)

var (
	storiesUrl = "https://i.instagram.com/api/v1/feed/reels_media/"
)

// A SuperAgent is an object storing all required request data
type SuperAgent struct {
	Client        *resty.Client
	Users         []User
	StoriesStruct []Stories
	AfterTs       int64
	QueryString   string
}

type User struct {
	Id   int64
	Name string
}

type Stories []struct {
	LatestReelMedia int `json:"latest_reel_media"`
	User            struct {
		Pk            int    `json:"pk"`
		Username      string `json:"username"`
		FullName      string `json:"full_name"`
		ProfilePicURL string `json:"profile_pic_url"`
	} `json:"user"`
	Items []struct {
		TakenAt        int `json:"taken_at"`
		MediaType      int `json:"media_type"`
		ImageVersions2 struct {
			Candidates []struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"candidates"`
		} `json:"image_versions2"`
		VideoVersions []struct {
			Type   int    `json:"type"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			URL    string `json:"url"`
			ID     string `json:"id"`
		} `json:"video_versions"`
		LinkText string `json:"link_text,omitempty"`
		StoryCta []struct {
			Links []struct {
				Weburi string `json:"webUri"`
			} `json:"links"`
		} `json:"story_cta,omitempty"`
	} `json:"items"`
}

// Used to create a new SuperAgent object
func Get(cookies string) *SuperAgent {
	if cookies == "" {
		log.Println("You have to specify instagram session id")
	}

	s := &SuperAgent{
		Client: &resty.Client{},
		//Query:  make(map[string]string),
	}
	s.Client = resty.New()
	s.Client.SetHeader("user-agent", userAgent)
	s.Client.SetHeader("x-ig-app-id", xIgAppId)
	s.Client.SetHeader("cookie", cookies) // Get cookies using browser's web inspector

	return s
}

// Used to set timestamp. All publications after specific timestamp will be ignored
func (s *SuperAgent) After(ts int64) *SuperAgent {
	s.AfterTs = ts
	return s
}

func (s *SuperAgent) Posts(names ...interface{}) *SuperAgent {
	pp.Println(names)
	for _, name := range names {
		pp.Printf("id: %v\n", name)

		switch v := reflect.ValueOf(name); v.Kind() {
		case reflect.String:
			//			log.Println("reflect.String")
			s.Users = append(s.Users, User{
				Name: v.String(),
			})
		case reflect.Slice:
			slice := makeSliceOfReflectValue(v)
			//			log.Println("slice:", slice)
			for _, value := range slice {
				//				pp.Printf("value: %v\n", value)
				user := reflect.ValueOf(value).String()
				s.Users = append(s.Users, User{
					Name: user,
				})
			}
		default:
			//			log.Println("default:")
			break
		}

	}
	//	pp.Printf("s.Users: %v\n", s.Users)
	s.getPosts()
	return s
}

func (s *SuperAgent) getPosts() {
	log.Println("Trying to get profile posts")

}

func (s *SuperAgent) Stories(ids ...interface{}) []byte {
	for _, id := range ids {
		switch v := reflect.ValueOf(id); v.Kind() {
		case reflect.Int64, reflect.Int:
			reel_id := strconv.FormatInt(v.Int(), 10)
			s.QueryString += "reel_ids=" + reel_id + "&"
			s.Users = append(s.Users, User{
				Id: v.Int(),
			})
		case reflect.Slice:
			slice := makeSliceOfReflectValue(v)
			//log.Println("slice:", slice)
			for _, value := range slice {
				//pp.Printf("value: %v\n", value)
				reel_id := strconv.FormatInt(reflect.ValueOf(value).Int(), 10)
				s.QueryString += "reel_ids=" + reel_id + "&"
			}
		default:
			log.Println("default:")
			break
		}

	}
	s.QueryString = strings.TrimSuffix(s.QueryString, "&")
	url := storiesUrl + "?" + s.QueryString
	pp.Printf("URL: %v\n", url)
	json := s.getStories(url)

	//	for k, v := range s.StoriesStruct {
	//		pp.Printf("\nk: %v v: %v\n", k, v[k].User.Username)
	//		pp.Printf("k: %v v: %v\n", k, v[k].Items[0].ImageVersions2.Candidates[0].URL)
	//	}

	return json
}

func (s *SuperAgent) getStories(url string) []byte {

	resp, err := s.Client.R().Get(url)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	result := gjson.GetBytes(resp.Body(), "reels_media")
	return []byte(result.Raw)

}

func (s *SuperAgent) getStories2() *SuperAgent {

	resp, err := s.Client.R().Get(storiesUrl)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	//log.Println("resp.String()", resp.String())
	result := gjson.GetBytes(resp.Body(), "reels_media")
	//log.Println("reflect:", reflect.TypeOf(result), reflect.TypeOf(result.String()))
	//log.Println("result:", result.String())
	result.ForEach(func(key, value gjson.Result) bool {
		//println("\n\nvalue.String()", value.String())
		var stories Stories
		err = json.Unmarshal([]byte(result.String()), &stories)
		if err != nil {
			log.Println(err.Error())
		}
		s.StoriesStruct = append(s.StoriesStruct, stories)
		pp.Println("ForEach: stories append", len(stories))
		return true // keep iterating
	})
	//pp.Println(s.StoriesStruct)
	return s
}

func (s *SuperAgent) do() []byte {
	resp, err := s.Client.R().Get(storiesUrl)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	//log.Printf("%v\n", resp)
	return resp.Body()

}

func makeSliceOfReflectValue(v reflect.Value) (slice []interface{}) {

	kind := v.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return slice
	}

	slice = make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		slice[i] = v.Index(i).Interface()
	}

	return slice
}
