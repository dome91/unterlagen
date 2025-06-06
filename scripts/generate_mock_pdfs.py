#!/usr/bin/env python3
"""
Script to generate mock PDFs for testing document management app
"""

import os
from datetime import datetime, timedelta
from fpdf import FPDF
import random

def create_output_dir():
    """Create testdata/mock_pdfs directory if it doesn't exist"""
    output_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'testdata', 'mock_pdfs')
    os.makedirs(output_dir, exist_ok=True)
    return output_dir

def generate_invoice_pdf(output_dir, invoice_num):
    """Generate a mock invoice PDF"""
    pdf = FPDF()
    pdf.add_page()
    pdf.set_font('Arial', 'B', 16)
    
    # Header
    pdf.cell(0, 10, f'INVOICE #{invoice_num:04d}', 0, 1, 'C')
    pdf.ln(10)
    
    # Company info
    pdf.set_font('Arial', 'B', 12)
    pdf.cell(0, 8, 'Acme Corporation', 0, 1)
    pdf.set_font('Arial', '', 10)
    pdf.cell(0, 6, '123 Business Street', 0, 1)
    pdf.cell(0, 6, 'City, State 12345', 0, 1)
    pdf.cell(0, 6, 'Phone: (555) 123-4567', 0, 1)
    pdf.ln(10)
    
    # Invoice details
    pdf.set_font('Arial', 'B', 10)
    pdf.cell(40, 6, 'Invoice Date:', 0, 0)
    pdf.set_font('Arial', '', 10)
    invoice_date = datetime.now() - timedelta(days=random.randint(1, 90))
    pdf.cell(0, 6, invoice_date.strftime('%Y-%m-%d'), 0, 1)
    
    pdf.set_font('Arial', 'B', 10)
    pdf.cell(40, 6, 'Due Date:', 0, 0)
    pdf.set_font('Arial', '', 10)
    due_date = invoice_date + timedelta(days=30)
    pdf.cell(0, 6, due_date.strftime('%Y-%m-%d'), 0, 1)
    pdf.ln(10)
    
    # Bill to
    pdf.set_font('Arial', 'B', 12)
    pdf.cell(0, 8, 'Bill To:', 0, 1)
    pdf.set_font('Arial', '', 10)
    companies = ['Tech Solutions Inc', 'Global Enterprises', 'Innovation Labs', 'Digital Partners']
    pdf.cell(0, 6, random.choice(companies), 0, 1)
    pdf.cell(0, 6, f'{random.randint(100, 999)} Customer Ave', 0, 1)
    pdf.cell(0, 6, f'City, State {random.randint(10000, 99999)}', 0, 1)
    pdf.ln(10)
    
    # Items table
    pdf.set_font('Arial', 'B', 10)
    pdf.cell(80, 8, 'Description', 1, 0, 'C')
    pdf.cell(30, 8, 'Quantity', 1, 0, 'C')
    pdf.cell(30, 8, 'Unit Price', 1, 0, 'C')
    pdf.cell(30, 8, 'Total', 1, 1, 'C')
    
    pdf.set_font('Arial', '', 10)
    items = [
        ('Software License', random.randint(1, 10), random.randint(100, 500)),
        ('Consulting Services', random.randint(10, 50), random.randint(50, 200)),
        ('Support Package', 1, random.randint(200, 1000))
    ]
    
    subtotal = 0
    for item, qty, price in items:
        total = qty * price
        subtotal += total
        pdf.cell(80, 6, item, 1, 0)
        pdf.cell(30, 6, str(qty), 1, 0, 'C')
        pdf.cell(30, 6, f'${price:.2f}', 1, 0, 'R')
        pdf.cell(30, 6, f'${total:.2f}', 1, 1, 'R')
    
    # Totals
    pdf.ln(5)
    tax = subtotal * 0.08
    total = subtotal + tax
    
    pdf.cell(140, 6, '', 0, 0)
    pdf.cell(20, 6, 'Subtotal:', 0, 0, 'R')
    pdf.cell(20, 6, f'${subtotal:.2f}', 0, 1, 'R')
    
    pdf.cell(140, 6, '', 0, 0)
    pdf.cell(20, 6, 'Tax (8%):', 0, 0, 'R')
    pdf.cell(20, 6, f'${tax:.2f}', 0, 1, 'R')
    
    pdf.cell(140, 6, '', 0, 0)
    pdf.set_font('Arial', 'B', 10)
    pdf.cell(20, 6, 'Total:', 0, 0, 'R')
    pdf.cell(20, 6, f'${total:.2f}', 0, 1, 'R')
    
    filename = f'invoice_{invoice_num:04d}.pdf'
    pdf.output(os.path.join(output_dir, filename))
    return filename

