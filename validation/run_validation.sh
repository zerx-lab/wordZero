#!/bin/bash
# WordZero Comprehensive Validation Script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/output"

echo "=========================================="
echo "WordZero Feature Validation Suite"
echo "=========================================="

# Create output directory
mkdir -p "${OUTPUT_DIR}"

# Copy validator to output directory
cp "${SCRIPT_DIR}/docx_validator.py" "${OUTPUT_DIR}/"

# Run the Go test
echo ""
echo "Running comprehensive Go tests..."
cd "${SCRIPT_DIR}"

export GOTOOLCHAIN=local
go run comprehensive_test.go

echo ""
echo "=========================================="
echo "Running detailed OOXML validation..."
echo "=========================================="

# Run Python validator on all generated documents
python3 "${OUTPUT_DIR}/docx_validator.py" "${OUTPUT_DIR}" 2>&1 || true

echo ""
echo "Validation complete!"
