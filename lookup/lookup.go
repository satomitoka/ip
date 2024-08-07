package lookup

import (
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/http"

    "github.com/oschwald/maxminddb-golang"
)

var (
    asnDB     *maxminddb.Reader
    countryDB *maxminddb.Reader
)

// ASNRecord 保存ASN数据库的查询结果
type ASNRecord struct {
    ASN    string `maxminddb:"asn"`
    Domain string `maxminddb:"domain"`
    Name   string `maxminddb:"name"`
}

// CountryRecord 保存国家数据库的查询结果
type CountryRecord struct {
    Continent     string `maxminddb:"continent"`
    ContinentName string `maxminddb:"continent_name"`
    Country       string `maxminddb:"country"`
    CountryName   string `maxminddb:"country_name"`
}

// Init 初始化日志文件和数据库
func Init() {
    var err error
    // 打开ASN数据库
    asnDB, err = maxminddb.Open("/data/ipinfo/db/asn.mmdb")
    if err != nil {
        log.Fatal("Error opening ASN database:", err)
    }

    // 打开国家数据库
    countryDB, err = maxminddb.Open("/data/ipinfo/db/country.mmdb")
    if err != nil {
        log.Fatal("Error opening country database:", err)
    } 
}

// GetIPHandler 获取IP地址的处理函数
func GetIPHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // 尝试从X-Forwarded-For头部取得IP
    fwdIP := r.Header.Get("X-Forwarded-For")
    if fwdIP == "" {
        fwdIP = r.Header.Get("X-Real-IP")
    }
    // 如果两个头部都没有，则从连接中获取IP
    if fwdIP == "" {
        ip, _, _ := net.SplitHostPort(r.RemoteAddr)
        fwdIP = ip
    }
    // 直接返回IP地址，不使用JSON格式化
    fmt.Fprintf(w, fwdIP)
}

// IPLookupHandler IP查询的处理函数
func IPLookupHandler(w http.ResponseWriter, r *http.Request) {
    // 允许跨站请求
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    // 从请求中获取User-Agent头部，即浏览器信息
    userAgent := r.Header.Get("User-Agent")   
    
    // 尝试从查询参数获取IP
    ipStr := r.URL.Query().Get("ip")
    if ipStr == "" {
        // 尝试从X-Forwarded-For头部取得IP
        fwdIP := r.Header.Get("X-Forwarded-For")
        if fwdIP == "" {
            fwdIP = r.Header.Get("X-Real-IP")
        }
        // 如果两个头部都没有，则从连接中获取IP
        if fwdIP == "" {
            ip, _, _ := net.SplitHostPort(r.RemoteAddr)
            fwdIP = ip
        }
        ipStr = fwdIP
    }

    ip := net.ParseIP(ipStr)
    if ip == nil {
        http.Error(w, "Invalid IP address", http.StatusBadRequest)
        return
    }

    // 查询ASN记录
    var asn ASNRecord
    err := asnDB.Lookup(ip, &asn)
    if err != nil {
        http.Error(w, fmt.Sprintf("ASN Lookup failed: %v", err), http.StatusInternalServerError)
        return
    }

    // 查询国家记录
    var country CountryRecord
    err = countryDB.Lookup(ip, &country)
    if err != nil {
        http.Error(w, fmt.Sprintf("Country Lookup failed: %v", err), http.StatusInternalServerError)
        return
    }

    // 整理响应数据
    responseData := struct {
        IP            string `json:"ip"`
        ASN           string `json:"asn"`
        Domain        string `json:"domain"`
        ISP           string `json:"isp"`
        ContinentCode string `json:"continent_code"`
        ContinentName string `json:"continent_name"`
        CountryCode   string `json:"country_code"`
        CountryName   string `json:"country_name"`
        UserAgent     string `json:"user_agent"`
    }{
        IP:            ipStr,
        ASN:           asn.ASN,
        Domain:        asn.Domain,
        ISP:           asn.Name,
        ContinentCode: country.Continent,
        ContinentName: country.ContinentName,
        CountryCode:   country.Country,
        CountryName:   country.CountryName,
        UserAgent:     userAgent,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(responseData)
}