// Package utils provides helper functions and utilities for the bot.
package utils

import (
	"botwa/types"
	"context"
	"fmt"
	"runtime"
	"time"
    "io/ioutil"
    "net/http"
    "os/exec"
    "strings"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"go.mau.fi/libsignal/logger"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// Global start time untuk runtime bot
var appStartTime = time.Now()

// ExecuteShell menjalankan perintah shell dan mengembalikan output
func ExecuteShell(command string) (string, error) {
    parts := strings.Fields(command)
    if len(parts) == 0 {
        return "", fmt.Errorf("command kosong")
    }

    cmd := exec.Command(parts[0], parts[1:]...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}

func infoMessage(msg types.Messages) string {
    return fmt.Sprintf(`
=============================
ℹ️ MESSAGE INFO DETAIL
=============================
🔗 Chat        : %v
👤 FromUser    : %s
🌐 FromServer  : %s
🤖 FromMe      : %v
🆔 ID          : %s
👥 IsGroup     : %v
👑 IsOwner     : %v
📨 Sender      : %v
👤 SenderUser  : %s
🌐 SenderServer: %s
🏷️ Pushname    : %s
🕒 Timestamp   : %s
🔣 Prefix      : %s
💻 Command     : %s
📦 Args        : %v
💬 Text        : %s
📝 Body        : %s
=============================
`, 
        msg.From,
        msg.FromUser,
        msg.FromServer,
        msg.FromMe,
        msg.ID,
        msg.IsGroup,
        msg.IsOwner,
        msg.Sender,
        msg.SenderUser,
        msg.SenderServer,
        msg.Pushname,
        msg.Timestamp.Format("2006-01-02 15:04:05"),
        msg.Prefix,
        msg.Command,
        msg.Args,
        msg.Text,
        msg.Body,
    )
}

// HandleCommand routes and processes user commands.
func HandleCommand(client *whatsmeow.Client, m types.Messages, evt *events.Message) {
	if m.Prefix == "" {
		return
	}

	switch m.Command {
	//--------CASE MENU-------//
case "menu":
	jid := evt.Info.Chat
	currentTime := time.Now().Format("02-Jan-2006 15:04:05")
	uptime := formatDuration(time.Since(appStartTime))

	menuText := fmt.Sprintf(`*📋 DAFTAR MENU BOT*  
🕒 Waktu: %s  
⏱️ Runtime Bot: %s  

*⚡ Command Utama:*  
• *.ping* – Cek status server dan bot  
• *.exec* – Command server bot  
• *.info* – Detail Pesan Anda

*📤 Promosi / JPM:*  
• *.jpmall* – Promosi All List  
• *.jpmvpn* – Promosi List VPN  
• *.jpmvps* – Promosi List VPS  
• *.jpmdor* – Promosi List Paket  

Silakan ketik salah satu perintah di atas.
`, currentTime, uptime)

	_, err := client.SendMessage(context.Background(), jid, &waProto.Message{
		Conversation: proto.String(menuText),
	})
	if err != nil {
		logger.Error("Failed to send menu: " + err.Error())
	}


case "info":
    msg := Serialize(evt)
    infoText := infoMessage(msg) // panggil fungsi untuk dapatkan detail

    // Kirim ke WhatsApp
    _, _ = client.SendMessage(context.Background(), msg.From, &waProto.Message{
        Conversation: proto.String(infoText),
    })

case "exec":
jid := evt.Info.Chat
msg := Serialize(evt)
if !msg.IsOwner {
    _, _ = client.SendMessage(context.Background(), msg.From, &waProto.Message{
        Conversation: proto.String("❌ Perintah ini hanya untuk owner."),
    })
    return
}

    if m.Text == "" {
        _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
            Conversation: proto.String("⚠️ Harap masukkan perintah setelah *.exec*."),
        })
        return
    }

    output, err := ExecuteShell(m.Text)
    if err != nil {
        output = fmt.Sprintf("❌ Error: %v\n%s", err, output)
    }

    if len(output) > 4000 { // batasi biar ga kepanjangan
        output = output[:4000] + "\n\n⚠️ Output terpotong..."
    }

    _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
        Conversation: proto.String(fmt.Sprintf("📌 Hasil eksekusi:\n\n%s", output)),
    })

