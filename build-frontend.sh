#!/bin/bash
set -e
cd "$(git rev-parse --show-toplevel)"
vp run build
