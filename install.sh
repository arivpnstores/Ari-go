#!/bin/bash
# =========================================
# Menu Install Golang + Ari-go + systemctl
# Author : Ari Setiawan 
# =========================================

GO_VERSION="1.22.0"
SERVICE_PATH="/etc/systemd/system/ari-go.service"

install_golang() {
    echo "=== Installing Dependencies ==="
    apt update -y
    apt upgrade -y
    apt install wget curl build-essential gcc -y

    echo "=== Installing Golang v$GO_VERSION ==="
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin

    go version
    echo "=== Golang installation complete! ==="
}

install_ari-go() {
    git clone https://github.com/arivpnstores/Ari-go.git
    cd Ari-go 
    echo "=== Installing Go Modules ==="
    go get
    echo "=== Ari-go installation complete! ==="
}

start_ari-go() {
    echo "=== Menjalankan Ari-go (Tekan CTRL+C untuk stop) ==="
    /usr/local/go/bin/go run main.go
}

setup_systemctl() {
    echo "=== Creating systemd service ==="
    cat > $SERVICE_PATH << EOF
[Unit]
Description=ari-go Go WhatsApp Bot
After=network.target

[Service]
Type=simple
WorkingDirectory=/root/Ari-go
ExecStart=/usr/local/go/bin/go run main.go
Restart=on-failure
User=root

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reexec
    systemctl daemon-reload
    systemctl enable ari-go.service
    systemctl start ari-go.service

    echo "=== systemctl service created & started ==="
    echo "Cek status: systemctl status ari-go.service"
    echo "Cek log   : journalctl -u ari-go.service -f"
}

AUTH_FILE="root/Ari-go/.auth"

# Cek apakah file auth ada
if [[ -f "$AUTH_FILE" ]]; then
    echo "Admin sudah terverifikasi."
else
    # Minta password
    read -p "   VERIFIKASI SCRIPT :   " ari
    if [[ $ari == "ScGolang30k" ]]; then
        echo "Terverifikasi!"
        # Buat file auth supaya tidak perlu input lagi
        mkdir -p "root/Ari-go"
        touch "$AUTH_FILE"
    else
        echo "Password salah!"
        exit 1
    fi
fi
# Menu
while true; do
    clear
    echo "==============================="
    echo "      MENU INSTALL BOT GO"
    echo "==============================="
    echo "1. Install Golang ($GO_VERSION)"
    echo "2. Install Ari-go"
    echo "3. Start Ari-go (Manual)"
    echo "4. Buat & Jalankan systemctl"
    echo "0. Keluar"
    echo "==============================="
    read -p "Pilih menu [0-4]: " choice

    case $choice in
        1) install_golang ;;
        2) install_ari-go ;;
        3) start_ari-go ;;
        4) setup_systemctl ;;
        0) echo "Keluar..."; exit 0 ;;
        *) echo "Pilihan tidak valid!" ;;
    esac
    read -p "Tekan Enter untuk kembali ke menu..."
done
