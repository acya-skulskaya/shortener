# Инкремент 17. Улучшение производительности по памяти при помощи анализа профилей pprof

Для снятия профиля использовался вариант конфигурации приложения с in-memory storage с подачей нагрузки на эндпоинт "POST /"

Для эксперимента поменяла в хранилище работу с обычного map с sync.RWMutex на sync.Map - потребление памяти заметно увеличилось - эксперимент не удачный

```shell
$ go tool pprof -top -diff_base=profiles/base.pprof profiles/result_sync_map.pprof
File: shortener
Type: inuse_space
Time: 2026-01-19 14:11:39 MSK
Duration: 60.01s, Total samples = 17430.12kB 
Showing nodes accounting for 8171.66kB, 46.88% of 17430.12kB total
      flat  flat%   sum%        cum   cum%
 4096.19kB 23.50% 23.50%  4096.19kB 23.50%  internal/sync.newEntryNode[go.shape.interface {},go.shape.interface {}] (inline)
 3584.55kB 20.57% 44.07%  3584.55kB 20.57%  internal/sync.newIndirectNode[go.shape.interface {},go.shape.interface {}] (inline)
-2069.27kB 11.87% 32.19%  5611.46kB 32.19%  github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory.(*Container).add
 1536.07kB  8.81% 41.01%  1536.07kB  8.81%  encoding/json.(*decodeState).literalStore
 1024.02kB  5.87% 46.88%  1024.02kB  5.87%  github.com/acya-skulskaya/shortener/internal/helpers.RandStringRunes
  512.17kB  2.94% 49.82%   512.17kB  2.94%  internal/profile.(*Profile).postDecode
 -512.06kB  2.94% 46.88%  -512.06kB  2.94%  strings.(*Builder).grow
         0     0% 46.88%  1536.07kB  8.81%  encoding/json.(*Decoder).Decode
         0     0% 46.88%  1536.07kB  8.81%  encoding/json.(*decodeState).object
         0     0% 46.88%  1536.07kB  8.81%  encoding/json.(*decodeState).unmarshal
         0     0% 46.88%  1536.07kB  8.81%  encoding/json.(*decodeState).value
         0     0% 46.88%  8171.66kB 46.88%  github.com/acya-skulskaya/shortener/internal/middleware.CookieAuth.func1
         0     0% 46.88%  8171.66kB 46.88%  github.com/acya-skulskaya/shortener/internal/middleware.RequestLogger.func1
         0     0% 46.88%  1024.01kB  5.87%  github.com/acya-skulskaya/shortener/internal/middleware.getUserID
         0     0% 46.88%  6635.48kB 38.07%  github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory.(*InMemoryShortURLRepository).Store
         0     0% 46.88%   512.17kB  2.94%  github.com/go-chi/chi/v5.(*Mux).Mount.func1
         0     0% 46.88%  8171.66kB 46.88%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 46.88%  7147.65kB 41.01%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 46.88%  1536.07kB  8.81%  github.com/golang-jwt/jwt/v4.(*Parser).ParseUnverified
         0     0% 46.88%  1024.01kB  5.87%  github.com/golang-jwt/jwt/v4.(*Parser).ParseWithClaims
         0     0% 46.88%  1024.01kB  5.87%  github.com/golang-jwt/jwt/v4.ParseWithClaims
         0     0% 46.88%   512.17kB  2.94%  internal/profile.Parse
         0     0% 46.88%   512.17kB  2.94%  internal/profile.parseUncompressed
         0     0% 46.88%  7680.73kB 44.07%  internal/sync.(*HashTrieMap[go.shape.interface {},go.shape.interface {}]).Store (inline)
         0     0% 46.88%  7680.73kB 44.07%  internal/sync.(*HashTrieMap[go.shape.interface {},go.shape.interface {}]).Swap
         0     0% 46.88%  3584.55kB 20.57%  internal/sync.(*HashTrieMap[go.shape.interface {},go.shape.interface {}]).expand
         0     0% 46.88%  6635.48kB 38.07%  main.(*ShortUrlsService).apiPageMain
         0     0% 46.88%  8171.66kB 46.88%  net/http.(*conn).serve
         0     0% 46.88%  8171.66kB 46.88%  net/http.HandlerFunc.ServeHTTP
         0     0% 46.88%  8171.66kB 46.88%  net/http.serverHandler.ServeHTTP
         0     0% 46.88%   512.17kB  2.94%  net/http/pprof.collectProfile
         0     0% 46.88%   512.17kB  2.94%  net/http/pprof.handler.ServeHTTP
         0     0% 46.88%   512.17kB  2.94%  net/http/pprof.handler.serveDeltaProfile
         0     0% 46.88%  -512.06kB  2.94%  strings.(*Builder).Grow
         0     0% 46.88%  -512.06kB  2.94%  strings.Join
         0     0% 46.88%  7680.73kB 44.07%  sync.(*Map).Store (inline)
```

