#!/usr/bin/env python3
"""
Summary script to display all generated mock PDFs
"""

from pathlib import Path

def categorize_file(name, size):
    """Categorize a PDF file based on its name."""
    if name.startswith('invoice_'):
        return 'Invoices', (name, size)
    elif name.startswith('contract_'):
        return 'Contracts', (name, size)
    elif name.startswith('report_') and not name.startswith('monthly_report') and not name.startswith('annual_report'):
        return 'Reports', (name, size)
    elif name.startswith('manual_'):
        return 'Manuals', (name, size)
    elif name.startswith('receipt_'):
        return 'Receipts', (name, size)
    elif name.startswith('large_document'):
        return 'Large Documents', (name, size)
    elif name.startswith('confidential_'):
        return 'Confidential', (name, size)
    elif name.startswith('application_form'):
        return 'Forms', (name, size)
    elif name.startswith('financial_statement'):
        return 'Financial Statements', (name, size)
    elif name.startswith('presentation_'):
        return 'Presentations', (name, size)
    elif name.startswith('nda_'):
        return 'Legal Documents', (name, size)
    elif name.startswith('monthly_report'):
        return 'Monthly Reports', (name, size)
    elif name.startswith('weekly_status'):
        return 'Weekly Status', (name, size)
    elif name.startswith('quarterly_review'):
        return 'Quarterly Reviews', (name, size)
    elif name.startswith('annual_report') or name.startswith('tax_summary') or name.startswith('compliance_review'):
        return 'Annual Documents', (name, size)
    return None, None

def display_categories(categories):
    """Display categorized files."""
    for category, files in categories.items():
        if files:
            print(f"📁 {category} ({len(files)} files)")
            print(f"   {'-'*40}")
            for filename, size in files[:5]:  # Show first 5 files
                size_str = f"{size:,} bytes" if size < 1024 else f"{size/1024:.1f}KB"
                print(f"   • {filename:<35} ({size_str})")
            if len(files) > 5:
                print(f"   • ... and {len(files)-5} more files")
            print()

def display_file_statistics(pdf_files):
    """Display file size statistics."""
    sizes = [f.stat().st_size for f in pdf_files]
    total_size = sum(sizes)
    avg_size = total_size / len(sizes) if sizes else 0
    max_size = max(sizes) if sizes else 0
    min_size = min(sizes) if sizes else 0

    print("📊 File Size Statistics")
    print(f"   {'='*30}")
    print(f"   Total Size: {total_size/1024:.1f}KB")
    print(f"   Average Size: {avg_size/1024:.1f}KB")
    print(f"   Largest File: {max_size/1024:.1f}KB")
    print(f"   Smallest File: {min_size/1024:.1f}KB")
    print()

def display_testing_capabilities(categories):
    """Display testing capabilities."""
    print("🧪 Testing Capabilities")
    print(f"   {'='*30}")
    print(f"   ✓ Document type classification ({len([c for c in categories.values() if c])} types)")
    print("   ✓ Date-based filtering (monthly, weekly, quarterly, annual)")
    print("   ✓ Content search across varied document structures")
    print("   ✓ File size handling (1KB - 50KB range)")
    print("   ✓ Bulk processing scenarios")
    print("   ✓ Metadata extraction testing")
    print("   ✓ Performance benchmarking")
    print()

def main():
    # Get the mock PDFs directory
    script_dir = Path(__file__).parent
    pdf_dir = script_dir.parent / 'testdata' / 'mock_pdfs'

    if not pdf_dir.exists():
        print("Mock PDFs directory not found!")
        return

    # Get all PDF files
    pdf_files = list(pdf_dir.glob('*.pdf'))
    pdf_files.sort()

    print("📄 Mock PDF Test Documents Summary")
    print(f"{'='*50}")
    print(f"Directory: {pdf_dir}")
    print(f"Total PDFs: {len(pdf_files)}")
    print()

    # Categorize files
    categories = {
        'Invoices': [],
        'Contracts': [],
        'Reports': [],
        'Manuals': [],
        'Receipts': [],
        'Large Documents': [],
        'Confidential': [],
        'Forms': [],
        'Financial Statements': [],
        'Presentations': [],
        'Legal Documents': [],
        'Monthly Reports': [],
        'Weekly Status': [],
        'Quarterly Reviews': [],
        'Annual Documents': []
    }

    for pdf_file in pdf_files:
        name = pdf_file.name
        size = pdf_file.stat().st_size
        category, file_info = categorize_file(name, size)
        if category:
            categories[category].append(file_info)

    # Display categories
    display_categories(categories)

    # File size distribution
    display_file_statistics(pdf_files)

    # Testing capabilities
    display_testing_capabilities(categories)

    print("🚀 Ready for document management app testing!")

if __name__ == "__main__":
    main()
