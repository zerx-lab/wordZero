#!/usr/bin/env python3
"""
DOCX Open XML Specification Validator
This script validates DOCX files against the Open XML standard.
"""

import zipfile
import os
import sys
import json
from lxml import etree
from pathlib import Path

# OOXML namespaces
NAMESPACES = {
    'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main',
    'r': 'http://schemas.openxmlformats.org/officeDocument/2006/relationships',
    'wp': 'http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing',
    'a': 'http://schemas.openxmlformats.org/drawingml/2006/main',
    'pic': 'http://schemas.openxmlformats.org/drawingml/2006/picture',
    'c': 'http://schemas.openxmlformats.org/drawingml/2006/chart',
    'ct': 'http://schemas.openxmlformats.org/package/2006/content-types',
    'pr': 'http://schemas.openxmlformats.org/package/2006/relationships',
    'm': 'http://schemas.openxmlformats.org/officeDocument/2006/math',
    'v': 'urn:schemas-microsoft-com:vml',
    'o': 'urn:schemas-microsoft-com:office:office',
    'wpc': 'http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas',
    'mc': 'http://schemas.openxmlformats.org/markup-compatibility/2006',
    'wp14': 'http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing',
}

# Required files in a valid DOCX
REQUIRED_FILES = [
    '[Content_Types].xml',
    '_rels/.rels',
    'word/document.xml',
]

# Required relationship types
REQUIRED_RELATIONSHIP_TYPES = [
    'http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument',
]


class ValidationResult:
    def __init__(self, file_path):
        self.file_path = file_path
        self.is_valid = True
        self.errors = []
        self.warnings = []
        self.info = []
        self.document_stats = {}

    def add_error(self, message, component=None):
        self.is_valid = False
        self.errors.append({"component": component, "message": message})

    def add_warning(self, message, component=None):
        self.warnings.append({"component": component, "message": message})

    def add_info(self, message, component=None):
        self.info.append({"component": component, "message": message})

    def to_dict(self):
        return {
            "file": self.file_path,
            "valid": self.is_valid,
            "errors": self.errors,
            "warnings": self.warnings,
            "info": self.info,
            "stats": self.document_stats
        }

    def __str__(self):
        lines = [f"\n{'='*60}", f"Validation Report: {self.file_path}", '='*60]

        if self.is_valid:
            lines.append("✓ Document is VALID")
        else:
            lines.append("✗ Document is INVALID")

        if self.errors:
            lines.append(f"\nErrors ({len(self.errors)}):")
            for err in self.errors:
                comp = f"[{err['component']}] " if err['component'] else ""
                lines.append(f"  ✗ {comp}{err['message']}")

        if self.warnings:
            lines.append(f"\nWarnings ({len(self.warnings)}):")
            for warn in self.warnings:
                comp = f"[{warn['component']}] " if warn['component'] else ""
                lines.append(f"  ⚠ {comp}{warn['message']}")

        if self.info:
            lines.append(f"\nInfo ({len(self.info)}):")
            for info in self.info:
                comp = f"[{info['component']}] " if info['component'] else ""
                lines.append(f"  ℹ {comp}{info['message']}")

        if self.document_stats:
            lines.append("\nDocument Statistics:")
            for key, value in self.document_stats.items():
                lines.append(f"  • {key}: {value}")

        lines.append('='*60)
        return '\n'.join(lines)