// --------CASE JPM-------//
case "jpmdor":
jid := evt.Info.Chat
msg := Serialize(evt)
if !msg.IsOwner {
    _, _ = client.SendMessage(context.Background(), msg.From, &waProto.Message{
        Conversation: proto.String("❌ Perintah ini hanya untuk owner."),
    })
    return
}

    // Ambil text dari raw GitHub
    resp, err := http.Get("https://raw.githubusercontent.com/arivpnstores/izin/main/listvpn.txt")
    if err != nil {
        _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
            Conversation: proto.String("❌ Gagal mengambil file dari GitHub."),
        })
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
            Conversation: proto.String("❌ Gagal membaca isi file."),
        })
        return
    }

    messageText := string(body) // isi file dari GitHub

    allGroups, err := client.GetJoinedGroups()
    if err != nil {
        logger.Error("Gagal mengambil grup: " + err.Error())
        return
    }

    totalSent := 0

    // Info awal
    _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
        Conversation: proto.String(fmt.Sprintf("Memproses *jpm* ke %d grup...", len(allGroups))),
    })

    for _, group := range allGroups {
        _, err := client.SendMessage(context.Background(), group.JID, &waProto.Message{
            Conversation: proto.String(messageText),
        })
        if err == nil {
            totalSent++
        }
        time.Sleep(20 * time.Second)
    }

    // Info akhir
    _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
        Conversation: proto.String(fmt.Sprintf("*JPM Selesai ✅*\nTotal grup yang berhasil dikirimi pesan: %d", totalSent)),
    })

	case "jpmvpn":
jid := evt.Info.Chat
msg := Serialize(evt)
if !msg.IsOwner {
    _, _ = client.SendMessage(context.Background(), msg.From, &waProto.Message{
        Conversation: proto.String("❌ Perintah ini hanya untuk owner."),
    })
    return
}
    // Ambil text dari raw GitHub
    resp, err := http.Get("https://raw.githubusercontent.com/arivpnstores/izin/main/listvpn.txt")
    if err != nil {
        _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
            Conversation: proto.String("❌ Gagal mengambil file dari GitHub."),
        })
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
            Conversation: proto.String("❌ Gagal membaca isi file."),
        })
        return
    }

    messageText := string(body) // isi file dari GitHub

    allGroups, err := client.GetJoinedGroups()
    if err != nil {
        logger.Error("Gagal mengambil grup: " + err.Error())
        return
    }

    totalSent := 0

    // Info awal
    _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
        Conversation: proto.String(fmt.Sprintf("Memproses *jpm* ke %d grup...", len(allGroups))),
    })

    for _, group := range allGroups {
        _, err := client.SendMessage(context.Background(), group.JID, &waProto.Message{
            Conversation: proto.String(messageText),
        })
        if err == nil {
            totalSent++
        }
        time.Sleep(20 * time.Second)
    }

    // Info akhir
    _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
        Conversation: proto.String(fmt.Sprintf("*JPM Selesai ✅*\nTotal grup yang berhasil dikirimi pesan: %d", totalSent)),
    })

	case "jpmvps":
jid := evt.Info.Chat
msg := Serialize(evt)
if !msg.IsOwner {
    _, _ = client.SendMessage(context.Background(), msg.From, &waProto.Message{
        Conversation: proto.String("❌ Perintah ini hanya untuk owner."),
    })
    return
}
    // Ambil text dari raw GitHub
    resp, err := http.Get("https://raw.githubusercontent.com/arivpnstores/izin/main/listvps.txt")
    if err != nil {
        _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
            Conversation: proto.String("❌ Gagal mengambil file dari GitHub."),
        })
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
            Conversation: proto.String("❌ Gagal membaca isi file."),
        })
        return
    }

    messageText := string(body) // isi file dari GitHub

    allGroups, err := client.GetJoinedGroups()
    if err != nil {
        logger.Error("Gagal mengambil grup: " + err.Error())
        return
    }

    totalSent := 0

    // Info awal
    _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
        Conversation: proto.String(fmt.Sprintf("Memproses *jpm* ke %d grup...", len(allGroups))),
    })

    for _, group := range allGroups {
        _, err := client.SendMessage(context.Background(), group.JID, &waProto.Message{
            Conversation: proto.String(messageText),
        })
        if err == nil {
            totalSent++
        }
        time.Sleep(20 * time.Second)
    }

    // Info akhir
    _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
        Conversation: proto.String(fmt.Sprintf("*JPM Selesai ✅*\nTotal grup yang berhasil dikirimi pesan: %d", totalSent)),
    })

	case "jpmall":
