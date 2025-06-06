#!/usr/bin/env python3
"""
Script to generate PDFs with specific dates and metadata for testing date-based filtering
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

def generate_monthly_reports(output_dir):
    """Generate monthly reports for the past 12 months"""
    generated_files = []
    
    for i in range(12):
        try:
            pdf = FPDF()
            pdf.add_page()
            
            # Calculate date for each month
            report_date = datetime.now() - timedelta(days=30 * i)
            month_year = report_date.strftime("%B_%Y")
            
            pdf.set_font('Arial', 'B', 16)
            pdf.cell(0, 15, f'MONTHLY SALES REPORT', 0, 1, 'C')
            pdf.cell(0, 10, f'{report_date.strftime("%B %Y")}', 0, 1, 'C')
            pdf.ln(10)
            
            pdf.set_font('Arial', 'B', 12)
            pdf.cell(0, 8, 'PERFORMANCE SUMMARY', 0, 1)
            pdf.set_font('Arial', '', 10)
            
            # Random sales data
            sales = random.randint(50000, 200000)
            growth = random.randint(-10, 25)
            customers = random.randint(100, 500)
            
            pdf.cell(0, 6, f'Total Sales: ${sales:,}', 0, 1)
            pdf.cell(0, 6, f'Growth Rate: {growth}%', 0, 1)
            pdf.cell(0, 6, f'New Customers: {customers}', 0, 1)
            pdf.cell(0, 6, f'Report Generated: {datetime.now().strftime("%Y-%m-%d")}', 0, 1)
            
            filename = f'monthly_report_{month_year}.pdf'
            pdf.output(os.path.join(output_dir, filename))
            generated_files.append(filename)
            print(f"Generated: {filename}")
            
        except Exception as e:
            print(f"Error generating monthly report {i}: {e}")
    
    return generated_files

def generate_weekly_status_reports(output_dir):
    """Generate weekly status reports for the past 8 weeks"""
    generated_files = []
    
    for i in range(8):
        try:
            pdf = FPDF()
            pdf.add_page()
            
            # Calculate date for each week
            report_date = datetime.now() - timedelta(weeks=i)
            week_str = f"Week_{report_date.strftime('%Y_%m_%d')}"
            
            pdf.set_font('Arial', 'B', 16)
            pdf.cell(0, 15, 'WEEKLY STATUS REPORT', 0, 1, 'C')
            pdf.cell(0, 10, f'Week ending {report_date.strftime("%B %d, %Y")}', 0, 1, 'C')
            pdf.ln(10)
            
            pdf.set_font('Arial', 'B', 12)
            pdf.cell(0, 8, 'PROJECT UPDATES', 0, 1)
            pdf.set_font('Arial', '', 10)
            
            projects = ['Alpha Project', 'Beta Initiative', 'Gamma System', 'Delta Platform']
            for project in random.sample(projects, 3):
                status = random.choice(['On Track', 'Behind Schedule', 'Completed', 'In Review'])
                pdf.cell(0, 6, f'{project}: {status}', 0, 1)
            
            pdf.ln(5)
            pdf.set_font('Arial', 'B', 12)
            pdf.cell(0, 8, 'TEAM METRICS', 0, 1)
            pdf.set_font('Arial', '', 10)
            
            pdf.cell(0, 6, f'Tasks Completed: {random.randint(15, 45)}', 0, 1)
            pdf.cell(0, 6, f'Issues Resolved: {random.randint(5, 20)}', 0, 1)
            pdf.cell(0, 6, f'Team Utilization: {random.randint(75, 95)}%', 0, 1)
            
            filename = f'weekly_status_{week_str}.pdf'
            pdf.output(os.path.join(output_dir, filename))
            generated_files.append(filename)
            print(f"Generated: {filename}")
            
        except Exception as e:
            print(f"Error generating weekly report {i}: {e}")
    
    return generated_files

def generate_quarterly_reviews(output_dir):
    """Generate quarterly review documents for the past 2 years"""
    generated_files = []
    
    quarters = ['Q1', 'Q2', 'Q3', 'Q4']
    current_year = datetime.now().year
    
    for year in [current_year - 1, current_year]:
        for quarter in quarters:
            try:
                pdf = FPDF()
                pdf.add_page()
                
                pdf.set_font('Arial', 'B', 16)
                pdf.cell(0, 15, f'{quarter} {year} QUARTERLY REVIEW', 0, 1, 'C')
                pdf.ln(10)
                
                pdf.set_font('Arial', 'B', 12)
                pdf.cell(0, 8, 'EXECUTIVE SUMMARY', 0, 1)
                pdf.set_font('Arial', '', 10)
                
                revenue = random.randint(1000000, 5000000)
                profit = random.randint(100000, 800000)
                growth = random.randint(-5, 30)
                
                pdf.cell(0, 6, f'Revenue: ${revenue:,}', 0, 1)
                pdf.cell(0, 6, f'Net Profit: ${profit:,}', 0, 1)
                pdf.cell(0, 6, f'YoY Growth: {growth}%', 0, 1)
                pdf.ln(8)
                
                pdf.set_font('Arial', 'B', 12)
                pdf.cell(0, 8, 'KEY ACHIEVEMENTS', 0, 1)
                pdf.set_font('Arial', '', 10)
                
                achievements = [
                    'Successful product launch in new market',
                    'Strategic partnership agreement signed',
                    'Team expansion and talent acquisition',
                    'Technology infrastructure improvements',
                    'Customer satisfaction score improvement'
                ]
                
                for achievement in random.sample(achievements, 3):
                    pdf.cell(0, 6, f'- {achievement}', 0, 1)
                
                filename = f'quarterly_review_{quarter}_{year}.pdf'
                pdf.output(os.path.join(output_dir, filename))
                generated_files.append(filename)
                print(f"Generated: {filename}")
                
            except Exception as e:
                print(f"Error generating quarterly review {quarter} {year}: {e}")
    
    return generated_files

def generate_annual_documents(output_dir):
    """Generate annual documents for the past 3 years"""
    generated_files = []
    
    current_year = datetime.now().year
    doc_types = ['annual_report', 'tax_summary', 'compliance_review']
    
    for year in range(current_year - 2, current_year + 1):
        for doc_type in doc_types:
            try:
                pdf = FPDF()
                pdf.add_page()
                
                title_map = {
                    'annual_report': 'ANNUAL BUSINESS REPORT',
                    'tax_summary': 'ANNUAL TAX SUMMARY',
                    'compliance_review': 'COMPLIANCE REVIEW REPORT'
                }
                
                pdf.set_font('Arial', 'B', 16)
                pdf.cell(0, 15, f'{title_map[doc_type]} {year}', 0, 1, 'C')
                pdf.ln(10)
                
                pdf.set_font('Arial', '', 12)
                pdf.cell(0, 8, f'Fiscal Year: {year}', 0, 1)
                pdf.cell(0, 8, f'Document Type: {doc_type.replace("_", " ").title()}', 0, 1)
                pdf.cell(0, 8, f'Generated: {datetime.now().strftime("%Y-%m-%d")}', 0, 1)
                pdf.ln(8)
                
                if doc_type == 'annual_report':
                    pdf.set_font('Arial', 'B', 12)
                    pdf.cell(0, 8, 'ANNUAL HIGHLIGHTS', 0, 1)
                    pdf.set_font('Arial', '', 10)
                    
                    annual_revenue = random.randint(5000000, 20000000)
                    employees = random.randint(50, 300)
                    markets = random.randint(3, 15)
                    
                    pdf.cell(0, 6, f'Total Annual Revenue: ${annual_revenue:,}', 0, 1)
                    pdf.cell(0, 6, f'Employee Count: {employees}', 0, 1)
                    pdf.cell(0, 6, f'Markets Served: {markets}', 0, 1)
                
                elif doc_type == 'tax_summary':
                    pdf.set_font('Arial', 'B', 12)
                    pdf.cell(0, 8, 'TAX INFORMATION', 0, 1)
                    pdf.set_font('Arial', '', 10)
                    
                    taxable_income = random.randint(1000000, 8000000)
                    tax_paid = int(taxable_income * 0.21)
                    
                    pdf.cell(0, 6, f'Taxable Income: ${taxable_income:,}', 0, 1)
                    pdf.cell(0, 6, f'Federal Tax Paid: ${tax_paid:,}', 0, 1)
                    pdf.cell(0, 6, f'Effective Tax Rate: 21%', 0, 1)
                
                else:  # compliance_review
                    pdf.set_font('Arial', 'B', 12)
                    pdf.cell(0, 8, 'COMPLIANCE STATUS', 0, 1)
                    pdf.set_font('Arial', '', 10)
                    
                    compliance_areas = ['Financial Reporting', 'Data Privacy', 'Environmental', 'Safety Standards']
                    for area in compliance_areas:
                        status = random.choice(['Compliant', 'Under Review', 'Remediation Required'])
                        pdf.cell(0, 6, f'{area}: {status}', 0, 1)
                
                filename = f'{doc_type}_{year}.pdf'
                pdf.output(os.path.join(output_dir, filename))
                generated_files.append(filename)
                print(f"Generated: {filename}")
                
            except Exception as e:
                print(f"Error generating {doc_type} for {year}: {e}")
    
    return generated_files

def main():
    """Generate all dated PDFs"""
    output_dir = create_output_dir()
    print(f"Creating dated mock PDFs in: {output_dir}")
    
    all_generated_files = []
    
    # Generate different time-based document series
    print("\nGenerating monthly reports...")
    monthly_files = generate_monthly_reports(output_dir)
    all_generated_files.extend(monthly_files)
    
    print("\nGenerating weekly status reports...")
    weekly_files = generate_weekly_status_reports(output_dir)
    all_generated_files.extend(weekly_files)
    
    print("\nGenerating quarterly reviews...")
    quarterly_files = generate_quarterly_reviews(output_dir)
    all_generated_files.extend(quarterly_files)
    
    print("\nGenerating annual documents...")
    annual_files = generate_annual_documents(output_dir)
    all_generated_files.extend(annual_files)
    
    print(f"\nTotal dated files generated: {len(all_generated_files)}")
    print(f"Output directory: {output_dir}")
    
    # Print summary by type
    print("\nSummary by document type:")
    print(f"- Monthly reports: {len(monthly_files)}")
    print(f"- Weekly status reports: {len(weekly_files)}")
    print(f"- Quarterly reviews: {len(quarterly_files)}")
    print(f"- Annual documents: {len(annual_files)}")

if __name__ == "__main__":
    main()