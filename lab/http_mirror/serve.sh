#!/bin/bash
cd "$(dirname "$0")/tmp/fextralife_eldenring"
python -m http.server 8000