jid := evt.Info.Chat
msg := Serialize(evt)
if !msg.IsOwner {
    _, _ = client.SendMessage(context.Background(), msg.From, &waProto.Message{
        Conversation: proto.String("❌ Perintah ini hanya untuk owner."),
    })
    return
}

    // Ambil text dari raw GitHub
    resp, err := http.Get("https://raw.githubusercontent.com/arivpnstores/izin/main/list.txt")
    if err != nil {
        _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
            Conversation: proto.String("❌ Gagal mengambil file dari GitHub."),
        })
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
            Conversation: proto.String("❌ Gagal membaca isi file."),
        })
        return
    }

    messageText := string(body) // isi file dari GitHub

    allGroups, err := client.GetJoinedGroups()
    if err != nil {
        logger.Error("Gagal mengambil grup: " + err.Error())
        return
    }

    totalSent := 0

    // Info awal
    _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
        Conversation: proto.String(fmt.Sprintf("Memproses *jpm* ke %d grup...", len(allGroups))),
    })

    for _, group := range allGroups {
        _, err := client.SendMessage(context.Background(), group.JID, &waProto.Message{
            Conversation: proto.String(messageText),
        })
        if err == nil {
            totalSent++
        }
        time.Sleep(20 * time.Second)
    }

    // Info akhir
    _, _ = client.SendMessage(context.Background(), jid, &waProto.Message{
        Conversation: proto.String(fmt.Sprintf("*JPM Selesai ✅*\nTotal grup yang berhasil dikirimi pesan: %d", totalSent)),
    })
		
		//--------CASE PING-------//
	case "ping", "uptime":
		jid := evt.Info.Chat
		start := time.Now()
		// Ambil info sistem
		platform := runtime.GOOS
		totalRam := getTotalMemory()
		totalDisk := getTotalDiskSpace()
		cpuCount := runtime.NumCPU()
		uptimeVps := getUptime()
		botUptime := formatDuration(time.Since(appStartTime))
		latency := time.Since(start).Seconds()

		// Format pesan
		msg := fmt.Sprintf(`*🔴 INFORMATION SERVER*

• Platform : %s
• Total Ram : %s
• Total Disk : %s
• Total Cpu : %d Core
• Runtime VPS : %s

*🔵 INFORMATION GOLANG BOT*

• Respon Speed : %.4f detik
• Runtime Bot : %s`,
			platform,
			totalRam,
			totalDisk,
			cpuCount,
			uptimeVps,
			latency,
			botUptime,
		)

		// Kirim pesan ke WhatsApp
		_, err := client.SendMessage(context.Background(), jid, &waProto.Message{
			Conversation: proto.String(msg),
		})
		if err != nil {
			logger.Error("Failed to send uptime reply: " + err.Error())
		}
	}
}

// Fungsi bantu untuk format waktu
func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())

	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%d hari ", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%d jam ", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%d menit ", minutes)
	}
	if seconds > 0 {
		result += fmt.Sprintf("%d detik", seconds)
	}
	return result
}

// RAM total
func getTotalMemory() string {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "Unknown"
	}
	return fmt.Sprintf("%.2f GB", float64(v.Total)/1e9)
}

// Disk total
func getTotalDiskSpace() string {
	d, err := disk.Usage("/")
	if err != nil {
		return "Unknown"
	}
	return fmt.Sprintf("%.2f GB", float64(d.Total)/1e9)
}

// Uptime VPS
func getUptime() string {
	uptimeSec, err := host.Uptime()
	if err != nil {
		return "Unknown"
	}
	return formatDuration(time.Duration(uptimeSec) * time.Second)
}
