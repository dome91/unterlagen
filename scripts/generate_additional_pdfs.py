#!/usr/bin/env python3
"""
Script to generate additional mock PDFs with various characteristics for testing
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

def generate_large_document_pdf(output_dir, doc_num):
    """Generate a large multi-page document for testing pagination"""
    try:
        pdf = FPDF()
        
        # Generate 50+ pages
        for page in range(1, 52):
            pdf.add_page()
            pdf.set_font('Arial', 'B', 16)
            pdf.cell(0, 15, f'TECHNICAL SPECIFICATION DOCUMENT', 0, 1, 'C')
            pdf.cell(0, 10, f'Page {page} of 51', 0, 1, 'C')
            pdf.ln(10)
            
            pdf.set_font('Arial', 'B', 12)
            pdf.cell(0, 8, f'Section {page}: Technical Requirements', 0, 1)
            pdf.set_font('Arial', '', 10)
            
            # Add content to make it realistic
            for i in range(30):
                content = f'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Page {page}, line {i+1}.'
                if i % 5 == 0:
                    pdf.set_font('Arial', 'B', 10)
                    pdf.cell(0, 6, f'Subsection {page}.{i//5 + 1}', 0, 1)
                    pdf.set_font('Arial', '', 10)
                pdf.cell(0, 5, content, 0, 1)
        
        filename = f'large_document_{doc_num:03d}.pdf'
        pdf.output(os.path.join(output_dir, filename))
        return filename
    except Exception as e:
        print(f"Error generating large document {doc_num}: {e}")
        return None

def generate_password_protected_pdf(output_dir, doc_num):
    """Generate a password-protected PDF (simulated with special filename)"""
    try:
        pdf = FPDF()
        pdf.add_page()
        pdf.set_font('Arial', 'B', 16)
        
        pdf.cell(0, 15, 'CONFIDENTIAL DOCUMENT', 0, 1, 'C')
        pdf.ln(10)
        
        pdf.set_font('Arial', '', 12)
        pdf.cell(0, 8, 'This document contains sensitive information.', 0, 1)
        pdf.cell(0, 8, 'Access requires proper authorization.', 0, 1)
        pdf.ln(5)
        
        pdf.set_font('Arial', 'B', 11)
        pdf.cell(0, 8, 'CLASSIFICATION: RESTRICTED', 0, 1)
        pdf.set_font('Arial', '', 10)
        
        sensitive_content = [
            'Employee Records and Salary Information',
            'Strategic Business Plans for Q1-Q4',
            'Customer Database and Contact Information',
            'Financial Projections and Budget Allocations',
            'Proprietary Technology Specifications'
        ]
        
        for content in sensitive_content:
            pdf.cell(0, 6, f'- {content}', 0, 1)
        
        filename = f'confidential_protected_{doc_num:03d}.pdf'
        pdf.output(os.path.join(output_dir, filename))
        return filename
    except Exception as e:
        print(f"Error generating protected document {doc_num}: {e}")
        return None

def generate_form_pdf(output_dir, form_num):
    """Generate a fillable form PDF"""
    try:
        pdf = FPDF()
        pdf.add_page()
        pdf.set_font('Arial', 'B', 16)
        
        pdf.cell(0, 15, 'EMPLOYEE APPLICATION FORM', 0, 1, 'C')
        pdf.ln(10)
        
        pdf.set_font('Arial', 'B', 12)
        pdf.cell(0, 8, 'PERSONAL INFORMATION', 0, 1)
        pdf.set_font('Arial', '', 10)
        
        # Form fields
        fields = [
            ('Full Name:', '_' * 40),
            ('Email Address:', '_' * 40),
            ('Phone Number:', '_' * 25),
            ('Address:', '_' * 50),
            ('City, State, ZIP:', '_' * 35),
            ('Date of Birth:', '_' * 15),
            ('Social Security Number:', '_' * 15)
        ]
        
        for label, field in fields:
            pdf.cell(60, 8, label, 0, 0)
            pdf.cell(0, 8, field, 0, 1)
            pdf.ln(2)
        
        pdf.ln(5)
        pdf.set_font('Arial', 'B', 12)
        pdf.cell(0, 8, 'EMPLOYMENT HISTORY', 0, 1)
        pdf.set_font('Arial', '', 10)
        
        for i in range(3):
            pdf.cell(0, 6, f'Previous Employer {i+1}:', 0, 1)
            pdf.cell(20, 6, '', 0, 0)
            pdf.cell(50, 6, 'Company Name: ' + '_' * 25, 0, 0)
            pdf.cell(50, 6, 'Position: ' + '_' * 20, 0, 1)
            pdf.cell(20, 6, '', 0, 0)
            pdf.cell(50, 6, 'Start Date: ' + '_' * 12, 0, 0)
            pdf.cell(50, 6, 'End Date: ' + '_' * 12, 0, 1)
            pdf.ln(3)
        
        # Signature section
        pdf.ln(10)
        pdf.set_font('Arial', 'B', 12)
        pdf.cell(0, 8, 'APPLICANT SIGNATURE', 0, 1)
        pdf.set_font('Arial', '', 10)
        pdf.ln(5)
        
        pdf.cell(80, 6, 'Signature: ' + '_' * 30, 0, 0)
        pdf.cell(50, 6, 'Date: ' + '_' * 15, 0, 1)
        
        filename = f'application_form_{form_num:03d}.pdf'
        pdf.output(os.path.join(output_dir, filename))
        return filename
    except Exception as e:
        print(f"Error generating form {form_num}: {e}")
        return None

def generate_financial_statement_pdf(output_dir, stmt_num):
    """Generate a financial statement PDF"""
    try:
        pdf = FPDF()
        pdf.add_page()
        pdf.set_font('Arial', 'B', 16)
        
        companies = ['TechCorp Industries', 'Global Solutions Ltd', 'Innovation Partners Inc', 'Digital Dynamics Corp']
        company = random.choice(companies)
        
        pdf.cell(0, 15, f'{company.upper()}', 0, 1, 'C')
        pdf.cell(0, 10, 'FINANCIAL STATEMENT', 0, 1, 'C')
        
        stmt_date = datetime.now() - timedelta(days=random.randint(30, 365))
        quarter = f'Q{((stmt_date.month - 1) // 3) + 1} {stmt_date.year}'
        pdf.cell(0, 8, f'For the period ending {quarter}', 0, 1, 'C')
        pdf.ln(10)
        
        # Assets section
        pdf.set_font('Arial', 'B', 12)
        pdf.cell(0, 8, 'ASSETS', 0, 1)
        pdf.set_font('Arial', 'B', 10)
        pdf.cell(120, 6, 'Current Assets', 0, 1)
        pdf.set_font('Arial', '', 10)
        
        current_assets = [
            ('Cash and Cash Equivalents', random.randint(500000, 2000000)),
            ('Accounts Receivable', random.randint(300000, 1500000)),
            ('Inventory', random.randint(200000, 800000)),
            ('Prepaid Expenses', random.randint(50000, 200000))
        ]
        
        total_current = 0
        for asset, amount in current_assets:
            total_current += amount
            pdf.cell(100, 5, f'  {asset}', 0, 0)
            pdf.cell(40, 5, f'${amount:,}', 0, 1, 'R')
        
        pdf.set_font('Arial', 'B', 10)
        pdf.cell(100, 6, 'Total Current Assets', 0, 0)
        pdf.cell(40, 6, f'${total_current:,}', 0, 1, 'R')
        pdf.ln(3)
        
        # Fixed Assets
        pdf.cell(120, 6, 'Fixed Assets', 0, 1)
        pdf.set_font('Arial', '', 10)
        
        fixed_assets = [
            ('Property, Plant & Equipment', random.randint(1000000, 5000000)),
            ('Less: Accumulated Depreciation', -random.randint(200000, 800000)),
            ('Intangible Assets', random.randint(500000, 2000000))
        ]
        
        total_fixed = 0
        for asset, amount in fixed_assets:
            total_fixed += amount
            pdf.cell(100, 5, f'  {asset}', 0, 0)
            pdf.cell(40, 5, f'${amount:,}', 0, 1, 'R')
        
        pdf.set_font('Arial', 'B', 10)
        pdf.cell(100, 6, 'Total Fixed Assets', 0, 0)
        pdf.cell(40, 6, f'${total_fixed:,}', 0, 1, 'R')
        pdf.ln(2)
        
        pdf.cell(100, 8, 'TOTAL ASSETS', 1, 0)
        pdf.cell(40, 8, f'${total_current + total_fixed:,}', 1, 1, 'R')
        pdf.ln(5)
        
        # Liabilities section
        pdf.set_font('Arial', 'B', 12)
        pdf.cell(0, 8, 'LIABILITIES & EQUITY', 0, 1)
        pdf.set_font('Arial', 'B', 10)
        pdf.cell(120, 6, 'Current Liabilities', 0, 1)
        pdf.set_font('Arial', '', 10)
        
        liabilities = [
            ('Accounts Payable', random.randint(100000, 500000)),
            ('Short-term Debt', random.randint(200000, 800000)),
            ('Accrued Expenses', random.randint(50000, 200000))
        ]
        
        total_liabilities = 0
        for liability, amount in liabilities:
            total_liabilities += amount
            pdf.cell(100, 5, f'  {liability}', 0, 0)
            pdf.cell(40, 5, f'${amount:,}', 0, 1, 'R')
        
        equity = (total_current + total_fixed) - total_liabilities
        pdf.ln(3)
        pdf.set_font('Arial', 'B', 10)
        pdf.cell(100, 6, 'Total Liabilities', 0, 0)
        pdf.cell(40, 6, f'${total_liabilities:,}', 0, 1, 'R')
        pdf.cell(100, 6, 'Shareholders Equity', 0, 0)
        pdf.cell(40, 6, f'${equity:,}', 0, 1, 'R')
        pdf.ln(2)
        
        pdf.cell(100, 8, 'TOTAL LIABILITIES & EQUITY', 1, 0)
        pdf.cell(40, 8, f'${total_liabilities + equity:,}', 1, 1, 'R')
        
        filename = f'financial_statement_{stmt_num:03d}_{quarter.replace(" ", "_")}.pdf'
        pdf.output(os.path.join(output_dir, filename))
        return filename
    except Exception as e:
        print(f"Error generating financial statement {stmt_num}: {e}")
        return None

def generate_presentation_pdf(output_dir, pres_num):
    """Generate a presentation-style PDF"""
    try:
        pdf = FPDF()
        
        # Title slide
        pdf.add_page()
        pdf.set_font('Arial', 'B', 20)
        pdf.ln(40)
        pdf.cell(0, 20, 'QUARTERLY BUSINESS REVIEW', 0, 1, 'C')
        pdf.set_font('Arial', '', 14)
        pdf.cell(0, 10, f'Presentation #{pres_num}', 0, 1, 'C')
        pdf.cell(0, 10, datetime.now().strftime('%B %Y'), 0, 1, 'C')
        
        # Slide 2 - Agenda
        pdf.add_page()
        pdf.set_font('Arial', 'B', 18)
        pdf.cell(0, 15, 'AGENDA', 0, 1, 'C')
        pdf.ln(10)
        
        pdf.set_font('Arial', '', 12)
        agenda_items = [
            '1. Executive Summary',
            '2. Financial Performance',
            '3. Market Analysis',
            '4. Operational Updates',
            '5. Strategic Initiatives',
            '6. Q&A Session'
        ]
        
        for item in agenda_items:
            pdf.cell(0, 8, item, 0, 1)
            pdf.ln(2)
        
        # Slide 3 - Financial Performance
        pdf.add_page()
        pdf.set_font('Arial', 'B', 18)
        pdf.cell(0, 15, 'FINANCIAL PERFORMANCE', 0, 1, 'C')
        pdf.ln(10)
        
        pdf.set_font('Arial', 'B', 12)
        pdf.cell(0, 8, 'Key Metrics:', 0, 1)
        pdf.set_font('Arial', '', 11)
        
        metrics = [
            f'- Revenue: ${random.randint(5, 50)}M (+{random.randint(5, 25)}% YoY)',
            f'- Gross Margin: {random.randint(35, 65)}%',
            f'- Operating Income: ${random.randint(1, 10)}M',
            f'- Customer Growth: +{random.randint(100, 500)} new customers'
        ]
        
        for metric in metrics:
            pdf.cell(0, 8, metric, 0, 1)
        
        filename = f'presentation_{pres_num:03d}.pdf'
        pdf.output(os.path.join(output_dir, filename))
        return filename
    except Exception as e:
        print(f"Error generating presentation {pres_num}: {e}")
        return None

def generate_legal_document_pdf(output_dir, legal_num):
    """Generate a legal document PDF"""
    try:
        pdf = FPDF()
        pdf.add_page()
        pdf.set_font('Arial', 'B', 14)
        
        pdf.cell(0, 15, 'NON-DISCLOSURE AGREEMENT', 0, 1, 'C')
        pdf.ln(5)
        
        pdf.set_font('Arial', '', 10)
        pdf.cell(0, 6, f'Document #: NDA-{legal_num:04d}', 0, 1)
        pdf.cell(0, 6, f'Date: {datetime.now().strftime("%B %d, %Y")}', 0, 1)
        pdf.ln(8)
        
        pdf.set_font('Arial', 'B', 11)
        pdf.cell(0, 8, 'PARTIES TO THIS AGREEMENT:', 0, 1)
        pdf.set_font('Arial', '', 10)
        
        pdf.cell(0, 6, 'Disclosing Party: Acme Corporation', 0, 1)
        pdf.cell(0, 6, 'Receiving Party: [PARTY NAME]', 0, 1)
        pdf.ln(8)
        
        pdf.set_font('Arial', 'B', 11)
        pdf.cell(0, 8, 'TERMS AND CONDITIONS:', 0, 1)
        pdf.set_font('Arial', '', 10)
        
        legal_terms = [
            '1. CONFIDENTIAL INFORMATION: For purposes of this Agreement, "Confidential',
            '   Information" means all non-public, proprietary information disclosed by the',
            '   Disclosing Party to the Receiving Party.',
            '',
            '2. OBLIGATIONS: The Receiving Party agrees to:',
            '   a) Hold all Confidential Information in strict confidence',
            '   b) Not disclose Confidential Information to third parties',
            '   c) Use Confidential Information solely for evaluation purposes',
            '',
            '3. TERM: This Agreement shall remain in effect for a period of two (2) years',
            '   from the date of execution.',
            '',
            '4. GOVERNING LAW: This Agreement shall be governed by the laws of [STATE].'
        ]
        
        for term in legal_terms:
            if term:
                pdf.cell(0, 5, term, 0, 1)
            else:
                pdf.ln(3)
        
        pdf.ln(15)
        pdf.set_font('Arial', 'B', 10)
        pdf.cell(0, 8, 'SIGNATURES:', 0, 1)
        pdf.ln(8)
        
        pdf.set_font('Arial', '', 9)
        pdf.cell(90, 6, 'Disclosing Party:', 0, 0)
        pdf.cell(90, 6, 'Receiving Party:', 0, 1)
        pdf.ln(12)
        pdf.cell(90, 6, '_' * 30, 0, 0)
        pdf.cell(90, 6, '_' * 30, 0, 1)
        pdf.cell(90, 6, 'Signature', 0, 0)
        pdf.cell(90, 6, 'Signature', 0, 1)
        
        filename = f'nda_{legal_num:04d}.pdf'
        pdf.output(os.path.join(output_dir, filename))
        return filename
    except Exception as e:
        print(f"Error generating legal document {legal_num}: {e}")
        return None

def main():
    """Generate additional mock PDFs"""
    output_dir = create_output_dir()
    print(f"Creating additional mock PDFs in: {output_dir}")
    
    generated_files = []
    
    # Generate different types of documents
    for i in range(1, 3):  # 2 large documents
        filename = generate_large_document_pdf(output_dir, i)
        if filename:
            generated_files.append(filename)
            print(f"Generated: {filename}")
    
    for i in range(1, 4):  # 3 protected documents
        filename = generate_password_protected_pdf(output_dir, i)
        if filename:
            generated_files.append(filename)
            print(f"Generated: {filename}")
    
    for i in range(1, 4):  # 3 forms
        filename = generate_form_pdf(output_dir, i)
        if filename:
            generated_files.append(filename)
            print(f"Generated: {filename}")
    
    for i in range(1, 5):  # 4 financial statements
        filename = generate_financial_statement_pdf(output_dir, i)
        if filename:
            generated_files.append(filename)
            print(f"Generated: {filename}")
    
    for i in range(1, 3):  # 2 presentations
        filename = generate_presentation_pdf(output_dir, i)
        if filename:
            generated_files.append(filename)
            print(f"Generated: {filename}")
    
    for i in range(1, 4):  # 3 legal documents
        filename = generate_legal_document_pdf(output_dir, i)
        if filename:
            generated_files.append(filename)
            print(f"Generated: {filename}")
    
    print(f"\nTotal additional files generated: {len(generated_files)}")
    print(f"Output directory: {output_dir}")

if __name__ == "__main__":
    main()