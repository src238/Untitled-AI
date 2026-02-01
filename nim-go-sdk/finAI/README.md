# Liminal AI Sales Finder Agent

## Overview
The **Liminal AI Sales Finder Agent** is an intelligent assistant designed to help Liminal account holders make better purchasing decisions. By analysing a user’s historical transaction data, the agent identifies cheaper or more suitable alternative products and presents clear, actionable recommendations directly within Liminal’s chat interface.

The goal is to reduce unnecessary spending, reduce manual product research, and provide users with confidence that they are making informed and ethical purchasing decisions.

## What Software Are We Creating?
We are designing a **Sales Finder AI Agent** that analyses past purchases made by a Liminal user and determines whether more cost-effective or better-suited alternatives exist.

The system is split into four core components:

### 1. Transaction Analysis
Extracts transaction context and interprets the gathered data from a user’s previous transactions, including product type, category, and price.

### 2. Competitor Finder
The AI agent then scrapes the web to search for potential viable alternative products or providers that come at a lower cost to the account holder while still fulfilling the key requirements / aspects of the original product.
### 3. Product Comparison
Compares the original purchase against the potential alternatives, weighing differences in price, value, drawbacks and potential benefits.

### 4. Chat Output Agent
Sends recommendations and insights to Liminal’s chat system in a clear, readable, and user-friendly format, Providing the alternative, it's price, how much the account holder will save and a link to the product.

## Project Aim
This project is designed to align with Liminal’s values of **user consent, privacy, and transparency**. The agent only analyses authorised transaction data and provides explainable recommendations, empowering account holders to make the most efficient use of their money, especially in a day and age of a bleak cost of living.

## Goals
Our primary goals for this project are:

- Develop a **professional, modular system** that can be adapted or extended by a second-party developer
- Minimise the need for manual product research by users, reducing uncertainty and decision fatigue
- Help users feel confident that their purchases are **cost-effective, reliable, and ethically sourced**
- This project aims to produce a proof of concept for a potential framework that liminal could build upon when creating the foundation for the next age of finances and private banking.

## Features
The Liminal AI Sales Finder Agent currently includes:

- Analysis of recent transaction history
- Alternative product suggestions based on past purchases
- A recent transactions list for contextual recommendations
- A chatbot interface that allows users to communicate directly with the agent through Liminal and perform transactional operations
- Insight generation that explains *why* an alternative may be beneficial

## How Does This Help Customers?
This tool empowers Liminal users to make smarter financial decisions by:

- Identifying cheaper or more suitable alternatives quickly
- Saving time by eliminating tedious manual research
- Increasing transparency around spending choices
- Helping users maximise value while providing a secure, safe and trustworthy system that works for their interests, not against them

Overall, it enables users to spend more efficiently while staying informed and in control of their finances, a service which is now a luxury reserved for the super rich of the modern world.

## Challenges Encountered
During development, we encountered several technical issues, including:

- Difficulties connecting to required APIs
- Network and Wi-Fi instability
- Environment variable misconfiguration (incorrectly naming the `.env` file), which caused significant debugging delays

These challenges provided valuable experience in debugging, configuration management, and development workflows.

## Limitations and Future Improvements
While functional, the current version of the project has limitations:

- Limited time for user testing and feedback integration
- Limited testing with real world data
- Alternative suggestions are based on recent insights rather than long-term learning
- Product matching could be improved with richer datasets and broader competitor coverage
- 

Future improvements could include:

- More advanced pattern analysis over longer transaction histories
- Improved recommendation accuracy using additional data sources
- Enhanced explainability and customisation options for users
- Improved project structure to allow for long term maintenance and development is required
- Adding further capabilities to the filter and search features for transactions