def generate_contract_pdf(output_dir, contract_num):
    """Generate a mock contract PDF"""
    pdf = FPDF()
    pdf.add_page()
    pdf.set_font('Arial', 'B', 16)
    
    # Title
    pdf.cell(0, 15, 'SERVICE AGREEMENT', 0, 1, 'C')
    pdf.ln(10)
    
    # Contract details
    pdf.set_font('Arial', 'B', 12)
    pdf.cell(0, 8, f'Contract #: SA-{contract_num:04d}', 0, 1)
    
    contract_date = datetime.now() - timedelta(days=random.randint(30, 365))
    pdf.cell(0, 8, f'Date: {contract_date.strftime("%B %d, %Y")}', 0, 1)
    pdf.ln(5)
    
    # Parties
    pdf.set_font('Arial', 'B', 11)
    pdf.cell(0, 8, 'PARTIES:', 0, 1)
    pdf.set_font('Arial', '', 10)
    
    pdf.cell(0, 6, 'Provider: Acme Corporation', 0, 1)
    pdf.cell(0, 6, '123 Business Street, City, State 12345', 0, 1)
    pdf.ln(3)
    
    clients = ['TechStart Inc.', 'Global Solutions LLC', 'Innovation Partners', 'Digital Dynamics Corp.']
    pdf.cell(0, 6, f'Client: {random.choice(clients)}', 0, 1)
    pdf.cell(0, 6, f'{random.randint(500, 999)} Client Boulevard, City, State {random.randint(10000, 99999)}', 0, 1)
    pdf.ln(8)
    
    # Terms
    pdf.set_font('Arial', 'B', 11)
    pdf.cell(0, 8, 'TERMS AND CONDITIONS:', 0, 1)
    pdf.set_font('Arial', '', 10)
    
    terms = [
        '1. Services: Provider agrees to deliver software development and consulting services as outlined in the attached Statement of Work.',
        '2. Duration: This agreement shall remain in effect for a period of 12 months from the date of execution.',
        '3. Payment: Client agrees to pay the fees as specified in the payment schedule. Late payments may incur a 1.5% monthly service charge.',
        '4. Confidentiality: Both parties agree to maintain the confidentiality of proprietary information shared during the course of this agreement.',
        '5. Intellectual Property: All work products created under this agreement shall be owned by the Client upon full payment.',
        '6. Termination: Either party may terminate this agreement with 30 days written notice.'
    ]
    
    for term in terms:
        pdf.cell(0, 6, term[:90], 0, 1)
        if len(term) > 90:
            pdf.cell(10, 6, '', 0, 0)
            pdf.cell(0, 6, term[90:], 0, 1)
        pdf.ln(2)
    
    # Signatures
    pdf.ln(15)
    pdf.set_font('Arial', 'B', 11)
    pdf.cell(0, 8, 'SIGNATURES:', 0, 1)
    pdf.ln(10)
    
    pdf.set_font('Arial', '', 10)
    pdf.cell(90, 6, 'Provider:', 0, 0)
    pdf.cell(90, 6, 'Client:', 0, 1)
    pdf.ln(15)
    
    pdf.cell(90, 6, '_' * 30, 0, 0)
    pdf.cell(90, 6, '_' * 30, 0, 1)
    pdf.cell(90, 6, 'John Smith, CEO', 0, 0)
    pdf.cell(90, 6, 'Jane Doe, Director', 0, 1)
    pdf.cell(90, 6, 'Acme Corporation', 0, 0)
    pdf.cell(90, 6, random.choice(clients), 0, 1)
    
    filename = f'contract_SA_{contract_num:04d}.pdf'
    pdf.output(os.path.join(output_dir, filename))
    return filename

