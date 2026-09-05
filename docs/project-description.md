# Project Description

## Project Name

หมูไม่เด้ง

## Group Members 

TODO

## Problem Description

Expenses are difficult to track. Categorizing all daily expenses is tedious, and expenses are done on many different channels.

## Target Customers

- College Student

## Scenario (use-case & description)

### Register

Register by Google OAuth, same operation as login?

### Login

Login by Google OAuth

### Logout

Logout..............

### Manage bank expense hook

Add / Edit / Remove bank expense hook, which can be chosen from a variety of bank providers (KBank, SCB, etc.). When a certain bank provider hook is added, the app watches the screenshot on your phone and automatically uploads when there's a new receipt image.
After the receipt image is uploaded to the system, we extract the metadata. Then notify the user that there's a new expense record, suggest a category, and let the user confirm or choose another category.
After the user chooses a category, the system inserts an expense record.

### Manage income record

Add / Edit / Remove income record.

### Manage expense record

Add / Edit / Remove expense record. As a fallback method. 

### Manage category

Add / Edit / Remove category.

### View income/expense dashboard

User can view a donut chart of certain type at specific time range. The available filters are
- Type: income (default) or expense
- Time range: current month (default), last month, current year, last year, and custom time range (any start date and end date range).
- Bank provider: all (default) or filter by some bank providers on checkboxes.

User can view bar chart for comparing many time intervals. The available filters are
- Time interval: every week, month, 3 months, year (should we allow custom?)
- Bank provider: all (default) or filter by some bank providers on checkboxes.
- The bar chart will plot every time interval, each interval has 2 bar: income and expense. It shows at most X interval (up to time interval?).

## Functional Requirements

## Non-functional Requirements

