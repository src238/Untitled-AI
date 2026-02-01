# Liminal AI Sales Finder Agent

## Overview
The **Liminal AI Sales Finder Agent** is an intelligent assistant designed to help Liminal account holders make better purchasing decisions. By analysing a user’s historical transaction data, the agent identifies cheaper or more suitable alternative products and presents clear, actionable recommendations directly within Liminal’s chat interface.

The goal is to reduce unnecessary spending, eliminate manual product research, and provide users with confidence that they are making informed and ethical purchasing decisions.

## What Software Are We Creating?
We are designing a **Sales Finder AI Agent** that analyses past purchases made by a Liminal user and determines whether more cost-effective or better-suited alternatives exist.

The system is split into four core components:

### 1. Receipt Reader
Extracts and interprets data from a user’s previous transactions, including product type, category, and price.

### 2. Competitor Finder
Searches for similar products offered by alternative companies or providers.

### 3. Product Comparator
Compares the original purchase against available alternatives, highlighting differences in price, value, and potential benefits.

### 4. Chat Output Agent
Sends recommendations and insights to Liminal’s chat system in a clear, readable, and user-friendly format.

## Project Aim
This project is designed to align with Liminal’s values of **user consent, privacy, and transparency**. The agent only analyses authorised transaction data and provides explainable recommendations.

Additionally, the project explores how **AI-driven analysis can be integrated into blockchain-based banking systems**. By automating sales comparison and pattern analysis, the system demonstrates how artificial intelligence can improve efficiency for both users and financial platforms, allowing greater focus on security, decentralisation, and ethical spending.

## Goals
Our primary goals for this project are:

- Develop a **professional, modular system** that can be adapted or extended by a second-party developer
- Minimise the need for manual product research by users, reducing uncertainty and decision fatigue
- Help users feel confident that their purchases are **cost-effective, reliable, and ethically sourced**
- Gain a deeper understanding of **sales patterns within blockchain-based banking systems** and apply this knowledge to intelligent recommendation logic

## Features
The Liminal AI Sales Finder Agent currently includes:

- Analysis of recent transaction history
- Alternative product suggestions based on past purchases
- A recent transactions list for contextual recommendations
- A chatbot interface that allows users to communicate directly with the agent through Liminal
- Insight generation that explains *why* an alternative may be beneficial

## How Does This Help Customers?
This tool empowers Liminal users to make smarter financial decisions by:

- Identifying cheaper or more suitable alternatives quickly
- Saving time by eliminating manual research
- Increasing transparency around spending choices
- Helping users maximise value while maintaining trust and safety

Overall, it enables users to spend more efficiently while staying informed and in control.

## Challenges Encountered
During development, we encountered several technical issues, including:

- Difficulties connecting to required APIs
- Network and Wi-Fi instability
- Environment variable misconfiguration (incorrectly naming the `.env` file), which caused significant debugging delays

These challenges provided valuable experience in debugging, configuration management, and development workflows.

## Limitations and Future Improvements
While functional, the current version of the project has limitations:

- Limited time for user testing and feedback integration
- Alternative suggestions are based on recent insights rather than long-term learning
- Product matching could be improved with richer datasets and broader competitor coverage

Future improvements could include:

- More advanced pattern analysis over longer transaction histories
- Improved recommendation accuracy using additional data sources
- Enhanced explainability and customisation options for users