def generate_report_pdf(output_dir, report_num):
    """Generate a mock business report PDF"""
    pdf = FPDF()
    pdf.add_page()
    pdf.set_font('Arial', 'B', 16)
    
    # Title
    report_types = ['Quarterly Business Report', 'Annual Financial Summary', 'Project Status Report', 'Market Analysis Report']
    title = random.choice(report_types)
    pdf.cell(0, 15, title.upper(), 0, 1, 'C')
    
    # Date and quarter
    report_date = datetime.now() - timedelta(days=random.randint(1, 180))
    quarter = f'Q{((report_date.month - 1) // 3) + 1} {report_date.year}'
    
    pdf.set_font('Arial', '', 12)
    pdf.cell(0, 8, f'Period: {quarter}', 0, 1, 'C')
    pdf.cell(0, 8, f'Report Date: {report_date.strftime("%B %d, %Y")}', 0, 1, 'C')
    pdf.ln(10)
    
    # Executive Summary
    pdf.set_font('Arial', 'B', 12)
    pdf.cell(0, 8, 'EXECUTIVE SUMMARY', 0, 1)
    pdf.set_font('Arial', '', 10)
    
    summary_text = [
        'This report provides a comprehensive overview of our business performance for the reporting period.',
        'Key highlights include strong revenue growth, improved operational efficiency, and successful',
        'completion of strategic initiatives. Our team has demonstrated exceptional performance across',
        'all major business units, positioning us well for continued growth in the upcoming quarters.'
    ]
    
    for line in summary_text:
        pdf.cell(0, 6, line, 0, 1)
    pdf.ln(8)
    
    # Key Metrics
    pdf.set_font('Arial', 'B', 12)
    pdf.cell(0, 8, 'KEY PERFORMANCE METRICS', 0, 1)
    pdf.set_font('Arial', '', 10)
    
    # Create a simple table
    pdf.set_font('Arial', 'B', 10)
    pdf.cell(60, 8, 'Metric', 1, 0, 'C')
    pdf.cell(40, 8, 'Current Period', 1, 0, 'C')
    pdf.cell(40, 8, 'Previous Period', 1, 0, 'C')
    pdf.cell(30, 8, 'Change', 1, 1, 'C')
    
    pdf.set_font('Arial', '', 10)
    metrics = [
        ('Revenue', f'${random.randint(500, 2000)}K', f'${random.randint(400, 1800)}K'),
        ('Gross Margin', f'{random.randint(35, 65)}%', f'{random.randint(30, 60)}%'),
        ('Customer Acquisition', f'{random.randint(50, 200)}', f'{random.randint(40, 180)}'),
        ('Employee Satisfaction', f'{random.randint(75, 95)}%', f'{random.randint(70, 90)}%')
    ]
    
    for metric, current, previous in metrics:
        change = random.choice(['+5%', '+12%', '+8%', '-2%', '+15%'])
        pdf.cell(60, 6, metric, 1, 0)
        pdf.cell(40, 6, current, 1, 0, 'C')
        pdf.cell(40, 6, previous, 1, 0, 'C')
        pdf.cell(30, 6, change, 1, 1, 'C')
    
    pdf.ln(10)
    
    # Conclusion
    pdf.set_font('Arial', 'B', 12)
    pdf.cell(0, 8, 'CONCLUSION', 0, 1)
    pdf.set_font('Arial', '', 10)
    
    conclusion_text = [
        'The results demonstrate strong business performance and operational excellence.',
        'Moving forward, we will continue to focus on sustainable growth, customer satisfaction,',
        'and strategic market expansion to maintain our competitive advantage.'
    ]
    
    for line in conclusion_text:
        pdf.cell(0, 6, line, 0, 1)
    
    filename = f'report_{report_num:04d}_{quarter.replace(" ", "_")}.pdf'
    pdf.output(os.path.join(output_dir, filename))
    return filename

def generate_manual_pdf(output_dir, manual_num):
    """Generate a mock user manual PDF"""
    try:
        pdf = FPDF()
        pdf.add_page()
        pdf.set_font('Arial', 'B', 18)
        
        # Title
        products = ['Software Installation Guide', 'User Manual v2.1', 'Quick Start Guide', 'API Documentation']
        title = random.choice(products)
        pdf.cell(0, 15, title.upper(), 0, 1, 'C')
        pdf.ln(10)
        
        # Version info
        pdf.set_font('Arial', '', 12)
        pdf.cell(0, 8, f'Version: {random.randint(1, 5)}.{random.randint(0, 9)}.{random.randint(0, 9)}', 0, 1)
        pdf.cell(0, 8, f'Last Updated: {datetime.now().strftime("%B %Y")}', 0, 1)
        pdf.ln(10)
        
        # Introduction section
        pdf.set_font('Arial', 'B', 12)
        pdf.cell(0, 8, '1. INTRODUCTION', 0, 1)
        pdf.set_font('Arial', '', 10)
        
        intro_text = [
            'Welcome to our comprehensive user guide. This document will help you get started',
            'with our software solution and make the most of its powerful features.',
            '',
            'Our platform is designed to streamline your workflow and increase productivity',
            'through intuitive design and robust functionality.'
        ]
        
        for line in intro_text:
            if line:
                pdf.cell(0, 6, line, 0, 1)
            else:
                pdf.ln(3)
        
        filename = f'manual_{manual_num:04d}.pdf'
        pdf.output(os.path.join(output_dir, filename))
        return filename
    except Exception as e:
        print(f"Error generating manual {manual_num}: {e}")
        return None

