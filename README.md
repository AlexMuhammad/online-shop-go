# Online Shop Fastcampus

### Introduction

The Online Shop Fastcampus API is designed to manage an online shop with essential features such as listing products, processing orders, and managing product inventory. This document outlines the key requirements and functionalities to help track the API development.

### Features

1. Public Features

- List ProductsEndpoint to retrieve a list of all available products.
- Checkout ProductEndpoint to process an order for a selected product.
- Confirm PaymentEndpoint to confirm payment and finalize the order.
- View Detail OrderEndpoint to retrieve detailed information about an order.

2. Protected Features (Admin Only)

- Add ProductEndpoint to add a new product to the catalog (authentication required).
- Update ProductEndpoint to update existing product details (authentication required).
- Delete ProductEndpoint to delete a product from the catalog (authentication required).
