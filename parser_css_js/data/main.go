package main

import (
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path"
    "regexp"
    "strings"
)

func main() {
    cssURL := "https://10web-site.ai/124/wp-content/plugins/ai-builder-demo-plugin-master/assets/css/fonts.css"
    fontsDir := "./fonts"

    if err := os.MkdirAll(fontsDir, os.ModePerm); err != nil {
        fmt.Printf("Ошибка создания папки fonts: %v\n", err)
        return
    }

    resp, err := http.Get(cssURL)
    if err != nil {
        fmt.Printf("Ошибка загрузки CSS: %v\n", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Ошибка: статус HTTP %d\n", resp.StatusCode)
        return
    }

    cssData, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Ошибка чтения CSS: %v\n", err)
        return
    }

    fontRegex := regexp.MustCompile(url\(['"]?([^'"]+\.(ttf|woff|sfd))['"]?\))
    matches := fontRegex.FindAllStringSubmatch(string(cssData), -1)

    if len(matches) == 0 {
        fmt.Println("Шрифты не найдены.")
        return
    }

    baseURL, err := url.Parse(cssURL)
    if err != nil {
        fmt.Printf("Ошибка парсинга базового URL: %v\n", err)
        return
    }

    downloadedFiles := make(map[string]bool)

    for _, match := range matches {
        fontPath := match[1]
        
        if strings.HasPrefix(fontPath, "../") {
            fontPath = strings.TrimPrefix(fontPath, "../")
        }

        fontURL := fmt.Sprintf("%s://%s%s/%s", 
            baseURL.Scheme, 
            baseURL.Host, 
            path.Dir(path.Dir(baseURL.Path)),
            fontPath)

        fileName := path.Base(fontPath)

        if downloadedFiles[fileName] {
            continue
        }
        downloadedFiles[fileName] = true

        fmt.Printf("Загрузка шрифта: %s\n", fontURL)

        resp, err := http.Get(fontURL)
        if err != nil {
            fmt.Printf("Ошибка загрузки файла %s: %v\n", fontURL, err)
            continue
        }

        if resp.StatusCode != http.StatusOK {
            fmt.Printf("Ошибка загрузки файла %s: статус HTTP %d\n", fontURL, resp.StatusCode)
            resp.Body.Close()
            continue
        }

        outFile, err := os.Create(path.Join(fontsDir, fileName))
        if err != nil {
            fmt.Printf("Ошибка создания файла %s: %v\n", fileName, err)
            resp.Body.Close()
            continue
        }

        _, err = io.Copy(outFile, resp.Body)
        outFile.Close()
        resp.Body.Close()

        if err != nil {
            fmt.Printf("Ошибка сохранения файла %s: %v\n", fileName, err)
            continue
        }

        fmt.Printf("Файл %s успешно загружен.\n", fileName)
    }

    fmt.Println("Парсинг завершен.")
}
