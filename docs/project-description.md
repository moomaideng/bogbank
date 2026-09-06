# Project Description

## Project Name

TBD, not finalized by the team yet.

## Group Members 

- Ashira Aungsumal 6732041921
- Kittichet Arayasujin 6732005321
- Kittichon Chaonawig 6731305421
- Pasin Thanyakasikol 6732025921
- Peeravas Piboolvorakul 6732033921

## Problem Description

Expenses are difficult to track. Categorizing all daily expenses is tedious, and expenses are done on many different channels.

## Target Customers

- College Student

## Scenario (use-case & description)

### Sign In (Register + Login)

Sign in with Google OAuth. First-time sign-in automatically creates a new account (register and login are the same operation).

### Sign Out

Sign out ends the current session.

### Manage bank expense hook

Add / Edit / Remove bank expense hook, which can be chosen from a variety of bank providers (KBank, SCB, etc.). When a bank hook is added, the user grants media-library permission so the native app can watch for new receipt images from that provider (e.g. bank-specific albums or the Screenshots album).
When a new receipt image is detected, the app uploads it and the system extracts metadata. Then notify the user that there's a new expense record, suggest a category, and let the user confirm or choose another category.
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

User can view bar chart for comparing time intervals. The available filters are
- Time interval: every week, month, 3 months, or year.
- Bank provider: all (default) or filter by some bank providers on checkboxes.
- The bar chart plots each interval as 2 bars: income and expense. Each interval option has its own fixed trailing window.

## Functional Requirements

### Authentication

- FR1: The system shall allow users to sign in using Google OAuth.
- FR2: The system shall automatically create a new account on a user's first successful sign-in.
- FR3: The system shall allow a signed-in user to sign out, ending their session.

### Bank Expense Hook
- FR4: The system shall allow a user to add, edit, and remove bank expense hooks, selecting from supported bank providers (e.g., KBank, SCB).
- FR5: For a configured bank hook, after the user grants media-library permission, the system shall detect newly added receipt images associated with that bank provider on the device and upload them for processing.
- FR6: The system shall extract metadata (e.g., amount, merchant, date) from an uploaded receipt image.
- FR7: The system shall notify the user of a new expense record candidate and suggest a category based on the extracted metadata.
- FR8: The system shall allow the user to confirm the suggested category or choose a different one before the expense record is created.

### Income Records
- FR9: The system shall allow a user to add, edit, and remove income records manually.

### Expense Records
- FR10: The system shall allow a user to add, edit, and remove expense records manually, as a fallback when gallery-watch receipt ingestion is not used or fails.

### Category Management
- FR11: The system shall allow a user to add, edit, and remove categories used to classify income and expense records.

### Dashboard
- FR12: The system shall display a donut chart of income or expense records, filterable by type (income/expense), time range (current month, last month, current year, last year, or a custom range), and bank provider (all, or a checkbox-selected subset).
- FR13: The system shall display a bar chart comparing income and expense across a fixed set of time interval options (week, month, 3 months, year), each with a fixed trailing window.

## Non-functional Requirements

### Security & Privacy
- NFR1: The system shall not store user passwords. Authentication shall rely solely on Google OAuth tokens.
- NFR2: The system shall not access the device media gallery until the user explicitly grants media-library permission when configuring a bank hook, and shall ingest only newly detected receipt images after that grant.
- NFR3: The system shall encrypt receipt images and financial records at rest.
- NFR4: The system shall handle personal and financial data in compliance with Thailand's Personal Data Protection Act (PDPA).

### Usability
- NFR5: The system shall provide a native mobile app UI optimized for phones, since target users primarily use phones.
- NFR6: The system shall minimize manual steps for expense confirmation (e.g., single-tap category confirm or change).

### Reliability
- NFR7: The system shall keep manual expense/income entry available at all times as a fallback if gallery-watch receipt ingestion does not apply or fails.

### Performance
- NFR8: The system shall respond to dashboard chart requests within an acceptable time for typical (up to 1-year) date ranges.