class DOCXValidator:
    def __init__(self, file_path):
        self.file_path = file_path
        self.result = ValidationResult(file_path)
        self.zf = None
        self.content_types = {}
        self.relationships = {}

    def validate(self):
        """Run all validation checks"""
        try:
            # Check if file exists
            if not os.path.exists(self.file_path):
                self.result.add_error(f"File does not exist: {self.file_path}")
                return self.result

            # Check file extension
            if not self.file_path.lower().endswith('.docx'):
                self.result.add_warning("File does not have .docx extension")

            # Open as ZIP
            try:
                self.zf = zipfile.ZipFile(self.file_path, 'r')
            except zipfile.BadZipFile:
                self.result.add_error("File is not a valid ZIP archive")
                return self.result

            # Run all validation checks
            self._validate_structure()
            self._validate_content_types()
            self._validate_relationships()
            self._validate_document_xml()
            self._validate_styles()
            self._validate_headers_footers()
            self._validate_images()
            self._validate_tables()
            self._validate_footnotes()
            self._validate_numbering()
            self._collect_statistics()

            self.zf.close()

        except Exception as e:
            self.result.add_error(f"Unexpected error during validation: {str(e)}")

        return self.result

    def _validate_structure(self):
        """Validate required OOXML structure"""
        file_list = self.zf.namelist()

        for required in REQUIRED_FILES:
            if required not in file_list:
                self.result.add_error(f"Missing required file: {required}", "Structure")

        self.result.add_info(f"Package contains {len(file_list)} files", "Structure")

    def _validate_content_types(self):
        """Validate [Content_Types].xml"""
        try:
            content = self.zf.read('[Content_Types].xml')
            tree = etree.fromstring(content)

            # Check for required content types
            defaults = tree.findall('.//{%s}Default' % NAMESPACES['ct'])
            overrides = tree.findall('.//{%s}Override' % NAMESPACES['ct'])

            self.result.add_info(f"Found {len(defaults)} default types, {len(overrides)} overrides", "ContentTypes")

            # Check for document.xml content type
            has_doc_type = False
            for override in overrides:
                part = override.get('PartName')
                if part == '/word/document.xml':
                    has_doc_type = True
                    break

            if not has_doc_type:
                self.result.add_warning("No explicit content type for /word/document.xml", "ContentTypes")

        except Exception as e:
            self.result.add_error(f"Failed to parse [Content_Types].xml: {str(e)}", "ContentTypes")

    def _validate_relationships(self):
        """Validate relationship files"""
        # Root relationships
        try:
            content = self.zf.read('_rels/.rels')
            tree = etree.fromstring(content)

            rels = tree.findall('.//{%s}Relationship' % NAMESPACES['pr'])

            has_office_doc = False
            for rel in rels:
                rel_type = rel.get('Type')
                if 'officeDocument' in rel_type:
                    has_office_doc = True
                    break

            if not has_office_doc:
                self.result.add_error("Missing officeDocument relationship", "Relationships")

            self.result.add_info(f"Found {len(rels)} root relationships", "Relationships")

        except Exception as e:
            self.result.add_error(f"Failed to parse _rels/.rels: {str(e)}", "Relationships")

        # Document relationships
        try:
            if 'word/_rels/document.xml.rels' in self.zf.namelist():
                content = self.zf.read('word/_rels/document.xml.rels')
                tree = etree.fromstring(content)

                rels = tree.findall('.//{%s}Relationship' % NAMESPACES['pr'])
                self.result.add_info(f"Found {len(rels)} document relationships", "Relationships")

                # Store relationships for later validation
                for rel in rels:
                    self.relationships[rel.get('Id')] = {
                        'type': rel.get('Type'),
                        'target': rel.get('Target')
                    }
        except Exception as e:
            self.result.add_warning(f"Failed to parse document relationships: {str(e)}", "Relationships")

    def _validate_document_xml(self):
        """Validate word/document.xml"""
        try:
            content = self.zf.read('word/document.xml')
            tree = etree.fromstring(content)

            # Check root element
            root_tag = etree.QName(tree.tag)
            if root_tag.localname != 'document':
                self.result.add_error(f"Root element should be 'document', got '{root_tag.localname}'", "Document")

            # Check for body
            body = tree.find('.//{%s}body' % NAMESPACES['w'])
            if body is None:
                self.result.add_error("Missing body element", "Document")
            else:
                self.result.add_info("Found body element", "Document")

            # Check for section properties (required for valid document)
            sectPr = tree.find('.//{%s}sectPr' % NAMESPACES['w'])
            if sectPr is None:
                self.result.add_warning("Missing section properties (sectPr)", "Document")
            else:
                self.result.add_info("Found section properties", "Document")

            # Check paragraphs
            paragraphs = tree.findall('.//{%s}p' % NAMESPACES['w'])
            self.result.document_stats['paragraphs'] = len(paragraphs)

            # Check runs
            runs = tree.findall('.//{%s}r' % NAMESPACES['w'])
            self.result.document_stats['runs'] = len(runs)

            # Check for text content
            texts = tree.findall('.//{%s}t' % NAMESPACES['w'])
            total_text = ''.join([t.text or '' for t in texts])
            self.result.document_stats['text_length'] = len(total_text)

        except etree.XMLSyntaxError as e:
            self.result.add_error(f"Invalid XML in document.xml: {str(e)}", "Document")
        except Exception as e:
            self.result.add_error(f"Failed to validate document.xml: {str(e)}", "Document")

    def _validate_styles(self):
        """Validate word/styles.xml"""
        if 'word/styles.xml' not in self.zf.namelist():
            self.result.add_warning("Missing styles.xml", "Styles")
            return

        try:
            content = self.zf.read('word/styles.xml')
            tree = etree.fromstring(content)

            # Check root element
            root_tag = etree.QName(tree.tag)
            if root_tag.localname != 'styles':
                self.result.add_error(f"Root element should be 'styles', got '{root_tag.localname}'", "Styles")

            # Count styles
            styles = tree.findall('.//{%s}style' % NAMESPACES['w'])
            self.result.document_stats['styles'] = len(styles)

            # Check for Normal style
            has_normal = False
            for style in styles:
                style_id = style.get('{%s}styleId' % NAMESPACES['w'])
                if style_id == 'Normal':
                    has_normal = True
                    break

            if not has_normal:
                self.result.add_warning("Missing 'Normal' style definition", "Styles")

            self.result.add_info(f"Found {len(styles)} style definitions", "Styles")

        except etree.XMLSyntaxError as e:
            self.result.add_error(f"Invalid XML in styles.xml: {str(e)}", "Styles")
        except Exception as e:
            self.result.add_error(f"Failed to validate styles.xml: {str(e)}", "Styles")

    def _validate_headers_footers(self):
        """Validate headers and footers"""
        headers = [f for f in self.zf.namelist() if f.startswith('word/header') and f.endswith('.xml')]
        footers = [f for f in self.zf.namelist() if f.startswith('word/footer') and f.endswith('.xml')]

        self.result.document_stats['headers'] = len(headers)
        self.result.document_stats['footers'] = len(footers)

        for header in headers:
            try:
                content = self.zf.read(header)
                tree = etree.fromstring(content)
                root_tag = etree.QName(tree.tag)
                if root_tag.localname != 'hdr':
                    self.result.add_error(f"Header {header} has invalid root element", "Headers")
            except Exception as e:
                self.result.add_error(f"Failed to parse {header}: {str(e)}", "Headers")

        for footer in footers:
            try:
                content = self.zf.read(footer)
                tree = etree.fromstring(content)
                root_tag = etree.QName(tree.tag)
                if root_tag.localname != 'ftr':
                    self.result.add_error(f"Footer {footer} has invalid root element", "Footers")
            except Exception as e:
                self.result.add_error(f"Failed to parse {footer}: {str(e)}", "Footers")

    def _validate_images(self):
        """Validate embedded images"""
        media_files = [f for f in self.zf.namelist() if f.startswith('word/media/')]

        self.result.document_stats['images'] = len(media_files)

        for media in media_files:
            try:
                data = self.zf.read(media)
                if len(data) == 0:
                    self.result.add_error(f"Empty media file: {media}", "Images")

                # Check image magic bytes
                if media.lower().endswith(('.png', '.jpg', '.jpeg', '.gif')):
                    is_valid = False
                    if data[:8] == b'\x89PNG\r\n\x1a\n':  # PNG
                        is_valid = True
                    elif data[:2] == b'\xff\xd8':  # JPEG
                        is_valid = True
                    elif data[:6] in (b'GIF87a', b'GIF89a'):  # GIF
                        is_valid = True

                    if not is_valid:
                        self.result.add_warning(f"Image format may not match extension: {media}", "Images")

            except Exception as e:
                self.result.add_error(f"Failed to validate {media}: {str(e)}", "Images")

        if media_files:
            self.result.add_info(f"Found {len(media_files)} media files", "Images")

    def _validate_tables(self):
        """Validate tables in document"""
        try:
            content = self.zf.read('word/document.xml')
            tree = etree.fromstring(content)

            tables = tree.findall('.//{%s}tbl' % NAMESPACES['w'])
            self.result.document_stats['tables'] = len(tables)

            for i, table in enumerate(tables):
                # Check table properties
                tblPr = table.find('./{%s}tblPr' % NAMESPACES['w'])

                # Check rows
                rows = table.findall('./{%s}tr' % NAMESPACES['w'])
                if len(rows) == 0:
                    self.result.add_warning(f"Table {i+1} has no rows", "Tables")

                # Check cell structure
                for j, row in enumerate(rows):
                    cells = row.findall('./{%s}tc' % NAMESPACES['w'])
                    if len(cells) == 0:
                        self.result.add_warning(f"Table {i+1} Row {j+1} has no cells", "Tables")

                    # Check for cell merging issues
                    for k, cell in enumerate(cells):
                        tcPr = cell.find('./{%s}tcPr' % NAMESPACES['w'])
                        if tcPr is not None:
                            gridSpan = tcPr.find('./{%s}gridSpan' % NAMESPACES['w'])
                            vMerge = tcPr.find('./{%s}vMerge' % NAMESPACES['w'])
                            # These are valid elements, just track them

            if tables:
                self.result.add_info(f"Found {len(tables)} tables", "Tables")

        except Exception as e:
            self.result.add_error(f"Failed to validate tables: {str(e)}", "Tables")

    def _validate_footnotes(self):
        """Validate footnotes and endnotes"""
        for file_name, file_type in [('word/footnotes.xml', 'Footnotes'), ('word/endnotes.xml', 'Endnotes')]:
            if file_name in self.zf.namelist():
                try:
                    content = self.zf.read(file_name)
                    tree = etree.fromstring(content)

                    expected_root = 'footnotes' if 'footnote' in file_name else 'endnotes'
                    root_tag = etree.QName(tree.tag)
                    if root_tag.localname != expected_root:
                        self.result.add_error(f"Invalid root element in {file_name}", file_type)

                    notes_tag = 'footnote' if 'footnote' in file_name else 'endnote'
                    notes = tree.findall('.//{%s}%s' % (NAMESPACES['w'], notes_tag))
                    self.result.document_stats[file_type.lower()] = len(notes)

                    self.result.add_info(f"Found {len(notes)} {file_type.lower()}", file_type)

                except Exception as e:
                    self.result.add_error(f"Failed to parse {file_name}: {str(e)}", file_type)

    def _validate_numbering(self):
        """Validate numbering.xml (lists)"""
        if 'word/numbering.xml' not in self.zf.namelist():
            return

        try:
            content = self.zf.read('word/numbering.xml')
            tree = etree.fromstring(content)

            root_tag = etree.QName(tree.tag)
            if root_tag.localname != 'numbering':
                self.result.add_error("Invalid root element in numbering.xml", "Numbering")

            # Check abstract numbering definitions
            abstractNum = tree.findall('.//{%s}abstractNum' % NAMESPACES['w'])
            num = tree.findall('.//{%s}num' % NAMESPACES['w'])

            self.result.document_stats['numbering_definitions'] = len(abstractNum)
            self.result.add_info(f"Found {len(abstractNum)} abstract numbering definitions, {len(num)} instances", "Numbering")

        except Exception as e:
            self.result.add_error(f"Failed to validate numbering.xml: {str(e)}", "Numbering")

    def _collect_statistics(self):
        """Collect additional document statistics"""
        # Calculate approximate word count
        try:
            content = self.zf.read('word/document.xml')
            tree = etree.fromstring(content)
            texts = tree.findall('.//{%s}t' % NAMESPACES['w'])
            total_text = ' '.join([t.text or '' for t in texts])
            words = total_text.split()
            self.result.document_stats['word_count'] = len(words)
        except:
            pass


