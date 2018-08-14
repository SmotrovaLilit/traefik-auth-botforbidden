package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

//список взят состраниц: https://yandex.ru/support/webmaster/robot-workings/check-yandex-robots.html#robot-in-logs
// и https://support.google.com/webmasters/answer/1061943?hl=ru
var botUserAgents = []string{
	"Mozilla/5.0 (compatible; YandexAccessibilityBot/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexAdNet/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexBlogs/0.99; robot; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots) ",
	"Mozilla/5.0 (compatible; YandexBot/3.0; MirrorDetector; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexCalendar/1.0; +http://yandex.com/bots) ",
	"Mozilla/5.0 (compatible; YandexCatalog/3.0; +http://yandex.com/bots) ",
	"Mozilla/5.0 (compatible; YandexDirect/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexDirectDyn/1.0; +http://yandex.com/bots",
	"Mozilla/5.0 (compatible; YaDirectFetcher/1.0; Dyatel; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexFavicons/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexFavicons/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexForDomain/1.0; +http://yandex.com/bots) ",
	"Mozilla/5.0 (compatible; YandexImages/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexImageResizer/2.0; +http://yandex.com/bots) ",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B411 Safari/600.1.4 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B411 Safari/600.1.4 (compatible; YandexMobileBot/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexMarket/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexMedia/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexMetrika/2.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexNews/4.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexOntoDB/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexOntoDBAPI/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexPagechecker/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexSearchShop/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexSitelinks; Dyatel; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexSpravBot/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexTurbo/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexVertis/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexVerticals/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexVideo/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexVideoParser/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexWebmaster/2.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36 (compatible; YandexScreenshotBot/3.0; +http://yandex.com/bots) ",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36 (compatible; YandexMedianaBot/1.0; +http://yandex.com/bots)",
	"APIs-Google (+https://developers.google.com/webmasters/APIs-Google.html)",
	"Mediapartners-Google",
	"Mozilla/5.0 (Linux; Android 5.0; SM-G920A) AppleWebKit (KHTML, like Gecko) Chrome Mobile Safari (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
	"AdsBot-Google (+http://www.google.com/adsbot.html)",
	"Googlebot-Image/1.0",
	"Googlebot-News",
	"Googlebot-Video/1.0",
	"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Safari/537.36",
	"Googlebot/2.1 (+http://www.google.com/bot.html)",
	"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"(compatible; Mediapartners-Google/2.1; +http://www.google.com/bot.html)",
	"AdsBot-Google-Mobile-Apps",
}

func TestIfBotRequest(t *testing.T) {
	for _, userAgent := range botUserAgents {
		for _, typeR := range []string{"POST", "GET"} {
			t.Run(userAgent, func(t *testing.T) {
				request, err := http.NewRequest(typeR, "/", nil)
				if err != nil {
					t.Error(err)
					return
				}
				request.Header.Add("User-Agent", userAgent)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(handlerCheckBot)
				handler.ServeHTTP(rr, request)
				if err != nil {
					t.Error(err)
					return
				}

				if rr.Code != 401 {
					t.Errorf("type request: %s, got status code %d, want %d", typeR, rr.Code, 401)
				}
			})
		}

	}
}

var notBotUserAgents = []string{
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36",
	"Mozilla/5.0 (Linux; Android 4.4.2; XMP-6250 Build/HAWK) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Safari/537.36 ADAPI/2.0 (UUID:9e7df0ed-2a5c-4a19-bec7-2cc54800f99d) RK3188-ADAPI/1.2.84.533 (MODEL:XMP-6250)",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; Media Center PC 6.0; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 2.0.50727; .NET4.0C; .NET4.0E)",
	"Opera/9.80 (J2ME/MIDP; Opera Mini/4.2/28.3492; U; en) Presto/2.8.119 Version/11.10",
	"Mozilla/5.0 (Windows NT 6.2; WOW64; rv:52.0) Gecko/20100101 Firefox/52.0",
	"Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_3 like Mac OS X) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.0 Mobile/14G60 Safari/602.1",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
}

func TestIfNotBotRequest(t *testing.T) {
	for _, userAgent := range notBotUserAgents {
		for _, typeR := range []string{"POST", "GET"} {
			t.Run(userAgent, func(t *testing.T) {
				request, err := http.NewRequest(typeR, "/", nil)
				if err != nil {
					t.Error(err)
					return
				}
				request.Header.Add("User-Agent", userAgent)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(handlerCheckBot)
				handler.ServeHTTP(rr, request)
				if err != nil {
					t.Error(err)
					return
				}

				if rr.Code != 200 {
					t.Errorf("type request: %s, got status code %d, want %d", typeR, rr.Code, 200)
				}
			})
		}

	}
}
