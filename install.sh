#!/usr/bin/env bash

set -e  # Abort on errors

# Define messages for English and Polish
declare -A messages
messages[en_start]="Starting installation..."
messages[pl_start]="Rozpoczynam instalację..."
messages[en_abort]="Installation aborted. Cleaning up..."
messages[pl_abort]="Instalacja przerwana. Czyszczenie..."
messages[en_cleaned]="Cleanup complete."
messages[pl_cleaned]="Czyszczenie zakończone."
messages[en_error]="Error:"
messages[pl_error]="Błąd:"
messages[en_requires_root]="This script must be run as root. Please use 'sudo'."
messages[pl_requires_root]="Ten skrypt musi być uruchomiony jako root. Użyj 'sudo'."
messages[en_missing_dep]="is not installed. Please install it first."
messages[pl_missing_dep]="nie jest zainstalowany. Proszę najpierw to zainstalować."
messages[en_dir_exists]="Directory already exists. Skipping creation."
messages[pl_dir_exists]="Katalog już istnieje. Pomijam tworzenie."
messages[en_dir_created]="Directory created successfully."
messages[pl_dir_created]="Katalog utworzony pomyślnie."
messages[en_installed]="Program installed successfully. You can now edit your subscribed channels in"
messages[pl_installed]="Program został pomyślnie zainstalowany. Możesz teraz edytować swoje subskrybowane kanały w"
messages[en_reminder]="Before you run the program for the first time, ensure you have the packages installed: mpv, yt-dlp, and youtube-dl."
messages[pl_reminder]="Przed pierwszym uruchomieniem programu upewnij się, że masz zainstalowane pakiety: mpv, yt-dlp i youtube-dl."

# Detect system language
LANGUAGE=$(locale | grep LANG= | cut -d= -f2 | cut -d_ -f1)
if [[ "$LANGUAGE" != "pl" ]]; then
    LANGUAGE="en"
fi

# Cleanup function to handle Ctrl+C
cleanup() {
    echo
    echo "${messages[${LANGUAGE}_abort]}"
    [[ -n "$TEMP_DIR" ]] && sudo rm -rf "$TEMP_DIR"
    echo "${messages[${LANGUAGE}_cleaned]}"
    exit 1
}

# Error handling function
error_exit() {
    echo "${messages[${LANGUAGE}_error]} $1" >&2
    [[ -n "$TEMP_DIR" ]] && sudo rm -rf "$TEMP_DIR"
    exit 1
}

# Trap Ctrl+C
trap cleanup INT

echo "${messages[${LANGUAGE}_start]}"

# Ensure script is run as root
if [[ $EUID -ne 0 ]]; then
    error_exit "${messages[${LANGUAGE}_requires_root]}"
fi

# Variables
TARGET_DIR="/etc/yt-cli"
BIN_PATH="/usr/local/bin/yt-cli"

# Check if necessary tools are installed
for cmd in mpv yt-dlp youtube-dl; do
    if ! command -v "$cmd" &>/dev/null; then
        error_exit "$cmd ${messages[${LANGUAGE}_missing_dep]}"
    fi
done

# Create a temporary directory for installation
TEMP_DIR=$(mktemp -d) || error_exit "Failed to create temporary directory."

# Copy files to temporary directory
cp channels.go "$TEMP_DIR" || error_exit "Failed to copy channels.go."
cp framework.go "$TEMP_DIR" || error_exit "Failed to copy framework.go."
cp main.go "$TEMP_DIR" || error_exit "Failed to copy main.go."
cp go.mod "$TEMP_DIR" || error_exit "Failed to copy go.mod."
cp go.sum "$TEMP_DIR" || error_exit "Failed to copy go.sum."
cp yt-cli "$TEMP_DIR" || error_exit "Failed to copy yt-cli."

# Move to target directory
if [[ -d "$TARGET_DIR" ]]; then
    echo "${messages[${LANGUAGE}_dir_exists]}"
else
    sudo mkdir "$TARGET_DIR" || error_exit "Failed to create directory $TARGET_DIR."
    echo "${messages[${LANGUAGE}_dir_created]}"
fi

sudo mv "$TEMP_DIR"/* "$TARGET_DIR" || error_exit "Failed to move files to $TARGET_DIR."

# Set permissions and install binary
sudo chmod +x "$TARGET_DIR/yt-cli" || error_exit "Failed to set permissions for yt-cli."
sudo ln -sf "$TARGET_DIR/yt-cli" "$BIN_PATH" || error_exit "Failed to link yt-cli to $BIN_PATH."

# Cleanup temporary directory
sudo rm -rf "$TEMP_DIR"

echo "${messages[${LANGUAGE}_installed]} $TARGET_DIR/channels.go"
echo "${messages[${LANGUAGE}_reminder]}"
