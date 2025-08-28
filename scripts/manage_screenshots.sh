#!/bin/bash

# Quell- und Zielverzeichnis
SOURCE_DIR="/mnt/c/Users/denkhaus/Pictures/Screenshots"
TARGET_DIR="/home/denkhaus/dev/gomodules/agents/screenshots"

# Prüfen ob Quellverzeichnis existiert
if [ ! -d "$SOURCE_DIR" ]; then
    echo "Quellverzeichnis $SOURCE_DIR existiert nicht"
    exit 1
fi

# Prüfen ob Zielverzeichnis existiert, falls nicht erstellen
if [ ! -d "$TARGET_DIR" ]; then
    mkdir -p "$TARGET_DIR"
fi

# Kopiere alle Screenshots vom Quell- zum Zielverzeichnis
cp "$SOURCE_DIR"/* "$TARGET_DIR"/ 2>/dev/null

# In das Zielverzeichnis wechseln
cd "$TARGET_DIR" || exit 1

# Alle Dateien nach Änderungsdatum sortieren (neueste zuerst) und in Array speichern
mapfile -t sorted_files < <(ls -t *.png 2>/dev/null)

# Wenn keine PNG-Dateien vorhanden sind, Skript beenden
if [ ${#sorted_files[@]} -eq 0 ]; then
    echo "Keine PNG-Dateien im Zielverzeichnis gefunden"
    exit 0
fi

# Die neueste Datei umbenennen
if [ -f "${sorted_files[0]}" ]; then
    mv "${sorted_files[0]}" "screenshot_current.png"
fi

# Aktualisierte Liste nach dem Umbenennen
mapfile -t sorted_files < <(ls -t *.png 2>/dev/null)

# Zähler für verbleibende Dateien
count=0

# Durch alle Dateien iterieren und alle außer die neuesten 3 löschen
for file in "${sorted_files[@]}"; do
    count=$((count + 1))
    if [ $count -gt 3 ]; then
        rm "$file"
    fi
done

echo "Screenshots erfolgreich verwaltet: Behalte die neuesten 3 Dateien"
