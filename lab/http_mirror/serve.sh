#!/bin/bash
cd "$(dirname "$0")/tmp/fextralife_eldenring/https/eldenring.wiki.fextralife.com"
python -m http.server 3000
