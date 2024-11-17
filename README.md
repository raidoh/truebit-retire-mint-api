# TrueBit Mint and Retire Calculator

## Overview
A simple calculator for truebit os (truebit.io) token minting and retiring operations.
This is inspired by [lpmythbuster's implementation](https://replit.com/@lpmythbuster/truebit#index.js).

## Getting Started

### Prerequisites
- Go installed on your system
- `.env` file configured

### Running the Application

1. Configure the environment variables in `.env` file
2. Start the server:

   go run cmd/server/main.go

3. Wait for the "Server starting on localhost:XXXX" message

### Using the Calculator

You can access the calculator in two ways:

1. **Via cURL**:

   curl -H "Accept: application/json" http://localhost:XXXX/truebit | jq


2. **Via Web Browser**:

   Simply navigate to `http://localhost:XXXX/truebit` in your preferred browser