Вернула обратно работу с map с sync.RWMutex, но добавила задание размера мапы при старте, это уменьшило успальзование памяти на 13мб

```shell
$ go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
File: shortener
Type: inuse_space
Time: 2026-01-19 14:11:39 MSK
Duration: 60s, Total samples = 17430.12kB 
Showing nodes accounting for -13333.94kB, 76.50% of 17430.12kB total
      flat  flat%   sum%        cum   cum%
-7189.54kB 41.25% 41.25% -7189.54kB 41.25%  github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory.(*Container).add
-4608.28kB 26.44% 67.69% -11285.81kB 64.75%  main.(*ShortUrlsService).apiPageMain
-1536.07kB  8.81% 76.50% -1536.07kB  8.81%  encoding/json.(*decodeState).literalStore
 -512.06kB  2.94% 79.44%  -512.06kB  2.94%  strings.(*Builder).grow
  512.01kB  2.94% 76.50%   512.01kB  2.94%  github.com/acya-skulskaya/shortener/internal/helpers.RandStringRunes
         0     0% 76.50% -1536.07kB  8.81%  encoding/json.(*Decoder).Decode
         0     0% 76.50% -1536.07kB  8.81%  encoding/json.(*decodeState).object
         0     0% 76.50% -1536.07kB  8.81%  encoding/json.(*decodeState).unmarshal
         0     0% 76.50% -1536.07kB  8.81%  encoding/json.(*decodeState).value
         0     0% 76.50% -13333.94kB 76.50%  github.com/acya-skulskaya/shortener/internal/middleware.CookieAuth.func1
         0     0% 76.50% -13333.94kB 76.50%  github.com/acya-skulskaya/shortener/internal/middleware.RequestLogger.func1
         0     0% 76.50% -2048.13kB 11.75%  github.com/acya-skulskaya/shortener/internal/middleware.getUserID
         0     0% 76.50% -6677.53kB 38.31%  github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory.(*InMemoryShortURLRepository).Store
         0     0% 76.50% -13333.94kB 76.50%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 76.50% -11285.81kB 64.75%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 76.50% -1536.07kB  8.81%  github.com/golang-jwt/jwt/v4.(*Parser).ParseUnverified
         0     0% 76.50% -2048.13kB 11.75%  github.com/golang-jwt/jwt/v4.(*Parser).ParseWithClaims
         0     0% 76.50% -2048.13kB 11.75%  github.com/golang-jwt/jwt/v4.ParseWithClaims
         0     0% 76.50% -13333.94kB 76.50%  net/http.(*conn).serve
         0     0% 76.50% -13333.94kB 76.50%  net/http.HandlerFunc.ServeHTTP
         0     0% 76.50% -13333.94kB 76.50%  net/http.serverHandler.ServeHTTP
         0     0% 76.50%  -512.06kB  2.94%  strings.(*Builder).Grow
         0     0% 76.50%  -512.06kB  2.94%  strings.Join

```

## Вывод
Удалось снизить потребление памяти для указанного эндпоинта с 17Мб до 4Мб, что составляет 76%