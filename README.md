# OnionScraper

<div align="center">

![galileoff-OnionScraper Intro](assets/galileoff.png)

# 🧅 galileoff-OnionScraper

![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=for-the-badge&logo=go)
![Tor](https://img.shields.io/badge/Network-Tor-7D4698?style=for-the-badge&logo=tor-browser&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Maintained](https://img.shields.io/badge/Maintained-Yes-blue?style=for-the-badge)

**Tor Ağı İçin Gelişmiş Web Kazıma, Ekran Görüntüsü Alma ve Analiz Aracı**

*Siber Vatan Programı Yıldız CTI Ekibi görevi kapsamında geliştirilmiştir.*

[Özellikler](#-özellikler) • [Kurulum](#-kurulum) • [Kullanım](#-kullanım) • [Çıktı Yapısı](#-çıktı-yapısı) • [İletişim](#-destek)

</div>

---

## 📖 Hakkında

**galileoff-OnionScraper**, Tor ağı üzerindeki `.onion` uzantılı siteleri analiz etmek, veri toplamak ve belgelemek için tasarlanmış profesyonel bir araçtır. Standart `http` kütüphanelerinin aksine, **Chromedp** altyapısını bir Tor Proxy (SOCKS5) üzerinden yönlendirerek gerçek bir tarayıcı gibi davranır. Bu sayede JavaScript gerektiren modern dark web sitelerini, Captcha (otomatik olmayan) korumalarını ve dinamik içerikleri sorunsuz bir şekilde işleyebilir.

> [!TIP]
> **Estetik ve İşlevsel**: Araç, siber güvenlik uzmanlarının kullanım alışkanlıklarına uygun olarak gelişmiş bir **ASCII sanat arayüzü**, **renkli loglar** ve **canlı ilerleme çubukları** ile donatılmıştır.

## ✨ Teknik Özellikler

| Özellik | Açıklama |
| :--- | :--- |
| **🧅 Akıllı Tor Entegrasyonu** | Sistem (9050) ve Tor Browser (9150) portlarını otomatik algılar ve bağlanır. |
| **🔎 Tarayıcı Önceliği** | Öncelikle `msedge.exe` (Edge) arar, bulamazsa Chrome kullanarak sayfaları render eder. |
| **🏷️ Modüler Sınıflandırma** | `rules.yaml` kurallarına göre siteleri **Market, Forum, Fidye Yazılım, Silah** vb. olarak otomatik etiketler. |
| **🛡️ Gelişmiş Gizlilik** | WebRTC kapatma, DNS sızıntı koruması ve dinamik User-Agent rotasyonu sağlar. |
| **📸 Tam Ekran Görüntüsü** | Sitelerin render edilmiş son halini yüksek kaliteli `.png` olarak kaydeder. |
| **🔗 Güvenli Link Haritası** | Sayfa içindeki linkleri çıkarır ve yanlış tıklamaları önlemek için güvenli formatta (`[.]onion`) raporlar. |
| **⚡ Performans Yönetimi** | İhtiyaca göre **3 (Düşük)**, **5 (Orta)** veya **10 (Yüksek)** worker ile eşzamanlı tarama yapabilir. |

## 🛠 Kurulum

### Ön Gereksinimler

1.  **Go**: v1.23 veya üzeri.
2.  **Tor Bağlantısı**:
    *   **Yöntem 1 (Önerilen):** Tor Browser'ı açın ve açık bırakın (Port 9150).
    *   **Yöntem 2:** Tor servisini sistem servisi olarak başlatın (Port 9050).
3.  **Tarayıcı**: Microsoft Edge (Önerilen) veya Google Chrome.

### Hızlı Kurulum

```bash
# 1. Projeyi klonlayın
git clone https://github.com/galile0ff/galileoff-OnionScraper.git

# 2. Proje dizinine girin
cd galileoff-OnionScraper

# 3. Bağımlılıkları yükleyin
go mod tidy
```

## 🚀 Kullanım

Projeyi başlatmak için:

```bash
go run main.go
```

### 🎮 Etkileşimli Arayüz

Program sizi adım adım yönlendiren renkli bir menüye sahiptir:

1.  **Bağlantı Kontrolü**: Başlangıçta Tor IP adresinizi ve bağlantı durumunuzu test eder.
2.  **Hedef Seçimi**: `config/` klasöründeki dosyaları listeler.
3.  **Performans Ayarı**: Sistem gücünüze göre 3, 5 veya 10 "Worker" (Köle) seçebilirsiniz.
4.  **Canlı Takip**: Tarama sırasında işlem durumunu canlı bir ilerleme çubuğu ile izleyebilirsiniz.

### ⚙️ Yapılandırma Formatı (ÖNEMLİ)

Taranacak siteleri `config/targets.yaml` dosyasına ekleyin veya dosya uzantısı `.yaml` olacak şekilde yeni taranacak URL'lerin olduğu dosyayı `config/` klasörüne koyun, içerikte **hangi satıra nasıl URL koyduğunuzun bir önemi yoktur.** Örnekte olduğu gibi olabilir:

```text
http://exampleonionaddress.onion
darkmarketv2.onionhttp://forumxyz.onion
```

> [!WARNING]
> Dosyayı standart YAML formatında (örn: `- url: ...`) **YAZMAYINIZ**. Düz metin dosyayı gibi kullanınız. Program satır satır okuma yapar ve birleşik olan linkleri sizin için ayırıp tarama yapabilir.

## � Gelişmiş Yapılandırma

OnionScraper, tarama davranışını özelleştirmeniz için iki temel dosyaya daha sahiptir.

### 1. Sınıflandırma Kuralları (`config/rules.yaml`)
Programın siteleri nasıl etiketleyeceğini (örn: `[MARKET]`, `[FORUM]`) belirleyen kurallar bu dosyada tanımlanır. Kendi kurallarınızı ekleyebilirsiniz:

```yaml
categories:
  - id: "yeni_kategori"
    name: "Özel Kategori Adı"
    tag: "[ÖZEL-ETİKET]" 
    keywords:
      high:
        - "Kesin Eşleşme Kelimesi"
      medium:
        - "Olası Kelime 1"
        - "Olası Kelime 2"
    structure_rules:
      - selector: ".class-adi" # CSS Seçici ile kontrol
```

### 2. User-Agent Havuzu (`config/user_agents.json`)
Gizliliği artırmak için kullanılan tarayıcı kimlikleri burada bulunur. Listeyi güncel tutarak parmak izinizi değiştirebilirsiniz:

```json
[
  {
    "name": "Tor Browser - Windows",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/115.0",
    "headers": {
      "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
      "Accept-Language": "en-US,en;q=0.5"
    }
  }
]
```

## �📂 Çıktı Yapısı

Sonuçlar, seçtiğiniz config dosyasının adıyla bir dosyada toplanır (Örn: `targets` klasörü). Her site için ayrı klasör açılmaz, tüm veriler URL tabanlı isimlendirilerek düzenli bir şekilde saklanır.

```text
targets/
├── scan_result.log                     # Detaylı işlem ve hata günlüğü
├── links.txt                           # Tüm sitelerden toplanan linkler (Alt linklerde eklenir)
├── http_exampleonion_onion.html        # 1. Sitenin kaynak kodu
├── http_exampleonion_onion.png         # 1. Sitenin ekran görüntüsü
├── http_galileoff_onion.html          # 2. Sitenin kaynak kodu
└── http_galileoff_onion.png           # 2. Sitenin ekran görüntüsü
```

### links.txt Örneği
Linkler güvenlik amacıyla "defanged" formatta kaydedilir:
```text
================================================================================
  KAYNAK ADRES: [MARKET] http://exampleonion.onion
  BULUNAN LİNK SAYISI: 12
================================================================================
  [+] [LOGIN?]          http://auth[.]onion
  [+] [FORUM?]          http://community[.]onion
```

## 🏗 Proje Ağacı

```bash
.
├── 📂 config/           # Yapılandırma dosyaları
│   ├── rules.yaml       # Örnek sınıflandırma kuralları (Etiketleme için)
│   ├── targets.yaml     # Örnek hedef site listesi (Düz metin olarak linkler eklenebilir)
│   └── user_agents.json # Örnek User-Agent havuzu
├── 📂 internal/         # Uygulama çekirdek modülleri
│   ├── 📂 classifier/   # İçerik analiz ve etiketleme motoru
│   ├── 📂 config/       # Dosya okuma işlemleri
│   ├── 📂 network/      # Tor bağlantısı ve IP kontrolü
│   ├── 📂 report/       # Loglama ve dosya yazma işlemleri
│   ├── 📂 scanner/      # Chromedp motoru ve ekran görüntüsü
│   ├── 📂 ui/           # ASCII sanatları, menüler ve canlı ilerleme çubuğu
│   └── 📂 utils/        # Link ayıklama ve metin işleme
├── main.go              # Ana giriş noktası
└── README.md            # Dokümantasyon
```

## ☕ Destek

Bu proje açık kaynaklıdır ve topluluk desteğiyle geliştirilebilir. Eğer işinize yaradıysa:

<div align="center">
<a href="https://www.buymeacoffee.com/galileoff" target="_blank">
<img src="https://cdn.buymeacoffee.com/buttons/v2/default-red.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" >
</a>
</div>

## 🤝 Katkıda Bulunma

1. Fork'layın
2. Branch oluşturun (`git checkout -b feature/yeniozellik`)
3. Commit'leyin (`git commit -m 'Yeni özellik: X eklendi'`)
4. Push'layın (`git push origin feature/yeniozellik`)
5. Pull Request açın

## 📄 Lisans

Bu proje **MIT Lisansı** ile lisanslanmıştır. Detaylar için `LICENSE` dosyasına bakınız.

---

<div align="center">
Developed with 🧡 by <a href="https://github.com/galile0ff">galile0ff</a>
</div>
