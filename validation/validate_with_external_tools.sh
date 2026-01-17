#!/bin/bash
# Validate DOCX files using external tools (not self-validation)

set -e

OUTPUT_DIR="/home/user/wordZero/validation/output/issues_check"
EXTRACT_DIR="/home/user/wordZero/validation/output/extracted"

echo "=============================================="
echo "External OOXML Validation"
echo "=============================================="

# Create extraction directory
rm -rf "$EXTRACT_DIR"
mkdir -p "$EXTRACT_DIR"

# Function to validate a single docx file
validate_docx() {
    local docx_file="$1"
    local name=$(basename "$docx_file" .docx)
    local extract_path="$EXTRACT_DIR/$name"

    echo ""
    echo "--- Validating: $name.docx ---"

    # 1. Check if it's a valid ZIP
    if ! unzip -t "$docx_file" > /dev/null 2>&1; then
        echo "  [FAIL] Not a valid ZIP archive"
        return 1
    fi
    echo "  [OK] Valid ZIP archive"

    # 2. Extract contents
    mkdir -p "$extract_path"
    unzip -q -o "$docx_file" -d "$extract_path"

    # 3. Check required files exist
    local required_files=(
        "[Content_Types].xml"
        "_rels/.rels"
        "word/document.xml"
    )

    for f in "${required_files[@]}"; do
        if [ -f "$extract_path/$f" ]; then
            echo "  [OK] Found: $f"
        else
            echo "  [FAIL] Missing: $f"
        fi
    done

    # 4. Validate XML with xmllint
    echo "  XML Validation:"
    for xml_file in $(find "$extract_path" -name "*.xml" -type f); do
        local rel_path="${xml_file#$extract_path/}"
        if xmllint --noout "$xml_file" 2>/dev/null; then
            echo "    [OK] $rel_path"
        else
            echo "    [FAIL] $rel_path - Invalid XML"
        fi
    done

    # 5. Check with python-docx
    echo "  python-docx validation:"
    if python3 -c "from docx import Document; Document('$docx_file')" 2>/dev/null; then
        echo "    [OK] python-docx can read file"
    else
        echo "    [WARN] python-docx cannot read file"
    fi

    # 6. Specific checks based on file name
    case "$name" in
        *issue78*)
            echo "  Tab character check:"
            if grep -q $'\t' "$extract_path/word/document.xml" 2>/dev/null; then
                echo "    [OK] Tab characters found in XML"
            elif grep -q '<w:tab/>' "$extract_path/word/document.xml" 2>/dev/null; then
                echo "    [OK] Tab elements (<w:tab/>) found"
            else
                echo "    [WARN] No tabs found - may be issue #78"
            fi
            ;;
        *issue76*rendered*)
            echo "  Header/Footer check:"
            if ls "$extract_path/word/header"*.xml 2>/dev/null | head -1 > /dev/null; then
                echo "    [OK] Header files exist"
            else
                echo "    [FAIL] No header files - issue #76"
            fi
            if ls "$extract_path/word/footer"*.xml 2>/dev/null | head -1 > /dev/null; then
                echo "    [OK] Footer files exist"
            else
                echo "    [FAIL] No footer files - issue #76"
            fi
            ;;
        *issue91*result*)
            echo "  Image replacement check:"
            local img_count=$(ls "$extract_path/word/media/"* 2>/dev/null | wc -l)
            echo "    Images in media folder: $img_count"
            if [ "$img_count" -ge 2 ]; then
                echo "    [OK] Multiple images found (fix working)"
            else
                echo "    [WARN] Expected 2 images, found $img_count"
            fi
            # Check for placeholder remnants
            if grep -q '{{#image' "$extract_path/word/document.xml" 2>/dev/null; then
                echo "    [FAIL] Image placeholder not replaced"
            fi
            ;;
        *issue88*result*)
            echo "  Bold tag check in loop content:"
            # Count bold tags
            local bold_count=$(grep -o '<w:b/>\|<w:b>' "$extract_path/word/document.xml" 2>/dev/null | wc -l)
            echo "    Bold tags found: $bold_count"
            # Check if loop content has bold
            if grep -A5 "Item1\|Item2" "$extract_path/word/document.xml" 2>/dev/null | grep -q '<w:b'; then
                echo "    [FAIL] Loop content is bold - issue #88 not fixed"
            else
                echo "    [OK] Loop content not bold (fix working)"
            fi
            ;;
    esac
}

# Run Go test to generate files
echo "Generating test files..."
cd /home/user/wordZero/validation
export GOTOOLCHAIN=local
go run comprehensive_issue_check.go 2>&1 | grep -E "Created|ERROR"

echo ""
echo "=============================================="
echo "Validating generated files..."
echo "=============================================="

# Validate all generated docx files
for docx in "$OUTPUT_DIR"/*.docx; do
    if [ -f "$docx" ]; then
        validate_docx "$docx"
    fi
done

echo ""
echo "=============================================="
echo "Detailed XML Analysis"
echo "=============================================="

# Show specific XML content for key issues
echo ""
echo "--- Issue #78: Tab characters in document.xml ---"
if [ -f "$EXTRACT_DIR/issue78_original/word/document.xml" ]; then
    echo "Original document tabs:"
    grep -o '<w:t[^>]*>[^<]*</w:t>' "$EXTRACT_DIR/issue78_original/word/document.xml" | head -5
fi
if [ -f "$EXTRACT_DIR/issue78_reopened/word/document.xml" ]; then
    echo "Reopened document tabs:"
    grep -o '<w:t[^>]*>[^<]*</w:t>' "$EXTRACT_DIR/issue78_reopened/word/document.xml" | head -5
fi

echo ""
echo "--- Issue #88: Bold formatting check ---"
if [ -f "$EXTRACT_DIR/issue88_result/word/document.xml" ]; then
    echo "Checking for <w:b> near loop content:"
    grep -B2 -A2 "Item1" "$EXTRACT_DIR/issue88_result/word/document.xml" 2>/dev/null | head -20
fi

echo ""
echo "=============================================="
echo "Validation Complete"
echo "=============================================="