def validate_file(file_path):
    """Validate a single DOCX file"""
    validator = DOCXValidator(file_path)
    return validator.validate()


def validate_directory(dir_path):
    """Validate all DOCX files in a directory"""
    results = []
    path = Path(dir_path)

    for docx_file in path.glob('**/*.docx'):
        result = validate_file(str(docx_file))
        results.append(result)

    return results


def main():
    import argparse

    parser = argparse.ArgumentParser(description='DOCX Open XML Validator')
    parser.add_argument('path', help='DOCX file or directory to validate')
    parser.add_argument('--json', action='store_true', help='Output as JSON')
    parser.add_argument('--quiet', action='store_true', help='Only show errors')

    args = parser.parse_args()

    if os.path.isdir(args.path):
        results = validate_directory(args.path)
    else:
        results = [validate_file(args.path)]

    if args.json:
        output = [r.to_dict() for r in results]
        print(json.dumps(output, indent=2))
    else:
        for result in results:
            if args.quiet and result.is_valid:
                continue
            print(result)

    # Exit with error code if any validation failed
    if any(not r.is_valid for r in results):
        sys.exit(1)

    # Summary
    if not args.json and len(results) > 1:
        valid_count = sum(1 for r in results if r.is_valid)
        print(f"\nSummary: {valid_count}/{len(results)} files valid")


if __name__ == '__main__':
    main()
