#!/bin/bash

SCRIPT_PATH="scripts/manage_screenshots.sh"

echo "Starte Screenshot-Management alle 10 Sekunden..."
echo "Drücken Sie STRG+C um das Script zu beenden"
watch -n10 --color $SCRIPT_PATH