def generate_receipt_pdf(output_dir, receipt_num):
    """Generate a mock receipt PDF"""
    try:
        pdf = FPDF()
        pdf.add_page()
        pdf.set_font('Arial', 'B', 14)
        
        # Header
        pdf.cell(0, 10, 'PURCHASE RECEIPT', 0, 1, 'C')
        pdf.ln(5)
        
        # Store info
        stores = ['TechMart Electronics', 'Office Supplies Plus', 'Digital Solutions Store', 'Business Equipment Co.']
        store = random.choice(stores)
        
        pdf.set_font('Arial', 'B', 12)
        pdf.cell(0, 8, store, 0, 1, 'C')
        pdf.set_font('Arial', '', 10)
        pdf.cell(0, 6, f'{random.randint(100, 999)} Retail Street', 0, 1, 'C')
        pdf.cell(0, 6, f'City, State {random.randint(10000, 99999)}', 0, 1, 'C')
        pdf.ln(10)
        
        # Receipt details
        purchase_date = datetime.now() - timedelta(days=random.randint(1, 30))
        pdf.set_font('Arial', 'B', 10)
        pdf.cell(40, 6, 'Receipt #:', 0, 0)
        pdf.set_font('Arial', '', 10)
        pdf.cell(0, 6, f'R{receipt_num:06d}', 0, 1)
        
        pdf.set_font('Arial', 'B', 10)
        pdf.cell(40, 6, 'Date:', 0, 0)
        pdf.set_font('Arial', '', 10)
        pdf.cell(0, 6, purchase_date.strftime('%Y-%m-%d'), 0, 1)
        pdf.ln(8)
        
        # Items
        pdf.set_font('Arial', 'B', 10)
        pdf.cell(80, 6, 'Item', 1, 0, 'C')
        pdf.cell(20, 6, 'Qty', 1, 0, 'C')
        pdf.cell(30, 6, 'Price', 1, 0, 'C')
        pdf.cell(30, 6, 'Total', 1, 1, 'C')
        
        pdf.set_font('Arial', '', 9)
        items = [
            ('Wireless Mouse', 1, 29.99),
            ('USB Cable', 2, 12.99),
            ('Notebook', 1, 8.50)
        ]
        
        subtotal = 0
        for item, qty, price in items:
            total = qty * price
            subtotal += total
            pdf.cell(80, 5, item, 1, 0)
            pdf.cell(20, 5, str(qty), 1, 0, 'C')
            pdf.cell(30, 5, f'${price:.2f}', 1, 0, 'R')
            pdf.cell(30, 5, f'${total:.2f}', 1, 1, 'R')
        
        # Totals
        pdf.ln(3)
        tax = subtotal * 0.08
        total = subtotal + tax
        
        pdf.cell(100, 5, '', 0, 0)
        pdf.cell(30, 5, 'Subtotal:', 0, 0, 'R')
        pdf.cell(30, 5, f'${subtotal:.2f}', 0, 1, 'R')
        
        pdf.cell(100, 5, '', 0, 0)
        pdf.cell(30, 5, 'Tax (8%):', 0, 0, 'R')
        pdf.cell(30, 5, f'${tax:.2f}', 0, 1, 'R')
        
        pdf.set_font('Arial', 'B', 9)
        pdf.cell(100, 5, '', 0, 0)
        pdf.cell(30, 5, 'Total:', 0, 0, 'R')
        pdf.cell(30, 5, f'${total:.2f}', 0, 1, 'R')
        
        filename = f'receipt_{receipt_num:06d}.pdf'
        pdf.output(os.path.join(output_dir, filename))
        return filename
    except Exception as e:
        print(f"Error generating receipt {receipt_num}: {e}")
        return None

def main():
    """Generate all mock PDFs"""
    output_dir = create_output_dir()
    print(f"Creating mock PDFs in: {output_dir}")
    
    generated_files = []
    
    # Generate various types of documents
    for i in range(1, 6):  # 5 invoices
        filename = generate_invoice_pdf(output_dir, i)
        generated_files.append(filename)
        print(f"Generated: {filename}")
    
    for i in range(1, 4):  # 3 contracts
        filename = generate_contract_pdf(output_dir, i)
        generated_files.append(filename)
        print(f"Generated: {filename}")
    
    for i in range(1, 4):  # 3 reports
        filename = generate_report_pdf(output_dir, i)
        generated_files.append(filename)
        print(f"Generated: {filename}")
    
    for i in range(1, 3):  # 2 manuals
        filename = generate_manual_pdf(output_dir, i)
        if filename:
            generated_files.append(filename)
            print(f"Generated: {filename}")
    
    for i in range(1, 8):  # 7 receipts
        filename = generate_receipt_pdf(output_dir, i)
        if filename:
            generated_files.append(filename)
            print(f"Generated: {filename}")
    
    print(f"\nTotal files generated: {len(generated_files)}")
    print(f"Output directory: {output_dir}")

if __name__ == "__main__":
    main()